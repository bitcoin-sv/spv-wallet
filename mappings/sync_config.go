package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
)

func MapToSyncConfigContract(sc *bux.SyncConfig) *buxmodels.SyncConfig {
	return &buxmodels.SyncConfig{
		Broadcast:        sc.Broadcast,
		BroadcastInstant: sc.BroadcastInstant,
		PaymailP2P:       sc.PaymailP2P,
		SyncOnChain:      sc.SyncOnChain,
	}
}

func MapToSyncConfigBux(sc *buxmodels.SyncConfig) *bux.SyncConfig {
	return &bux.SyncConfig{
		Broadcast:        sc.Broadcast,
		BroadcastInstant: sc.BroadcastInstant,
		PaymailP2P:       sc.PaymailP2P,
		SyncOnChain:      sc.SyncOnChain,
	}
}
