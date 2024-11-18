package engine

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/bitcoin-sv/go-sdk/chainhash"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

const maxBumpHeight = 64

// BUMPs represents a slice of BUMPs - BSV Unified Merkle Paths
type BUMPs []*trx.MerklePath

// MerklePath wraps trx.MerklePath from which is a BSV Unified Merkle Path
type MerklePath struct{ trx.MerklePath }

// MerklePaths represents a slice of MerklePath
type MerklePaths []*trx.MerklePath

// BUMP represents BUMP (BSV Unified Merkle Path) format
type BUMP struct {
	BlockHeight uint64       `json:"blockHeight,string"`
	Path        [][]BUMPLeaf `json:"path"`
}

// BUMPLeaf represents each BUMP path element
type BUMPLeaf struct {
	Offset    uint64 `json:"offset,string"`
	Hash      string `json:"hash,omitempty"`
	TxID      bool   `json:"txid,omitempty"`
	Duplicate bool   `json:"duplicate,omitempty"`
}

// CalculateMergedBUMPSDK calculates Merged BUMP from a slice of BUMPs
func CalculateMergedBUMPSDK(bumps []trx.MerklePath) (*trx.MerklePath, error) {
	if len(bumps) == 0 {
		return nil, errors.New("no BUMPs provided")
	}

	// Initialize merged BUMP as a copy of the first element in the list
	mergedBump := bumps[0]

	// Check block height consistency and maximum bump height constraints
	blockHeight := mergedBump.BlockHeight
	bumpHeight := len(mergedBump.Path)
	if bumpHeight > maxBumpHeight {
		return nil, spverrors.Newf("BUMP cannot be higher than %d", maxBumpHeight)
	}

	for _, b := range bumps {
		if bumpHeight != len(b.Path) {
			allBumps := ""
			for _, bb := range bumps {
				allBumps += fmt.Sprintf("%+v\n", bb)
			}
			return nil, spverrors.Newf("merged BUMP cannot be obtained from Merkle Proofs of different heights: %s", allBumps)
		}
		if b.BlockHeight != blockHeight {
			return nil,
				spverrors.Newf("cannot merge BUMPs from different blocks")
		}
		if len(b.Path) == 0 {
			return nil,
				spverrors.Newf("empty BUMP given")
		}
	}

	bump := BUMP{BlockHeight: blockHeight}
	bump.Path = make([][]BUMPLeaf, bumpHeight)
	bump.allNodes = make([]map[uint64]bool, bumpHeight)
	for i := range bump.allNodes {
		bump.allNodes[i] = make(map[uint64]bool, 0)
	}

	merkleRoot, err := bumps[0].calculateMerkleRoot()
	if err != nil {
		return nil, err
	}

	for _, b := range bumps {
		mr, err := b.calculateMerkleRoot()
		if err != nil {
			return nil, err
		}

		if merkleRoot != mr {
			return nil, spverrors.Newf("different merkle roots in BUMPs")
		}

		err = bump.add(b)
		if err != nil {
			return nil, err
		}
	}

	for _, p := range bump.Path {
		sort.Slice(p, func(i, j int) bool {
			return p[i].Offset < p[j].Offset
		})
	}

	return &bump, nil
}

func (bump *BUMP) add(b BUMP) error {
	if len(bump.Path) != len(b.Path) {
		return spverrors.Newf("cannot merge BUMPs of different heights")
	}

	for i := range b.Path {
		for _, v := range b.Path[i] {
			_, value := bump.allNodes[i][v.Offset]
			if !value {
				bump.Path[i] = append(bump.Path[i], v)
				bump.allNodes[i][v.Offset] = true
				continue
			}
			if i == 0 && value && v.TxID {
				for j := range bump.Path[i] {
					if bump.Path[i][j].Offset == v.Offset {
						bump.Path[i][j] = v
					}
				}
			}
		}
	}

	return nil
}

func (bump *BUMP) calculateMerkleRoot() (string, error) {
	merkleRoot := ""

	for _, bumpPathElement := range bump.Path[0] {
		if bumpPathElement.TxID {
			calcMerkleRoot, err := calculateMerkleRoot(bumpPathElement, bump)
			if err != nil {
				return "", err
			}

			if merkleRoot == "" {
				merkleRoot = calcMerkleRoot
				continue
			}

			if calcMerkleRoot != merkleRoot {
				return "", spverrors.Newf("different merkle roots for the same block")
			}
		}
	}
	return merkleRoot, nil
}

