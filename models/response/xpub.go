package response

// Xpub is a model that represents a xpub.
type Xpub struct {
	// Model is a common model that contains common fields for all models.
	Model

	// ID is a hash of the xpub.
	ID string `json:"id" example:"bb8593f85ef8056a77026ad415f02128f3768906de53e9e8bf8749fe2d66cf50"`
	// CurrentBalance is a xpub's current balance.
	CurrentBalance uint64 `json:"currentBalance" example:"1234"`
	// NextInternalNum is the index derivation number use to generate NEXT internal xPub (internal xPub are used for change destinations).
	NextInternalNum uint32 `json:"nextInternalNum" example:"0"`
	// NextExternalNum is the index derivation number use to generate NEXT external xPub (external xPub are used for address destinations).
	NextExternalNum uint32 `json:"nextExternalNum" example:"0"`
}
