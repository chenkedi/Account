package services

import (
	"account/internal/business/models"
	"account/internal/data/repository"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ImportService handles bill import operations
type ImportService struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
	categoryRepo    *repository.CategoryRepository
	logger          *zap.Logger
}

// NewImportService creates a new ImportService
func NewImportService(
	transactionRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	categoryRepo *repository.CategoryRepository,
	logger *zap.Logger,
) *ImportService {
	return &ImportService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		logger:          logger,
	}
}

// ParseRequest contains the parameters for parsing a file
type ParseRequest struct {
	Source   models.ImportSource `json:"source" binding:"required"`
	FileName string               `json:"file_name"`
	File     io.Reader            `json:"-"`
}

// ParseFile parses an uploaded file and returns preview data
func (s *ImportService) ParseFile(userID uuid.UUID, req *ParseRequest) (*models.ImportPreview, error) {
	var transactions []models.ParsedTransaction
	var err error

	switch req.Source {
	case models.ImportSourceAlipay:
		transactions, err = s.parseAlipayCSV(req.File)
	case models.ImportSourceWeChat:
		transactions, err = s.parseWeChatCSV(req.File)
	case models.ImportSourceBank:
		transactions, err = s.parseBankCSV(req.File)
	case models.ImportSourceGeneric:
		transactions, err = s.parseGenericCSV(req.File)
	default:
		return nil, fmt.Errorf("unsupported import source: %s", req.Source)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	// Enhance transactions with account/category suggestions
	validCount := 0
	duplicateCount := 0
	for i := range transactions {
		tx := &transactions[i]
		tx.Source = req.Source
		tx.LineNumber = i + 1

		// Check for duplicates
		isDup, err := s.checkForDuplicate(userID, tx)
		if err != nil {
			s.logger.Warn("Failed to check duplicate", zap.Error(err))
		}
		tx.IsDuplicate = isDup

		// Determine if can be imported
		tx.CanBeImported = !tx.IsDuplicate && tx.Amount > 0 && !tx.TransactionDate.IsZero()

		if tx.IsDuplicate {
			duplicateCount++
		}
		if tx.CanBeImported {
			validCount++
		}
	}

	// Get user's accounts and categories for suggestions
	accounts, _ := s.accountRepo.GetAll(userID)
	categories, _ := s.categoryRepo.GetAll(userID)

	accountSuggestions := make(map[string][]models.Account)
	for _, tx := range transactions {
		if tx.AccountName != "" {
			// Find matching accounts
			matches := s.findMatchingAccounts(tx.AccountName, accounts)
			if len(matches) > 0 {
				accountSuggestions[tx.AccountName] = matches
			}
		}
	}

	preview := &models.ImportPreview{
		JobID:             uuid.New(),
		Source:            req.Source,
		TotalRows:         len(transactions),
		ValidRows:         validCount,
		DuplicateRows:     duplicateCount,
		Transactions:      transactions,
		AccountSuggestions: accountSuggestions,
		Categories:        categories,
	}

	return preview, nil
}

// ExecuteImportRequest contains the data needed to execute an import
type ExecuteImportRequest struct {
	JobID        uuid.UUID               `json:"job_id" binding:"required"`
	Transactions []models.ParsedTransaction `json:"transactions" binding:"required"`
}

// ExecuteImport executes the actual import of transactions
func (s *ImportService) ExecuteImport(userID uuid.UUID, req *ExecuteImportRequest) (*models.ImportResult, error) {
	result := &models.ImportResult{
		JobID:        req.JobID,
		TotalRows:    len(req.Transactions),
		ImportedIDs:  make([]uuid.UUID, 0),
		Errors:       make([]models.ImportError, 0),
	}

	for _, parsedTx := range req.Transactions {
		if !parsedTx.CanBeImported || parsedTx.SelectedAccountID == nil {
			result.SkippedRows++
			continue
		}

		// Create transaction
		createReq := &CreateTransactionRequest{
			AccountID:       *parsedTx.SelectedAccountID,
			CategoryID:      parsedTx.SelectedCategoryID,
			Type:            parsedTx.Type,
			Amount:          parsedTx.Amount,
			Currency:        parsedTx.Currency,
			Note:            parsedTx.Note,
			TransactionDate: parsedTx.TransactionDate,
		}

		if createReq.Currency == "" {
			createReq.Currency = "CNY"
		}

		tx, err := s.createTransactionInternal(userID, createReq)
		if err != nil {
			result.FailedRows++
			result.Errors = append(result.Errors, models.ImportError{
				LineNumber: parsedTx.LineNumber,
				Error:      err.Error(),
			})
			continue
		}

		result.ImportedRows++
		result.ImportedIDs = append(result.ImportedIDs, tx.ID)
	}

	return result, nil
}

// Internal method to create transaction (simplified version without full service dependencies)
func (s *ImportService) createTransactionInternal(userID uuid.UUID, req *CreateTransactionRequest) (*models.Transaction, error) {
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	account, err := s.accountRepo.GetByID(req.AccountID, userID)
	if err != nil {
		return nil, fmt.Errorf("invalid account: %w", err)
	}

	if req.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(*req.CategoryID, userID)
		if err != nil {
			req.CategoryID = nil // Ignore invalid category
		}
	}

	transaction, err := s.transactionRepo.Create(
		userID,
		req.AccountID,
		req.CategoryID,
		req.Type,
		req.Amount,
		req.Currency,
		req.Note,
		req.TransactionDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update account balance
	switch req.Type {
	case models.TransactionTypeIncome:
		account.Balance += req.Amount
	case models.TransactionTypeExpense:
		account.Balance -= req.Amount
	}
	_, _ = s.accountRepo.Update(account, userID)

	return transaction, nil
}

// parseAlipayCSV parses Alipay CSV format
func (s *ImportService) parseAlipayCSV(r io.Reader) ([]models.ParsedTransaction, error) {
	reader := csv.NewReader(r)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("no data in file")
	}

	// Alipay CSV usually has a header row at index 0 or 1
	// Find the actual header by looking for known column names
	headerIndex := 0
	for i, record := range records {
		if len(record) > 0 && (strings.Contains(record[0], "交易时间") || strings.Contains(record[0], "时间")) {
			headerIndex = i
			break
		}
		if i > 5 { // Don't search too far
			break
		}
	}

	transactions := make([]models.ParsedTransaction, 0)

	for i := headerIndex + 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 3 {
			continue
		}

		tx := models.ParsedTransaction{
			RawData: make(map[string]string),
			Currency: "CNY",
		}

		// Map fields by index (Alipay format is somewhat consistent)
		// Common Alipay columns: 交易时间, 交易分类, 交易对方, 金额, 收/支, 支付方式
		for j, field := range record {
			if headerIndex < len(records) && j < len(records[headerIndex]) {
				tx.RawData[records[headerIndex][j]] = field
			}
		}

		// Parse transaction date
		for _, key := range []string{"交易时间", "时间", "日期"} {
			if val, ok := tx.RawData[key]; ok && val != "" {
				tx.TransactionDate, _ = s.parseChineseDate(val)
				break
			}
		}

		// Parse amount
		for _, key := range []string{"金额", "资金变动", "消费金额"} {
			if val, ok := tx.RawData[key]; ok && val != "" {
				tx.Amount, _ = s.parseAmount(val)
				break
			}
		}

		// Determine type (income/expense)
		for _, key := range []string{"收/支", "类型", "收支类型"} {
			if val, ok := tx.RawData[key]; ok {
				if strings.Contains(val, "收") || strings.Contains(val, "收入") || strings.Contains(val, "+") {
					tx.Type = models.TransactionTypeIncome
				} else if strings.Contains(val, "支") || strings.Contains(val, "支出") || strings.Contains(val, "-") {
					tx.Type = models.TransactionTypeExpense
				}
			}
		}

		// Fallback: infer type from amount sign
		if tx.Type == "" {
			if tx.Amount < 0 {
				tx.Type = models.TransactionTypeExpense
				tx.Amount = math.Abs(tx.Amount)
			} else if tx.Amount > 0 {
				tx.Type = models.TransactionTypeExpense // Default to expense for Alipay
			}
		}

		// Get counterparty/merchant
		for _, key := range []string{"交易对方", "对方", "商户名称", "名称"} {
			if val, ok := tx.RawData[key]; ok && val != "" {
				tx.Counterparty = val
				break
			}
		}

		// Get account name
		for _, key := range []string{"支付方式", "账户", "付款方式"} {
			if val, ok := tx.RawData[key]; ok && val != "" {
				tx.AccountName = val
				break
			}
		}

		// Get note
		for _, key := range []string{"备注", "说明", "商品说明"} {
			if val, ok := tx.RawData[key]; ok && val != "" {
				tx.Note = val
				break
			}
		}

		// Combine for category hint
		tx.CategoryHint = strings.Join([]string{tx.Counterparty, tx.Note}, " ")

		if tx.Amount > 0 && !tx.TransactionDate.IsZero() {
			transactions = append(transactions, tx)
		}
	}

	return transactions, nil
}

