package services

import (
	"account/internal/business/models"
	"account/internal/data/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
	categoryRepo    *repository.CategoryRepository
}

func NewTransactionService(
	transactionRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	categoryRepo *repository.CategoryRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
	}
}

type CreateTransactionRequest struct {
	AccountID       uuid.UUID                `json:"account_id" binding:"required"`
	CategoryID      *uuid.UUID               `json:"category_id"`
	Type            models.TransactionType   `json:"type" binding:"required"`
	Amount          float64                  `json:"amount" binding:"required,min=0"`
	Currency        string                   `json:"currency"`
	Note            string                   `json:"note"`
	TransactionDate time.Time                `json:"transaction_date" binding:"required"`
}

type UpdateTransactionRequest struct {
	AccountID       uuid.UUID                `json:"account_id"`
	CategoryID      *uuid.UUID               `json:"category_id"`
	Type            models.TransactionType   `json:"type"`
	Amount          float64                  `json:"amount"`
	Currency        string                   `json:"currency"`
	Note            string                   `json:"note"`
	TransactionDate time.Time                `json:"transaction_date"`
}

type TransactionListRequest struct {
	Limit  int    `form:"limit,default=50"`
	Offset int    `form:"offset,default=0"`
}

type StatsRequest struct {
	StartDate time.Time `form:"start_date" binding:"required"`
	EndDate   time.Time `form:"end_date" binding:"required"`
}

type StatsResponse struct {
	IncomeTotal  float64                `json:"income_total"`
	ExpenseTotal float64                `json:"expense_total"`
	NetTotal     float64                `json:"net_total"`
	StartDate    time.Time              `json:"start_date"`
	EndDate      time.Time              `json:"end_date"`
}

