package models

import (
	"time"

	"github.com/google/uuid"
)

// ImportSource represents the source type of the imported file
type ImportSource string

const (
	ImportSourceAlipay   ImportSource = "alipay"
	ImportSourceWeChat   ImportSource = "wechat"
	ImportSourceBank     ImportSource = "bank"
	ImportSourceGeneric  ImportSource = "generic"
)

// ImportStatus represents the status of an import job
type ImportStatus string

const (
	ImportStatusPending   ImportStatus = "pending"
	ImportStatusParsing   ImportStatus = "parsing"
	ImportStatusPreview   ImportStatus = "preview"
	ImportStatusImporting ImportStatus = "importing"
	ImportStatusCompleted ImportStatus = "completed"
	ImportStatusFailed    ImportStatus = "failed"
)

// ImportJob represents a file import job
type ImportJob struct {
	ID          uuid.UUID       `db:"id" json:"id"`
	UserID      uuid.UUID       `db:"user_id" json:"user_id"`
	Source      ImportSource    `db:"source" json:"source"`
	FileName    string          `db:"file_name" json:"file_name"`
	FileSize    int64           `db:"file_size" json:"file_size"`
	Status      ImportStatus    `db:"status" json:"status"`
	TotalRows   int             `db:"total_rows" json:"total_rows"`
	ImportedRows int            `db:"imported_rows" json:"imported_rows"`
	ErrorMsg    string          `db:"error_msg" json:"error_msg,omitempty"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
}

// ParsedTransaction represents a transaction parsed from an imported file
type ParsedTransaction struct {
	// Raw data from the file
	RawData map[string]string `json:"raw_data"`

	// Parsed fields
	TransactionDate time.Time       `json:"transaction_date"`
	Type            TransactionType `json:"type"`
	Amount          float64         `json:"amount"`
	Currency        string          `json:"currency"`
	Note            string          `json:"note"`

	// Account matching hints
	AccountName     string `json:"account_name,omitempty"`
	AccountNumber   string `json:"account_number,omitempty"`
	Counterparty    string `json:"counterparty,omitempty"`

	// Category matching hints
	CategoryHint    string `json:"category_hint,omitempty"`

	// Metadata
	Source          ImportSource `json:"source"`
	LineNumber      int          `json:"line_number"`
	IsDuplicate     bool         `json:"is_duplicate"`
	CanBeImported   bool         `json:"can_be_imported"`
	ImportWarning   string       `json:"import_warning,omitempty"`

	// Selected account/category for import (set by user in preview)
	SelectedAccountID  *uuid.UUID `json:"selected_account_id,omitempty"`
	SelectedCategoryID *uuid.UUID `json:"selected_category_id,omitempty"`
}

// ImportPreview represents the preview data before actual import
type ImportPreview struct {
	JobID           uuid.UUID           `json:"job_id"`
	Source          ImportSource        `json:"source"`
	TotalRows       int                 `json:"total_rows"`
	ValidRows       int                 `json:"valid_rows"`
	DuplicateRows   int                 `json:"duplicate_rows"`
	Transactions    []ParsedTransaction `json:"transactions"`
	AccountSuggestions map[string][]Account `json:"account_suggestions,omitempty"` // key: account name hint
	Categories      []Category          `json:"categories,omitempty"`
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	JobID           uuid.UUID       `json:"job_id"`
	TotalRows       int             `json:"total_rows"`
	ImportedRows    int             `json:"imported_rows"`
	SkippedRows     int             `json:"skipped_rows"`
	FailedRows      int             `json:"failed_rows"`
	ImportedIDs     []uuid.UUID     `json:"imported_ids,omitempty"`
	Errors          []ImportError    `json:"errors,omitempty"`
}

// ImportError represents an error during import
type ImportError struct {
	LineNumber int    `json:"line_number"`
	Error      string `json:"error"`
}