// parseWeChatCSV parses WeChat CSV format
func (s *ImportService) parseWeChatCSV(r io.Reader) ([]models.ParsedTransaction, error) {
	reader := csv.NewReader(r)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("no data in file")
	}

	// Find header row
	headerIndex := 0
	for i, record := range records {
		if len(record) > 0 && (strings.Contains(record[0], "交易时间") || strings.Contains(record[0], "时间")) {
			headerIndex = i
			break
		}
		if i > 5 {
			break
		}
	}

	transactions := make([]models.ParsedTransaction, 0)

	for i := headerIndex + 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 3 {
			continue
		}

		tx := models.ParsedTransaction{
			RawData: make(map[string]string),
			Currency: "CNY",
		}

		// Store raw data
		for j, field := range record {
			if headerIndex < len(records) && j < len(records[headerIndex]) {
				tx.RawData[records[headerIndex][j]] = field
			}
		}

		// Parse transaction date
		for _, key := range []string{"交易时间", "时间", "日期"} {
			if val, ok := tx.RawData[key]; ok && val != "" {
				tx.TransactionDate, _ = s.parseChineseDate(val)
				break
			}
		}

		// Parse amount
		for _, key := range []string{"金额", "资金变动", "消费金额", "收/支金额"} {
			if val, ok := tx.RawData[key]; ok && val != "" {
				tx.Amount, _ = s.parseAmount(val)
				break
			}
		}

		// Determine type
		for _, key := range []string{"收/支", "类型", "收支类型", "交易类型"} {
			if val, ok := tx.RawData[key]; ok {
				if strings.Contains(val, "收") || strings.Contains(val, "收入") || strings.Contains(val, "转账") && strings.Contains(val, "存入") {
					tx.Type = models.TransactionTypeIncome
				} else if strings.Contains(val, "支") || strings.Contains(val, "支出") || strings.Contains(val, "消费") || strings.Contains(val, "转账") && strings.Contains(val, "转出") {
					tx.Type = models.TransactionTypeExpense
				}
			}
		}

		// Fallback type detection
		if tx.Type == "" {
			for _, key := range []string{"交易类型", "类型"} {
				if val, ok := tx.RawData[key]; ok {
					if strings.Contains(val, "转账") {
						// Need more info for transfers
						tx.Type = models.TransactionTypeExpense
					} else if strings.Contains(val, "红包") || strings.Contains(val, "退款") {
						tx.Type = models.TransactionTypeIncome
					}
				}
			}
		}

		if tx.Type == "" && tx.Amount > 0 {
			tx.Type = models.TransactionTypeExpense // Default to expense
		}

		// Get counterparty
		for _, key := range []string{"交易对方", "对方", "商户名称", "名称", "联系人"} {
			if val, ok := tx.RawData[key]; ok && val != "" {
				tx.Counterparty = val
				break
			}
		}

		// Get account
		for _, key := range []string{"支付方式", "账户", "付款方式", "收/付款方式"} {
			if val, ok := tx.RawData[key]; ok && val != "" {
				tx.AccountName = val
				break
			}
		}

		// Get note
		for _, key := range []string{"备注", "说明", "商品", "商品说明"} {
			if val, ok := tx.RawData[key]; ok && val != "" {
				tx.Note = val
				break
			}
		}

		tx.CategoryHint = strings.Join([]string{tx.Counterparty, tx.Note}, " ")

		if tx.Amount > 0 && !tx.TransactionDate.IsZero() {
			transactions = append(transactions, tx)
		}
	}

	return transactions, nil
}

