package models

import (
	"time"

	"github.com/google/uuid"
)

type SyncState struct {
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	DeviceID    string    `db:"device_id" json:"device_id"`
	LastSyncAt  time.Time `db:"last_sync_at" json:"last_sync_at"`
	SyncToken   string    `db:"sync_token" json:"sync_token"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type SyncPullRequest struct {
	DeviceID   string    `json:"device_id" binding:"required"`
	LastSyncAt time.Time `json:"last_sync_at"`
}

type SyncPullResponse struct {
	Accounts     []Account     `json:"accounts"`
	Categories   []Category    `json:"categories"`
	Transactions []Transaction `json:"transactions"`
	CurrentSyncAt time.Time   `json:"current_sync_at"`
}

type SyncPushRequest struct {
	DeviceID     string        `json:"device_id" binding:"required"`
	Accounts     []Account     `json:"accounts"`
	Categories   []Category    `json:"categories"`
	Transactions []Transaction `json:"transactions"`
	LastSyncAt   time.Time     `json:"last_sync_at"`
}

type SyncPushResponse struct {
	Success       bool      `json:"success"`
	CurrentSyncAt time.Time `json:"current_sync_at"`
}
