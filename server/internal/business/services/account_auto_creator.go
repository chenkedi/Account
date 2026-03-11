package services

import (
	"strings"

	"account/internal/business/models"
	"account/internal/data/repository"
	"github.com/google/uuid"
)

// AccountAutoCreator 自动账户创建器
type AccountAutoCreator struct {
	accountRepo *repository.AccountRepository
}

// NewAccountAutoCreator 创建自动账户创建器
func NewAccountAutoCreator(accountRepo *repository.AccountRepository) *AccountAutoCreator {
	return &AccountAutoCreator{
		accountRepo: accountRepo,
	}
}

// CreateAccountsFromHints 根据账户线索批量创建账户
func (c *AccountAutoCreator) CreateAccountsFromHints(
	userID uuid.UUID,
	hints []models.AccountHint,
	existingAccounts []models.Account,
) ([]models.Account, []models.Account, error) {
	var createdAccounts []models.Account
	var skippedAccounts []models.Account

	// 去重处理：相同账户线索只处理一次
	seenHints := make(map[string]bool)

	for _, hint := range hints {
		// 生成唯一标识
		hintKey := c.generateHintKey(hint)
		if seenHints[hintKey] {
			continue
		}
		seenHints[hintKey] = true

		// 检查是否已存在相似账户
		existing := c.findSimilarAccount(hint, existingAccounts)
		if existing != nil {
			// 更新现有账户信息（如余额）
			c.updateExistingAccount(existing, hint)
			skippedAccounts = append(skippedAccounts, *existing)
			continue
		}

		// 创建新账户
		newAccount := c.buildAccountFromHint(hint, userID)
		created, err := c.accountRepo.Create(
			newAccount.UserID,
			newAccount.Name,
			newAccount.Type,
			newAccount.Currency,
			newAccount.TailNumber,
			newAccount.Balance,
		)
		if err != nil {
			continue // 记录错误但继续处理其他账户
		}

		createdAccounts = append(createdAccounts, *created)
		existingAccounts = append(existingAccounts, *created) // 更新已存在列表
	}

	return createdAccounts, skippedAccounts, nil
}

// generateHintKey 生成账户线索的唯一键
func (c *AccountAutoCreator) generateHintKey(hint models.AccountHint) string {
	parts := []string{
		string(hint.Source),
		hint.BankName,
		hint.AccountNumber,
		hint.CardType,
	}
	return strings.Join(parts, "|")
}

// findSimilarAccount 查找相似的现有账户
func (c *AccountAutoCreator) findSimilarAccount(
	hint models.AccountHint,
	existingAccounts []models.Account,
) *models.Account {
	for _, account := range existingAccounts {
		similarity := c.calculateSimilarity(hint, account)
		if similarity >= 0.85 { // 85%相似度阈值
			return &account
		}
	}
	return nil
}

// calculateSimilarity 计算账户线索与现有账户的相似度
// 优先级：尾号匹配(50%) > 名称相似度(30%) > 类型匹配(20%)
func (c *AccountAutoCreator) calculateSimilarity(
	hint models.AccountHint,
	account models.Account,
) float64 {
	var score float64

	// 1. 账号尾号匹配（权重：50%）- 最高优先级
	if hint.AccountNumber != "" && account.TailNumber != "" {
		if normalizeTailNumber(hint.AccountNumber) == normalizeTailNumber(account.TailNumber) {
			score += 0.50
		}
	}

	// 2. 账户名称相似度（权重：30%）
	nameSim := c.stringSimilarity(
		strings.ToLower(hint.AccountName),
		strings.ToLower(account.Name),
	)
	score += nameSim * 0.30

	// 3. 账户类型匹配（权重：20%）
	if c.accountTypesMatch(hint, account) {
		score += 0.20
	}

	return score
}

// stringSimilarity 计算两个字符串的相似度
func (c *AccountAutoCreator) stringSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	if strings.Contains(s1, s2) || strings.Contains(s2, s1) {
		return 0.8
	}

	// 检查关键词重叠
	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)

	if len(words1) == 0 || len(words2) == 0 {
		return 0
	}

	matches := 0
	for _, w1 := range words1 {
		for _, w2 := range words2 {
			if w1 == w2 {
				matches++
				break
			}
		}
	}

	return float64(matches) / float64(max(len(words1), len(words2)))
}

