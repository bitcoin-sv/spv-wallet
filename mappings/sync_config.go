package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
)

// MapToSyncConfigContract will map the sync-config model from bux to the bux-models contract
func MapToSyncConfigContract(sc *bux.SyncConfig) *buxmodels.SyncConfig {
	if sc == nil {
		return nil
	}

	return &buxmodels.SyncConfig{
		Broadcast:        sc.Broadcast,
		BroadcastInstant: sc.BroadcastInstant,
		PaymailP2P:       sc.PaymailP2P,
		SyncOnChain:      sc.SyncOnChain,
	}
}

// MapToSyncConfigBux will map the sync-config model from bux-models to the bux contract
func MapToSyncConfigBux(sc *buxmodels.SyncConfig) *bux.SyncConfig {
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