// parseBankCSV parses generic bank CSV format
func (s *ImportService) parseBankCSV(r io.Reader) ([]models.ParsedTransaction, error) {
	reader := csv.NewReader(r)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("no data in file")
	}

	// Find header row
	headerIndex := 0
	for i, record := range records {
		recordStr := strings.Join(record, " ")
		if strings.Contains(recordStr, "日期") || strings.Contains(recordStr, "date") ||
		   strings.Contains(recordStr, "摘要") || strings.Contains(recordStr, "description") {
			headerIndex = i
			break
		}
		if i > 5 {
			break
		}
	}

	transactions := make([]models.ParsedTransaction, 0)

	for i := headerIndex + 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 3 {
			continue
		}

		tx := models.ParsedTransaction{
			RawData: make(map[string]string),
			Currency: "CNY",
		}

		// Store raw data with header keys
		for j, field := range record {
			key := fmt.Sprintf("col%d", j)
			if headerIndex < len(records) && j < len(records[headerIndex]) {
				key = records[headerIndex][j]
			}
			tx.RawData[key] = field
		}

		// Try to parse date from common column names
		dateFound := false
		for key, val := range tx.RawData {
			if (strings.Contains(key, "日期") || strings.Contains(key, "date") ||
			    strings.Contains(key, "time") || strings.Contains(key, "时间")) && val != "" {
				if t, err := s.parseChineseDate(val); err == nil {
					tx.TransactionDate = t
					dateFound = true
					break
				} else if t, err := time.Parse("2006-01-02", val); err == nil {
					tx.TransactionDate = t
					dateFound = true
					break
				} else if t, err := time.Parse("2006/01/02", val); err == nil {
					tx.TransactionDate = t
					dateFound = true
					break
				}
			}
		}

		if !dateFound {
			continue
		}

		// Try to parse amount
		amountFound := false
		for key, val := range tx.RawData {
			if (strings.Contains(key, "金额") || strings.Contains(key, "amount") ||
			    strings.Contains(key, "支出") || strings.Contains(key, "收入") ||
			    strings.Contains(key, "debit") || strings.Contains(key, "credit")) && val != "" {
				if amt, err := s.parseAmount(val); err == nil && amt > 0 {
					tx.Amount = amt
					// Determine type from column name
					if strings.Contains(key, "支出") || strings.Contains(key, "debit") || strings.Contains(key, "out") {
						tx.Type = models.TransactionTypeExpense
					} else if strings.Contains(key, "收入") || strings.Contains(key, "credit") || strings.Contains(key, "in") {
						tx.Type = models.TransactionTypeIncome
					}
					amountFound = true
					break
				}
			}
		}

		// Try separate income/expense columns
		if !amountFound {
			var incomeAmt, expenseAmt float64
			for key, val := range tx.RawData {
				if (strings.Contains(key, "收入") || strings.Contains(key, "credit")) && val != "" {
					incomeAmt, _ = s.parseAmount(val)
				}
				if (strings.Contains(key, "支出") || strings.Contains(key, "debit")) && val != "" {
					expenseAmt, _ = s.parseAmount(val)
				}
			}
			if incomeAmt > 0 {
				tx.Amount = incomeAmt
				tx.Type = models.TransactionTypeIncome
				amountFound = true
			} else if expenseAmt > 0 {
				tx.Amount = expenseAmt
				tx.Type = models.TransactionTypeExpense
				amountFound = true
			}
		}

		if !amountFound || tx.Type == "" {
			continue
		}

		// Get description/note
		for key, val := range tx.RawData {
			if (strings.Contains(key, "摘要") || strings.Contains(key, "description") ||
			    strings.Contains(key, "备注") || strings.Contains(key, "note")) && val != "" {
				tx.Note = val
				break
			}
		}

		// Get counterparty
		for key, val := range tx.RawData {
			if (strings.Contains(key, "对方") || strings.Contains(key, "merchant") ||
			    strings.Contains(key, "payee") || strings.Contains(key, "收款人")) && val != "" {
				tx.Counterparty = val
				break
			}
		}

		// Get account info
		for key, val := range tx.RawData {
			if (strings.Contains(key, "账户") || strings.Contains(key, "account") ||
			    strings.Contains(key, "卡号")) && val != "" {
				tx.AccountName = val
				break
			}
		}

		tx.CategoryHint = strings.Join([]string{tx.Counterparty, tx.Note}, " ")

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// parseGenericCSV parses generic CSV with flexible column mapping
func (s *ImportService) parseGenericCSV(r io.Reader) ([]models.ParsedTransaction, error) {
	reader := csv.NewReader(r)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("no data in file")
	}

	// Use first row as header
	headerIndex := 0
	transactions := make([]models.ParsedTransaction, 0)

	for i := headerIndex + 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 2 {
			continue
		}

		tx := models.ParsedTransaction{
			RawData: make(map[string]string),
			Currency: "CNY",
		}

		// Store raw data
		for j, field := range record {
			key := fmt.Sprintf("col%d", j)
			if headerIndex < len(records) && j < len(records[headerIndex]) {
				key = records[headerIndex][j]
			}
			tx.RawData[key] = field
		}

		// Try to find date in any column
		for key, val := range tx.RawData {
			if val == "" {
				continue
			}
			if t, err := s.parseChineseDate(val); err == nil {
				tx.TransactionDate = t
				break
			} else if t, err := time.Parse("2006-01-02", val); err == nil {
				tx.TransactionDate = t
				break
			} else if t, err := time.Parse("2006/01/02", val); err == nil {
				tx.TransactionDate = t
				break
			}
		}

		// Try to find amount in any column
		for key, val := range tx.RawData {
			if val == "" {
				continue
			}
			if amt, err := s.parseAmount(val); err == nil && amt > 0 {
				tx.Amount = amt
				// Guess type from sign or column name
				if strings.HasPrefix(strings.TrimSpace(val), "-") {
					tx.Type = models.TransactionTypeExpense
					tx.Amount = math.Abs(tx.Amount)
				} else if strings.Contains(key, "支出") || strings.Contains(key, "expense") {
					tx.Type = models.TransactionTypeExpense
				} else if strings.Contains(key, "收入") || strings.Contains(key, "income") {
					tx.Type = models.TransactionTypeIncome
				} else {
					tx.Type = models.TransactionTypeExpense // Default to expense
				}
				break
			}
		}

		// Use other columns as note
		var notes []string
		for key, val := range tx.RawData {
			if val == "" {
				continue
			}
			// Skip columns that are already used for date/amount
			if _, err := s.parseAmount(val); err == nil {
				continue
			}
			if _, err := time.Parse("2006-01-02", val); err == nil {
				continue
			}
			if _, err := s.parseChineseDate(val); err == nil {
				continue
			}
			notes = append(notes, val)
		}
		tx.Note = strings.Join(notes, " | ")

		if tx.Amount > 0 && !tx.TransactionDate.IsZero() && tx.Type != "" {
			transactions = append(transactions, tx)
		}
	}

	return transactions, nil
}

