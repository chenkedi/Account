package repository

import (
	"account/internal/business/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SyncRepository struct {
	db             *sqlx.DB
	accountRepo    *AccountRepository
	categoryRepo   *CategoryRepository
	transactionRepo *TransactionRepository
}

func NewSyncRepository(
	db *sqlx.DB,
	accountRepo *AccountRepository,
	categoryRepo *CategoryRepository,
	transactionRepo *TransactionRepository,
) *SyncRepository {
	return &SyncRepository{
		db:              db,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		transactionRepo: transactionRepo,
	}
}

func (r *SyncRepository) GetSyncState(userID uuid.UUID, deviceID string) (*models.SyncState, error) {
	var syncState models.SyncState

	query := `
		SELECT user_id, device_id, last_sync_at, sync_token, created_at, updated_at
		FROM sync_state
		WHERE user_id = $1 AND device_id = $2
	`

	err := r.db.Get(&syncState, query, userID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sync state: %w", err)
	}

	return &syncState, nil
}

func (r *SyncRepository) UpsertSyncState(userID uuid.UUID, deviceID string, syncToken string) (*models.SyncState, error) {
	now := time.Now().UTC()
	syncState := &models.SyncState{
		UserID:     userID,
		DeviceID:   deviceID,
		LastSyncAt: now,
		SyncToken:  syncToken,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	query := `
		INSERT INTO sync_state (user_id, device_id, last_sync_at, sync_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, device_id) DO UPDATE
		SET last_sync_at = EXCLUDED.last_sync_at,
		    sync_token = EXCLUDED.sync_token,
		    updated_at = EXCLUDED.updated_at
		RETURNING created_at
	`

	err := r.db.QueryRow(query,
		syncState.UserID, syncState.DeviceID, syncState.LastSyncAt, syncState.SyncToken,
		syncState.CreatedAt, syncState.UpdatedAt,
	).Scan(&syncState.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert sync state: %w", err)
	}

	return syncState, nil
}

func (r *SyncRepository) GetChangesSince(userID uuid.UUID, since time.Time) (*models.SyncPullResponse, error) {
	accounts, err := r.accountRepo.GetModifiedSince(userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get account changes: %w", err)
	}

	categories, err := r.categoryRepo.GetModifiedSince(userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get category changes: %w", err)
	}

	transactions, err := r.transactionRepo.GetModifiedSince(userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction changes: %w", err)
	}

	return &models.SyncPullResponse{
		Accounts:     accounts,
		Categories:   categories,
		Transactions: transactions,
		CurrentSyncAt: time.Now().UTC(),
	}, nil
}

func (r *SyncRepository) ApplyChanges(userID uuid.UUID, req *models.SyncPushRequest) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if len(req.Accounts) > 0 {
		if err := r.accountRepo.CreateMany(req.Accounts); err != nil {
			return fmt.Errorf("failed to apply account changes: %w", err)
		}
	}

	if len(req.Categories) > 0 {
		if err := r.categoryRepo.CreateMany(req.Categories); err != nil {
			return fmt.Errorf("failed to apply category changes: %w", err)
		}
	}

	if len(req.Transactions) > 0 {
		if err := r.transactionRepo.CreateMany(req.Transactions); err != nil {
			return fmt.Errorf("failed to apply transaction changes: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
