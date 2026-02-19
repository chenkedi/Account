package services

import (
	"account/internal/business/models"
	"account/internal/data/repository"
	"fmt"

	"github.com/google/uuid"
)

type AccountService struct {
	accountRepo  *repository.AccountRepository
	categoryRepo *repository.CategoryRepository
}

func NewAccountService(accountRepo *repository.AccountRepository, categoryRepo *repository.CategoryRepository) *AccountService {
	return &AccountService{
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
	}
}

type CreateAccountRequest struct {
	Name     string                `json:"name" binding:"required"`
	Type     models.AccountType    `json:"type" binding:"required"`
	Currency string                `json:"currency"`
	Balance  float64               `json:"balance"`
}

type UpdateAccountRequest struct {
	Name     string                `json:"name"`
	Type     models.AccountType    `json:"type"`
	Currency string                `json:"currency"`
	Balance  float64               `json:"balance"`
}

func (s *AccountService) CreateAccount(userID uuid.UUID, req *CreateAccountRequest) (*models.Account, error) {
	if req.Currency == "" {
		req.Currency = "CNY"
	}

	if !isValidAccountType(req.Type) {
		return nil, fmt.Errorf("invalid account type: %s", req.Type)
	}

	account, err := s.accountRepo.Create(userID, req.Name, req.Type, req.Currency, req.Balance)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (s *AccountService) GetAccount(id uuid.UUID, userID uuid.UUID) (*models.Account, error) {
	account, err := s.accountRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) GetAllAccounts(userID uuid.UUID) ([]models.Account, error) {
	accounts, err := s.accountRepo.GetAll(userID)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *AccountService) UpdateAccount(id uuid.UUID, userID uuid.UUID, req *UpdateAccountRequest) (*models.Account, error) {
	account, err := s.accountRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		account.Name = req.Name
	}
	if req.Type != "" {
		if !isValidAccountType(req.Type) {
			return nil, fmt.Errorf("invalid account type: %s", req.Type)
		}
		account.Type = req.Type
	}
	if req.Currency != "" {
		account.Currency = req.Currency
	}
	if req.Balance != 0 {
		account.Balance = req.Balance
	}

	updatedAccount, err := s.accountRepo.Update(account, userID)
	if err != nil {
		return nil, err
	}

	return updatedAccount, nil
}

func (s *AccountService) DeleteAccount(id uuid.UUID, userID uuid.UUID) error {
	return s.accountRepo.Delete(id, userID)
}

func isValidAccountType(t models.AccountType) bool {
	switch t {
	case models.AccountTypeBank, models.AccountTypeCash, models.AccountTypeAlipay,
		models.AccountTypeWeChat, models.AccountTypeCredit, models.AccountTypeInvestment,
		models.AccountTypeOther:
		return true
	default:
		return false
	}
}