// parseChineseDate parses Chinese date formats like "2024-01-15 14:30:00" or "2024/01/15"
func (s *ImportService) parseChineseDate(str string) (time.Time, error) {
	str = strings.TrimSpace(str)

	// Try various formats
	formats := []string{
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"2006-01-02 15:04",
		"2006/01/02 15:04",
		"2006-01-02",
		"2006/01/02",
		"2006.01.02",
		"20060102",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, str, time.Local); err == nil {
			return t.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", str)
}

// parseAmount parses amount strings with currency symbols and commas
func (s *ImportService) parseAmount(str string) (float64, error) {
	str = strings.TrimSpace(str)

	// Remove currency symbols and spaces
	replacer := strings.NewReplacer(
		"¥", "", "￥", "", "CNY", "", "元", "",
		"$", "", "USD", "", "€", "", "£", "",
		" ", "", ",", "", "+", "",
	)
	str = replacer.Replace(str)

	// Handle negative numbers in parentheses, e.g., (100.00)
	if strings.HasPrefix(str, "(") && strings.HasSuffix(str, ")") {
		str = "-" + str[1:len(str)-1]
	}

	return strconv.ParseFloat(str, 64)
}

// checkForDuplicate checks if a transaction already exists
func (s *ImportService) checkForDuplicate(userID uuid.UUID, tx *models.ParsedTransaction) (bool, error) {
	// Simple duplicate check: same date, amount, and similar note
	existing, err := s.transactionRepo.GetByDateRangeAndAmount(
		userID,
		tx.TransactionDate.Add(-24*time.Hour),
		tx.TransactionDate.Add(24*time.Hour),
		tx.Amount-0.01,
		tx.Amount+0.01,
	)
	if err != nil {
		return false, err
	}

	return len(existing) > 0, nil
}

// findMatchingAccounts finds accounts that match the given name hint
func (s *ImportService) findMatchingAccounts(nameHint string, accounts []models.Account) []models.Account {
	nameHint = strings.ToLower(nameHint)
	matches := make([]models.Account, 0)

	for _, account := range accounts {
		accountName := strings.ToLower(account.Name)
		if strings.Contains(accountName, nameHint) || strings.Contains(nameHint, accountName) {
			matches = append(matches, account)
		}
	}

	return matches
}
