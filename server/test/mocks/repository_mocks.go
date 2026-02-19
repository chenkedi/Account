package mocks

import (
	"account/internal/business/models"
	"account/internal/data/repository"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// User model for mocks
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(email, password string) (*models.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) VerifyPassword(user *models.User, password string) bool {
	args := m.Called(user, password)
	return args.Bool(0)
}

// MockAccountRepository is a mock implementation of AccountRepository
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(userID uuid.UUID, name string, accountType models.AccountType, currency string, balance float64) (*models.Account, error) {
	args := m.Called(userID, name, accountType, currency, balance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByID(id uuid.UUID, userID uuid.UUID) (*models.Account, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) GetAll(userID uuid.UUID) ([]models.Account, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Account), args.Error(1)
}

func (m *MockAccountRepository) Update(account *models.Account, userID uuid.UUID) (*models.Account, error) {
	args := m.Called(account, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAccountRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockAccountRepository) GetModifiedSince(userID uuid.UUID, since time.Time) ([]models.Account, error) {
	args := m.Called(userID, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Account), args.Error(1)
}

func (m *MockAccountRepository) CreateMany(accounts []models.Account) error {
	args := m.Called(accounts)
	return args.Error(0)
}

// MockCategoryRepository is a mock implementation of CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(userID uuid.UUID, name string, categoryType models.CategoryType, parentID *uuid.UUID, icon string) (*models.Category, error) {
	args := m.Called(userID, name, categoryType, parentID, icon)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByID(id uuid.UUID, userID uuid.UUID) (*models.Category, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetAll(userID uuid.UUID) ([]models.Category, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByType(userID uuid.UUID, categoryType models.CategoryType) ([]models.Category, error) {
	args := m.Called(userID, categoryType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(category *models.Category, userID uuid.UUID) (*models.Category, error) {
	args := m.Called(category, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockCategoryRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockCategoryRepository) CreateDefaultCategories(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetModifiedSince(userID uuid.UUID, since time.Time) ([]models.Category, error) {
	args := m.Called(userID, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoryRepository) CreateMany(categories []models.Category) error {
	args := m.Called(categories)
	return args.Error(0)
}

// MockTransactionRepository is a mock implementation of TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(tx *models.Transaction) (*models.Transaction, error) {
	args := m.Called(tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByID(id uuid.UUID, userID uuid.UUID) (*models.Transaction, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetAll(userID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByDateRange(userID uuid.UUID, startDate, endDate time.Time) ([]models.Transaction, error) {
	args := m.Called(userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(tx *models.Transaction) (*models.Transaction, error) {
	args := m.Called(tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Delete(id uuid.UUID, userID uuid.UUID) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetModifiedSince(userID uuid.UUID, since time.Time) ([]models.Transaction, error) {
	args := m.Called(userID, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) CreateMany(transactions []models.Transaction) error {
	args := m.Called(transactions)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetStatsByDateRange(userID uuid.UUID, startDate, endDate time.Time) (float64, float64, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

func (m *MockTransactionRepository) GetCategoryStats(userID uuid.UUID, startDate, endDate time.Time) ([]repository.CategoryStats, error) {
	args := m.Called(userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.CategoryStats), args.Error(1)
}

func (m *MockTransactionRepository) GetMonthlyTrend(userID uuid.UUID, startDate, endDate time.Time) ([]repository.MonthlyStats, error) {
	args := m.Called(userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.MonthlyStats), args.Error(1)
}

// MockSyncRepository is a mock implementation of SyncRepository
type MockSyncRepository struct {
	mock.Mock
}

func (m *MockSyncRepository) GetSyncState(userID uuid.UUID, deviceID string) (*models.SyncState, error) {
	args := m.Called(userID, deviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SyncState), args.Error(1)
}

func (m *MockSyncRepository) UpsertSyncState(userID uuid.UUID, deviceID string, syncToken string) (*models.SyncState, error) {
	args := m.Called(userID, deviceID, syncToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SyncState), args.Error(1)
}

func (m *MockSyncRepository) GetChangesSince(userID uuid.UUID, since time.Time) (*models.SyncPullResponse, error) {
	args := m.Called(userID, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SyncPullResponse), args.Error(1)
}

func (m *MockSyncRepository) ApplyChanges(userID uuid.UUID, req *models.SyncPushRequest) error {
	args := m.Called(userID, req)
	return args.Error(0)
}

// SyncState is a placeholder for import
type SyncState struct {
	UserID    uuid.UUID
	DeviceID  string
	LastSyncAt time.Time
	SyncToken string
	CreatedAt time.Time
	UpdatedAt time.Time
}
