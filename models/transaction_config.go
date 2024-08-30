package models

import "time"

// TransactionConfig is a model that represents a transaction config.
type TransactionConfig struct {
	// ChangeDestinations is a slice of change destinations.
	ChangeDestinations []*Destination `json:"change_destinations"`
	// ChangeStrategy is a change strategy.
	ChangeStrategy string `json:"change_destinations_strategy"`
	// ChangeMinimumSatoshis is a minimum satoshis for change.
	ChangeMinimumSatoshis uint64 `json:"change_minimum_satoshis" example:"0"`
	// ChangeNumberOfDestinations is a number of change destinations.
	ChangeNumberOfDestinations int `json:"change_number_of_destinations" example:"1"`
	// ChangeSatoshis is a change satoshis.
	ChangeSatoshis uint64 `json:"change_satoshis" example:"49"`
	// ExpiresAt is a time when transaction expires.
	ExpiresIn time.Duration `json:"expires_in" example:"1000" swaggertype:"string"`
	// Fee is a fee amount.
	Fee uint64 `json:"fee" example:"1"`
	// FeeUnit is a pointer to a fee unit object.
	FeeUnit *FeeUnit `json:"fee_unit"`
	// FromUtxos is a slice of from utxos used to build transaction.
	FromUtxos []*UtxoPointer `json:"from_utxos"`
	// IncludeUtxos is a slice of utxos to include in transaction.
	IncludeUtxos []*UtxoPointer `json:"include_utxos"`
	// Inputs is a slice of transaction inputs.
	Inputs []*TransactionInput `json:"inputs"`
	// Outputs is a slice of transaction outputs.
	Outputs []*TransactionOutput `json:"outputs"`
	// SendAllTo is a pointer to a transaction output object.
	SendAllTo *TransactionOutput `json:"send_all_to"`
	// Sync contains sync configuration.
	Sync *SyncConfig `json:"sync"`
}

// TransactionInput is a model that represents a transaction input.
type TransactionInput struct {
	// Utxo is a pointer to a utxo object.
	Utxo `json:",inline"`
	// Destination is a pointer to a destination object.
	Destination Destination `json:"destination"`
}

// TransactionOutput is a model that represents a transaction output.
type TransactionOutput struct {
	// OpReturn is a pointer to a op return object.
	OpReturn *OpReturn `json:"op_return,omitempty"`
	// PaymailP4 is a pointer to a paymail p4 object.
	PaymailP4 *PaymailP4 `json:"paymail_p4,omitempty"`
	// Satoshis is a satoshis amount.
	Satoshis uint64 `json:"satoshis" example:"50"`
	// Script is a transaction output string representation of script.
	Script string `json:"script" example:"76a91433ba3607a902bc022164bcb6e993f27bd040241c88ac"`
	// ScriptType is a transaction output script type.
	Scripts []*ScriptOutput `json:"scripts,omitempty"`
	// To is a transaction output destination address.
	To string `json:"to" example:"1MB8MfCyA5mGt3UBhxYr1exBfsFWgL1gCm"`
	// UseForChange is a flag that indicates if this output should be used for change.
	UseForChange bool `json:"use_for_change" example:"false"`
}

// MapProtocol is a model that represents a map protocol.
type MapProtocol struct {
	// App is a map protocol app.
	App string `json:"app,omitempty"`
	// Keys is a map protocol keys.
	Keys map[string]interface{} `json:"keys,omitempty"`
	// Type is a map protocol type.
	Type string `json:"type,omitempty"`
}

// OpReturn is a model that represents a op return.
type OpReturn struct {
	// Hex is a full hex of op return.
	Hex string `json:"hex,omitempty"`
	// HexParts is a slice of splitted hex parts.
	HexParts []string `json:"hex_parts,omitempty"`
	// Map is a pointer to a map protocol object.
	Map *MapProtocol `json:"map,omitempty"`
	// StringParts is a slice of string parts.
	StringParts []string `json:"string_parts,omitempty"`
}

// PaymailP4 is a model that represents a paymail p4.
type PaymailP4 struct {
	// Alias is a paymail p4 alias.
	Alias string `json:"alias,omitempty"`
	// Domain is a paymail p4 domain.
	Domain string `json:"domain,omitempty"`
	// FromPaymail is a paymail p4 from paymail.
	FromPaymail string `json:"from_paymail,omitempty"`
	// Note is a paymail p4 note.
	Note string `json:"note,omitempty"`
	// PubKey is a paymail p4 pub key.
	PubKey string `json:"pub_key,omitempty"`
	// ReceiveEndpoint is a paymail p4 receive endpoint.
	ReceiveEndpoint string `json:"receive_endpoint,omitempty"`
	// ReferenceID is a paymail p4 reference id.
	ReferenceID string `json:"reference_id,omitempty"`
	// ResolutionType is a paymail p4 resolution type.
	ResolutionType string `json:"resolution_type,omitempty"`
}

// ScriptOutput is a model that represents a script output.
type ScriptOutput struct {
	// Address is a script output address.
	Address string `json:"address,omitempty"`
	// Satoshis is a script output satoshis.
	Satoshis uint64 `json:"satoshis,omitempty"`
	// Script is a script output script.
	Script string `json:"script"`
	// ScriptType is a script output script type.
	ScriptType string `json:"script_type"`
}
