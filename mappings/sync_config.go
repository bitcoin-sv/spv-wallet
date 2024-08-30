package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// MapToSyncConfigContract will map the sync-config model from spv-wallet to the spv-wallet-models contract
func MapToSyncConfigContract(sc *engine.SyncConfig) *response.SyncConfig {
	if sc == nil {
		return nil
	}

	return &response.SyncConfig{
		Broadcast:        sc.Broadcast,
		BroadcastInstant: sc.BroadcastInstant,
		PaymailP2P:       sc.PaymailP2P,
		SyncOnChain:      sc.SyncOnChain,
	}
}

// MapSyncConfigModelToEngine will map the sync-config model from spv-wallet-models to the spv-wallet contract
func MapSyncConfigModelToEngine(sc *response.SyncConfig) *engine.SyncConfig {
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
