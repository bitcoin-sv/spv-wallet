package engine

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	magic "github.com/bitcoinschema/go-map"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/mrz1836/go-cachestore"
)

// TransactionConfig is the configuration used to start a transaction
type TransactionConfig struct {
	ChangeDestinations         []*Destination       `json:"change_destinations" toml:"change_destinations" yaml:"change_destinations" bson:"change_destinations"`
	ChangeDestinationsStrategy ChangeStrategy       `json:"change_destinations_strategy" toml:"change_destinations_strategy" yaml:"change_destinations_strategy" bson:"change_destinations_strategy"`
	ChangeMinimumSatoshis      uint64               `json:"change_minimum_satoshis" toml:"change_minimum_satoshis" yaml:"change_minimum_satoshis" bson:"change_minimum_satoshis"`
	ChangeNumberOfDestinations int                  `json:"change_number_of_destinations" toml:"change_number_of_destinations" yaml:"change_number_of_destinations" bson:"change_number_of_destinations"`
	ChangeSatoshis             uint64               `json:"change_satoshis" toml:"change_satoshis" yaml:"change_satoshis" bson:"change_satoshis"` // The satoshis used for change
	ExpiresIn                  time.Duration        `json:"expires_in" toml:"expires_in" yaml:"expires_in" bson:"expires_in"`                     // The expiration time for the draft and utxos
	Fee                        uint64               `json:"fee" toml:"fee" yaml:"fee" bson:"fee"`                                                 // The fee used for the transaction (auto generated)
	FeeUnit                    *utils.FeeUnit       `json:"fee_unit" toml:"fee_unit" yaml:"fee_unit" bson:"fee_unit"`                             // Fee unit to use (overrides chainstate if set)
	FromUtxos                  []*UtxoPointer       `json:"from_utxos" toml:"from_utxos" yaml:"from_utxos" bson:"from_utxos"`                     // Use these specific utxos for the transaction
	IncludeUtxos               []*UtxoPointer       `json:"include_utxos" toml:"include_utxos" yaml:"include_utxos" bson:"include_utxos"`         // Include these utxos for the transaction, among others necessary if more is needed for fees
	Inputs                     []*TransactionInput  `json:"inputs" toml:"inputs" yaml:"inputs" bson:"inputs"`                                     // All transaction inputs
	Outputs                    []*TransactionOutput `json:"outputs" toml:"outputs" yaml:"outputs" bson:"outputs"`                                 // All transaction outputs
	SendAllTo                  *TransactionOutput   `json:"send_all_to,omitempty" toml:"send_all_to" yaml:"send_all_to" bson:"send_all_to"`       // Send ALL utxos to the output
	Sync                       *SyncConfig          `json:"sync" toml:"sync" yaml:"sync" bson:"sync"`                                             // Sync config for broadcasting and on-chain sync
	// Future ideas:
	// Conditions (utxo strategy, chain limit, split utxos)
	// NlockTime uint32
}

// TransactionInput is an input on the transaction config
type TransactionInput struct {
	Utxo
	Destination Destination `json:"destination" toml:"destination" yaml:"destination" bson:"destination"`
}

// MapProtocol is a specific MAP protocol interface for an op_return
type MapProtocol struct {
	App  string                 `json:"app,omitempty"`  // Application name
	Keys map[string]interface{} `json:"keys,omitempty"` // Keys to set
	Type string                 `json:"type,omitempty"` // Type of action
}

// OpReturn is the op_return definition for the output
type OpReturn struct {
	Hex         string       `json:"hex,omitempty"`          // Full hex
	HexParts    []string     `json:"hex_parts,omitempty"`    // Hex into parts
	Map         *MapProtocol `json:"map,omitempty"`          // MAP protocol
	StringParts []string     `json:"string_parts,omitempty"` // String parts
}

