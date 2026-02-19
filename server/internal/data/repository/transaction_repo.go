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
	ErrTransactionNotFound = errors.New("transaction not found")
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(
	userID uuid.UUID,
	accountID uuid.UUID,
	categoryID *uuid.UUID,
	transactionType models.TransactionType,
	amount float64,
	currency string,
	note string,
	transactionDate time.Time,
) (*models.Transaction, error) {
	now := time.Now().UTC()
	transaction := &models.Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		AccountID:       accountID,
		CategoryID:      categoryID,
		Type:            transactionType,
		Amount:          amount,
		Currency:        currency,
		Note:            note,
		TransactionDate: transactionDate.UTC(),
		CreatedAt:       now,
		UpdatedAt:       now,
		LastModifiedAt:  now,
		Version:         1,
		IsDeleted:       false,
	}

	query := `
		INSERT INTO transactions (id, user_id, account_id, category_id, type, amount, currency, note, transaction_date, created_at, updated_at, last_modified_at, version, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.Exec(query,
		transaction.ID, transaction.UserID, transaction.AccountID, transaction.CategoryID,
		transaction.Type, transaction.Amount, transaction.Currency, transaction.Note,
		transaction.TransactionDate, transaction.CreatedAt, transaction.UpdatedAt,
		transaction.LastModifiedAt, transaction.Version, transaction.IsDeleted,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

func (r *TransactionRepository) GetByID(id uuid.UUID, userID uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction

	query := `
		SELECT id, user_id, account_id, category_id, type, amount, currency, note, transaction_date, created_at, updated_at, last_modified_at, version, is_deleted
		FROM transactions
		WHERE id = $1 AND user_id = $2 AND is_deleted = false
	`

	err := r.db.Get(&transaction, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

func (r *TransactionRepository) GetAll(userID uuid.UUID, limit int, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT id, user_id, account_id, category_id, type, amount, currency, note, transaction_date, created_at, updated_at, last_modified_at, version, is_deleted
		FROM transactions
		WHERE user_id = $1 AND is_deleted = false
		ORDER BY transaction_date DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.Select(&transactions, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetByDateRange(userID uuid.UUID, start, end time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := `
		SELECT id, user_id, account_id, category_id, type, amount, currency, note, transaction_date, created_at, updated_at, last_modified_at, version, is_deleted
		FROM transactions
		WHERE user_id = $1 AND is_deleted = false
		AND transaction_date >= $2 AND transaction_date <= $3
		ORDER BY transaction_date DESC, created_at DESC
	`

	err := r.db.Select(&transactions, query, userID, start.UTC(), end.UTC())
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by date range: %w", err)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetByAccount(userID uuid.UUID, accountID uuid.UUID, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT id, user_id, account_id, category_id, type, amount, currency, note, transaction_date, created_at, updated_at, last_modified_at, version, is_deleted
		FROM transactions
		WHERE user_id = $1 AND account_id = $2 AND is_deleted = false
		ORDER BY transaction_date DESC, created_at DESC
		LIMIT $3
	`

	err := r.db.Select(&transactions, query, userID, accountID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by account: %w", err)
	}

	return transactions, nil
}

func (r *TransactionRepository) Update(transaction *models.Transaction, userID uuid.UUID) (*models.Transaction, error) {
	now := time.Now().UTC()
	transaction.UpdatedAt = now
	transaction.LastModifiedAt = now
	transaction.Version++

	query := `
		UPDATE transactions
		SET account_id = $1, category_id = $2, type = $3, amount = $4, currency = $5, note = $6, transaction_date = $7, updated_at = $8, last_modified_at = $9, version = $10
		WHERE id = $11 AND user_id = $12 AND is_deleted = false
	`

	result, err := r.db.Exec(query,
		transaction.AccountID, transaction.CategoryID, transaction.Type, transaction.Amount,
		transaction.Currency, transaction.Note, transaction.TransactionDate,
		transaction.UpdatedAt, transaction.LastModifiedAt, transaction.Version,
		transaction.ID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return nil, ErrTransactionNotFound
	}

	return transaction, nil
}

