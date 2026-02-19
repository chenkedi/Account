package sync

import (
	"account/internal/business/models"
	"account/internal/data/repository"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SyncEngine struct {
	syncRepo      *repository.SyncRepository
	accountRepo   *repository.AccountRepository
	categoryRepo  *repository.CategoryRepository
	transactionRepo *repository.TransactionRepository
	lwwStrategy   *LWWStrategy
	logger        *zap.Logger
	notifier      *SyncNotifier
	mu            sync.RWMutex
}

func NewSyncEngine(
	syncRepo *repository.SyncRepository,
	accountRepo *repository.AccountRepository,
	categoryRepo *repository.CategoryRepository,
	transactionRepo *repository.TransactionRepository,
	logger *zap.Logger,
) *SyncEngine {
	return &SyncEngine{
		syncRepo:        syncRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		transactionRepo: transactionRepo,
		lwwStrategy:     NewLWWStrategy(),
		logger:          logger,
		notifier:        NewSyncNotifier(),
	}
}

// PullResult contains all changes since last sync
type PullResult struct {
	Accounts     []models.Account
	Categories   []models.Category
	Transactions []models.Transaction
	CurrentSyncAt time.Time
}

// PushResult contains the result of a push operation
type PushResult struct {
	Success       bool
	CurrentSyncAt time.Time
}

// Pull gets all changes for a user since their last sync
func (e *SyncEngine) Pull(userID uuid.UUID, deviceID string, lastSyncAt time.Time) (*PullResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	e.logger.Debug("Processing pull request",
		zap.String("user_id", userID.String()),
		zap.String("device_id", deviceID),
		zap.Time("last_sync_at", lastSyncAt),
	)

	since := lastSyncAt
	if since.IsZero() {
		since = time.Unix(0, 0)
	}

	changes, err := e.syncRepo.GetChangesSince(userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get changes: %w", err)
	}

	return &PullResult{
		Accounts:     changes.Accounts,
		Categories:   changes.Categories,
		Transactions: changes.Transactions,
		CurrentSyncAt: changes.CurrentSyncAt,
	}, nil
}

// Push applies changes from a client to the server
func (e *SyncEngine) Push(userID uuid.UUID, req *models.SyncPushRequest) (*PushResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.logger.Debug("Processing push request",
		zap.String("user_id", userID.String()),
		zap.String("device_id", req.DeviceID),
		zap.Int("accounts_count", len(req.Accounts)),
		zap.Int("categories_count", len(req.Categories)),
		zap.Int("transactions_count", len(req.Transactions)),
	)

	// Get current server state to resolve conflicts
	serverChanges, err := e.syncRepo.GetChangesSince(userID, req.LastSyncAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get server state: %w", err)
	}

	// Resolve conflicts using LWW strategy
	resolvedAccounts := e.lwwStrategy.MergeAccountLists(serverChanges.Accounts, req.Accounts)
	resolvedCategories := e.lwwStrategy.MergeCategoryLists(serverChanges.Categories, req.Categories)
	resolvedTransactions := e.lwwStrategy.MergeTransactionLists(serverChanges.Transactions, req.Transactions)

	// Apply the resolved changes
	pushReq := &models.SyncPushRequest{
		DeviceID:     req.DeviceID,
		Accounts:     resolvedAccounts,
		Categories:   resolvedCategories,
		Transactions: resolvedTransactions,
		LastSyncAt:   req.LastSyncAt,
	}

	if err := e.syncRepo.ApplyChanges(userID, pushReq); err != nil {
		return nil, fmt.Errorf("failed to apply changes: %w", err)
	}

	// Update sync state
	now := time.Now().UTC()
	_, err = e.syncRepo.UpsertSyncState(userID, req.DeviceID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to update sync state: %w", err)
	}

	// Notify other devices
	e.notifier.Notify(userID, req.DeviceID)

	return &PushResult{
		Success:       true,
		CurrentSyncAt: now,
	}, nil
}

// GetNotifier returns the sync notifier for WebSocket notifications
func (e *SyncEngine) GetNotifier() *SyncNotifier {
	return e.notifier
}

// SyncNotifier manages real-time sync notifications
type SyncNotifier struct {
	mu       sync.RWMutex
	channels map[uuid.UUID]map[string]chan struct{}
}

func NewSyncNotifier() *SyncNotifier {
	return &SyncNotifier{
		channels: make(map[uuid.UUID]map[string]chan struct{}),
	}
}

// Subscribe registers a device for sync notifications
func (n *SyncNotifier) Subscribe(userID uuid.UUID, deviceID string) <-chan struct{} {
	n.mu.Lock()
	defer n.mu.Unlock()

	if _, ok := n.channels[userID]; !ok {
		n.channels[userID] = make(map[string]chan struct{})
	}

	ch := make(chan struct{}, 1)
	n.channels[userID][deviceID] = ch
	return ch
}

// Unsubscribe removes a device from notifications
func (n *SyncNotifier) Unsubscribe(userID uuid.UUID, deviceID string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if userChannels, ok := n.channels[userID]; ok {
		if ch, ok := userChannels[deviceID]; ok {
			close(ch)
			delete(userChannels, deviceID)
		}
		if len(userChannels) == 0 {
			delete(n.channels, userID)
		}
	}
}

// Notify sends a sync notification to all other devices of the same user
func (n *SyncNotifier) Notify(userID uuid.UUID, excludeDeviceID string) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if userChannels, ok := n.channels[userID]; ok {
		for deviceID, ch := range userChannels {
			if deviceID != excludeDeviceID {
				select {
				case ch <- struct{}{}:
				default:
					// Channel full, skip
				}
			}
		}
	}
}