// accountTypesMatch 检查账户类型是否匹配
func (c *AccountAutoCreator) accountTypesMatch(
	hint models.AccountHint,
	account models.Account,
) bool {
	// 映射线索类型到账户类型
	var expectedType models.AccountType

	switch hint.Source {
	case models.ImportSourceAlipay, models.ImportSourceWeChat, models.ImportSourceJD:
		expectedType = models.AccountTypeAlipay // 或 AccountTypeEWallet
	case models.ImportSourceBank:
		if hint.CardType == "信用卡" || hint.CardType == "贷记卡" {
			expectedType = models.AccountTypeCredit
		} else {
			expectedType = models.AccountTypeBank
		}
	}

	return account.Type == expectedType
}

// updateExistingAccount 更新现有账户信息
func (c *AccountAutoCreator) updateExistingAccount(
	account *models.Account,
	hint models.AccountHint,
) {
	// 如果账单中有余额信息且比现有数据新，则更新余额
	if hint.Balance != 0 && account.Balance != hint.Balance {
		account.Balance = hint.Balance
	}

	// 如果现有账户没有尾号但hint有，则更新尾号
	if account.TailNumber == "" && hint.AccountNumber != "" {
		account.TailNumber = normalizeTailNumber(hint.AccountNumber)
	}
}

// buildAccountFromHint 根据线索构建账户对象
func (c *AccountAutoCreator) buildAccountFromHint(
	hint models.AccountHint,
	userID uuid.UUID,
) models.Account {
	// 确定账户类型
	var accountType models.AccountType
	switch hint.Source {
	case models.ImportSourceAlipay:
		accountType = models.AccountTypeAlipay
	case models.ImportSourceWeChat:
		accountType = models.AccountTypeWeChat
	case models.ImportSourceJD:
		// 京东支付通常关联银行卡，根据卡类型判断
		if hint.CardType == "信用卡" || hint.CardType == "贷记卡" {
			accountType = models.AccountTypeCredit
		} else {
			accountType = models.AccountTypeBank
		}
	case models.ImportSourceBank:
		if hint.CardType == "信用卡" || hint.CardType == "贷记卡" {
			accountType = models.AccountTypeCredit
		} else {
			accountType = models.AccountTypeBank
		}
	default:
		accountType = models.AccountTypeOther
	}

	// 构建账户名称
	nameParts := []string{}
	if hint.BankName != "" {
		nameParts = append(nameParts, hint.BankName)
	}
	if hint.CardType != "" {
		nameParts = append(nameParts, hint.CardType)
	}
	if hint.AccountNumber != "" {
		nameParts = append(nameParts, "("+hint.AccountNumber+")")
	}

	name := strings.Join(nameParts, "")
	if name == "" {
		name = "未知账户"
	}

	return models.Account{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           name,
		Type:           accountType,
		TailNumber:     normalizeTailNumber(hint.AccountNumber),
		Currency:       "CNY",
		Balance:        hint.Balance,
		IsDeleted:      false,
	}
}

// max 返回较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// normalizeTailNumber 标准化账户尾号（去除"尾号"、"****"等前缀，只保留数字）
func normalizeTailNumber(tailNumber string) string {
	// 去除常见前缀
	tailNumber = strings.TrimSpace(tailNumber)
	tailNumber = strings.ReplaceAll(tailNumber, "尾号", "")
	tailNumber = strings.ReplaceAll(tailNumber, "卡号", "")
	tailNumber = strings.ReplaceAll(tailNumber, "账号", "")
	tailNumber = strings.ReplaceAll(tailNumber, "****", "")
	tailNumber = strings.ReplaceAll(tailNumber, "***", "")
	tailNumber = strings.ReplaceAll(tailNumber, "**", "")
	tailNumber = strings.ReplaceAll(tailNumber, "*", "")

	// 去除空白
	tailNumber = strings.TrimSpace(tailNumber)

	// 只保留数字
	var result strings.Builder
	for _, r := range tailNumber {
		if r >= '0' && r <= '9' {
			result.WriteRune(r)
		}
	}

	return result.String()
}
