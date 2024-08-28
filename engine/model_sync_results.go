package engine

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// SyncResults is the results from all sync attempts (broadcast or sync)
type SyncResults struct {
	Results []*SyncResult `json:"results"` // Each result of a sync task
}

// Sync actions for syncing transactions
const (
	syncActionBroadcast = "broadcast" // Broadcast a transaction into the mempool
	syncActionP2P       = "p2p"       // Notify all paymail providers associated to the transaction
	syncActionSync      = "sync"      // Get on-chain data about the transaction (IE: block hash, height, etc)
)

// SyncResult is the complete attempt/result to sync (multiple providers and strategies)
type SyncResult struct {
	Action        string    `json:"action"`         // type: broadcast, sync etc
	ExecutedAt    time.Time `json:"executed_at"`    // Time it was executed
	StatusMessage string    `json:"status_message"` // Success or failure message
}

// Scan will scan the value into Struct, implements sql.Scanner interface
func (t *SyncResults) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	byteValue, err := utils.ToByteArray(value)
	if err != nil || bytes.Equal(byteValue, []byte("")) || bytes.Equal(byteValue, []byte("\"\"")) {
		return nil
	}

	err = json.Unmarshal(byteValue, &t)
	return spverrors.Wrapf(err, "failed to parse SyncResults from JSON")
}

// Value return json value, implement driver.Valuer interface
func (t SyncResults) Value() (driver.Value, error) {
	marshal, err := json.Marshal(t)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to convert SyncResults to JSON")
	}

	return string(marshal), nil
}
