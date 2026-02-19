package api

import (
	"account/internal/api/handlers"
	"account/internal/business/models"
	"account/internal/business/services"
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

func setupAccountTestRouter() (*gin.Engine, *mocks.MockAccountRepository, *mocks.MockCategoryRepository, uuid.UUID) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	accountRepo := new(mocks.MockAccountRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	logger := zap.NewNop()

	accountService := services.NewAccountService(accountRepo, categoryRepo)
	accountHandler := handlers.NewAccountHandler(accountService, logger)

	testUserID := uuid.New()

	// Mock auth middleware
	router.Use(func(c *gin.Context) {
		c.Set("user_id", testUserID)
		c.Set("email", "test@example.com")
		c.Next()
	})

	router.POST("/api/v1/accounts", accountHandler.CreateAccount)
	router.GET("/api/v1/accounts", accountHandler.GetAllAccounts)
	router.GET("/api/v1/accounts/:id", accountHandler.GetAccount)
	router.PUT("/api/v1/accounts/:id", accountHandler.UpdateAccount)
	router.DELETE("/api/v1/accounts/:id", accountHandler.DeleteAccount)

	return router, accountRepo, categoryRepo, testUserID
}

func TestCreateAccount_Success(t *testing.T) {
	router, accountRepo, _, testUserID := setupAccountTestRouter()

	accountID := uuid.New()

	createReq := services.CreateAccountRequest{
		Name:     "Test Account",
		Type:     models.AccountTypeCash,
		Currency: "CNY",
		Balance:  1000.0,
	}

	expectedAccount := &models.Account{
		ID:             accountID,
		UserID:         testUserID,
		Name:           createReq.Name,
		Type:           createReq.Type,
		Currency:       createReq.Currency,
		Balance:        createReq.Balance,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
		Version:        1,
	}

	accountRepo.On("Create", testUserID, createReq.Name, createReq.Type, createReq.Currency, createReq.Balance).Return(expectedAccount, nil)

	jsonBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Account
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAccount.ID, response.ID)
	assert.Equal(t, expectedAccount.Name, response.Name)

	accountRepo.AssertExpectations(t)
}

func TestGetAccount_Success(t *testing.T) {
	router, accountRepo, _, testUserID := setupAccountTestRouter()

	accountID := uuid.New()

	expectedAccount := &models.Account{
		ID:      accountID,
		UserID:  testUserID,
		Name:    "Test Account",
		Type:    models.AccountTypeCash,
		Balance: 1000.0,
	}

	accountRepo.On("GetByID", accountID, testUserID).Return(expectedAccount, nil)

	req := httptest.NewRequest("GET", "/api/v1/accounts/"+accountID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Account
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAccount.ID, response.ID)

	accountRepo.AssertExpectations(t)
}

func TestGetAllAccounts_Success(t *testing.T) {
	router, accountRepo, _, testUserID := setupAccountTestRouter()

	expectedAccounts := []models.Account{
		{
			ID:      uuid.New(),
			UserID:  testUserID,
			Name:    "Cash",
			Type:    models.AccountTypeCash,
			Balance: 500.0,
		},
		{
			ID:      uuid.New(),
			UserID:  testUserID,
			Name:    "Bank",
			Type:    models.AccountTypeBank,
			Balance: 5000.0,
		},
	}

	accountRepo.On("GetAll", testUserID).Return(expectedAccounts, nil)

	req := httptest.NewRequest("GET", "/api/v1/accounts", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Account
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)

	accountRepo.AssertExpectations(t)
}

func TestUpdateAccount_Success(t *testing.T) {
	router, accountRepo, _, testUserID := setupAccountTestRouter()

	accountID := uuid.New()

	updateReq := services.UpdateAccountRequest{
		Name:    "Updated Account",
		Type:    models.AccountTypeBank,
		Balance: 2000.0,
	}

	updatedAccount := &models.Account{
		ID:             accountID,
		UserID:         testUserID,
		Name:           updateReq.Name,
		Type:           updateReq.Type,
		Balance:        updateReq.Balance,
		LastModifiedAt: time.Now(),
	}

	accountRepo.On("GetByID", accountID, testUserID).Return(updatedAccount, nil)
	accountRepo.On("Update", mock.AnythingOfType("*models.Account"), testUserID).Return(updatedAccount, nil)

	jsonBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/accounts/"+accountID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Account
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updateReq.Name, response.Name)

	accountRepo.AssertExpectations(t)
}

func TestDeleteAccount_Success(t *testing.T) {
	router, accountRepo, _, testUserID := setupAccountTestRouter()

	accountID := uuid.New()

	accountRepo.On("Delete", accountID, testUserID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/v1/accounts/"+accountID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	accountRepo.AssertExpectations(t)
}