// TransactionOutput is an output on the transaction config
type TransactionOutput struct {
	OpReturn     *OpReturn       `json:"op_return,omitempty" toml:"op_return" yaml:"op_return" bson:"op_return,omitempty"`                     // Add op_return data as an output
	PaymailP4    *PaymailP4      `json:"paymail_p4,omitempty" toml:"paymail_p4" yaml:"paymail_p4" bson:"paymail_p4,omitempty"`                 // Additional information for P4 or Paymail
	Satoshis     uint64          `json:"satoshis" toml:"satoshis" yaml:"satoshis" bson:"satoshis"`                                             // Set the specific satoshis to send (when applicable)
	Script       string          `json:"script,omitempty" toml:"script" yaml:"script" bson:"script,omitempty"`                                 // custom (non-standard) script output
	Scripts      []*ScriptOutput `json:"scripts" toml:"scripts" yaml:"scripts" bson:"scripts"`                                                 // Add script outputs
	To           string          `json:"to,omitempty" toml:"to" yaml:"to" bson:"to,omitempty"`                                                 // To address, paymail, handle
	UseForChange bool            `json:"use_for_change,omitempty" toml:"use_for_change" yaml:"use_for_change" bson:"use_for_change,omitempty"` // if set, no change destinations will be created, but all outputs flagged will get the change
}

// PaymailPayloadFormat is the format of the paymail payload
type PaymailPayloadFormat uint32

// Types of Paymail payload formats
const (
	BasicPaymailPayloadFormat PaymailPayloadFormat = iota
	BeefPaymailPayloadFormat
)

func (format PaymailPayloadFormat) String() string {
	switch format {
	case BasicPaymailPayloadFormat:
		return "BasicPaymailPayloadFormat"

	case BeefPaymailPayloadFormat:
		return "BeefPaymailPayloadFormat"

	default:
		return fmt.Sprintf("%d", uint32(format))
	}
}

// PaymailP4 paymail configuration for the p2p payments on this output
type PaymailP4 struct {
	Alias           string               `json:"alias" toml:"alias" yaml:"alias" bson:"alias,omitempty"`                                                       // Alias of the paymail {alias}@domain.com
	Domain          string               `json:"domain" toml:"domain" yaml:"domain" bson:"domain,omitempty"`                                                   // Domain of the paymail alias@{domain.com}
	FromPaymail     string               `json:"from_paymail,omitempty" toml:"from_paymail" yaml:"from_paymail" bson:"from_paymail,omitempty"`                 // From paymail address: alias@domain.com
	Note            string               `json:"note,omitempty" toml:"note" yaml:"note" bson:"note,omitempty"`                                                 // Friendly readable note to the paymail receiver
	PubKey          string               `json:"pub_key,omitempty" toml:"pub_key" yaml:"pub_key" bson:"pub_key,omitempty"`                                     // Used for validating the signature
	ReceiveEndpoint string               `json:"receive_endpoint,omitempty" toml:"receive_endpoint" yaml:"receive_endpoint" bson:"receive_endpoint,omitempty"` // P2P endpoint when notifying
	ReferenceID     string               `json:"reference_id,omitempty" toml:"reference_id" yaml:"reference_id" bson:"reference_id,omitempty"`                 // Reference ID saved from P2P request
	ResolutionType  string               `json:"resolution_type" toml:"resolution_type" yaml:"resolution_type" bson:"resolution_type,omitempty"`               // Type of address resolution (basic vs p2p)
	Format          PaymailPayloadFormat `json:"format,omitempty" toml:"format" yaml:"format" bson:"format,omitempty"`                                         // Use beef format for the transaction
}

// Types of resolution methods
const (
	// ResolutionTypeBasic is for the "deprecated" way to resolve an address from a Paymail
	ResolutionTypeBasic = "basic_resolution"

	// ResolutionTypeP2P is the current way to resolve a Paymail (prior to P4)
	ResolutionTypeP2P = "p2p"
)

// ChangeStrategy strategy to use for change
type ChangeStrategy string

// Types of change destination strategies
const (
	// ChangeStrategyDefault is a strategy that divides the satoshis among the change destinations
	ChangeStrategyDefault ChangeStrategy = "default"

	// ChangeStrategyRandom is a strategy randomizing the output of satoshis among the change destinations
	ChangeStrategyRandom ChangeStrategy = "random"

	// ChangeStrategyNominations is a strategy using coin nominations for the outputs (10, 25, 50, 100, 250 etc.)
	ChangeStrategyNominations ChangeStrategy = "nominations"
)

// ScriptOutput is the actual script record (could be several for one output record)
type ScriptOutput struct {
	Address    string `json:"address,omitempty"`  // Hex encoded locking script
	Satoshis   uint64 `json:"satoshis,omitempty"` // Number of satoshis for that output
	Script     string `json:"script"`             // Hex encoded locking script
	ScriptType string `json:"script_type"`        // The type of output
}

