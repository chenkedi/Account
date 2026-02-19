package api

import (
	"account/internal/api/handlers"
	"account/internal/api/middleware"
	"account/internal/business/models"
	"account/internal/data/repository"
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
	"go.uber.org/zap"
)

func setupAuthTestRouter() (*gin.Engine, *mocks.MockUserRepository, *mocks.MockCategoryRepository, *auth.TokenManager) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	userRepo := new(mocks.MockUserRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	logger := zap.NewNop()
	tokenMgr := auth.NewTokenManager("test-secret-key", 24*time.Hour)

	authHandler := handlers.NewAuthHandler(userRepo, categoryRepo, tokenMgr, logger)

	router.POST("/api/v1/auth/register", authHandler.Register)
	router.POST("/api/v1/auth/login", authHandler.Login)

	return router, userRepo, categoryRepo, tokenMgr
}

func TestRegister_Success(t *testing.T) {
	router, userRepo, categoryRepo, _ := setupAuthTestRouter()

	testUserID := uuid.New()
	testEmail := "test@example.com"
	testPassword := "password123"

	userRepo.On("Create", testEmail, testPassword).Return(&models.User{
		ID:        testUserID,
		Email:     testEmail,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	categoryRepo.On("CreateDefaultCategories", testUserID).Return(nil)

	reqBody := map[string]string{
		"email":    testEmail,
		"password": testPassword,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["access_token"])
	assert.Equal(t, "Bearer", response["token_type"])

	userRepo.AssertExpectations(t)
	categoryRepo.AssertExpectations(t)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	router, userRepo, _, _ := setupAuthTestRouter()

	testEmail := "test@example.com"
	testPassword := "password123"

	userRepo.On("Create", testEmail, testPassword).Return(nil, repository.ErrUserAlreadyExists)

	reqBody := map[string]string{
		"email":    testEmail,
		"password": testPassword,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	userRepo.AssertExpectations(t)
}

func TestRegister_InvalidRequest(t *testing.T) {
	router, _, _, _ := setupAuthTestRouter()

	// Missing password
	reqBody := map[string]string{
		"email": "test@example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	router, userRepo, _, tokenMgr := setupAuthTestRouter()

	testUserID := uuid.New()
	testEmail := "test@example.com"
	testPassword := "password123"

	testUser := &models.User{
		ID:           testUserID,
		Email:        testEmail,
		PasswordHash: "hashed-password",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userRepo.On("GetByEmail", testEmail).Return(testUser, nil)
	userRepo.On("VerifyPassword", testUser, testPassword).Return(true)

	reqBody := map[string]string{
		"email":    testEmail,
		"password": testPassword,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["access_token"])

	userRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	router, userRepo, _, _ := setupAuthTestRouter()

	testEmail := "test@example.com"
	testPassword := "password123"

	testUser := &models.User{
		ID:           uuid.New(),
		Email:        testEmail,
		PasswordHash: "hashed-password",
	}

	userRepo.On("GetByEmail", testEmail).Return(testUser, nil)
	userRepo.On("VerifyPassword", testUser, testPassword).Return(false)

	reqBody := map[string]string{
		"email":    testEmail,
		"password": testPassword,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	userRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	router, userRepo, _, _ := setupAuthTestRouter()

	testEmail := "test@example.com"
	testPassword := "password123"

	userRepo.On("GetByEmail", testEmail).Return(nil, repository.ErrUserNotFound)

	reqBody := map[string]string{
		"email":    testEmail,
		"password": testPassword,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	userRepo.AssertExpectations(t)
}


// Add auth middleware test
func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tokenMgr := auth.NewTokenManager("test-secret-key", 24*time.Hour)
	testUserID := uuid.New()
	testEmail := "test@example.com"

	// Create a protected route
	router.Use(middleware.AuthMiddleware(tokenMgr))
	router.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		email, _ := c.Get("email")
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"email":   email,
		})
	})

	// Test with valid token
	token, _, err := tokenMgr.GenerateToken(testUserID, testEmail)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Test without token
	req = httptest.NewRequest("GET", "/protected", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test with invalid token
	req = httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
