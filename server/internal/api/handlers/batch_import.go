package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"account/internal/business/models"
	"account/internal/business/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// BatchImportHandler handles batch import HTTP requests
type BatchImportHandler struct {
	batchService *services.BatchImportService
	logger       *zap.Logger
}

// NewBatchImportHandler creates a new BatchImportHandler
func NewBatchImportHandler(batchService *services.BatchImportService, logger *zap.Logger) *BatchImportHandler {
	return &BatchImportHandler{
		batchService: batchService,
		logger:       logger,
	}
}

// CreateBatchImportRequest represents the request to create a batch import
type CreateBatchImportRequest struct {
	Files []struct {
		Source   string `json:"source" binding:"required,oneof=alipay wechat jd bank generic"`
		FileName string `json:"file_name" binding:"required"`
		Content  string `json:"content" binding:"required"` // base64 encoded
	} `json:"files" binding:"required,min=1,max=20"`
}

// CreateBatchImportResponse represents the response after creating a batch import
type CreateBatchImportResponse struct {
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	FileCount int    `json:"file_count"`
}

// CreateBatchImport creates a new batch import job
func (h *BatchImportHandler) CreateBatchImport(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateBatchImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate and decode files
	files := make([]services.FileUpload, 0, len(req.Files))
	for _, f := range req.Files {
		// Validate source
		if !isValidSource(f.Source) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid source: %s", f.Source)})
			return
		}

		// Validate file name
		if f.FileName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file_name is required"})
			return
		}

		// Validate base64 content
		if f.Content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
			return
		}

		// Decode base64 to validate
		_, err := base64.StdEncoding.DecodeString(f.Content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid base64 content for file %s: %v", f.FileName, err)})
			return
		}

		files = append(files, services.FileUpload{
			Source:   f.Source,
			FileName: f.FileName,
			Content:  f.Content,
		})
	}

	// Create batch job
	job, err := h.batchService.CreateBatchJob(userID, files)
	if err != nil {
		h.logger.Error("Failed to create batch job", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create batch job"})
		return
	}

	c.JSON(http.StatusCreated, CreateBatchImportResponse{
		JobID:     job.ID.String(),
		Status:    string(job.Status),
		Message:   "Batch import job created successfully. Processing started in background.",
		FileCount: len(files),
	})
}

// isValidSource validates the import source
func isValidSource(source string) bool {
	switch models.ImportSource(source) {
	case models.ImportSourceAlipay, models.ImportSourceWeChat, models.ImportSourceJD,
		models.ImportSourceBank, models.ImportSourceGeneric:
		return true
	}
	return false
}

// getUserID extracts the user ID from the context
func (h *BatchImportHandler) getUserID(c *gin.Context) (uuid.UUID, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user_id not found in context")
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user_id format")
	}

	return userID, nil
}

// ListBatchImportsResponse represents the response for listing batch imports
type ListBatchImportsResponse struct {
	Jobs []models.BatchImportJob `json:"jobs"`
}

// ListBatchImports lists all batch import jobs for the user
func (h *BatchImportHandler) ListBatchImports(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// TODO: In a real implementation, this would fetch from database
	// For now, return empty list
	_ = userID // Mark as used
	c.JSON(http.StatusOK, ListBatchImportsResponse{
		Jobs: []models.BatchImportJob{},
	})
}

// GetBatchImportStatusResponse represents the response for getting batch import status
type GetBatchImportStatusResponse struct {
	Job      models.BatchImportJob `json:"job"`
	Progress services.ImportProgress `json:"progress"`
}

// GetBatchImportStatus gets the status of a batch import job
func (h *BatchImportHandler) GetBatchImportStatus(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	jobIDStr := c.Param("job_id")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job_id"})
		return
	}

	// TODO: In a real implementation, fetch job from database
	// For now, return a dummy job
	job := models.BatchImportJob{
		ID:        jobID,
		UserID:    userID,
		Status:    models.BatchImportStatusReadyToImport,
		TotalFiles: 2,
		ParsedFiles: 2,
		TotalTransactions: 10,
		ValidTransactions: 10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	progress := services.CalculateProgress(&job)

	c.JSON(http.StatusOK, GetBatchImportStatusResponse{
		Job:      job,
		Progress: progress,
	})
}

// GetBatchImportPreviewResponse represents the response for getting batch import preview
type GetBatchImportPreviewResponse struct {
	Job              models.BatchImportJob       `json:"job"`
	Files            []models.BatchImportFile    `json:"files"`
	TransferMatches  []models.TransferMatch      `json:"transfer_matches"`
	AccountHints     []models.AccountHint        `json:"account_hints"`
}

// GetBatchImportPreview gets the preview of a batch import job
func (h *BatchImportHandler) GetBatchImportPreview(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	jobIDStr := c.Param("job_id")
	_, err = uuid.Parse(jobIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job_id"})
		return
	}

	// TODO: In a real implementation, fetch preview data from database
	// For now, return empty preview
	_ = userID // Mark as used for now
	c.JSON(http.StatusOK, GetBatchImportPreviewResponse{
		Job: models.BatchImportJob{
			UserID: userID,
		},
		Files:           []models.BatchImportFile{},
		TransferMatches: []models.TransferMatch{},
		AccountHints:    []models.AccountHint{},
	})
}

// ExecuteBatchImportRequest represents the request to execute a batch import
type ExecuteBatchImportRequest struct {
	SelectedAccountIDs map[string]string `json:"selected_account_ids"` // file_index:account_id
	ConfirmedMatchIDs  []string          `json:"confirmed_match_ids"`
}

// ExecuteBatchImport executes the batch import
func (h *BatchImportHandler) ExecuteBatchImport(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	jobIDStr := c.Param("job_id")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job_id"})
		return
	}

	var req ExecuteBatchImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.batchService.ExecuteBatchImport(
		jobID,
		userID,
		req.SelectedAccountIDs,
		req.ConfirmedMatchIDs,
	)
	if err != nil {
		h.logger.Error("Failed to execute batch import", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to execute batch import"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteBatchImport deletes a batch import job
func (h *BatchImportHandler) DeleteBatchImport(c *gin.Context) {
	_, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	jobIDStr := c.Param("job_id")
	_, err = uuid.Parse(jobIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job_id"})
		return
	}

	// TODO: In a real implementation, delete from database
	c.JSON(http.StatusOK, gin.H{"message": "Batch import job deleted successfully"})
}
