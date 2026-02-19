package handlers

import (
	"account/internal/business/models"
	"account/internal/business/services"
	"bytes"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ImportHandler handles import-related HTTP requests
type ImportHandler struct {
	importService *services.ImportService
	logger        *zap.Logger
}

// NewImportHandler creates a new ImportHandler
func NewImportHandler(importService *services.ImportService, logger *zap.Logger) *ImportHandler {
	return &ImportHandler{
		importService: importService,
		logger:        logger,
	}
}

// UploadAndParseRequest is the request for uploading and parsing a file
type UploadAndParseRequest struct {
	Source string `form:"source" binding:"required,oneof=alipay wechat bank generic"`
}

// UploadAndParse handles file upload and initial parsing
func (h *ImportHandler) UploadAndParse(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get source from form data
	var req UploadAndParseRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}
	defer file.Close()

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		h.logger.Error("Failed to read uploaded file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	// Parse the file
	parseReq := &services.ParseRequest{
		Source:   models.ImportSource(req.Source),
		FileName: header.Filename,
		File:     bytes.NewReader(content),
	}

	preview, err := h.importService.ParseFile(userID, parseReq)
	if err != nil {
		h.logger.Error("Failed to parse file", zap.Error(err), zap.String("source", req.Source))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, preview)
}

// ExecuteImportRequest is the request to execute import
type ExecuteImportRequest struct {
	JobID        string                  `json:"job_id" binding:"required"`
	Transactions []models.ParsedTransaction `json:"transactions" binding:"required"`
}

// ExecuteImport executes the import with user-selected mappings
func (h *ImportHandler) ExecuteImport(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ExecuteImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse job ID
	jobUUID, err := parseUUID(req.JobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job_id"})
		return
	}

	// Build execute request
	executeReq := &services.ExecuteImportRequest{
		JobID:        jobUUID,
		Transactions: req.Transactions,
	}

	// Execute import
	result, err := h.importService.ExecuteImport(userID, executeReq)
	if err != nil {
		h.logger.Error("Failed to execute import", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ParseTemplateRequest requests a template for a specific source
type ParseTemplateRequest struct {
	Source string `form:"source" binding:"required,oneof=alipay wechat bank generic"`
}

// GetTemplateInfo returns information about the expected file format
func (h *ImportHandler) GetTemplateInfo(c *gin.Context) {
	source := c.Query("source")
	if source == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source parameter required"})
		return
	}

	var templateInfo map[string]interface{}

	switch source {
	case "alipay":
		templateInfo = map[string]interface{}{
			"source": "alipay",
			"description": "Alipay CSV export format. Download from Alipay app: My -> Bills -> ... -> Export Bill",
			"required_columns": []string{"交易时间", "金额", "收/支"},
			"optional_columns": []string{"交易对方", "备注", "支付方式"},
			"file_extensions": []string{".csv"},
		}
	case "wechat":
		templateInfo = map[string]interface{}{
			"source": "wechat",
			"description": "WeChat Pay CSV export format. Download from WeChat: Me -> Services -> Wallet -> Bill -> ... -> Export Bill",
			"required_columns": []string{"交易时间", "金额", "收/支"},
			"optional_columns": []string{"交易对方", "备注", "支付方式"},
			"file_extensions": []string{".csv"},
		}
	case "bank":
		templateInfo = map[string]interface{}{
			"source": "bank",
			"description": "Bank statement CSV format. Supports most Chinese bank export formats.",
			"required_columns": []string{"日期", "金额"},
			"optional_columns": []string{"摘要", "备注", "对方账户", "账户名"},
			"file_extensions": []string{".csv"},
		}
	case "generic":
		templateInfo = map[string]interface{}{
			"source": "generic",
			"description": "Generic CSV format. The parser will try to automatically detect columns.",
			"required_columns": []string{"date (or similar)", "amount (or similar)"},
			"optional_columns": []string{"description", "note", "category", "account"},
			"file_extensions": []string{".csv"},
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source"})
		return
	}

	c.JSON(http.StatusOK, templateInfo)
}

// GetSupportedSources returns all supported import sources
func (h *ImportHandler) GetSupportedSources(c *gin.Context) {
	sources := []map[string]interface{}{
		{
			"id":          "alipay",
			"name":        "支付宝",
			"description": "导入支付宝账单",
			"icon":        "alipay",
		},
		{
			"id":          "wechat",
			"name":        "微信支付",
			"description": "导入微信支付账单",
			"icon":        "wechat",
		},
		{
			"id":          "bank",
			"name":        "银行流水",
			"description": "导入银行对账单",
			"icon":        "bank",
		},
		{
			"id":          "generic",
			"name":        "通用CSV",
			"description": "导入通用CSV格式",
			"icon":        "file",
		},
	}

	c.JSON(http.StatusOK, gin.H{"sources": sources})
}

// parseUUID safely parses a UUID string
func parseUUID(s string) (interface{}, error) {
	// This is a placeholder - use google/uuid in real implementation
	// For now, just return a dummy value
	return s, nil
}

// Helper function to get file extension
func getFileExtension(filename string) string {
	return filepath.Ext(filename)
}