type CategoryStatsItem struct {
	CategoryID      string  `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	CategoryType    string  `json:"category_type"`
	TotalAmount     float64 `json:"total_amount"`
	TransactionCount int    `json:"transaction_count"`
	Percentage      float64 `json:"percentage"`
}

type MonthlyStatsItem struct {
	Year         int     `json:"year"`
	Month        int     `json:"month"`
	IncomeTotal  float64 `json:"income_total"`
	ExpenseTotal float64 `json:"expense_total"`
	NetTotal     float64 `json:"net_total"`
}

type DetailedStatsResponse struct {
	Summary       *StatsResponse         `json:"summary"`
	ByCategory    []CategoryStatsItem    `json:"by_category"`
	MonthlyTrend  []MonthlyStatsItem     `json:"monthly_trend"`
}


func (s *TransactionService) CreateTransaction(userID uuid.UUID, req *CreateTransactionRequest) (*models.Transaction, error) {
	if !isValidTransactionType(req.Type) {
		return nil, fmt.Errorf("invalid transaction type: %s", req.Type)
	}

	if req.Currency == "" {
		req.Currency = "CNY"
	}

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
			return nil, fmt.Errorf("invalid category: %w", err)
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

	if err := s.updateAccountBalance(account, req.Type, req.Amount); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) GetTransaction(id uuid.UUID, userID uuid.UUID) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *TransactionService) GetAllTransactions(userID uuid.UUID, req *TransactionListRequest) ([]models.Transaction, error) {
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 200 {
		req.Limit = 200
	}

	transactions, err := s.transactionRepo.GetAll(userID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *TransactionService) GetTransactionsByDateRange(userID uuid.UUID, start, end time.Time) ([]models.Transaction, error) {
	transactions, err := s.transactionRepo.GetByDateRange(userID, start, end)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *TransactionService) GetTransactionsByAccount(userID uuid.UUID, accountID uuid.UUID, limit int) ([]models.Transaction, error) {
	transactions, err := s.transactionRepo.GetByAccount(userID, accountID, limit)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *TransactionService) UpdateTransaction(id uuid.UUID, userID uuid.UUID, req *UpdateTransactionRequest) (*models.Transaction, error) {
	transaction, err := s.transactionRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	oldAccount, err := s.accountRepo.GetByID(transaction.AccountID, userID)
	if err != nil {
		return nil, fmt.Errorf("invalid old account: %w", err)
	}

	if req.AccountID != uuid.Nil && req.AccountID != transaction.AccountID {
		_, err := s.accountRepo.GetByID(req.AccountID, userID)
		if err != nil {
			return nil, fmt.Errorf("invalid new account: %w", err)
		}
	}

	if req.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(*req.CategoryID, userID)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %w", err)
		}
	}

	oldType := transaction.Type
	oldAmount := transaction.Amount

	if req.Type != "" {
		if !isValidTransactionType(req.Type) {
			return nil, fmt.Errorf("invalid transaction type: %s", req.Type)
		}
		transaction.Type = req.Type
	}
	if req.AccountID != uuid.Nil {
		transaction.AccountID = req.AccountID
	}
	if req.CategoryID != nil {
		transaction.CategoryID = req.CategoryID
	}
	if req.Amount > 0 {
		transaction.Amount = req.Amount
	}
	if req.Currency != "" {
		transaction.Currency = req.Currency
	}
	if req.Note != "" {
		transaction.Note = req.Note
	}
	if !req.TransactionDate.IsZero() {
		transaction.TransactionDate = req.TransactionDate.UTC()
	}

	updatedTransaction, err := s.transactionRepo.Update(transaction, userID)
	if err != nil {
		return nil, err
	}

	if err := s.reverseAccountBalance(oldAccount, oldType, oldAmount); err != nil {
		return nil, err
	}

	newAccount, err := s.accountRepo.GetByID(updatedTransaction.AccountID, userID)
	if err != nil {
		return nil, fmt.Errorf("invalid new account: %w", err)
	}

	if err := s.updateAccountBalance(newAccount, updatedTransaction.Type, updatedTransaction.Amount); err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}

func (s *TransactionService) DeleteTransaction(id uuid.UUID, userID uuid.UUID) error {
	transaction, err := s.transactionRepo.GetByID(id, userID)
	if err != nil {
		return err
	}

	account, err := s.accountRepo.GetByID(transaction.AccountID, userID)
	if err != nil {
		return fmt.Errorf("invalid account: %w", err)
	}

	if err := s.transactionRepo.Delete(id, userID); err != nil {
		return err
	}

	if err := s.reverseAccountBalance(account, transaction.Type, transaction.Amount); err != nil {
		return err
	}

	return nil
}

func (s *TransactionService) GetStats(userID uuid.UUID, req *StatsRequest) (*StatsResponse, error) {
	income, expense, err := s.transactionRepo.GetStatsByDateRange(userID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	return &StatsResponse{
		IncomeTotal:  income,
		ExpenseTotal: expense,
		NetTotal:     income - expense,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
	}, nil
}

func (s *TransactionService) GetDetailedStats(userID uuid.UUID, req *StatsRequest) (*DetailedStatsResponse, error) {
	// Get summary
	income, expense, err := s.transactionRepo.GetStatsByDateRange(userID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	summary := &StatsResponse{
		IncomeTotal:  income,
		ExpenseTotal: expense,
		NetTotal:     income - expense,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
	}

	// Get category stats
	categoryStats, err := s.transactionRepo.GetCategoryStats(userID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	// Calculate percentages
	var categoryItems []CategoryStatsItem
	for _, cs := range categoryStats {
		var total float64
		if cs.CategoryType == "income" {
			total = income
		} else {
			total = expense
		}

		var percentage float64
		if total > 0 {
			percentage = (cs.TotalAmount / total) * 100
		}

		categoryItems = append(categoryItems, CategoryStatsItem{
			CategoryID:      cs.CategoryID.String(),
			CategoryName:    cs.CategoryName,
			CategoryType:    cs.CategoryType,
			TotalAmount:     cs.TotalAmount,
			TransactionCount: cs.TransactionCount,
			Percentage:      percentage,
		})
	}

	// Get monthly trend
	monthlyStats, err := s.transactionRepo.GetMonthlyTrend(userID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	var monthlyItems []MonthlyStatsItem
	for _, ms := range monthlyStats {
		monthlyItems = append(monthlyItems, MonthlyStatsItem{
			Year:         ms.Year,
			Month:        ms.Month,
			IncomeTotal:  ms.IncomeTotal,
			ExpenseTotal: ms.ExpenseTotal,
			NetTotal:     ms.IncomeTotal - ms.ExpenseTotal,
		})
	}

	return &DetailedStatsResponse{
		Summary:      summary,
		ByCategory:   categoryItems,
		MonthlyTrend: monthlyItems,
	}, nil
}

func (s *TransactionService) updateAccountBalance(account *models.Account, tType models.TransactionType, amount float64) error {
	switch tType {
	case models.TransactionTypeIncome:
		account.Balance += amount
	case models.TransactionTypeExpense:
		account.Balance -= amount
	}

	_, err := s.accountRepo.Update(account, account.UserID)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	return nil
}

func (s *TransactionService) reverseAccountBalance(account *models.Account, tType models.TransactionType, amount float64) error {
	switch tType {
	case models.TransactionTypeIncome:
		account.Balance -= amount
	case models.TransactionTypeExpense:
		account.Balance += amount
	}

	_, err := s.accountRepo.Update(account, account.UserID)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	return nil
}

func isValidTransactionType(t models.TransactionType) bool {
	switch t {
	case models.TransactionTypeIncome, models.TransactionTypeExpense, models.TransactionTypeTransfer:
		return true
	default:
		return false
	}
}
