package api

import (
	"account/internal/api/handlers"
	"account/internal/api/middleware"
	"account/internal/business/models"
	"account/internal/business/services"
	"account/pkg/auth"
	"account/test/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupSyncTestRouter() (*gin.Engine, *mocks.MockSyncRepository, *mocks.MockAccountRepository, *mocks.MockCategoryRepository, *mocks.MockTransactionRepository, *auth.TokenManager, uuid.UUID) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	syncRepo := new(mocks.MockSyncRepository)
	accountRepo := new(mocks.MockAccountRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	transactionRepo := new(mocks.MockTransactionRepository)
	logger := zap.NewNop()
	tokenMgr := auth.NewTokenManager("test-secret-key", 24*time.Hour)

	syncService := services.NewSyncService(syncRepo, accountRepo, categoryRepo, transactionRepo)
	syncHandler := handlers.NewSyncHandler(syncService, logger)

	testUserID := uuid.New()

	// Protected routes
	router.Use(func(c *gin.Context) {
		c.Set("user_id", testUserID)
		c.Set("email", "test@example.com")
		c.Next()
	})
	router.Use(middleware.AuthMiddleware(tokenMgr))
	router.POST("/api/v1/sync/pull", syncHandler.Pull)
	router.POST("/api/v1/sync/push", syncHandler.Push)

	return router, syncRepo, accountRepo, categoryRepo, transactionRepo, tokenMgr, testUserID
}

func TestPull_Success(t *testing.T) {
	router, syncRepo, _, _, _, _, testUserID := setupSyncTestRouter()

	deviceID := "device-123"
	lastSyncAt := time.Now().Add(-24 * time.Hour)

	testAccountID := uuid.New()
	testCategoryID := uuid.New()
	testTransactionID := uuid.New()

	expectedResponse := &models.SyncPullResponse{
		Accounts: []models.Account{
			{
				ID:             testAccountID,
				UserID:         testUserID,
				Name:           "Test Account",
				Type:           models.AccountTypeCash,
				LastModifiedAt: time.Now(),
			},
		},
		Categories: []models.Category{
			{
				ID:             testCategoryID,
				UserID:         testUserID,
				Name:           "Test Category",
				Type:           models.CategoryTypeExpense,
				LastModifiedAt: time.Now(),
			},
		},
		Transactions: []models.Transaction{
			{
				ID:              testTransactionID,
				UserID:          testUserID,
				AccountID:       testAccountID,
				Note:            "Test Transaction",
				Amount:          100.0,
				Type:            models.TransactionTypeExpense,
				LastModifiedAt:  time.Now(),
			},
		},
		CurrentSyncAt: time.Now(),
	}

	syncRepo.On("GetChangesSince", testUserID, mock.AnythingOfType("time.Time")).Return(expectedResponse, nil)

	reqBody := models.SyncPullRequest{
		DeviceID:   deviceID,
		LastSyncAt: lastSyncAt,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/sync/pull", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SyncPullResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Accounts, 1)
	assert.Len(t, response.Categories, 1)
	assert.Len(t, response.Transactions, 1)

	syncRepo.AssertExpectations(t)
}

func TestPull_MissingDeviceID(t *testing.T) {
	router, _, _, _, _, _, _ := setupSyncTestRouter()

	reqBody := map[string]interface{}{
		"last_sync_at": time.Now(),
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/sync/pull", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPush_Success(t *testing.T) {
	router, syncRepo, _, _, _, _, testUserID := setupSyncTestRouter()

	deviceID := "device-123"
	testAccountID := uuid.New()

	pushRequest := models.SyncPushRequest{
		DeviceID: deviceID,
		Accounts: []models.Account{
			{
				ID:             testAccountID,
				UserID:         testUserID,
				Name:           "Test Account",
				Type:           models.AccountTypeCash,
				LastModifiedAt: time.Now(),
			},
		},
		Categories:   []models.Category{},
		Transactions: []models.Transaction{},
		LastSyncAt:   time.Now().Add(-24 * time.Hour),
	}

	syncRepo.On("ApplyChanges", testUserID, mock.AnythingOfType("*models.SyncPushRequest")).Return(nil)
	syncRepo.On("UpsertSyncState", testUserID, deviceID, "").Return(&models.SyncState{}, nil)

	jsonBody, _ := json.Marshal(pushRequest)

	req := httptest.NewRequest("POST", "/api/v1/sync/push", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SyncPushResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotZero(t, response.CurrentSyncAt)

	syncRepo.AssertExpectations(t)
}

func TestPush_MissingDeviceID(t *testing.T) {
	router, _, _, _, _, _, _ := setupSyncTestRouter()

	reqBody := map[string]interface{}{
		"accounts":     []interface{}{},
		"categories":   []interface{}{},
		"transactions": []interface{}{},
		"last_sync_at": time.Now(),
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/sync/push", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