// Scan will scan the value into Struct, implements sql.Scanner interface
func (t *TransactionConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	byteValue, err := utils.ToByteArray(value)
	if err != nil || bytes.Equal(byteValue, []byte("")) || bytes.Equal(byteValue, []byte("\"\"")) {
		return nil
	}

	return json.Unmarshal(byteValue, &t)
}

// Value return json value, implement driver.Valuer interface
func (t TransactionConfig) Value() (driver.Value, error) {
	marshal, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return string(marshal), nil
}

// processOutput will inspect the output to determine how to process
func (t *TransactionOutput) processOutput(ctx context.Context, cacheStore cachestore.ClientInterface,
	paymailClient paymail.ClientInterface, defaultFromSender string, checkSatoshis bool,
) error {
	// Convert known handle formats ($handcash or 1relayx)
	if strings.Contains(t.To, handleHandcashPrefix) ||
		(len(t.To) < handleMaxLength && len(t.To) > 1 && t.To[:1] == handleRelayPrefix) {

		// Convert the handle and check if it's changed (becomes a paymail address)
		if p := paymail.ConvertHandle(t.To, false); p != t.To {
			t.To = p
		}
	}

	// Check for Paymail, Bitcoin Address or OP Return
	if len(t.To) > 0 && strings.Contains(t.To, "@") { // Paymail output
		if checkSatoshis && t.Satoshis <= 0 {
			return spverrors.ErrOutputValueTooLow
		}
		return t.processPaymailOutput(ctx, cacheStore, paymailClient, defaultFromSender)
	} else if len(t.To) > 0 { // Standard Bitcoin Address
		if checkSatoshis && t.Satoshis <= 0 {
			return spverrors.ErrOutputValueTooLow
		}
		return t.processAddressOutput()
	} else if t.OpReturn != nil { // OP_RETURN output
		return t.processOpReturnOutput()
	} else if t.Script != "" { // Custom script output
		return t.processScriptOutput()
	}

	// No value set in either ToPaymail or ToAddress
	return spverrors.ErrOutputValueNotRecognized
}

// processPaymailOutput will detect how to process the Paymail output given
func (t *TransactionOutput) processPaymailOutput(ctx context.Context, cacheStore cachestore.ClientInterface,
	paymailClient paymail.ClientInterface, fromPaymail string,
) error {
	// Standardize the paymail address (break into parts)
	alias, domain, paymailAddress := paymail.SanitizePaymail(t.To)
	if len(paymailAddress) == 0 {
		return spverrors.ErrPaymailAddressIsInvalid
	}

	// Set the sanitized version of the paymail address provided
	t.To = paymailAddress

	// Start setting the Paymail information (nil check might not be needed)
	if t.PaymailP4 == nil {
		t.PaymailP4 = &PaymailP4{
			Alias:  alias,
			Domain: domain,
		}
	} else {
		t.PaymailP4.Alias = alias
		t.PaymailP4.Domain = domain
	}

	// Get the capabilities for the domain
	capabilities, err := getCapabilities(
		ctx, cacheStore, paymailClient, domain,
	)
	if err != nil {
		return err
	}

	// Does the provider support P2P?
	success, p2pDestinationURL, p2pSubmitTxURL, format := hasP2P(capabilities)
	if success {
		return t.processPaymailViaP2P(
			paymailClient, p2pDestinationURL, p2pSubmitTxURL, fromPaymail, format,
		)
	}

	return fmt.Errorf("paymail provider does not support P2P")
}

// processPaymailViaP2P will process the output for P2P Paymail resolution
func (t *TransactionOutput) processPaymailViaP2P(client paymail.ClientInterface, p2pDestinationURL, p2pSubmitTxURL string, fromPaymail string, format PaymailPayloadFormat) error {
	// todo: this is a hack since paymail providers will complain if satoshis are empty (SendToAll has 0 satoshi)
	satoshis := t.Satoshis
	if satoshis <= 0 {
		satoshis = 100
	}

	// Get the outputs and destination information from the Paymail provider
	destinationInfo, err := startP2PTransaction(
		client, t.PaymailP4.Alias, t.PaymailP4.Domain,
		p2pDestinationURL, satoshis,
	)
	if err != nil {
		return err
	}

	// split the total output satoshis across all the paymail outputs given
	outputValues, err := utils.SplitOutputValues(satoshis, len(destinationInfo.Outputs))
	if err != nil {
		return err
	}

	// Loop all received P2P outputs and build scripts
	for index, out := range destinationInfo.Outputs {
		t.Scripts = append(
			t.Scripts,
			&ScriptOutput{
				Address:    out.Address,
				Satoshis:   outputValues[index],
				Script:     out.Script,
				ScriptType: utils.ScriptTypePubKeyHash,
			},
		)
	}

	// Set the remaining P2P information
	t.PaymailP4.ReceiveEndpoint = p2pSubmitTxURL
	t.PaymailP4.ReferenceID = destinationInfo.Reference
	t.PaymailP4.ResolutionType = ResolutionTypeP2P
	t.PaymailP4.FromPaymail = fromPaymail
	t.PaymailP4.Format = format

	return nil
}

