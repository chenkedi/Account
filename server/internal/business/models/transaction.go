package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

type Transaction struct {
	ID              uuid.UUID         `db:"id" json:"id"`
	UserID          uuid.UUID         `db:"user_id" json:"user_id"`
	AccountID       uuid.UUID         `db:"account_id" json:"account_id"`
	CategoryID      *uuid.UUID        `db:"category_id" json:"category_id"`
	Type            TransactionType   `db:"type" json:"type"`
	Amount          float64           `db:"amount" json:"amount"`
	Currency        string            `db:"currency" json:"currency"`
	Note            string            `db:"note" json:"note"`
	TransactionDate time.Time         `db:"transaction_date" json:"transaction_date"`
	CreatedAt       time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at" json:"updated_at"`
	LastModifiedAt  time.Time         `db:"last_modified_at" json:"last_modified_at"`
	Version         int               `db:"version" json:"version"`
	IsDeleted       bool              `db:"is_deleted" json:"is_deleted"`
}

func (t *TransactionType) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid transaction type: %v", value)
	}
	*t = TransactionType(str)
	return nil
}

func (t TransactionType) Value() (driver.Value, error) {
	return string(t), nil
}

type TransferLink struct {
	ID                uuid.UUID `db:"id" json:"id"`
	FromTransactionID uuid.UUID `db:"from_transaction_id" json:"from_transaction_id"`
	ToTransactionID   uuid.UUID `db:"to_transaction_id" json:"to_transaction_id"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
}
