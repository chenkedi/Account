package handlers

import (
	"account/internal/business/services"
	"account/internal/data/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
	logger             *zap.Logger
}

func NewTransactionHandler(transactionService *services.TransactionService, logger *zap.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		logger:             logger,
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req services.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.CreateTransaction(userID, &req)
	if err != nil {
		h.logger.Error("Failed to create transaction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	transaction, err := h.transactionService.GetTransaction(id, userID)
	if err != nil {
		if err == repository.ErrTransactionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		h.logger.Error("Failed to get transaction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req services.TransactionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Limit = 50
		req.Offset = 0
	}

	transactions, err := h.transactionService.GetAllTransactions(userID, &req)
	if err != nil {
		h.logger.Error("Failed to get transactions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransactionsByDateRange(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var start, end time.Time
	var parseErr error

	startStr := c.Query("start_date")
	if startStr != "" {
		start, parseErr = time.Parse(time.RFC3339, startStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use RFC3339"})
			return
		}
	} else {
		start = time.Now().AddDate(0, -1, 0)
	}

	endStr := c.Query("end_date")
	if endStr != "" {
		end, parseErr = time.Parse(time.RFC3339, endStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use RFC3339"})
			return
		}
	} else {
		end = time.Now()
	}

	transactions, err := h.transactionService.GetTransactionsByDateRange(userID, start, end)
	if err != nil {
		h.logger.Error("Failed to get transactions by date range", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	var req services.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionService.UpdateTransaction(id, userID, &req)
	if err != nil {
		if err == repository.ErrTransactionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		h.logger.Error("Failed to update transaction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	if err := h.transactionService.DeleteTransaction(id, userID); err != nil {
		if err == repository.ErrTransactionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		h.logger.Error("Failed to delete transaction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *TransactionHandler) GetStats(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var start, end time.Time
	var parseErr error

	startStr := c.Query("start_date")
	if startStr != "" {
		start, parseErr = time.Parse(time.RFC3339, startStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use RFC3339"})
			return
		}
	} else {
		start = time.Now().AddDate(0, -1, 0)
	}

	endStr := c.Query("end_date")
	if endStr != "" {
		end, parseErr = time.Parse(time.RFC3339, endStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use RFC3339"})
			return
		}
	} else {
		end = time.Now()
	}

	req := &services.StatsRequest{
		StartDate: start,
		EndDate:   end,
	}

	stats, err := h.transactionService.GetStats(userID, req)
	if err != nil {
		h.logger.Error("Failed to get stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *TransactionHandler) GetDetailedStats(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var start, end time.Time
	var parseErr error

	startStr := c.Query("start_date")
	if startStr != "" {
		start, parseErr = time.Parse(time.RFC3339, startStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use RFC3339"})
			return
		}
	} else {
		// Default to last 6 months
		start = time.Now().AddDate(0, -6, 0)
	}

	endStr := c.Query("end_date")
	if endStr != "" {
		end, parseErr = time.Parse(time.RFC3339, endStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use RFC3339"})
			return
		}
	} else {
		end = time.Now()
	}

	req := &services.StatsRequest{
		StartDate: start,
		EndDate:   end,
	}

	stats, err := h.transactionService.GetDetailedStats(userID, req)
	if err != nil {
		h.logger.Error("Failed to get detailed stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
