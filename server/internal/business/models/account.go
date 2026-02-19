package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeBank     AccountType = "bank"
	AccountTypeCash     AccountType = "cash"
	AccountTypeAlipay   AccountType = "alipay"
	AccountTypeWeChat   AccountType = "wechat"
	AccountTypeCredit   AccountType = "credit"
	AccountTypeInvestment AccountType = "investment"
	AccountTypeOther    AccountType = "other"
)

type Account struct {
	ID             uuid.UUID   `db:"id" json:"id"`
	UserID         uuid.UUID   `db:"user_id" json:"user_id"`
	Name           string      `db:"name" json:"name"`
	Type           AccountType `db:"type" json:"type"`
	Currency       string      `db:"currency" json:"currency"`
	Balance        float64     `db:"balance" json:"balance"`
	CreatedAt      time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at" json:"updated_at"`
	LastModifiedAt time.Time   `db:"last_modified_at" json:"last_modified_at"`
	Version        int         `db:"version" json:"version"`
	IsDeleted      bool        `db:"is_deleted" json:"is_deleted"`
}

func (a *AccountType) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid account type: %v", value)
	}
	*a = AccountType(str)
	return nil
}

func (a AccountType) Value() (driver.Value, error) {
	return string(a), nil
}
