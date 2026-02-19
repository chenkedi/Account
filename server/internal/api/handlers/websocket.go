package handlers

import (
	"account/internal/sync"
	"account/pkg/auth"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow any origin for development, restrict in production
		return true
	},
}

type WebSocketHandler struct {
	syncNotifier *sync.SyncNotifier
	tokenMgr     *auth.TokenManager
	logger       *zap.Logger
	mu           sync.Mutex
	connections  map[uuid.UUID]map[string]*WebSocketConnection
}

type WebSocketConnection struct {
	conn     *websocket.Conn
	userID   uuid.UUID
	deviceID string
	send     chan []byte
}

// WebSocket message types
const (
	MessageTypeSyncAvailable = "sync_available"
	MessageTypePing          = "ping"
	MessageTypePong          = "pong"
)

type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type SyncAvailableData struct {
	Timestamp time.Time `json:"timestamp"`
}

func NewWebSocketHandler(
	syncNotifier *sync.SyncNotifier,
	tokenMgr *auth.TokenManager,
	logger *zap.Logger,
) *WebSocketHandler {
	return &WebSocketHandler{
		syncNotifier: syncNotifier,
		tokenMgr:     tokenMgr,
		logger:       logger,
		connections:  make(map[uuid.UUID]map[string]*WebSocketConnection),
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get token from query parameter
	token := c.Query("token")
	deviceID := c.Query("device_id")

	if token == "" || deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token and device_id required"})
		return
	}

	// Validate token
	claims, err := h.tokenMgr.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade to WebSocket", zap.Error(err))
		return
	}

	// Create connection object
	wsConn := &WebSocketConnection{
		conn:     conn,
		userID:   claims.UserID,
		deviceID: deviceID,
		send:     make(chan []byte, 256),
	}

	// Register connection
	h.registerConnection(wsConn)

	// Subscribe to sync notifications
	syncCh := h.syncNotifier.Subscribe(claims.UserID, deviceID)

	// Start goroutines
	go wsConn.writePump()
	go wsConn.readPump(h, syncCh)

	h.logger.Info("WebSocket connection established",
		zap.String("user_id", claims.UserID.String()),
		zap.String("device_id", deviceID),
	)
}

func (h *WebSocketHandler) registerConnection(conn *WebSocketConnection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.connections[conn.userID]; !ok {
		h.connections[conn.userID] = make(map[string]*WebSocketConnection)
	}
	h.connections[conn.userID][conn.deviceID] = conn
}

func (h *WebSocketHandler) unregisterConnection(conn *WebSocketConnection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if userConnections, ok := h.connections[conn.userID]; ok {
		if c, ok := userConnections[conn.deviceID]; ok && c == conn {
			close(conn.send)
			delete(userConnections, conn.deviceID)
		}
		if len(userConnections) == 0 {
			delete(h.connections, conn.userID)
		}
	}

	h.syncNotifier.Unsubscribe(conn.userID, conn.deviceID)
}

func (c *WebSocketConnection) readPump(h *WebSocketHandler, syncCh <-chan struct{}) {
	defer func() {
		h.unregisterConnection(c)
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-syncCh:
			// Send sync available message
			msg := WebSocketMessage{
				Type: MessageTypeSyncAvailable,
				Data: SyncAvailableData{
					Timestamp: time.Now().UTC(),
				},
			}
			data, err := json.Marshal(msg)
			if err == nil {
				c.send <- data
			}
		default:
			// Read message from client
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					h.logger.Error("WebSocket read error", zap.Error(err))
				}
				return
			}

			// Handle message from client
			var wsMsg WebSocketMessage
			if err := json.Unmarshal(message, &wsMsg); err == nil {
				if wsMsg.Type == MessageTypePing {
					// Respond with pong
					pongMsg := WebSocketMessage{Type: MessageTypePong}
					pongData, _ := json.Marshal(pongMsg)
					c.send <- pongData
				}
			}
		}
	}
}

func (c *WebSocketConnection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Channel closed
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Queue pending messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
