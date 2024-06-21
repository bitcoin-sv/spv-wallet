package engine

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/libsv/go-bc"
	"github.com/libsv/go-bt/v2"
)

const maxBumpHeight = 64

// BUMPs represents a slice of BUMPs - BSV Unified Merkle Paths
type BUMPs []*BUMP

// BUMP represents BUMP (BSV Unified Merkle Path) format
type BUMP struct {
	BlockHeight uint64       `json:"blockHeight,string"`
	Path        [][]BUMPLeaf `json:"path"`
	// private field for storing already used offsets to avoid duplicate nodes
	allNodes []map[uint64]bool
}

// BUMPLeaf represents each BUMP path element
type BUMPLeaf struct {
	Offset    uint64 `json:"offset,string"`
	Hash      string `json:"hash,omitempty"`
	TxID      bool   `json:"txid,omitempty"`
	Duplicate bool   `json:"duplicate,omitempty"`
}

// CalculateMergedBUMP calculates Merged BUMP from a slice of BUMPs
func CalculateMergedBUMP(bumps []BUMP) (*BUMP, error) {
	if len(bumps) == 0 || bumps == nil {
		return nil, nil
	}

	blockHeight := bumps[0].BlockHeight
	bumpHeight := len(bumps[0].Path)
	if bumpHeight > maxBumpHeight {
		return nil,
			fmt.Errorf("BUMP cannot be higher than %d", maxBumpHeight)
	}

	for _, b := range bumps {
		if bumpHeight != len(b.Path) {
			return nil,
				errors.New("Merged BUMP cannot be obtained from Merkle Proofs of different heights")
		}
		if b.BlockHeight != blockHeight {
			return nil,
				errors.New("BUMPs have different block heights. Cannot merge BUMPs from different blocks")
		}
		if len(b.Path) == 0 {
			return nil,
				errors.New("Empty BUMP given")
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
			return nil, errors.New("BUMPs have different merkle roots")
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
		return errors.New("BUMPs with different heights cannot be merged")
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
				return "", errors.New("different merkle roots for the same block")
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
			return "", errors.New("could not find pair")
		}

		leftNode, rightNode := prepareNodes(baseLeaf, offset, *leafInPair, newOffset)

		str, err := bc.MerkleTreeParentStr(leftNode, rightNode)
		if err != nil {
			return "", err
		}
		calculatedHash = str

		offset = offset / 2

		baseLeaf = BUMPLeaf{
			Hash:   calculatedHash,
			Offset: offset,
		}
	}

	return calculatedHash, nil
}

func findLeafByOffset(offset uint64, bumpLeaves []BUMPLeaf) *BUMPLeaf {
	for _, bumpTx := range bumpLeaves {
		if bumpTx.Offset == offset {
			return &bumpTx
		}
	}
	return nil
}

func getOffsetPair(offset uint64) uint64 {
	if offset%2 == 0 {
		return offset + 1
	}
	return offset - 1
}

func getParentOffset(offset uint64) uint64 {
	return getOffsetPair(offset / 2)
}

func prepareNodes(baseLeaf BUMPLeaf, offset uint64, leafInPair BUMPLeaf, newOffset uint64) (string, string) {
	var baseLeafHash, pairLeafHash string

	if baseLeaf.Duplicate {
		baseLeafHash = leafInPair.Hash
	} else {
		baseLeafHash = baseLeaf.Hash
	}

	if leafInPair.Duplicate {
		pairLeafHash = baseLeaf.Hash
	} else {
		pairLeafHash = leafInPair.Hash
	}

	if newOffset > offset {
		return baseLeafHash, pairLeafHash
	}
	return pairLeafHash, baseLeafHash
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

// Hex returns BUMP in hex format
func (bump *BUMP) Hex() string {
	return bump.bytesBuffer().String()
}

func (bump *BUMP) bytesBuffer() *bytes.Buffer {
	var buff bytes.Buffer
	buff.WriteString(hex.EncodeToString(bt.VarInt(bump.BlockHeight).Bytes()))

	height := len(bump.Path)
	buff.WriteString(leadingZeroInt(height))

	for i := 0; i < height; i++ {
		nodes := bump.Path[i]

		nLeafs := len(nodes)
		buff.WriteString(hex.EncodeToString(bt.VarInt(nLeafs).Bytes()))
		for _, n := range nodes {
			buff.WriteString(hex.EncodeToString(bt.VarInt(n.Offset).Bytes()))
			buff.WriteString(fmt.Sprintf("%02x", flags(n.TxID, n.Duplicate)))
			decodedHex, _ := hex.DecodeString(n.Hash)
			buff.WriteString(hex.EncodeToString(bt.ReverseBytes(decodedHex)))
		}
	}
	return &buff
}

// In case the offset or height is less than 10, they must be written with a leading zero
func leadingZeroInt(i int) string {
	return fmt.Sprintf("%02x", i)
}

func flags(txID, duplicate bool) byte {
	var (
		dataFlag      byte = 0o0
		duplicateFlag byte = 0o1
		txIDFlag      byte = 0o2
	)

	if duplicate {
		return duplicateFlag
	}
	if txID {
		return txIDFlag
	}
	return dataFlag
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

	return json.Unmarshal(byteValue, &bump)
}

// Value return json value, implement driver.Valuer interface
func (bump BUMP) Value() (driver.Value, error) {
	if reflect.DeepEqual(bump, BUMP{}) {
		return nil, nil
	}
	marshal, err := json.Marshal(bump)
	if err != nil {
		return nil, err
	}

	return string(marshal), nil
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

	return json.Unmarshal(byteValue, &bumps)
}

// Value return json value, implement driver.Valuer interface
func (bumps BUMPs) Value() (driver.Value, error) {
	if reflect.DeepEqual(bumps, BUMPs{}) {
		return nil, nil
	}
	marshal, err := json.Marshal(bumps)
	if err != nil {
		return nil, err
	}

	return string(marshal), nil
}

func bcBumpToBUMP(bcBump *bc.BUMP) BUMP {
	path := make([][]BUMPLeaf, len(bcBump.Path))
	for i := range bcBump.Path {
		path[i] = make([]BUMPLeaf, len(bcBump.Path[i]))
		for j, source := range bcBump.Path[i] {
			leaf := BUMPLeaf{}

			// All fields in bc.leaf are pointers, so we need to use SafeAssign to avoid dereferencing nil pointers
			utils.SafeAssign(&leaf.Offset, source.Offset)
			utils.SafeAssign(&leaf.Hash, source.Hash)
			utils.SafeAssign(&leaf.TxID, source.Txid)
			utils.SafeAssign(&leaf.Duplicate, source.Duplicate)

			path[i][j] = leaf
		}
	}
	return BUMP{
		BlockHeight: bcBump.BlockHeight,
		Path:        path,
	}
}

func sortAndAddToPath(txIDPath1 BUMPLeaf, offset uint64, txIDPath2 BUMPLeaf, pairOffset uint64) [][]BUMPLeaf {
	path := make([][]BUMPLeaf, 0)
	txIDPath := make([]BUMPLeaf, 2)

	if offset < pairOffset {
		txIDPath[0] = txIDPath1
		txIDPath[1] = txIDPath2
	} else {
		txIDPath[0] = txIDPath2
		txIDPath[1] = txIDPath1
	}

	path = append(path, txIDPath)
	return path
}

func createLeaf(offset uint64, node string) BUMPLeaf {
	leaf := BUMPLeaf{Offset: offset}

	isDuplicate := node == "*"
	if !isDuplicate {
		leaf.Hash = node
	} else {
		leaf.Duplicate = true
	}

	return leaf
}
