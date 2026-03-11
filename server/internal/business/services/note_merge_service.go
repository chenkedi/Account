package services

import (
	"strings"
)

// NoteMergeService 备注合并服务
type NoteMergeService struct{}

// NewNoteMergeService 创建备注合并服务
func NewNoteMergeService() *NoteMergeService {
	return &NoteMergeService{}
}

// MergeSuggestion 合并建议
type MergeSuggestion struct {
	PrimaryNote    string `json:"primary_note"`    // 主备注（支付平台详情）
	SecondaryNote  string `json:"secondary_note"`  // 辅助备注（银行摘要）
	MergedNote     string `json:"merged_note"`     // 合并后的建议
	MergeStrategy  string `json:"merge_strategy"`  // 合并策略
	Reason         string `json:"reason"`          // 合并理由
}

// SuggestMerge 为匹配的交易对建议备注合并
func (s *NoteMergeService) SuggestMerge(match TransferMatch) MergeSuggestion {
	// 确定哪边是银行，哪边是支付平台
	var bankNote, appNote string

	// 简单判断：包含"支付宝"、"微信"、"京东"等关键词的是平台备注
	isPlatformNote := func(note string) bool {
		platformKeywords := []string{"支付宝", "微信", "京东", "美团", "滴滴", "拼多多"}
		lowerNote := strings.ToLower(note)
		for _, kw := range platformKeywords {
			if strings.Contains(lowerNote, strings.ToLower(kw)) {
				return true
			}
		}
		return false
	}

	// 判断转出和转入哪个是平台，哪个是银行
	if isPlatformNote(match.OutTx.Note) {
		appNote = match.OutTx.Note
		bankNote = match.InTx.Note
	} else if isPlatformNote(match.InTx.Note) {
		appNote = match.InTx.Note
		bankNote = match.OutTx.Note
	} else {
		// 无法判断，使用较长的作为平台备注（通常平台备注更详细）
		if len(match.OutTx.Note) > len(match.InTx.Note) {
			appNote = match.OutTx.Note
			bankNote = match.InTx.Note
		} else {
			appNote = match.InTx.Note
			bankNote = match.OutTx.Note
		}
	}

	return s.createMergeSuggestion(bankNote, appNote)
}

// createMergeSuggestion 创建合并建议
func (s *NoteMergeService) createMergeSuggestion(bankNote, appNote string) MergeSuggestion {
	// 清理备注
	bankNote = strings.TrimSpace(bankNote)
	appNote = strings.TrimSpace(appNote)

	// 如果应用备注为空，使用银行备注
	if appNote == "" {
		return MergeSuggestion{
			PrimaryNote:   bankNote,
			SecondaryNote: "",
			MergedNote:    bankNote,
			MergeStrategy: "bank_only",
			Reason:        "无平台详情，使用银行摘要",
		}
	}

	// 如果银行备注为空，使用应用备注
	if bankNote == "" {
		return MergeSuggestion{
			PrimaryNote:   appNote,
			SecondaryNote: "",
			MergedNote:    appNote,
			MergeStrategy: "app_only",
			Reason:        "无银行摘要，使用平台详情",
		}
	}

	// 如果两者相同，直接使用
	if bankNote == appNote {
		return MergeSuggestion{
			PrimaryNote:   appNote,
			SecondaryNote: bankNote,
			MergedNote:    appNote,
			MergeStrategy: "identical",
			Reason:        "两边备注相同",
		}
	}

	// 如果应用备注已包含银行备注，使用应用备注
	if strings.Contains(appNote, bankNote) {
		return MergeSuggestion{
			PrimaryNote:   appNote,
			SecondaryNote: bankNote,
			MergedNote:    appNote,
			MergeStrategy: "app_contains_bank",
			Reason:        "平台详情已包含银行摘要",
		}
	}

	// 智能合并策略
	merged := s.smartMerge(appNote, bankNote)

	return MergeSuggestion{
		PrimaryNote:   appNote,
		SecondaryNote: bankNote,
		MergedNote:    merged,
		MergeStrategy: "smart_merge",
		Reason:        "智能合并：平台详情为主，银行摘要补充",
	}
}

// smartMerge 智能合并两个备注
func (s *NoteMergeService) smartMerge(primary, secondary string) string {
	// 清理
	primary = strings.TrimSpace(primary)
	secondary = strings.TrimSpace(secondary)

	// 如果主要备注已经很详细，直接使用
	if len(primary) > 30 {
		// 但确保包含次要备注的关键信息
		keywords := s.extractKeywords(secondary)
		for _, kw := range keywords {
			if !strings.Contains(primary, kw) && len(kw) > 2 {
				primary = primary + " | " + secondary
				break
			}
		}
		return primary
	}

	// 合并两者，用分隔符
	if secondary != "" {
		return primary + " | " + secondary
	}
	return primary
}

// extractKeywords 提取关键词
func (s *NoteMergeService) extractKeywords(text string) []string {
	// 简单的关键词提取：去掉停用词，保留有意义的部分
	stopWords := map[string]bool{
		"的": true, "了": true, "在": true, "是": true, "我": true,
		"和": true, "与": true, "或": true, "等": true, "及": true,
	}

	var keywords []string
	words := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == ',' || r == '，' || r == '|' || r == '/'
	})

	for _, word := range words {
		word = strings.TrimSpace(word)
		if word != "" && !stopWords[word] && len(word) >= 2 {
			keywords = append(keywords, word)
		}
	}

	return keywords
}
