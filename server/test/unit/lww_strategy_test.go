package unit

import (
	"account/internal/business/models"
	"account/internal/sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLWWStrategy_ResolveAccount(t *testing.T) {
	strategy := sync.NewLWWStrategy()
	userID := uuid.New()
	accountID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		local          *models.Account
		remote         *models.Account
		expectedName   string
		expectNil      bool
	}{
		{
			name:         "local is nil, remote wins",
			local:        nil,
			remote:       createTestAccount(userID, accountID, "Remote", now.Add(-1*time.Hour)),
			expectedName: "Remote",
			expectNil:    false,
		},
		{
			name:         "remote is nil, local wins",
			local:        createTestAccount(userID, accountID, "Local", now.Add(-1*time.Hour)),
			remote:       nil,
			expectedName: "Local",
			expectNil:    false,
		},
		{
			name:         "both nil, returns nil",
			local:        nil,
			remote:       nil,
			expectedName: "",
			expectNil:    true,
		},
		{
			name:         "remote has later timestamp, remote wins",
			local:        createTestAccount(userID, accountID, "Local", now.Add(-2*time.Hour)),
			remote:       createTestAccount(userID, accountID, "Remote", now.Add(-1*time.Hour)),
			expectedName: "Remote",
			expectNil:    false,
		},
		{
			name:         "local has later timestamp, local wins",
			local:        createTestAccount(userID, accountID, "Local", now.Add(-1*time.Hour)),
			remote:       createTestAccount(userID, accountID, "Remote", now.Add(-2*time.Hour)),
			expectedName: "Local",
			expectNil:    false,
		},
		{
			name:         "same timestamp, local wins",
			local:        createTestAccount(userID, accountID, "Local", now.Truncate(time.Second)),
			remote:       createTestAccount(userID, accountID, "Remote", now.Truncate(time.Second)),
			expectedName: "Local",
			expectNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.ResolveAccount(tt.local, tt.remote)
			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expectedName, result.Name)
			}
		})
	}
}

func TestLWWStrategy_ResolveCategory(t *testing.T) {
	strategy := sync.NewLWWStrategy()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		local          *models.Category
		remote         *models.Category
		expectedName   string
		expectNil      bool
	}{
		{
			name:         "local is nil, remote wins",
			local:        nil,
			remote:       createTestCategory(userID, categoryID, "Remote", now.Add(-1*time.Hour)),
			expectedName: "Remote",
			expectNil:    false,
		},
		{
			name:         "remote is nil, local wins",
			local:        createTestCategory(userID, categoryID, "Local", now.Add(-1*time.Hour)),
			remote:       nil,
			expectedName: "Local",
			expectNil:    false,
		},
		{
			name:         "remote has later timestamp, remote wins",
			local:        createTestCategory(userID, categoryID, "Local", now.Add(-2*time.Hour)),
			remote:       createTestCategory(userID, categoryID, "Remote", now.Add(-1*time.Hour)),
			expectedName: "Remote",
			expectNil:    false,
		},
		{
			name:         "local has later timestamp, local wins",
			local:        createTestCategory(userID, categoryID, "Local", now.Add(-1*time.Hour)),
			remote:       createTestCategory(userID, categoryID, "Remote", now.Add(-2*time.Hour)),
			expectedName: "Local",
			expectNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.ResolveCategory(tt.local, tt.remote)
			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expectedName, result.Name)
			}
		})
	}
}

func TestLWWStrategy_ResolveTransaction(t *testing.T) {
	strategy := sync.NewLWWStrategy()
	userID := uuid.New()
	accountID := uuid.New()
	txID := uuid.New()
	now := time.Now()

	tests := []struct {
		name           string
		local          *models.Transaction
		remote         *models.Transaction
		expectedNote   string
		expectedAmount float64
		expectNil      bool
	}{
		{
			name:           "local is nil, remote wins",
			local:          nil,
			remote:         createTestTransaction(userID, accountID, txID, "Remote Note", 100.0, now.Add(-1*time.Hour)),
			expectedNote:   "Remote Note",
			expectedAmount: 100.0,
			expectNil:      false,
		},
		{
			name:           "remote is nil, local wins",
			local:          createTestTransaction(userID, accountID, txID, "Local Note", 100.0, now.Add(-1*time.Hour)),
			remote:         nil,
			expectedNote:   "Local Note",
			expectedAmount: 100.0,
			expectNil:      false,
		},
		{
			name:           "remote has later timestamp, remote wins",
			local:          createTestTransaction(userID, accountID, txID, "Local Note", 100.0, now.Add(-2*time.Hour)),
			remote:         createTestTransaction(userID, accountID, txID, "Remote Note", 200.0, now.Add(-1*time.Hour)),
			expectedNote:   "Remote Note",
			expectedAmount: 200.0,
			expectNil:      false,
		},
		{
			name:           "local has later timestamp, local wins",
			local:          createTestTransaction(userID, accountID, txID, "Local Note", 200.0, now.Add(-1*time.Hour)),
			remote:         createTestTransaction(userID, accountID, txID, "Remote Note", 100.0, now.Add(-2*time.Hour)),
			expectedNote:   "Local Note",
			expectedAmount: 200.0,
			expectNil:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.ResolveTransaction(tt.local, tt.remote)
			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expectedNote, result.Note)
				assert.Equal(t, tt.expectedAmount, result.Amount)
			}
		})
	}
}

