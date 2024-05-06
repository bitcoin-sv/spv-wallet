package filter

// UtxoFilter is a struct for handling request parameters for utxo search requests
type UtxoFilter struct {
	ModelFilter `json:",inline"`

	TransactionID *string `json:"transactionId,omitempty"`
	OutputIndex   *uint32 `json:"outputIndex,omitempty"`

	ID            *string    `json:"id,omitempty"`
	Satoshis      *uint64    `json:"satoshis,omitempty"`
	ScriptPubKey  *string    `json:"scriptPubKey,omitempty"`
	Type          *string    `json:"type,omitempty" enums:"pubkey,pubkeyhash,nulldata,multisig,nonstandard,scripthash,metanet,token_stas,token_sensible"`
	DraftID       *string    `json:"draftId,omitempty"`
	ReservedRange *TimeRange `json:"reservedRange,omitempty" swaggertype:"object,string"`
	SpendingTxID  *string    `json:"spendingTxId,omitempty"`
}

var validTypes = getEnumValues[UtxoFilter]("Type")

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *UtxoFilter) ToDbConditions() (map[string]interface{}, error) {
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_utxos.go
	applyIfNotNil(conditions, "transaction_id", d.TransactionID)
	applyIfNotNil(conditions, "output_index", d.OutputIndex)
	applyIfNotNil(conditions, "id", d.ID)
	applyIfNotNil(conditions, "satoshis", d.Satoshis)
	applyIfNotNil(conditions, "script_pub_key", d.ScriptPubKey)
	if err := checkAndApplyStrOption(conditions, "type", d.Type, validTypes...); err != nil {
		return nil, err
	}
	applyIfNotNil(conditions, "spending_tx_id", d.SpendingTxID)

	applyConditionsIfNotNil(conditions, "reserved_at", d.ReservedRange.ToDbConditions())

	return conditions, nil
}

// AdminUtxoFilter wraps the UtxoFilter providing additional fields for admin utxo search requests
type AdminUtxoFilter struct {
	UtxoFilter `json:",inline"`

	XpubID *string `json:"xpubId,omitempty"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *AdminUtxoFilter) ToDbConditions() (map[string]interface{}, error) {
	conditions, err := d.UtxoFilter.ToDbConditions()
	if err != nil {
		return nil, err
	}

	applyIfNotNil(conditions, "xpub_id", d.XpubID)

	return conditions, nil
}