// processAddressOutput will process an output for a standard Bitcoin Address Transaction
func (t *TransactionOutput) processAddressOutput() (err error) {
	// Create the script from the Bitcoin address
	var s *bscript.Script
	if s, err = bscript.NewP2PKHFromAddress(t.To); err != nil {
		return
	}

	// Append the script
	t.Scripts = append(
		t.Scripts,
		&ScriptOutput{
			Address:    t.To,
			Satoshis:   t.Satoshis,
			Script:     s.String(),
			ScriptType: utils.ScriptTypePubKeyHash,
		},
	)
	return
}

// processScriptOutput will process a custom bitcoin script output
func (t *TransactionOutput) processScriptOutput() (err error) {
	if t.Script == "" {
		return spverrors.ErrInvalidScriptOutput
	}

	// check whether go-bt parses the script correctly
	if _, err = bscript.NewFromHexString(t.Script); err != nil {
		return
	}

	// Append the script
	t.Scripts = append(
		t.Scripts,
		&ScriptOutput{
			Satoshis:   t.Satoshis,
			Script:     t.Script,
			ScriptType: utils.GetDestinationType(t.Script), // try to determine type
		},
	)

	return nil
}

// processOpReturnOutput will process an op_return output
func (t *TransactionOutput) processOpReturnOutput() (err error) {
	// Create the script from the Bitcoin address
	var script string
	if len(t.OpReturn.Hex) > 0 {
		// raw op_return output in hex
		var s *bscript.Script
		if s, err = bscript.NewFromHexString(t.OpReturn.Hex); err != nil {
			return
		}
		script = s.String()
	} else if len(t.OpReturn.HexParts) > 0 {
		// hex strings of the op_return output
		bytesArray := make([][]byte, 0)
		for _, h := range t.OpReturn.HexParts {
			var b []byte
			if b, err = hex.DecodeString(h); err != nil {
				return
			}
			bytesArray = append(bytesArray, b)
		}
		s := &bscript.Script{}
		_ = s.AppendOpcodes(bscript.OpFALSE, bscript.OpRETURN)
		if err = s.AppendPushDataArray(bytesArray); err != nil {
			return
		}
		script = s.String()
	} else if len(t.OpReturn.StringParts) > 0 {
		// strings for the op_return output
		bytesArray := make([][]byte, 0)
		for _, s := range t.OpReturn.StringParts {
			bytesArray = append(bytesArray, []byte(s))
		}
		s := &bscript.Script{}
		_ = s.AppendOpcodes(bscript.OpFALSE, bscript.OpRETURN)
		if err = s.AppendPushDataArray(bytesArray); err != nil {
			return
		}
		script = s.String()
	} else if t.OpReturn.Map != nil {
		// strings for the map op_return
		bytesArray := [][]byte{
			[]byte(magic.Prefix),
			[]byte(magic.Set),
			[]byte(magic.MapAppKey),
			[]byte(t.OpReturn.Map.App),
			[]byte(magic.MapTypeKey),
			[]byte(t.OpReturn.Map.Type),
		}
		if len(t.OpReturn.Map.Keys) > 0 {
			for key, value := range t.OpReturn.Map.Keys {
				bytesArray = append(bytesArray, []byte(key))
				bytesArray = append(bytesArray, []byte(value.(string)))
			}
		}
		s := &bscript.Script{}
		_ = s.AppendOpcodes(bscript.OpFALSE, bscript.OpRETURN)
		if err = s.AppendPushDataArray(bytesArray); err != nil {
			return
		}
		script = s.String()
	} else {
		return spverrors.ErrInvalidOpReturnOutput
	}

	// Append the script
	t.Scripts = append(
		t.Scripts,
		&ScriptOutput{
			Satoshis:   t.Satoshis,
			Script:     script,
			ScriptType: utils.ScriptTypeNullData,
		},
	)
	return
}
