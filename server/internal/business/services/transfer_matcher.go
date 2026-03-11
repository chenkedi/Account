package services

import (
	"math"
	"strings"
	"time"

	"account/internal/business/models"
	"github.com/google/uuid"
)

// TransferMatcher 跨文件转账匹配器
// 使用账户尾号匹配作为最高优先级的匹配策略
type TransferMatcher struct {
	timeWindow      time.Duration // 时间窗口（默认24小时）
	amountTolerance float64       // 金额容差（默认0.01）
}

// NewTransferMatcher 创建匹配器
func NewTransferMatcher() *TransferMatcher {
	return &TransferMatcher{
		timeWindow:      24 * time.Hour,
		amountTolerance: 0.01,
	}
}

// MatchResult 匹配结果
type MatchResult struct {
	Matches      []TransferMatch
	UnmatchedOut []TransactionInfo // 未匹配的转出
	UnmatchedIn  []TransactionInfo // 未匹配的转入
}

// TransferMatch 转账匹配对
type TransferMatch struct {
	ID           uuid.UUID
	OutTx        TransactionInfo
	InTx         TransactionInfo
	Confidence   float64  // 匹配置信度 0-1
	MatchFactors []string // 匹配依据
}

// TransactionInfo 交易信息（用于匹配）
type TransactionInfo struct {
	FileID             uuid.UUID
	TransactionIndex   int
	Date               time.Time
	Amount             float64
	Type               string // expense, income
	AccountName        string
	AccountNumber      string // 账户尾号
	Counterparty       string
	CounterpartyNumber string // 对方账户尾号
	Note               string
	RawData            map[string]string
}

// FindMatches 在所有文件中查找匹配的转账记录
// 核心算法：账户尾号匹配作为最高优先级（40%权重）
func (m *TransferMatcher) FindMatches(files []FileTransactions) *MatchResult {
	var matches []TransferMatch

	// 收集所有转出和转入记录
	var outTxs, inTxs []TransactionInfo

	for _, file := range files {
		for i, tx := range file.Transactions {
			info := TransactionInfo{
				FileID:             file.FileID,
				TransactionIndex:   i,
				Date:               tx.TransactionDate,
				Amount:             tx.Amount,
				Type:               string(tx.Type),
				AccountName:        tx.AccountName,
				AccountNumber:      tx.ParsedAccountNumber,
				Counterparty:       tx.Counterparty,
				CounterpartyNumber: tx.RelatedAccountNumber,
				Note:               tx.Note,
				RawData:            tx.RawData,
			}

			// 判断是否为转账类型
			if m.isTransferOut(tx) {
				info.Type = "expense"
				outTxs = append(outTxs, info)
			} else if m.isTransferIn(tx) {
				info.Type = "income"
				inTxs = append(inTxs, info)
			}
		}
	}

	// 匹配算法
	matchedOut := make(map[int]bool)
	matchedIn := make(map[int]bool)

	for i, outTx := range outTxs {
		if matchedOut[i] {
			continue
		}

		bestMatch := -1
		bestConfidence := 0.0
		var bestFactors []string

		for j, inTx := range inTxs {
			if matchedIn[j] {
				continue
			}

			confidence, factors := m.calculateMatchConfidence(outTx, inTx)

			if confidence > bestConfidence {
				bestConfidence = confidence
				bestMatch = j
				bestFactors = factors
			}
		}

		// 如果最佳匹配的置信度足够高，则确认为匹配
		if bestMatch >= 0 && bestConfidence >= 0.8 {
			match := TransferMatch{
				ID:           uuid.New(),
				OutTx:        outTx,
				InTx:         inTxs[bestMatch],
				Confidence:   bestConfidence,
				MatchFactors: bestFactors,
			}
			matches = append(matches, match)
			matchedOut[i] = true
			matchedIn[bestMatch] = true
		}
	}

	// 收集未匹配的记录
	var unmatchedOut, unmatchedIn []TransactionInfo
	for i, tx := range outTxs {
		if !matchedOut[i] {
			unmatchedOut = append(unmatchedOut, tx)
		}
	}
	for i, tx := range inTxs {
		if !matchedIn[i] {
			unmatchedIn = append(unmatchedIn, tx)
		}
	}

	return &MatchResult{
		Matches:      matches,
		UnmatchedOut: unmatchedOut,
		UnmatchedIn:  unmatchedIn,
	}
}

// isTransferOut 判断是否为转出记录
func (m *TransferMatcher) isTransferOut(tx models.ParsedTransaction) bool {
	// 类型为支出且包含转账关键词
	if tx.Type == models.TransactionTypeExpense {
		note := strings.ToLower(tx.Note + " " + tx.Counterparty)
		transferKeywords := []string{"转账", "转出", "汇款", "跨行", "行内"}
		for _, kw := range transferKeywords {
			if strings.Contains(note, kw) {
				return true
			}
		}
	}
	return false
}

