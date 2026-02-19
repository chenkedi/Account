package repository

import (
	"account/internal/business/models"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(userID uuid.UUID, name string, accountType models.AccountType, currency string, balance float64) (*models.Account, error) {
	now := time.Now().UTC()
	account := &models.Account{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           name,
		Type:           accountType,
		Currency:       currency,
		Balance:        balance,
		CreatedAt:      now,
		UpdatedAt:      now,
		LastModifiedAt: now,
		Version:        1,
		IsDeleted:      false,
	}

	query := `
		INSERT INTO accounts (id, user_id, name, type, currency, balance, created_at, updated_at, last_modified_at, version, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(query,
		account.ID, account.UserID, account.Name, account.Type, account.Currency,
		account.Balance, account.CreatedAt, account.UpdatedAt, account.LastModifiedAt,
		account.Version, account.IsDeleted,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (r *AccountRepository) GetByID(id uuid.UUID, userID uuid.UUID) (*models.Account, error) {
	var account models.Account

	query := `
		SELECT id, user_id, name, type, currency, balance, created_at, updated_at, last_modified_at, version, is_deleted
		FROM accounts
		WHERE id = $1 AND user_id = $2 AND is_deleted = false
	`

	err := r.db.Get(&account, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

func (r *AccountRepository) GetAll(userID uuid.UUID) ([]models.Account, error) {
	var accounts []models.Account

	query := `
		SELECT id, user_id, name, type, currency, balance, created_at, updated_at, last_modified_at, version, is_deleted
		FROM accounts
		WHERE user_id = $1 AND is_deleted = false
		ORDER BY name ASC
	`

	err := r.db.Select(&accounts, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, nil
}

func (r *AccountRepository) Update(account *models.Account, userID uuid.UUID) (*models.Account, error) {
	now := time.Now().UTC()
	account.UpdatedAt = now
	account.LastModifiedAt = now
	account.Version++

	query := `
		UPDATE accounts
		SET name = $1, type = $2, currency = $3, balance = $4, updated_at = $5, last_modified_at = $6, version = $7
		WHERE id = $8 AND user_id = $9 AND is_deleted = false
	`

	result, err := r.db.Exec(query,
		account.Name, account.Type, account.Currency, account.Balance,
		account.UpdatedAt, account.LastModifiedAt, account.Version,
		account.ID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (r *AccountRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	now := time.Now().UTC()

	query := `
		UPDATE accounts
		SET is_deleted = true, updated_at = $1, last_modified_at = $2, version = version + 1
		WHERE id = $3 AND user_id = $4 AND is_deleted = false
	`

	result, err := r.db.Exec(query, now, now, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return ErrAccountNotFound
	}

	return nil
}

func (r *AccountRepository) GetModifiedSince(userID uuid.UUID, since time.Time) ([]models.Account, error) {
	var accounts []models.Account

	query := `
		SELECT id, user_id, name, type, currency, balance, created_at, updated_at, last_modified_at, version, is_deleted
		FROM accounts
		WHERE user_id = $1 AND last_modified_at > $2
	`

	err := r.db.Select(&accounts, query, userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get modified accounts: %w", err)
	}

	return accounts, nil
}

func (r *AccountRepository) CreateMany(accounts []models.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	query := `
		INSERT INTO accounts (id, user_id, name, type, currency, balance, created_at, updated_at, last_modified_at, version, is_deleted)
		VALUES (:id, :user_id, :name, :type, :currency, :balance, :created_at, :updated_at, :last_modified_at, :version, :is_deleted)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    type = EXCLUDED.type,
		    currency = EXCLUDED.currency,
		    balance = EXCLUDED.balance,
		    updated_at = EXCLUDED.updated_at,
		    last_modified_at = EXCLUDED.last_modified_at,
		    version = EXCLUDED.version,
		    is_deleted = EXCLUDED.is_deleted
		WHERE accounts.last_modified_at <= EXCLUDED.last_modified_at
	`

	for _, account := range accounts {
		_, err := tx.NamedExec(query, account)
		if err != nil {
			return fmt.Errorf("failed to insert account %s: %w", account.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
