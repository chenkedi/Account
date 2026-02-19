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

func setupCategoryTestRouter() (*gin.Engine, *mocks.MockCategoryRepository, uuid.UUID) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	categoryRepo := new(mocks.MockCategoryRepository)
	logger := zap.NewNop()

	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)

	testUserID := uuid.New()

	// Mock auth middleware
	router.Use(func(c *gin.Context) {
		c.Set("user_id", testUserID)
		c.Set("email", "test@example.com")
		c.Next()
	})

	router.POST("/api/v1/categories", categoryHandler.CreateCategory)
	router.GET("/api/v1/categories", categoryHandler.GetAllCategories)
	router.GET("/api/v1/categories/type/:type", categoryHandler.GetCategoriesByType)
	router.GET("/api/v1/categories/:id", categoryHandler.GetCategory)
	router.PUT("/api/v1/categories/:id", categoryHandler.UpdateCategory)
	router.DELETE("/api/v1/categories/:id", categoryHandler.DeleteCategory)

	return router, categoryRepo, testUserID
}

func TestCreateCategory_Success(t *testing.T) {
	router, categoryRepo, testUserID := setupCategoryTestRouter()

	categoryID := uuid.New()

	createReq := services.CreateCategoryRequest{
		Name: "Food & Dining",
		Type: models.CategoryTypeExpense,
		Icon: "🍔",
	}

	expectedCategory := &models.Category{
		ID:             categoryID,
		UserID:         testUserID,
		Name:           createReq.Name,
		Type:           createReq.Type,
		Icon:           createReq.Icon,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
		Version:        1,
	}

	categoryRepo.On("Create", testUserID, createReq.Name, createReq.Type, (*uuid.UUID)(nil), createReq.Icon).Return(expectedCategory, nil)

	jsonBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/categories", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Category
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory.ID, response.ID)
	assert.Equal(t, expectedCategory.Name, response.Name)

	categoryRepo.AssertExpectations(t)
}

func TestGetCategory_Success(t *testing.T) {
	router, categoryRepo, testUserID := setupCategoryTestRouter()

	categoryID := uuid.New()

	expectedCategory := &models.Category{
		ID:     categoryID,
		UserID: testUserID,
		Name:   "Food & Dining",
		Type:   models.CategoryTypeExpense,
	}

	categoryRepo.On("GetByID", categoryID, testUserID).Return(expectedCategory, nil)

	req := httptest.NewRequest("GET", "/api/v1/categories/"+categoryID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Category
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory.ID, response.ID)

	categoryRepo.AssertExpectations(t)
}

func TestGetAllCategories_Success(t *testing.T) {
	router, categoryRepo, testUserID := setupCategoryTestRouter()

	expectedCategories := []models.Category{
		{
			ID:     uuid.New(),
			UserID: testUserID,
			Name:   "Food",
			Type:   models.CategoryTypeExpense,
		},
		{
			ID:     uuid.New(),
			UserID: testUserID,
			Name:   "Salary",
			Type:   models.CategoryTypeIncome,
		},
	}

	categoryRepo.On("GetAll", testUserID).Return(expectedCategories, nil)

	req := httptest.NewRequest("GET", "/api/v1/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Category
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)

	categoryRepo.AssertExpectations(t)
}

func TestGetCategoriesByType_Success(t *testing.T) {
	router, categoryRepo, testUserID := setupCategoryTestRouter()

	expectedCategories := []models.Category{
		{
			ID:     uuid.New(),
			UserID: testUserID,
			Name:   "Food",
			Type:   models.CategoryTypeExpense,
		},
		{
			ID:     uuid.New(),
			UserID: testUserID,
			Name:   "Transport",
			Type:   models.CategoryTypeExpense,
		},
	}

	categoryRepo.On("GetByType", testUserID, models.CategoryTypeExpense).Return(expectedCategories, nil)

	req := httptest.NewRequest("GET", "/api/v1/categories/type/expense", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Category
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)

	categoryRepo.AssertExpectations(t)
}

func TestUpdateCategory_Success(t *testing.T) {
	router, categoryRepo, testUserID := setupCategoryTestRouter()

	categoryID := uuid.New()

	updateReq := services.UpdateCategoryRequest{
		Name: "Updated Category",
		Icon: "💰",
	}

	updatedCategory := &models.Category{
		ID:             categoryID,
		UserID:         testUserID,
		Name:           updateReq.Name,
		Icon:           updateReq.Icon,
		LastModifiedAt: time.Now(),
	}

	categoryRepo.On("GetByID", categoryID, testUserID).Return(updatedCategory, nil)
	categoryRepo.On("Update", mock.AnythingOfType("*models.Category"), testUserID).Return(updatedCategory, nil)

	jsonBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/categories/"+categoryID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Category
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updateReq.Name, response.Name)

	categoryRepo.AssertExpectations(t)
}

func TestDeleteCategory_Success(t *testing.T) {
	router, categoryRepo, testUserID := setupCategoryTestRouter()

	categoryID := uuid.New()

	categoryRepo.On("Delete", categoryID, testUserID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/v1/categories/"+categoryID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	categoryRepo.AssertExpectations(t)
}
