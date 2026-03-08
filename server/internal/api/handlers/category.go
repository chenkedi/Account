package handlers

import (
	"account/internal/business/models"
	"account/internal/business/services"
	"account/internal/data/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
	logger          *zap.Logger
}

func NewCategoryHandler(categoryService *services.CategoryService, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		logger:          logger,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req services.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryService.CreateCategory(userID, &req)
	if err != nil {
		h.logger.Error("Failed to create category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	category, err := h.categoryService.GetCategory(id, userID)
	if err != nil {
		if err == repository.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		h.logger.Error("Failed to get category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	categories, err := h.categoryService.GetAllCategories(userID)
	if err != nil {
		h.logger.Error("Failed to get categories", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetCategoriesByType(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	typeParam := c.Param("type")
	categoryType := models.CategoryType(typeParam)

	categories, err := h.categoryService.GetCategoriesByType(userID, categoryType)
	if err != nil {
		h.logger.Error("Failed to get categories by type", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	// 使用中间结构来正确处理 parent_id 可以为 null 的情况
	type UpdateCategoryRequestBody struct {
		Name     string      `json:"name"`
		Type     string      `json:"type"`
		ParentID interface{} `json:"parent_id"` // 使用 interface{} 来接收 null 或 string
		Icon     string      `json:"icon"`
	}

	var reqBody UpdateCategoryRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建服务层请求
	var req services.UpdateCategoryRequest
	req.Name = reqBody.Name
	req.Type = models.CategoryType(reqBody.Type)
	req.Icon = reqBody.Icon

	// 处理 parent_id
	if reqBody.ParentID == nil {
		// 显式设置为 null
		req.ParentID = nil
	} else if parentIDStr, ok := reqBody.ParentID.(string); ok {
		if parentIDStr == "" {
			// 空字符串也视为 null
			req.ParentID = nil
		} else {
			// 解析 UUID
			if parentUUID, err := uuid.Parse(parentIDStr); err == nil {
				req.ParentID = &parentUUID
			}
		}
	}

	category, err := h.categoryService.UpdateCategory(id, userID, &req)
	if err != nil {
		if err == repository.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		h.logger.Error("Failed to update category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if err := h.categoryService.DeleteCategory(id, userID); err != nil {
		if err == repository.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		h.logger.Error("Failed to delete category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
