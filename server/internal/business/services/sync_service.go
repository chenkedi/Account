package services

import (
	"account/internal/business/models"
	"account/internal/data/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SyncService struct {
	syncRepo      *repository.SyncRepository
	accountRepo   *repository.AccountRepository
	categoryRepo  *repository.CategoryRepository
	transactionRepo *repository.TransactionRepository
}

func NewSyncService(
	syncRepo *repository.SyncRepository,
	accountRepo *repository.AccountRepository,
	categoryRepo *repository.CategoryRepository,
	transactionRepo *repository.TransactionRepository,
) *SyncService {
	return &SyncService{
		syncRepo:        syncRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *SyncService) PullChanges(userID uuid.UUID, deviceID string, lastSyncAt time.Time) (*models.SyncPullResponse, error) {
	since := lastSyncAt
	if since.IsZero() {
		since = time.Unix(0, 0)
	}

	changes, err := s.syncRepo.GetChangesSince(userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get changes: %w", err)
	}

	return changes, nil
}

func (s *SyncService) PushChanges(userID uuid.UUID, req *models.SyncPushRequest) (*models.SyncPushResponse, error) {
	if err := s.syncRepo.ApplyChanges(userID, req); err != nil {
		return nil, fmt.Errorf("failed to apply changes: %w", err)
	}

	_, err := s.syncRepo.UpsertSyncState(userID, req.DeviceID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to update sync state: %w", err)
	}

	return &models.SyncPushResponse{
		Success:       true,
		CurrentSyncAt: time.Now().UTC(),
	}, nil
}
