package engine

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

// SyncConfig is the configuration used for syncing a transaction (on-chain)
type SyncConfig struct {
	Broadcast        bool `json:"broadcast" toml:"broadcast" yaml:"broadcast"`                         // Transaction should be broadcasted
	BroadcastInstant bool `json:"broadcast_instant" toml:"broadcast_instant" yaml:"broadcast_instant"` // Transaction should be broadcasted instantly (ASAP)
	PaymailP2P       bool `json:"paymail_p2p" toml:"paymail_p2p" yaml:"paymail_p2p"`                   // Transaction will be sent to all related paymail providers if P2P is detected
	SyncOnChain      bool `json:"sync_on_chain" toml:"sync_on_chain" yaml:"sync_on_chain"`             // Transaction should be checked that it's on-chain
	// FUTURE IDEAS:
	// DelayToBroadcast time.Duration `json:"delay_to_broadcast" toml:"delay_to_broadcast" yaml:"delay_to_broadcast"` // Delay for broadcasting
	// Miner       string `json:"miner" toml:"miner" yaml:"miner"`  // Use a specific miner
	// miners: []miner{name, token, feeQuote}
	// default: miner
	// failover: miner
	// keep tx updated until x blocks?
}

// Scan will scan the value into Struct, implements sql.Scanner interface
func (t *SyncConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	byteValue, err := utils.ToByteArray(value)
	if err != nil || bytes.Equal(byteValue, []byte("")) || bytes.Equal(byteValue, []byte("\"\"")) {
		return nil
	}

	err = json.Unmarshal(byteValue, &t)
	if err != nil {
		return spverrors.Wrapf(err, "failed to parse SyncConfig from JSON")
	}
	return nil
}

// Value return json value, implement driver.Valuer interface
func (t *SyncConfig) Value() (driver.Value, error) {
	marshal, err := json.Marshal(t)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to convert SyncConfig to JSON")
	}

	return string(marshal), nil
}
