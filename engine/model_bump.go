package engine

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/bitcoin-sv/go-sdk/chainhash"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

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

// ToMerklePath converts BUMP to trx.MerklePath
func (b *BUMP) ToMerklePath() (*trx.MerklePath, error) {
	blockHeight, err := conv.Uint64ToUint32(b.BlockHeight)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to convert block height to uint32")
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
				return nil, spverrors.Wrapf(err, "failed to create chainhash from hex")
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
			tempLeaf := BUMPLeaf{}
			tempLeaf.Offset = leaf.Offset
			if leaf.Hash != nil {
				tempLeaf.Hash = leaf.Hash.String()
			}
			if leaf.Txid != nil {
				tempLeaf.TxID = *leaf.Txid
			}
			if leaf.Duplicate != nil {
				tempLeaf.Duplicate = *leaf.Duplicate
			}
			b.Path[i][j] = tempLeaf
		}
	}
	return b, nil
}

// IsEmpty checks if the BUMP is empty
func (bump BUMP) IsEmpty() bool {
	return reflect.DeepEqual(bump, BUMP{})
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

	err = json.Unmarshal(byteValue, bump)
	if err != nil {
		return fmt.Errorf("failed to parse BUMP from JSON: %w", err)
	}

	return nil
}

// Value return json value, implement driver.Valuer interface
func (bump BUMP) Value() (driver.Value, error) {
	if bump.IsEmpty() {
		return nil, nil
	}

	marshal, err := json.Marshal(bump)
	if err != nil {
		return nil, fmt.Errorf("failed to convert BUMP to JSON: %w", err)
	}

	return string(marshal), nil
}
