package handlers

import (
	"account/internal/business/models"
	"account/internal/data/repository"
	"account/pkg/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	userRepo     *repository.UserRepository
	categoryRepo *repository.CategoryRepository
	tokenMgr     *auth.TokenManager
	logger       *zap.Logger
}

func NewAuthHandler(userRepo *repository.UserRepository, categoryRepo *repository.CategoryRepository, tokenMgr *auth.TokenManager, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
		tokenMgr:     tokenMgr,
		logger:       logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.Create(req.Email, req.Password)
	if err != nil {
		if err == repository.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		h.logger.Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Create default categories for the new user
	if err := h.categoryRepo.CreateDefaultCategories(user.ID); err != nil {
		h.logger.Error("Failed to create default categories", zap.Error(err))
		// Continue without failing the registration
	}

	// Generate token
	token, expiry, err := h.tokenMgr.GenerateToken(user.ID, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(expiry).Seconds()),
		User:        *user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		h.logger.Error("Failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if !h.userRepo.VerifyPassword(user, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate token
	token, expiry, err := h.tokenMgr.GenerateToken(user.ID, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(expiry).Seconds()),
		User:        *user,
	})
}
