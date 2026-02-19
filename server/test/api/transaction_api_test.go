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

func setupTransactionTestRouter() (*gin.Engine, *mocks.MockTransactionRepository, *mocks.MockAccountRepository, uuid.UUID) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	txRepo := new(mocks.MockTransactionRepository)
	accountRepo := new(mocks.MockAccountRepository)
	logger := zap.NewNop()

	txService := services.NewTransactionService(txRepo, accountRepo, logger)
	txHandler := handlers.NewTransactionHandler(txService, logger)

	testUserID := uuid.New()

	// Mock auth middleware
	router.Use(func(c *gin.Context) {
		c.Set("user_id", testUserID)
		c.Set("email", "test@example.com")
		c.Next()
	})

	router.POST("/api/v1/transactions", txHandler.CreateTransaction)
	router.GET("/api/v1/transactions/:id", txHandler.GetTransaction)
	router.GET("/api/v1/transactions", txHandler.GetAllTransactions)
	router.PUT("/api/v1/transactions/:id", txHandler.UpdateTransaction)
	router.DELETE("/api/v1/transactions/:id", txHandler.DeleteTransaction)

	return router, txRepo, accountRepo, testUserID
}

func TestCreateTransaction_Success(t *testing.T) {
	router, txRepo, accountRepo, testUserID := setupTransactionTestRouter()

	accountID := uuid.New()
	txID := uuid.New()

	createReq := services.CreateTransactionRequest{
		AccountID:       accountID,
		Type:            models.TransactionTypeExpense,
		Amount:          100.50,
		Currency:        "CNY",
		Note:            "Test Transaction",
		TransactionDate: time.Now(),
	}

	testAccount := &models.Account{
		ID:      accountID,
		UserID:  testUserID,
		Balance: 1000.0,
	}

	expectedTx := &models.Transaction{
		ID:              txID,
		UserID:          testUserID,
		AccountID:       accountID,
		Type:            createReq.Type,
		Amount:          createReq.Amount,
		Currency:        createReq.Currency,
		Note:            createReq.Note,
		TransactionDate: createReq.TransactionDate,
	}

	accountRepo.On("GetByID", accountID, testUserID).Return(testAccount, nil)
	txRepo.On("Create", mock.AnythingOfType("*models.Transaction")).Return(expectedTx, nil)
	accountRepo.On("Update", mock.AnythingOfType("*models.Account")).Return(testAccount, nil)

	jsonBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Transaction
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTx.ID, response.ID)
	assert.Equal(t, expectedTx.Note, response.Note)

	accountRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestGetTransaction_Success(t *testing.T) {
	router, txRepo, _, testUserID := setupTransactionTestRouter()

	txID := uuid.New()
	accountID := uuid.New()

	expectedTx := &models.Transaction{
		ID:        txID,
		UserID:    testUserID,
		AccountID: accountID,
		Note:      "Test Transaction",
		Amount:    100.0,
		Type:      models.TransactionTypeExpense,
	}

	txRepo.On("GetByID", txID, testUserID).Return(expectedTx, nil)

	req := httptest.NewRequest("GET", "/api/v1/transactions/"+txID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Transaction
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTx.ID, response.ID)

	txRepo.AssertExpectations(t)
}

func TestGetTransaction_NotFound(t *testing.T) {
	router, txRepo, _, testUserID := setupTransactionTestRouter()

	txID := uuid.New()

	txRepo.On("GetByID", txID, testUserID).Return(nil, repository.ErrTransactionNotFound)

	req := httptest.NewRequest("GET", "/api/v1/transactions/"+txID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateTransaction_Success(t *testing.T) {
	router, txRepo, accountRepo, testUserID := setupTransactionTestRouter()

	txID := uuid.New()
	accountID := uuid.New()

	updateReq := services.UpdateTransactionRequest{
		Note: "Updated Note",
	}

	oldTx := &models.Transaction{
		ID:              txID,
		UserID:          testUserID,
		AccountID:       accountID,
		Note:            "Old Note",
		Amount:          100.0,
		Type:            models.TransactionTypeExpense,
		TransactionDate: time.Now(),
	}

	updatedTx := &models.Transaction{
		ID:              txID,
		UserID:          testUserID,
		AccountID:       accountID,
		Note:            updateReq.Note,
		Amount:          100.0,
		Type:            models.TransactionTypeExpense,
		TransactionDate: time.Now(),
	}

	testAccount := &models.Account{
		ID:      accountID,
		UserID:  testUserID,
		Balance: 1000.0,
	}

	txRepo.On("GetByID", txID, testUserID).Return(oldTx, nil)
	accountRepo.On("GetByID", accountID, testUserID).Return(testAccount, nil)
	txRepo.On("Update", mock.AnythingOfType("*models.Transaction")).Return(updatedTx, nil)
	accountRepo.On("Update", mock.AnythingOfType("*models.Account")).Return(testAccount, nil)

	jsonBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/transactions/"+txID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Transaction
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updateReq.Note, response.Note)

	txRepo.AssertExpectations(t)
}

func TestDeleteTransaction_Success(t *testing.T) {
	router, txRepo, accountRepo, testUserID := setupTransactionTestRouter()

	txID := uuid.New()
	accountID := uuid.New()

	oldTx := &models.Transaction{
		ID:              txID,
		UserID:          testUserID,
		AccountID:       accountID,
		Note:            "To Delete",
		Amount:          100.0,
		Type:            models.TransactionTypeExpense,
		TransactionDate: time.Now(),
	}

	testAccount := &models.Account{
		ID:      accountID,
		UserID:  testUserID,
		Balance: 1000.0,
	}

	txRepo.On("GetByID", txID, testUserID).Return(oldTx, nil)
	accountRepo.On("GetByID", accountID, testUserID).Return(testAccount, nil)
	txRepo.On("Delete", txID, testUserID).Return(nil)
	accountRepo.On("Update", mock.AnythingOfType("*models.Account")).Return(testAccount, nil)

	req := httptest.NewRequest("DELETE", "/api/v1/transactions/"+txID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	txRepo.AssertExpectations(t)
}

var ErrTransactionNotFound = &mockError{"transaction not found"}
