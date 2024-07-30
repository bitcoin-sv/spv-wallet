package models

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/common"
)

// BlockHeader is a model that represents a BSV block header.
type BlockHeader struct {
	// Model is a common model that contains common fields for all models.
	Model common.OldModel

	// ID is a block header id (hash).
	ID string `json:"id"`
	// Height is a block header height.
	Height uint32 `json:"height"`
	// Time is a block header time (timestamp).
	Time uint32 `json:"time"`
	// Nonce is a block header nonce.
	Nonce uint32 `json:"nonce"`
	// Version is a block header version.
	Version uint32 `json:"version"`
	// HashPreviousBlock is a block header hash of previous block.
	HashPreviousBlock string `json:"hash_previous_block"`
	// HashMerkleRoot is a block header hash merkle tree root.
	HashMerkleRoot string `json:"hash_merkle_root"`
	// Bits contains BSV block header bits no.
	Bits string `json:"bits"`
	// Synec is a time when block header was synced.
	Synced time.Time `json:"synced"`
}
