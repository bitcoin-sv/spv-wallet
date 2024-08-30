package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToOldSyncConfigContract will map the sync-config model from spv-wallet to the spv-wallet-models contract
func MapToOldSyncConfigContract(sc *engine.SyncConfig) *models.SyncConfig {
	if sc == nil {
		return nil
	}

	return &models.SyncConfig{
		Broadcast:        sc.Broadcast,
		BroadcastInstant: sc.BroadcastInstant,
		PaymailP2P:       sc.PaymailP2P,
		SyncOnChain:      sc.SyncOnChain,
	}
}

// MapOldSyncConfigModelToEngine will map the sync-config model from spv-wallet-models to the spv-wallet contract
func MapOldSyncConfigModelToEngine(sc *models.SyncConfig) *engine.SyncConfig {
	if sc == nil {
		return nil
	}

	return &engine.SyncConfig{
		Broadcast:        sc.Broadcast,
		BroadcastInstant: sc.BroadcastInstant,
		PaymailP2P:       sc.PaymailP2P,
		SyncOnChain:      sc.SyncOnChain,
	}
}
