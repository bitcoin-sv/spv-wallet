package models

// MerkleRoot holds the content of the synced Merkle root response
type MerkleRoot struct {
	MerkleRoot  string `json:"merkleRoot"`
	BlockHeight int    `json:"blockHeight"`
}

// MerkleRootsBHSResponse is a type that Block Header Service (BHS) returns when queried for merkleroots
type MerkleRootsBHSResponse = ExclusiveStartKeyPage[[]MerkleRoot]
