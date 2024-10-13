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

// MerkleRootsConfirmations is an API response for confirming
// merkle roots inclusion in the longest chain.
type MerkleRootsConfirmations struct {
	ConfirmationState MerkleRootConfirmationState `json:"confirmationState"`
	// BHS also returns Confirmations array - but it's not used in the code here
}
