package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	spvwalletmodels "github.com/bitcoin-sv/spv-wallet/models"
)

// MapToSyncConfigContract will map the sync-config model from spv-wallet to the spv-wallet-models contract
func MapToSyncConfigContract(sc *engine.SyncConfig) *spvwalletmodels.SyncConfig {
	if sc == nil {
		return nil
	}

	return &spvwalletmodels.SyncConfig{
		Broadcast:        sc.Broadcast,
		BroadcastInstant: sc.BroadcastInstant,
		PaymailP2P:       sc.PaymailP2P,
		SyncOnChain:      sc.SyncOnChain,
	}
}

// MapToSyncConfigSPV will map the sync-config model from spv-wallet-models to the spv-wallet contract
func MapToSyncConfigSPV(sc *spvwalletmodels.SyncConfig) *engine.SyncConfig {
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