func (r *TransactionRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	now := time.Now().UTC()

	query := `
		UPDATE transactions
		SET is_deleted = true, updated_at = $1, last_modified_at = $2, version = version + 1
		WHERE id = $3 AND user_id = $4 AND is_deleted = false
	`

	result, err := r.db.Exec(query, now, now, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

func (r *TransactionRepository) GetModifiedSince(userID uuid.UUID, since time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := `
		SELECT id, user_id, account_id, category_id, type, amount, currency, note, transaction_date, created_at, updated_at, last_modified_at, version, is_deleted
		FROM transactions
		WHERE user_id = $1 AND last_modified_at > $2
	`

	err := r.db.Select(&transactions, query, userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get modified transactions: %w", err)
	}

	return transactions, nil
}

func (r *TransactionRepository) CreateMany(transactions []models.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	query := `
		INSERT INTO transactions (id, user_id, account_id, category_id, type, amount, currency, note, transaction_date, created_at, updated_at, last_modified_at, version, is_deleted)
		VALUES (:id, :user_id, :account_id, :category_id, :type, :amount, :currency, :note, :transaction_date, :created_at, :updated_at, :last_modified_at, :version, :is_deleted)
		ON CONFLICT (id) DO UPDATE
		SET account_id = EXCLUDED.account_id,
		    category_id = EXCLUDED.category_id,
		    type = EXCLUDED.type,
		    amount = EXCLUDED.amount,
		    currency = EXCLUDED.currency,
		    note = EXCLUDED.note,
		    transaction_date = EXCLUDED.transaction_date,
		    updated_at = EXCLUDED.updated_at,
		    last_modified_at = EXCLUDED.last_modified_at,
		    version = EXCLUDED.version,
		    is_deleted = EXCLUDED.is_deleted
		WHERE transactions.last_modified_at <= EXCLUDED.last_modified_at
	`

	for _, transaction := range transactions {
		_, err := tx.NamedExec(query, transaction)
		if err != nil {
			return fmt.Errorf("failed to insert transaction %s: %w", transaction.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *TransactionRepository) GetStatsByDateRange(userID uuid.UUID, start, end time.Time) (incomeTotal float64, expenseTotal float64, err error) {
	query := `
		SELECT type, SUM(amount) as total
		FROM transactions
		WHERE user_id = $1 AND is_deleted = false
		AND type IN ('income', 'expense')
		AND transaction_date >= $2 AND transaction_date <= $3
		GROUP BY type
	`

	rows, err := r.db.Query(query, userID, start.UTC(), end.UTC())
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t string
		var total float64
		if err := rows.Scan(&t, &total); err != nil {
			return 0, 0, fmt.Errorf("failed to scan stats: %w", err)
		}
		if t == "income" {
			incomeTotal = total
		} else if t == "expense" {
			expenseTotal = total
		}
	}

	return incomeTotal, expenseTotal, nil
}

// CategoryStats represents statistics for a single category
type CategoryStats struct {
	CategoryID   uuid.UUID `db:"category_id"`
	CategoryName string    `db:"category_name"`
	CategoryType string    `db:"category_type"`
	TotalAmount  float64   `db:"total_amount"`
	TransactionCount int   `db:"transaction_count"`
}

// GetCategoryStats gets statistics grouped by category
func (r *TransactionRepository) GetCategoryStats(userID uuid.UUID, start, end time.Time) ([]CategoryStats, error) {
	var stats []CategoryStats

	query := `
		SELECT
			c.id as category_id,
			c.name as category_name,
			c.type as category_type,
			SUM(t.amount) as total_amount,
			COUNT(t.id) as transaction_count
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = $1
			AND t.is_deleted = false
			AND t.type IN ('income', 'expense')
			AND t.transaction_date >= $2
			AND t.transaction_date <= $3
		GROUP BY c.id, c.name, c.type
		ORDER BY total_amount DESC
	`

	err := r.db.Select(&stats, query, userID, start.UTC(), end.UTC())
	if err != nil {
		return nil, fmt.Errorf("failed to get category stats: %w", err)
	}

	return stats, nil
}

// MonthlyStats represents statistics for a single month
type MonthlyStats struct {
	Year         int     `db:"year"`
	Month        int     `db:"month"`
	IncomeTotal  float64 `db:"income_total"`
	ExpenseTotal float64 `db:"expense_total"`
}

// GetMonthlyTrend gets monthly income/expense trend
func (r *TransactionRepository) GetMonthlyTrend(userID uuid.UUID, start, end time.Time) ([]MonthlyStats, error) {
	var stats []MonthlyStats

	query := `
		SELECT
			EXTRACT(YEAR FROM transaction_date)::int as year,
			EXTRACT(MONTH FROM transaction_date)::int as month,
			SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as income_total,
			SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as expense_total
		FROM transactions
		WHERE user_id = $1
			AND is_deleted = false
			AND type IN ('income', 'expense')
			AND transaction_date >= $2
			AND transaction_date <= $3
		GROUP BY year, month
		ORDER BY year, month
	`

	err := r.db.Select(&stats, query, userID, start.UTC(), end.UTC())
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly trend: %w", err)
	}

	return stats, nil
}

// GetByDateRangeAndAmount finds transactions within a date range and amount range (for duplicate detection)
func (r *TransactionRepository) GetByDateRangeAndAmount(
	userID uuid.UUID,
	start, end time.Time,
	minAmount, maxAmount float64,
) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := `
		SELECT id, user_id, account_id, category_id, type, amount, currency, note, transaction_date, created_at, updated_at, last_modified_at, version, is_deleted
		FROM transactions
		WHERE user_id = $1 AND is_deleted = false
		AND transaction_date >= $2 AND transaction_date <= $3
		AND amount >= $4 AND amount <= $5
		ORDER BY transaction_date DESC, created_at DESC
	`

	err := r.db.Select(&transactions, query, userID, start.UTC(), end.UTC(), minAmount, maxAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for duplicate check: %w", err)
	}

	return transactions, nil
}