// calculateMerkleRoots will calculate one merkle root for tx in the BUMPLeaf
func calculateMerkleRoot(baseLeaf BUMPLeaf, bump *BUMP) (string, error) {
	calculatedHash := baseLeaf.Hash
	offset := baseLeaf.Offset

	for _, bLevel := range bump.Path {
		newOffset := getOffsetPair(offset)
		leafInPair := findLeafByOffset(newOffset, bLevel)
		if leafInPair == nil {
			return "", spverrors.Newf("could not find pair")
		}

		leftNode, rightNode := prepareNodes(baseLeaf, offset, *leafInPair, newOffset)

		str, err := utils.MerkleTreeParentStr(leftNode, rightNode)
		if err != nil {
			return "", spverrors.Wrapf(err, "failed to calculate merkle tree parent for %s and %s", leftNode, rightNode)
		}
		calculatedHash = str

		offset = offset / 2

		baseLeaf = BUMPLeaf{
			Hash:   calculatedHash,
			Offset: offset,
	// Iterate and merge each BUMP into the mergedBump
	for _, bump := range bumps[1:] {
		if bump.BlockHeight != blockHeight {
			return nil, spverrors.Newf("inconsistent block heights in BUMPs")
		}
		if err := mergedBump.Combine(&bump); err != nil {
			return nil, spverrors.Wrapf(err, "failed to combine BUMPs")
		}
	}

	return &mergedBump, nil
}

// Bytes returns BUMPs bytes
func (bumps *BUMPs) Bytes() []byte {
	var buff bytes.Buffer

	for _, bump := range *bumps {
		bytes, _ := hex.DecodeString(bump.Hex())
		buff.Write(bytes)
	}

	return buff.Bytes()
}



// Scan scan value into Json, implements sql.Scanner interface
func (bumps *BUMPs) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	byteValue, err := utils.ToByteArray(value)
	if err != nil || bytes.Equal(byteValue, []byte("")) || bytes.Equal(byteValue, []byte("\"\"")) {
		return nil
	}

	err = json.Unmarshal(byteValue, &bumps)
	return spverrors.Wrapf(err, "failed to parse BUMPs from JSON, data: %v", value)
}

// Value return json value, implement driver.Valuer interface
func (bumps BUMPs) Value() (driver.Value, error) {
	if reflect.DeepEqual(bumps, BUMPs{}) {
		return nil, nil
	}
	marshal, err := json.Marshal(bumps)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to convert BUMPs to JSON, data: %v", bumps)
	}

	return string(marshal), nil
}

// ToMerklePath converts BUMP to trx.MerklePath
func (b *BUMP) ToMerklePath() (*trx.MerklePath, error) {
	blockHeight, err := conv.Uint64ToUint32(b.BlockHeight)
	if err != nil {
		return nil, spverrors.Wrapf(err, "error in ToMerklePath: failed to convert block height to uint32")
	}
	mp := &trx.MerklePath{
		BlockHeight: blockHeight,
		Path:        make([][]*trx.PathElement, len(b.Path)),
	}
	for i, level := range b.Path {
		mp.Path[i] = make([]*trx.PathElement, len(level))
		for j, leaf := range level {
			hash, err := chainhash.NewHashFromHex(leaf.Hash)
			if err != nil {
				return nil, spverrors.Wrapf(err, "error in ToMerklePath: failed to create chainhash from hex")
			}
			mp.Path[i][j] = &trx.PathElement{
				Offset:    leaf.Offset,
				Hash:      hash,
				Txid:      &leaf.TxID,
				Duplicate: &leaf.Duplicate,
			}
		}
	}
	return mp, nil
}

// FromMerklePath converts trx.MerklePath to BUMP
func FromMerklePath(mp *trx.MerklePath) (*BUMP, error) {
	b := &BUMP{
		BlockHeight: uint64(mp.BlockHeight),
		Path:        make([][]BUMPLeaf, len(mp.Path)),
	}
	for i, level := range mp.Path {
		b.Path[i] = make([]BUMPLeaf, len(level))
		for j, leaf := range level {
			b.Path[i][j] = BUMPLeaf{
				Offset:    leaf.Offset,
				Hash:      leaf.Hash.String(),
				TxID:      leaf.Txid != nil && *leaf.Txid,
				Duplicate: leaf.Duplicate != nil && *leaf.Duplicate,
			}
		}
	}
	return b, nil
}

// Scan scan value into Json, implements sql.Scanner interface
func (bump *BUMP) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	byteValue, err := utils.ToByteArray(value)
	if err != nil || bytes.Equal(byteValue, []byte("")) || bytes.Equal(byteValue, []byte("\"\"")) {
		return nil
	}

	// Unmarshal into your BUMP struct
	err = json.Unmarshal(byteValue, bump)
	if err != nil {
		return fmt.Errorf("failed to parse BUMP from JSON: %w", err)
	}

	return nil
}

// IsEmpty returns true if BUMP is empty (all fields are zero values)
func (bump BUMP) IsEmpty() bool {
	return reflect.DeepEqual(bump, BUMP{})
}

// Value return json value, implement driver.Valuer interface
func (bump BUMP) Value() (driver.Value, error) {
	if bump.IsEmpty() {
		return nil, nil
	}

	// Marshal your BUMP struct
	marshal, err := json.Marshal(bump)
	if err != nil {
		return nil, fmt.Errorf("failed to convert BUMP to JSON: %w", err)
	}

	return string(marshal), nil
}
