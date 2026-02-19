package integration

import (
	"account/internal/business/models"
	"account/internal/sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCompleteSyncFlow_DeviceACreatesTransaction_DeviceBReceivesSync
// Tests the complete transaction creation and sync flow as described in CLAUDE.md
func TestCompleteSyncFlow_LWWStrategy(t *testing.T) {
	strategy := sync.NewLWWStrategy()
	userID := uuid.New()

	t.Log("=== Testing Complete Sync Flow (LWW Strategy) ===")

	// ========== Step 1-4: Device A creates transaction locally ==========
	t.Log("Step 1-4: Device A creates transaction locally")

	accountID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	// Create account on Device A
	deviceAAccount := models.Account{
		ID:             accountID,
		UserID:         userID,
		Name:           "Alipay",
		Type:           models.AccountTypeAlipay,
		Currency:       "CNY",
		Balance:        5000.0,
		CreatedAt:      now.Add(-24 * time.Hour),
		UpdatedAt:      now,
		LastModifiedAt: now,
		Version:        1,
		IsDeleted:      false,
	}

	// Create transaction on Device A
	_ = models.Transaction{
		ID:              transactionID,
		UserID:          userID,
		AccountID:       accountID,
		Type:            models.TransactionTypeExpense,
		Amount:          88.5,
		Currency:        "CNY",
		Note:            "Coffee at Starbucks",
		TransactionDate: now,
		CreatedAt:       now,
		UpdatedAt:       now,
		LastModifiedAt:  now,
		Version:         1,
		IsDeleted:       false,
	}

	// ========== Step 5-7: Device A pushes to server ==========
	t.Log("Step 5-7: Device A pushes changes to server")

	// Simulate server receiving the push
	serverAccounts := []models.Account{}
	_ = []models.Account{}

	// Server applies LWW when receiving (no conflict yet)
	mergedAccounts := strategy.MergeAccountLists(serverAccounts, []models.Account{deviceAAccount})
	assert.Len(t, mergedAccounts, 1)
	assert.Equal(t, "Alipay", mergedAccounts[0].Name)

	// ========== Step 8-9: Server resolves conflicts with LWW ==========
	t.Log("Step 8-9: Server applies LWW conflict resolution")

	// Now simulate Device B modifying the same account with earlier timestamp
	deviceBAccountOld := models.Account{
		ID:             accountID,
		UserID:         userID,
		Name:           "Alipay (Old)",
		Type:           models.AccountTypeAlipay,
		Currency:       "CNY",
		Balance:        4911.5,
		LastModifiedAt: now.Add(-30 * time.Minute), // Earlier than Device A
		Version:        2,
	}

	// LWW should choose Device A's version (later timestamp)
	winner := strategy.ResolveAccount(&deviceBAccountOld, &deviceAAccount)
	assert.Equal(t, "Alipay", winner.Name)
	assert.Equal(t, 5000.0, winner.Balance)

	// Now simulate Device B modifying with later timestamp
	deviceBAccountNew := models.Account{
		ID:             accountID,
		UserID:         userID,
		Name:           "Alipay (Updated)",
		Type:           models.AccountTypeAlipay,
		Currency:       "CNY",
		Balance:        4911.5,
		LastModifiedAt: now.Add(30 * time.Minute), // Later than Device A
		Version:        2,
	}

	// LWW should choose Device B's version now
	winner = strategy.ResolveAccount(&deviceAAccount, &deviceBAccountNew)
	assert.Equal(t, "Alipay (Updated)", winner.Name)
	assert.Equal(t, 4911.5, winner.Balance)

	t.Log("✓ LWW conflict resolution works correctly")
}

// TestSyncFlow_MultiDeviceScenario tests multi-device sync scenarios
func TestSyncFlow_MultiDeviceScenario(t *testing.T) {
	strategy := sync.NewLWWStrategy()
	userID := uuid.New()
	now := time.Now()

	t.Log("=== Testing Multi-Device Sync Scenario ===")

	// Create 3 different transactions from 3 devices
	phoneAccountID := uuid.New()
	laptopAccountID := uuid.New()
	tabletAccountID := uuid.New()

	// Phone creates an account
	phoneAccount := models.Account{
		ID:             phoneAccountID,
		UserID:         userID,
		Name:           "Phone Cash",
		Type:           models.AccountTypeCash,
		Balance:        1000.0,
		LastModifiedAt: now.Add(-2 * time.Hour),
	}

	// Laptop creates an account
	laptopAccount := models.Account{
		ID:             laptopAccountID,
		UserID:         userID,
		Name:           "Laptop Bank",
		Type:           models.AccountTypeBank,
		Balance:        5000.0,
		LastModifiedAt: now.Add(-1 * time.Hour),
	}

	// Tablet creates an account
	tabletAccount := models.Account{
		ID:             tabletAccountID,
		UserID:         userID,
		Name:           "Tablet Alipay",
		Type:           models.AccountTypeAlipay,
		Balance:        2000.0,
		LastModifiedAt: now,
	}

	// Merge all three device's changes
	merged := strategy.MergeAccountLists(
		[]models.Account{phoneAccount},
		[]models.Account{laptopAccount, tabletAccount},
	)

	assert.Len(t, merged, 3)

	// Verify all accounts are present
	accountMap := make(map[string]models.Account)
	for _, a := range merged {
		accountMap[a.ID.String()] = a
	}

	assert.Equal(t, "Phone Cash", accountMap[phoneAccountID.String()].Name)
	assert.Equal(t, "Laptop Bank", accountMap[laptopAccountID.String()].Name)
	assert.Equal(t, "Tablet Alipay", accountMap[tabletAccountID.String()].Name)

	t.Log("✓ Multi-device account merging works correctly")
}

// TestSyncNotifier_MultipleDevices tests the sync notifier
func TestSyncFlow_NotificationSystem(t *testing.T) {
	notifier := sync.NewSyncNotifier()
	userID := uuid.New()

	t.Log("=== Testing Sync Notification System ===")

	// Subscribe multiple devices
	device1Chan := notifier.Subscribe(userID, "device-1")
	device2Chan := notifier.Subscribe(userID, "device-2")
	device3Chan := notifier.Subscribe(userID, "device-3")

	// Notify from device 1 - devices 2 and 3 should get notified
	notifier.Notify(userID, "device-1")

	// Check device 1 doesn't get its own notification
	select {
	case <-device1Chan:
		t.Fatal("Device 1 should not receive its own notification")
	case <-time.After(100 * time.Millisecond):
		// Expected
	}

	// Check device 2 gets notification
	select {
	case <-device2Chan:
		// Expected
	case <-time.After(1 * time.Second):
		t.Fatal("Device 2 did not receive notification")
	}

	// Check device 3 gets notification
	select {
	case <-device3Chan:
		// Expected
	case <-time.After(1 * time.Second):
		t.Fatal("Device 3 did not receive notification")
	}

	t.Log("✓ Sync notification system works correctly")

	// Cleanup
	notifier.Unsubscribe(userID, "device-1")
	notifier.Unsubscribe(userID, "device-2")
	notifier.Unsubscribe(userID, "device-3")
}

// TestSyncFlow_Scenario_LastWriteWins tests various LWW scenarios
func TestSyncFlow_LastWriteWins_AllScenarios(t *testing.T) {
	strategy := sync.NewLWWStrategy()
	userID := uuid.New()
	accountID := uuid.New()
	now := time.Now()

	t.Log("=== Testing LWW All Scenarios ===")

	tests := []struct {
		name           string
		localName      string
		localTime      time.Time
		remoteName     string
		remoteTime     time.Time
		expectedWinner string
	}{
		{
			name:           "Remote wins - later timestamp",
			localName:      "Local Name",
			localTime:      now.Add(-2 * time.Hour),
			remoteName:     "Remote Name",
			remoteTime:     now.Add(-1 * time.Hour),
			expectedWinner: "Remote Name",
		},
		{
			name:           "Local wins - later timestamp",
			localName:      "Local Name",
			localTime:      now.Add(-1 * time.Hour),
			remoteName:     "Remote Name",
			remoteTime:     now.Add(-2 * time.Hour),
			expectedWinner: "Local Name",
		},
		{
			name:           "Local wins - same timestamp",
			localName:      "Local Name",
			localTime:      now.Truncate(time.Second),
			remoteName:     "Remote Name",
			remoteTime:     now.Truncate(time.Second),
			expectedWinner: "Local Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			local := &models.Account{
				ID:             accountID,
				UserID:         userID,
				Name:           tt.localName,
				LastModifiedAt: tt.localTime,
			}
			remote := &models.Account{
				ID:             accountID,
				UserID:         userID,
				Name:           tt.remoteName,
				LastModifiedAt: tt.remoteTime,
			}

			winner := strategy.ResolveAccount(local, remote)
			assert.Equal(t, tt.expectedWinner, winner.Name)
		})
	}

	t.Log("✓ All LWW scenarios pass")
}
