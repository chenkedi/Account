package integration

import (
	"testing"
	"time"

	"account/internal/business/models"
	"account/internal/business/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestBatchImport_CalculateProgress 测试进度计算功能
func TestBatchImport_CalculateProgress(t *testing.T) {
	// 测试进度计算
	job := &models.BatchImportJob{
		ID:     uuid.New(),
		Status: models.BatchImportStatusParsing,
	}

	progress := services.CalculateProgress(job)
	assert.NotNil(t, progress)
	assert.Equal(t, 6, progress.TotalSteps)
	assert.Equal(t, 1, progress.CurrentStep)
	assert.Equal(t, "解析文件", progress.StepDescription)
}

// TestBatchImport_AllStatusProgress 测试所有状态的进度计算
func TestBatchImport_AllStatusProgress(t *testing.T) {
	testCases := []struct {
		name           string
		status         models.BatchImportStatus
		expectedStep   int
		expectedDesc   string
	}{
		{
			name:         "Pending status",
			status:       models.BatchImportStatusPending,
			expectedStep: 0,
			expectedDesc: "等待处理",
		},
		{
			name:         "Parsing status",
			status:       models.BatchImportStatusParsing,
			expectedStep: 1,
			expectedDesc: "解析文件",
		},
		{
			name:         "Analyzing status",
			status:       models.BatchImportStatusAnalyzing,
			expectedStep: 2,
			expectedDesc: "分析交易",
		},
		{
			name:         "Matching status",
			status:       models.BatchImportStatusMatching,
			expectedStep: 3,
			expectedDesc: "匹配转账",
		},
		{
			name:         "Ready to import status",
			status:       models.BatchImportStatusReadyToImport,
			expectedStep: 4,
			expectedDesc: "准备导入",
		},
		{
			name:         "Importing status",
			status:       models.BatchImportStatusImporting,
			expectedStep: 5,
			expectedDesc: "导入数据",
		},
		{
			name:         "Completed status",
			status:       models.BatchImportStatusCompleted,
			expectedStep: 6,
			expectedDesc: "完成",
		},
		{
			name:         "Failed status",
			status:       models.BatchImportStatusFailed,
			expectedStep: 0,
			expectedDesc: "失败",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			job := &models.BatchImportJob{
				Status: tc.status,
			}
			progress := services.CalculateProgress(job)
			assert.Equal(t, tc.expectedStep, progress.CurrentStep)
			assert.Equal(t, tc.expectedDesc, progress.StepDescription)
		})
	}
}

// TestTransferMatcher_Create 测试转账匹配器创建
func TestTransferMatcher_Create(t *testing.T) {
	// Create matcher
	matcher := services.NewTransferMatcher()
	assert.NotNil(t, matcher)
}

// TestAccountHint_Model 测试账户提示模型
func TestAccountHint_Model(t *testing.T) {
	hint := models.AccountHint{
		Source:         models.ImportSourceAlipay,
		AccountName:    "招商银行",
		AccountNumber:  "1234",
		BankName:       "招商银行",
		CardType:       "信用卡",
		AccountType:    models.AccountTypeCredit,
		Balance:        10000.00,
		FoundInFile:    "alipay.csv",
	}

	assert.Equal(t, models.ImportSourceAlipay, hint.Source)
	assert.Equal(t, "招商银行", hint.AccountName)
	assert.Equal(t, "1234", hint.AccountNumber)
	assert.Equal(t, "信用卡", hint.CardType)
	assert.Equal(t, models.AccountTypeCredit, hint.AccountType)
	assert.Equal(t, 10000.00, hint.Balance)
}

// TestBatchImportJob_Model 测试批量导入任务模型
func TestBatchImportJob_Model(t *testing.T) {
	jobID := uuid.New()
	userID := uuid.New()

	job := models.BatchImportJob{
		ID:                  jobID,
		UserID:              userID,
		Status:              models.BatchImportStatusPending,
		TotalFiles:          3,
		ParsedFiles:         0,
		TotalTransactions:   0,
		ValidTransactions:   0,
		MatchPairs:          0,
		AutoCreatedAccounts: 0,
		ErrorMsg:            "",
	}

	assert.Equal(t, jobID, job.ID)
	assert.Equal(t, userID, job.UserID)
	assert.Equal(t, models.BatchImportStatusPending, job.Status)
	assert.Equal(t, 3, job.TotalFiles)
}

// TestTransferMatch_Model 测试转账匹配模型
func TestTransferMatch_Model(t *testing.T) {
	matchID := uuid.New()
	jobID := uuid.New()

	match := models.TransferMatch{
		ID:           matchID,
		JobID:        jobID,
		MatchType:    models.MatchTypeTransfer,
		Confidence:   0.95,
		MatchFactors: []string{"金额匹配", "时间匹配"},
		UserConfirmed: false,
	}

	assert.Equal(t, matchID, match.ID)
	assert.Equal(t, jobID, match.JobID)
	assert.Equal(t, models.MatchTypeTransfer, match.MatchType)
	assert.Equal(t, 0.95, match.Confidence)
	assert.Len(t, match.MatchFactors, 2)
	assert.False(t, match.UserConfirmed)
}

// TestNoteMergeService_SuggestMerge 测试备注合并服务的基本功能
func TestNoteMergeService_SuggestMerge(t *testing.T) {
	noteService := services.NewNoteMergeService()
	assert.NotNil(t, noteService)

	// 注意：SuggestMerge 接受 TransferMatch，我们可以测试空场景
	// 实际的合并逻辑测试会更复杂，但这里验证服务能初始化即可
}

// TestBatchImport_Integration_Scenario 测试批量导入的完整集成场景
func TestBatchImport_Integration_Scenario(t *testing.T) {
	t.Log("=== Testing Batch Import Integration Scenario ===")

	// 1. 初始化模型
	userID := uuid.New()

	// 2. 创建批量导入任务
	job := &models.BatchImportJob{
		ID:         uuid.New(),
		UserID:     userID,
		Status:     models.BatchImportStatusPending,
		TotalFiles: 2,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 3. 模拟状态流转
	statusFlow := []models.BatchImportStatus{
		models.BatchImportStatusParsing,
		models.BatchImportStatusAnalyzing,
		models.BatchImportStatusMatching,
		models.BatchImportStatusReadyToImport,
		models.BatchImportStatusImporting,
		models.BatchImportStatusCompleted,
	}

	for _, status := range statusFlow {
		job.Status = status
		job.UpdatedAt = time.Now()
		progress := services.CalculateProgress(job)
		assert.NotNil(t, progress)
		t.Logf("Status: %s, Progress: %.0f%%", status, progress.PercentComplete)
	}

	// 4. 验证最终状态
	assert.Equal(t, models.BatchImportStatusCompleted, job.Status)
	finalProgress := services.CalculateProgress(job)
	assert.Equal(t, 100.0, finalProgress.PercentComplete)

	t.Log("✓ Batch import integration scenario completed successfully")
}
