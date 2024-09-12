package response

// SyncConfig contains sync configuration flags.
type SyncConfig struct {
	// Broadcast is a flag that indicates whether to broadcast transaction or not.
	Broadcast bool `json:"broadcast"`
	// BroadcastInstant is a flag that indicates whether to broadcast transaction instantly or not.
	BroadcastInstant bool `json:"broadcastInstant"`
	// PaymailP2P is a flag that indicates whether to use paymail p2p or not.
	PaymailP2P bool `json:"paymailP2p"`
	// SyncOnChain is a flag that indicates whether to sync transaction on chain or not.
	SyncOnChain bool `json:"syncOnChain"`
}
