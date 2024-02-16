package mappings

import (
	"github.com/bitcoin-sv/bux"
	spvwalletmodels "github.com/bitcoin-sv/bux-models"
)

// MapToSyncConfigContract will map the sync-config model from spv-wallet to the spv-wallet-models contract
func MapToSyncConfigContract(sc *bux.SyncConfig) *spvwalletmodels.SyncConfig {
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
func MapToSyncConfigSPV(sc *spvwalletmodels.SyncConfig) *bux.SyncConfig {
	if sc == nil {
		return nil
	}

	return &bux.SyncConfig{
		Broadcast:        sc.Broadcast,
		BroadcastInstant: sc.BroadcastInstant,
		PaymailP2P:       sc.PaymailP2P,
		SyncOnChain:      sc.SyncOnChain,
	}
}
