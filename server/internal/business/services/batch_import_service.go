package services

import (
	"encoding/base64"
	"strings"
	"time"

	"account/internal/business/models"
	"account/internal/data/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// BatchImportService handles batch import of multiple files
type BatchImportService struct {
	importService     *ImportService
	accountRepo       *repository.AccountRepository
	transactionRepo   *repository.TransactionRepository
	categoryRepo      *repository.CategoryRepository
	logger            *zap.Logger
}

// NewBatchImportService creates a new BatchImportService
func NewBatchImportService(
	importService *ImportService,
	accountRepo *repository.AccountRepository,
	transactionRepo *repository.TransactionRepository,
	categoryRepo *repository.CategoryRepository,
	logger *zap.Logger,
) *BatchImportService {
	return &BatchImportService{
		importService:     importService,
		accountRepo:       accountRepo,
		transactionRepo:   transactionRepo,
		categoryRepo:      categoryRepo,
		logger:            logger,
	}
}

// FileUpload represents a file to be uploaded
type FileUpload struct {
	Source   string `json:"source"`
	FileName string `json:"file_name"`
	Content  string `json:"content"` // base64 encoded
}

// CreateBatchJob creates a new batch import job
func (s *BatchImportService) CreateBatchJob(userID uuid.UUID, files []FileUpload) (*models.BatchImportJob, error) {
	job := &models.BatchImportJob{
		ID:              uuid.New(),
		UserID:          userID,
		Status:          models.BatchImportStatusPending,
		TotalFiles:      len(files),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Start processing in background
	go s.processBatchJob(job, files)

	return job, nil
}

// processBatchJob processes all files in the batch
func (s *BatchImportService) processBatchJob(job *models.BatchImportJob, files []FileUpload) {
	job.Status = models.BatchImportStatusParsing
	job.UpdatedAt = time.Now()

	var allTransactions []models.ParsedTransaction
	var allAccountHints []models.AccountHint

	// Process each file
	for _, file := range files {
		s.logger.Info("Processing file", zap.String("filename", file.FileName))

		// Decode base64 content
		content, err := base64.StdEncoding.DecodeString(file.Content)
		if err != nil {
			s.logger.Error("Failed to decode file content", zap.Error(err))
			continue
		}

		// Parse file
		source := models.ImportSource(file.Source)
		parseReq := &ParseRequest{
			Source:   source,
			FileName: file.FileName,
			File:     strings.NewReader(string(content)),
		}

		preview, err := s.importService.ParseFile(job.UserID, parseReq)
		if err != nil {
			s.logger.Error("Failed to parse file", zap.Error(err))
			continue
		}

		// Collect transactions
		for _, tx := range preview.Transactions {
			tx.BatchJobID = &job.ID
			allTransactions = append(allTransactions, tx)
		}

		// Extract account hints
		hints := s.extractAccountHints(preview.Transactions, file.FileName)
		allAccountHints = append(allAccountHints, hints...)

		job.ParsedFiles++
	}

	job.TotalTransactions = len(allTransactions)
	job.ValidTransactions = len(allTransactions)

	// Analyze and match
	job.Status = models.BatchImportStatusAnalyzing
	job.UpdatedAt = time.Now()

	// Find transfer matches
	matchResult := s.findTransferMatches(allTransactions)
	job.MatchPairs = len(matchResult.Matches)

	// Auto-create accounts
	job.Status = models.BatchImportStatusMatching
	autoCreated, _, err := s.autoCreateAccounts(job.UserID, allAccountHints)
	if err != nil {
		s.logger.Error("Failed to auto-create accounts", zap.Error(err))
	}
	job.AutoCreatedAccounts = len(autoCreated)

	// Mark as ready to import
	job.Status = models.BatchImportStatusReadyToImport
	job.UpdatedAt = time.Now()

	s.logger.Info("Batch job processing completed",
		zap.String("job_id", job.ID.String()),
		zap.Int("transactions", job.TotalTransactions),
		zap.Int("matches", job.MatchPairs),
		zap.Int("auto_created_accounts", job.AutoCreatedAccounts),
	)
}

// extractAccountHints extracts account hints from transactions
func (s *BatchImportService) extractAccountHints(transactions []models.ParsedTransaction, fileName string) []models.AccountHint {
	var hints []models.AccountHint
	seen := make(map[string]bool)

	for _, tx := range transactions {
		if tx.ParsedAccountNumber == "" && tx.AccountName == "" {
			continue
		}

		key := tx.AccountName + "|" + tx.ParsedAccountNumber
		if seen[key] {
			continue
		}
		seen[key] = true

		hint := models.AccountHint{
			Source:          tx.Source,
			AccountName:     tx.AccountName,
			AccountNumber:   tx.ParsedAccountNumber,
			BankName:        tx.ParsedBankName,
			CardType:        tx.ParsedCardType,
			AccountType:     models.AccountType(tx.ParsedAccountType),
			Balance:         tx.ParsedBalance,
			FoundInFile:     fileName,
		}
		hints = append(hints, hint)
	}

	return hints
}

// findTransferMatches finds transfer matches across all transactions
func (s *BatchImportService) findTransferMatches(transactions []models.ParsedTransaction) *MatchResult {
	var outTxs, inTxs []TransactionInfo

	for i, tx := range transactions {
		fileID := uuid.New() // Generate a dummy file ID for matching
		if tx.BatchFileID != nil {
			fileID = *tx.BatchFileID
		}

		info := TransactionInfo{
			FileID:             fileID,
			TransactionIndex:   i,
			Date:               tx.TransactionDate,
			Amount:             tx.Amount,
			Type:               string(tx.Type),
			AccountName:        tx.AccountName,
			AccountNumber:      tx.ParsedAccountNumber,
			Counterparty:       tx.Counterparty,
			CounterpartyNumber: tx.RelatedAccountNumber,
			Note:               tx.Note,
			RawData:            tx.RawData,
		}

		if isTransferOut(tx) {
			outTxs = append(outTxs, info)
		} else if isTransferIn(tx) {
			inTxs = append(inTxs, info)
		}
	}

	// Use the TransferMatcher to find matches
	matcher := NewTransferMatcher()
	return matcher.FindMatches(convertToFileTransactions(outTxs, inTxs))
}

// isTransferOut checks if a transaction is a transfer out
func isTransferOut(tx models.ParsedTransaction) bool {
	if tx.Type == models.TransactionTypeExpense {
		note := strings.ToLower(tx.Note + " " + tx.Counterparty)
		transferKeywords := []string{"转账", "转出", "汇款", "跨行", "行内"}
		for _, kw := range transferKeywords {
			if strings.Contains(note, kw) {
				return true
			}
		}
	}
	return false
}

// isTransferIn checks if a transaction is a transfer in
func isTransferIn(tx models.ParsedTransaction) bool {
	if tx.Type == models.TransactionTypeIncome {
		note := strings.ToLower(tx.Note + " " + tx.Counterparty)
		inKeywords := []string{"转账", "转入", "汇款", "跨行", "来账"}
		for _, kw := range inKeywords {
			if strings.Contains(note, kw) {
				return true
			}
		}
	}
	return false
}

// convertToFileTransactions converts transaction info to FileTransactions
func convertToFileTransactions(outTxs, inTxs []TransactionInfo) []FileTransactions {
	fileMap := make(map[uuid.UUID]*FileTransactions)

	for _, tx := range outTxs {
		if _, ok := fileMap[tx.FileID]; !ok {
			fileMap[tx.FileID] = &FileTransactions{FileID: tx.FileID}
		}
	}

	for _, tx := range inTxs {
		if _, ok := fileMap[tx.FileID]; !ok {
			fileMap[tx.FileID] = &FileTransactions{FileID: tx.FileID}
		}
	}

	result := make([]FileTransactions, 0, len(fileMap))
	for _, ft := range fileMap {
		result = append(result, *ft)
	}
	return result
}

// autoCreateAccounts automatically creates accounts from hints
func (s *BatchImportService) autoCreateAccounts(userID uuid.UUID, hints []models.AccountHint) ([]models.Account, []models.Account, error) {
	// Get existing accounts
	existingAccounts, err := s.accountRepo.GetAll(userID)
	if err != nil {
		return nil, nil, err
	}

	creator := NewAccountAutoCreator(s.accountRepo)
	return creator.CreateAccountsFromHints(userID, hints, existingAccounts)
}

// ExecuteBatchImport executes the batch import with user selections
func (s *BatchImportService) ExecuteBatchImport(
	jobID uuid.UUID,
	userID uuid.UUID,
	selectedAccountIDs map[string]string,
	confirmedMatches []string,
) (*models.ImportResult, error) {
	result := &models.ImportResult{
		JobID:        jobID,
		ImportedIDs:  make([]uuid.UUID, 0),
		Errors:       make([]models.ImportError, 0),
	}

	// TODO: In a real implementation, we would:
	// 1. Retrieve the parsed transactions from storage
	// 2. Apply user account selections
	// 3. Process confirmed transfer matches
	// 4. Import all transactions with proper account mappings
	// 5. Create transfer transactions for matched pairs

	// For now, return a dummy result
	s.logger.Info("Executing batch import",
		zap.String("job_id", jobID.String()),
		zap.Int("selected_accounts", len(selectedAccountIDs)),
		zap.Int("confirmed_matches", len(confirmedMatches)),
	)

	return result, nil
}

// ImportProgress tracks the progress of a batch import
type ImportProgress struct {
	TotalSteps       int     `json:"total_steps"`
	CurrentStep      int     `json:"current_step"`
	StepDescription  string  `json:"step_description"`
	PercentComplete  float64 `json:"percent_complete"`
}

// CalculateProgress calculates the overall progress of the batch job
func CalculateProgress(job *models.BatchImportJob) ImportProgress {
	steps := map[models.BatchImportStatus]int{
		models.BatchImportStatusPending:          0,
		models.BatchImportStatusParsing:          1,
		models.BatchImportStatusAnalyzing:        2,
		models.BatchImportStatusMatching:         3,
		models.BatchImportStatusReadyToImport:    4,
		models.BatchImportStatusImporting:        5,
		models.BatchImportStatusCompleted:        6,
		models.BatchImportStatusFailed:           0,
	}

	currentStep := steps[job.Status]
	totalSteps := 6

	percentComplete := float64(currentStep) / float64(totalSteps) * 100

	stepDescriptions := map[models.BatchImportStatus]string{
		models.BatchImportStatusPending:          "等待处理",
		models.BatchImportStatusParsing:          "解析文件",
		models.BatchImportStatusAnalyzing:        "分析交易",
		models.BatchImportStatusMatching:         "匹配转账",
		models.BatchImportStatusReadyToImport:    "准备导入",
		models.BatchImportStatusImporting:        "导入数据",
		models.BatchImportStatusCompleted:        "完成",
		models.BatchImportStatusFailed:           "失败",
	}

	return ImportProgress{
		TotalSteps:      totalSteps,
		CurrentStep:     currentStep,
		StepDescription: stepDescriptions[job.Status],
		PercentComplete: percentComplete,
	}
}