func TestLWWStrategy_MergeAccountLists(t *testing.T) {
	strategy := sync.NewLWWStrategy()
	userID := uuid.New()

	account1ID := uuid.New()
	account2ID := uuid.New()
	account3ID := uuid.New()

	now := time.Now()

	localAccounts := []models.Account{
		*createTestAccount(userID, account1ID, "Local Account 1", now.Add(-2*time.Hour)),
		*createTestAccount(userID, account2ID, "Local Account 2", now.Add(-30*time.Minute)),
		*createTestAccount(userID, account3ID, "Local Only", now.Add(-1*time.Hour)),
	}

	remoteAccounts := []models.Account{
		*createTestAccount(userID, account1ID, "Remote Account 1", now.Add(-1*time.Hour)),
		*createTestAccount(userID, account2ID, "Remote Account 2", now.Add(-1*time.Hour)),
		{
			ID:             uuid.New(),
			UserID:         userID,
			Name:           "Remote Only",
			Type:           models.AccountTypeCash,
			Currency:       "CNY",
			Balance:        0,
			LastModifiedAt: now.Add(-30 * time.Minute),
		},
	}

	result := strategy.MergeAccountLists(localAccounts, remoteAccounts)

	assert.Len(t, result, 4)

	resultMap := make(map[string]models.Account)
	for _, a := range result {
		resultMap[a.ID.String()] = a
	}

	assert.Equal(t, "Remote Account 1", resultMap[account1ID.String()].Name)
	assert.Equal(t, "Local Account 2", resultMap[account2ID.String()].Name)
	assert.Equal(t, "Local Only", resultMap[account3ID.String()].Name)
}

func TestLWWStrategy_MergeCategoryLists(t *testing.T) {
	strategy := sync.NewLWWStrategy()
	userID := uuid.New()

	cat1ID := uuid.New()
	now := time.Now()

	localCategories := []models.Category{
		*createTestCategory(userID, cat1ID, "Local Category", now.Add(-2*time.Hour)),
	}

	remoteCategories := []models.Category{
		*createTestCategory(userID, cat1ID, "Remote Category", now.Add(-1*time.Hour)),
	}

	result := strategy.MergeCategoryLists(localCategories, remoteCategories)

	assert.Len(t, result, 1)
	assert.Equal(t, "Remote Category", result[0].Name)
}

func TestLWWStrategy_MergeTransactionLists(t *testing.T) {
	strategy := sync.NewLWWStrategy()
	userID := uuid.New()
	accountID := uuid.New()

	tx1ID := uuid.New()
	now := time.Now()

	localTransactions := []models.Transaction{
		*createTestTransaction(userID, accountID, tx1ID, "Local Note", 100.0, now.Add(-2*time.Hour)),
	}

	remoteTransactions := []models.Transaction{
		*createTestTransaction(userID, accountID, tx1ID, "Remote Note", 200.0, now.Add(-1*time.Hour)),
	}

	result := strategy.MergeTransactionLists(localTransactions, remoteTransactions)

	assert.Len(t, result, 1)
	assert.Equal(t, "Remote Note", result[0].Note)
	assert.Equal(t, 200.0, result[0].Amount)
}

func createTestAccount(userID, accountID uuid.UUID, name string, lastModified time.Time) *models.Account {
	return &models.Account{
		ID:             accountID,
		UserID:         userID,
		Name:           name,
		Type:           models.AccountTypeCash,
		Currency:       "CNY",
		Balance:        1000.0,
		CreatedAt:      time.Now().Add(-24 * time.Hour),
		UpdatedAt:      lastModified,
		LastModifiedAt: lastModified,
		Version:        1,
		IsDeleted:      false,
	}
}

func createTestCategory(userID, categoryID uuid.UUID, name string, lastModified time.Time) *models.Category {
	return &models.Category{
		ID:             categoryID,
		UserID:         userID,
		Name:           name,
		Type:           models.CategoryTypeExpense,
		CreatedAt:      time.Now().Add(-24 * time.Hour),
		UpdatedAt:      lastModified,
		LastModifiedAt: lastModified,
		Version:        1,
		IsDeleted:      false,
	}
}

func createTestTransaction(userID, accountID, txID uuid.UUID, note string, amount float64, lastModified time.Time) *models.Transaction {
	return &models.Transaction{
		ID:              txID,
		UserID:          userID,
		AccountID:       accountID,
		Type:            models.TransactionTypeExpense,
		Amount:          amount,
		Currency:        "CNY",
		Note:            note,
		TransactionDate: time.Now().Add(-1 * time.Hour),
		CreatedAt:       time.Now().Add(-24 * time.Hour),
		UpdatedAt:       lastModified,
		LastModifiedAt:  lastModified,
		Version:         1,
		IsDeleted:       false,
	}
}
