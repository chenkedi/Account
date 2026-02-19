package handlers

import (
	"account/internal/business/models"
	"account/internal/business/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SyncHandler struct {
	syncService *services.SyncService
	logger      *zap.Logger
}

func NewSyncHandler(syncService *services.SyncService, logger *zap.Logger) *SyncHandler {
	return &SyncHandler{
		syncService: syncService,
		logger:      logger,
	}
}

func (h *SyncHandler) Pull(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.SyncPullRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DeviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}

	changes, err := h.syncService.PullChanges(userID, req.DeviceID, req.LastSyncAt)
	if err != nil {
		h.logger.Error("Failed to pull changes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, changes)
}

func (h *SyncHandler) Push(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.SyncPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DeviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}

	response, err := h.syncService.PushChanges(userID, &req)
	if err != nil {
		h.logger.Error("Failed to push changes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}
