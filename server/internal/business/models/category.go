package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type Category struct {
	ID             uuid.UUID     `db:"id" json:"id"`
	UserID         uuid.UUID     `db:"user_id" json:"user_id"`
	Name           string        `db:"name" json:"name"`
	Type           CategoryType  `db:"type" json:"type"`
	ParentID       *uuid.UUID    `db:"parent_id" json:"parent_id"`
	Icon           string        `db:"icon" json:"icon"`
	CreatedAt      time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `db:"updated_at" json:"updated_at"`
	LastModifiedAt time.Time     `db:"last_modified_at" json:"last_modified_at"`
	Version        int           `db:"version" json:"version"`
	IsDeleted      bool          `db:"is_deleted" json:"is_deleted"`
}

func (c *CategoryType) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid category type: %v", value)
	}
	*c = CategoryType(str)
	return nil
}

func (c CategoryType) Value() (driver.Value, error) {
	return string(c), nil
}

// Default categories
var DefaultIncomeCategories = []string{
	"Salary", "Bonus", "Investment", "Gift", "Refund", "Reimbursement", "Other Income",
}

var DefaultExpenseCategories = []string{
	"Food & Dining", "Groceries", "Transportation", "Shopping", "Entertainment",
	"Healthcare", "Education", "Utilities", "Rent", "Travel", "Subscriptions",
	"Personal Care", "Gifts", "Other Expense",
}
