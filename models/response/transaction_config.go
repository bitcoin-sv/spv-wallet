package response

import "time"

// TransactionConfig is a model that represents a transaction config.
type TransactionConfig struct {
	// ChangeDestinations is a slice of change destinations.
	ChangeDestinations []*Destination `json:"changeDestinations"`
	// ChangeStrategy is a change strategy.
	ChangeStrategy string `json:"changeDestinationsStrategy"`
	// ChangeMinimumSatoshis is a minimum satoshis for change.
	ChangeMinimumSatoshis uint64 `json:"changeMinimumSatoshis" example:"0"`
	// ChangeNumberOfDestinations is a number of change destinations.
	ChangeNumberOfDestinations int `json:"changeNumberOfDestinations" example:"1"`
	// ChangeSatoshis is a change satoshis.
	ChangeSatoshis uint64 `json:"changeSatoshis" example:"49"`
	// ExpiresAt is a time when transaction expires.
	ExpiresIn time.Duration `json:"expiresIn" example:"1000" swaggertype:"string"`
	// Fee is a fee amount.
	Fee uint64 `json:"fee" example:"1"`
	// FeeUnit is a pointer to a fee unit object.
	FeeUnit *FeeUnit `json:"feeUnit"`
	// FromUtxos is a slice of from utxos used to build transaction.
	FromUtxos []*UtxoPointer `json:"fromUtxos"`
	// IncludeUtxos is a slice of utxos to include in transaction.
	IncludeUtxos []*UtxoPointer `json:"includeUtxos"`
	// Inputs is a slice of transaction inputs.
	Inputs []*TransactionInput `json:"inputs"`
	// Outputs is a slice of transaction outputs.
	Outputs []*TransactionOutput `json:"outputs"`
	// SendAllTo is a pointer to a transaction output object.
	SendAllTo *TransactionOutput `json:"sendAllTo"`
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
	OpReturn *OpReturn `json:"opReturn,omitempty"`
	// PaymailP4 is a pointer to a paymail p4 object.
	PaymailP4 *PaymailP4 `json:"paymailP4,omitempty"`
	// Satoshis is a satoshis amount.
	Satoshis uint64 `json:"satoshis" example:"50"`
	// Script is a transaction output string representation of script.
	Script string `json:"script" example:"76a91433ba3607a902bc022164bcb6e993f27bd040241c88ac"`
	// ScriptType is a transaction output script type.
	Scripts []*ScriptOutput `json:"scripts,omitempty"`
	// To is a transaction output destination address.
	To string `json:"to" example:"1MB8MfCyA5mGt3UBhxYr1exBfsFWgL1gCm"`
	// UseForChange is a flag that indicates if this output should be used for change.
	UseForChange bool `json:"useForChange" example:"false"`
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
	HexParts []string `json:"hexParts,omitempty"`
	// Map is a pointer to a map protocol object.
	Map *MapProtocol `json:"map,omitempty"`
	// StringParts is a slice of string parts.
	StringParts []string `json:"stringParts,omitempty"`
}

// PaymailP4 is a model that represents a paymail p4.
type PaymailP4 struct {
	// Alias is a paymail p4 alias.
	Alias string `json:"alias,omitempty"`
	// Domain is a paymail p4 domain.
	Domain string `json:"domain,omitempty"`
	// FromPaymail is a paymail p4 from paymail.
	FromPaymail string `json:"fromPaymail,omitempty"`
	// Note is a paymail p4 note.
	Note string `json:"note,omitempty"`
	// PubKey is a paymail p4 pub key.
	PubKey string `json:"pubKey,omitempty"`
	// ReceiveEndpoint is a paymail p4 receive endpoint.
	ReceiveEndpoint string `json:"receiveEndpoint,omitempty"`
	// ReferenceID is a paymail p4 reference id.
	ReferenceID string `json:"referenceId,omitempty"`
	// ResolutionType is a paymail p4 resolution type.
	ResolutionType string `json:"resolutionType,omitempty"`
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
	ScriptType string `json:"scriptType"`
}
