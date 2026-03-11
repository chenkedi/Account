package models

import (
	"time"

	"github.com/google/uuid"
)

// BatchImportStatus represents the status of a batch import job
type BatchImportStatus string

const (
	BatchImportStatusPending        BatchImportStatus = "pending"
	BatchImportStatusParsing        BatchImportStatus = "parsing"
	BatchImportStatusAnalyzing      BatchImportStatus = "analyzing"
	BatchImportStatusMatching       BatchImportStatus = "matching"
	BatchImportStatusReadyToImport  BatchImportStatus = "ready_to_import"
	BatchImportStatusImporting      BatchImportStatus = "importing"
	BatchImportStatusCompleted      BatchImportStatus = "completed"
	BatchImportStatusFailed         BatchImportStatus = "failed"
)

// FileImportStatus represents the status of a single file import
type FileImportStatus string

const (
	FileImportStatusPending   FileImportStatus = "pending"
	FileImportStatusParsing   FileImportStatus = "parsing"
	FileImportStatusParsed    FileImportStatus = "parsed"
	FileImportStatusAnalyzing FileImportStatus = "analyzing"
	FileImportStatusAnalyzed  FileImportStatus = "analyzed"
	FileImportStatusImported  FileImportStatus = "imported"
	FileImportStatusFailed    FileImportStatus = "failed"
)

// MatchType represents the type of match found
type MatchType string

const (
	MatchTypeTransfer   MatchType = "transfer"
	MatchTypeNoteMerge  MatchType = "note_merge"
	MatchTypeAccountLink MatchType = "account_link"
)

// BatchImportJob represents a batch import job for multiple files
type BatchImportJob struct {
	ID                  uuid.UUID         `db:"id" json:"id"`
	UserID              uuid.UUID         `db:"user_id" json:"user_id"`
	Status              BatchImportStatus `db:"status" json:"status"`
	TotalFiles          int               `db:"total_files" json:"total_files"`
	ParsedFiles         int               `db:"parsed_files" json:"parsed_files"`
	TotalTransactions   int               `db:"total_transactions" json:"total_transactions"`
	ValidTransactions   int               `db:"valid_transactions" json:"valid_transactions"`
	MatchPairs          int               `db:"match_pairs" json:"match_pairs"`
	AutoCreatedAccounts int               `db:"auto_created_accounts" json:"auto_created_accounts"`
	ErrorMsg            string            `db:"error_msg" json:"error_msg,omitempty"`
	CreatedAt           time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time         `db:"updated_at" json:"updated_at"`
}

// BatchImportFile represents a single file within a batch import job
type BatchImportFile struct {
	ID              uuid.UUID         `db:"id" json:"id"`
	JobID           uuid.UUID         `db:"job_id" json:"job_id"`
	Source          ImportSource      `db:"source" json:"source"`
	FileName        string            `db:"file_name" json:"file_name"`
	Status          FileImportStatus  `db:"status" json:"status"`
	ParsedContent   []ParsedTransaction `db:"parsed_content" json:"parsed_content"`
	AccountHints    []AccountHint     `db:"account_hints" json:"account_hints"`
	ParseErrors     []string          `db:"parse_errors" json:"parse_errors,omitempty"`
	CreatedAt       time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at" json:"updated_at"`
}

// AccountHint represents account information extracted from a parsed file
type AccountHint struct {
	Source          ImportSource `json:"source"`
	AccountName     string       `json:"account_name"`
	AccountNumber   string       `json:"account_number"`
	BankName        string       `json:"bank_name"`
	CardType        string       `json:"card_type"`
	AccountType     AccountType  `json:"account_type"`
	Balance         float64      `json:"balance"`
	FoundInFile     string       `json:"found_in_file"`
}

// TransferMatch represents a matched transfer between two transactions
type TransferMatch struct {
	ID              uuid.UUID           `json:"id"`
	JobID           uuid.UUID           `json:"job_id"`
	FromFileID      uuid.UUID           `json:"from_file_id"`
	FromTransaction ParsedTransaction   `json:"from_transaction"`
	ToFileID        uuid.UUID           `json:"to_file_id"`
	ToTransaction   ParsedTransaction   `json:"to_transaction"`
	MatchType       MatchType           `json:"match_type"`
	Confidence      float64             `json:"confidence"`
	MatchFactors    []string            `json:"match_factors"`
	UserConfirmed   bool                `json:"user_confirmed"`
	CreatedAt       time.Time           `json:"created_at"`
}

// NoteMergeSuggestion represents a suggestion for merging notes
type NoteMergeSuggestion struct {
	PrimaryNote     string   `json:"primary_note"`
	SecondaryNote   string   `json:"secondary_note"`
	MergedNote      string   `json:"merged_note"`
	MergeStrategy   string   `json:"merge_strategy"`
	Reason          string   `json:"reason"`
}