// isTransferIn 判断是否为转入记录
func (m *TransferMatcher) isTransferIn(tx models.ParsedTransaction) bool {
	// 类型为收入且包含转入关键词
	if tx.Type == models.TransactionTypeIncome {
		note := strings.ToLower(tx.Note + " " + tx.Counterparty)
		inKeywords := []string{"转账", "转入", "汇款", "跨行", "来账"}
		for _, kw := range inKeywords {
			if strings.Contains(note, kw) {
				return true
			}
		}
	}
	return false
}

// calculateMatchConfidence 计算匹配置信度
// 优先级排序：账户尾号匹配(40%) > 金额匹配(30%) > 时间匹配(20%) > 账户名称关联(10%)
func (m *TransferMatcher) calculateMatchConfidence(outTx, inTx TransactionInfo) (float64, []string) {
	var confidence float64
	var factors []string

	// ===== 最高优先级：账户尾号匹配（权重：40%）=====
	// 这是最强的关联证据，因为支付宝/微信/京东账单都会记录消费的账户尾号
	tailNumberMatched := false

	// 场景1：转出账户尾号 == 转入记录的对方账户尾号
	if outTx.AccountNumber != "" && inTx.CounterpartyNumber != "" {
		if normalizeTailNumber(outTx.AccountNumber) == normalizeTailNumber(inTx.CounterpartyNumber) {
			confidence += 0.40
			factors = append(factors, "账户尾号精确匹配(转出账户=转入对方账户)")
			tailNumberMatched = true
		}
	}

	// 场景2：转入账户尾号 == 转出记录的对方账户尾号
	if !tailNumberMatched && inTx.AccountNumber != "" && outTx.CounterpartyNumber != "" {
		if normalizeTailNumber(inTx.AccountNumber) == normalizeTailNumber(outTx.CounterpartyNumber) {
			confidence += 0.40
			factors = append(factors, "账户尾号精确匹配(转入账户=转出对方账户)")
			tailNumberMatched = true
		}
	}

	// 场景3：两边都记录了相同的尾号（同一账户的不同交易记录）
	if !tailNumberMatched && outTx.AccountNumber != "" && inTx.AccountNumber != "" {
		if normalizeTailNumber(outTx.AccountNumber) == normalizeTailNumber(inTx.AccountNumber) {
			// 这种情况可能是同一账户的进出，不一定是转账匹配
			// 给较低的权重，结合其他因素判断
			confidence += 0.15
			factors = append(factors, "同一账户尾号(可能为内部转账)")
		}
	}

	// ===== 第二优先级：金额匹配（权重：30%）=====
	amountDiff := math.Abs(outTx.Amount - inTx.Amount)
	if amountDiff <= m.amountTolerance {
		confidence += 0.30
		factors = append(factors, "金额完全匹配")
	} else if amountDiff <= 1.0 {
		// 小额差异可能是手续费
		confidence += 0.25
		factors = append(factors, "金额基本匹配(可能有手续费)")
	} else if amountDiff <= 5.0 {
		// 较大差异可能是跨行手续费
		confidence += 0.15
		factors = append(factors, "金额近似(可能有较高手续费)")
	}

	// ===== 第三优先级：时间匹配（权重：20%）=====
	timeDiff := inTx.Date.Sub(outTx.Date)
	if timeDiff >= 0 && timeDiff <= 30*time.Minute {
		confidence += 0.20
		factors = append(factors, "时间高度匹配(30分钟内)")
	} else if timeDiff > 0 && timeDiff <= 2*time.Hour {
		confidence += 0.18
		factors = append(factors, "时间较匹配(2小时内)")
	} else if timeDiff > 0 && timeDiff <= 24*time.Hour {
		confidence += 0.12
		factors = append(factors, "时间基本匹配(24小时内)")
	}

	// ===== 第四优先级：账户名称关联（权重：10%）=====
	if strings.Contains(inTx.Counterparty, outTx.AccountName) ||
		strings.Contains(outTx.AccountName, inTx.Counterparty) ||
		strings.Contains(inTx.Counterparty, outTx.AccountName) ||
		strings.Contains(outTx.Counterparty, inTx.AccountName) {
		confidence += 0.10
		factors = append(factors, "账户名称或银行名称关联")
	}

	return confidence, factors
}

// FileTransactions represents transactions from a single file
type FileTransactions struct {
	FileID       uuid.UUID
	Source       models.ImportSource
	FileName     string
	Transactions []models.ParsedTransaction
}
