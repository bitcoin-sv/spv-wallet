package chainmodels

// MerkleRootConfirmationState represents the state of each Merkle Root verification
// process and can be one of three values: Confirmed, Invalid and UnableToVerify.
type MerkleRootConfirmationState string

const (
	// MRConfirmed state occurs when Merkle Root is found in the longest chain.
	MRConfirmed MerkleRootConfirmationState = "CONFIRMED"
	// MRInvalid state occurs when Merkle Root is not found in the longest chain.
	MRInvalid MerkleRootConfirmationState = "INVALID"
	// MRUnableToVerify state occurs when Block Headers Service is behind in synchronization with the longest chain.
	MRUnableToVerify MerkleRootConfirmationState = "UNABLE_TO_VERIFY"
)

// MerkleRootConfirmation is a confirmation
// of merkle roots inclusion in the longest chain.
type MerkleRootConfirmation struct {
	Hash         string                      `json:"blockHash"`
	BlockHeight  uint64                      `json:"blockHeight"`
	MerkleRoot   string                      `json:"merkleRoot"`
	Confirmation MerkleRootConfirmationState `json:"confirmation"`
}

// MerkleRootsConfirmations is an API response for confirming
// merkle roots inclusion in the longest chain.
type MerkleRootsConfirmations struct {
	ConfirmationState MerkleRootConfirmationState `json:"confirmationState"`
	Confirmations     []MerkleRootConfirmation    `json:"confirmations"`
}
