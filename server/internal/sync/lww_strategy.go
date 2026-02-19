package sync

import (
	"account/internal/business/models"
	"time"
)

// LWWStrategy implements Last-Write-Wins conflict resolution
// The entity with the later last_modified_at timestamp always wins
type LWWStrategy struct{}

func NewLWWStrategy() *LWWStrategy {
	return &LWWStrategy{}
}

type LWWEntity interface {
	GetLastModifiedAt() time.Time
	GetID() string
	GetVersion() int
	IsDeleted() bool
}

// ResolveAccount resolves conflict between two account versions
// Returns the winning account
func (s *LWWStrategy) ResolveAccount(local, remote *models.Account) *models.Account {
	if local == nil {
		return remote
	}
	if remote == nil {
		return local
	}

	if remote.LastModifiedAt.After(local.LastModifiedAt) {
		return remote
	}
	return local
}

// ResolveCategory resolves conflict between two category versions
// Returns the winning category
func (s *LWWStrategy) ResolveCategory(local, remote *models.Category) *models.Category {
	if local == nil {
		return remote
	}
	if remote == nil {
		return local
	}

	if remote.LastModifiedAt.After(local.LastModifiedAt) {
		return remote
	}
	return local
}

// ResolveTransaction resolves conflict between two transaction versions
// Returns the winning transaction
func (s *LWWStrategy) ResolveTransaction(local, remote *models.Transaction) *models.Transaction {
	if local == nil {
		return remote
	}
	if remote == nil {
		return local
	}

	if remote.LastModifiedAt.After(local.LastModifiedAt) {
		return remote
	}
	return local
}

// MergeAccountLists merges two lists of accounts using LWW strategy
func (s *LWWStrategy) MergeAccountLists(local, remote []models.Account) []models.Account {
	accountMap := make(map[string]models.Account)

	for _, a := range local {
		accountMap[a.ID.String()] = a
	}

	for _, remoteAccount := range remote {
		key := remoteAccount.ID.String()
		if existing, exists := accountMap[key]; exists {
			winner := s.ResolveAccount(&existing, &remoteAccount)
			accountMap[key] = *winner
		} else {
			accountMap[key] = remoteAccount
		}
	}

	result := make([]models.Account, 0, len(accountMap))
	for _, a := range accountMap {
		result = append(result, a)
	}

	return result
}

// MergeCategoryLists merges two lists of categories using LWW strategy
func (s *LWWStrategy) MergeCategoryLists(local, remote []models.Category) []models.Category {
	categoryMap := make(map[string]models.Category)

	for _, c := range local {
		categoryMap[c.ID.String()] = c
	}

	for _, remoteCategory := range remote {
		key := remoteCategory.ID.String()
		if existing, exists := categoryMap[key]; exists {
			winner := s.ResolveCategory(&existing, &remoteCategory)
			categoryMap[key] = *winner
		} else {
			categoryMap[key] = remoteCategory
		}
	}

	result := make([]models.Category, 0, len(categoryMap))
	for _, c := range categoryMap {
		result = append(result, c)
	}

	return result
}

// MergeTransactionLists merges two lists of transactions using LWW strategy
func (s *LWWStrategy) MergeTransactionLists(local, remote []models.Transaction) []models.Transaction {
	transactionMap := make(map[string]models.Transaction)

	for _, t := range local {
		transactionMap[t.ID.String()] = t
	}

	for _, remoteTransaction := range remote {
		key := remoteTransaction.ID.String()
		if existing, exists := transactionMap[key]; exists {
			winner := s.ResolveTransaction(&existing, &remoteTransaction)
			transactionMap[key] = *winner
		} else {
			transactionMap[key] = remoteTransaction
		}
	}

	result := make([]models.Transaction, 0, len(transactionMap))
	for _, t := range transactionMap {
		result = append(result, t)
	}

	return result
}
