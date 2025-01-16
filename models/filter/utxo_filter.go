package filter

// UtxoFilter is a struct for handling request parameters for utxo search requests
type UtxoFilter struct {

	// ModelFilter is a struct for handling typical request parameters for search requests
	ModelFilter `json:",inline"`

	TransactionID *string `json:"transactionId,omitempty" example:"5e17858ea0ca4155827754ba82bdcfcce108d5bb5b47fbb3aa54bd14540683c6"`
	OutputIndex   *uint32 `json:"outputIndex,omitempty" example:"0"`

	ID            *string    `json:"id,omitempty" example:"fe4cbfee0258aa589cbc79963f7c204061fd67d987e32ee5049aa90ce14658ee"`
	Satoshis      *uint64    `json:"satoshis,omitempty" example:"1"`
	ScriptPubKey  *string    `json:"scriptPubKey,omitempty" example:"76a914a5f271385e75f57bcd9092592dede812f8c466d088ac"`
	Type          *string    `json:"type,omitempty" enums:"pubkey,pubkeyhash,nulldata,multisig,nonstandard,scripthash,metanet,token_stas,token_sensible"`
	DraftID       *string    `json:"draftId,omitempty" example:"89419d4c7c50810bfe5ff9df9ad5074b749959423782dc91a30f1058b9ad7ef7"`
	ReservedRange *TimeRange `json:"reservedRange,omitempty"` // ReservedRange specifies the time range when a UTXO was reserved.
	SpendingTxID  *string    `json:"spendingTxId,omitempty" example:"11a7746489a70e9c0170601c2be65558455317a984194eb2791b637f59f8cd6e"`
}

var validUtxoTypes = getEnumValues[UtxoFilter]("Type")

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *UtxoFilter) ToDbConditions() (map[string]interface{}, error) {
	if d == nil {
		return nil, nil
	}
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_utxos.go
	applyIfNotNil(conditions, "transaction_id", d.TransactionID)
	applyIfNotNil(conditions, "output_index", d.OutputIndex)
	applyIfNotNil(conditions, "id", d.ID)
	applyIfNotNil(conditions, "satoshis", d.Satoshis)
	applyIfNotNil(conditions, "draft_id", d.DraftID)
	applyIfNotNil(conditions, "script_pub_key", d.ScriptPubKey)
	if err := checkAndApplyStrOption(conditions, "type", d.Type, validUtxoTypes...); err != nil {
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
	if d == nil {
		return nil, nil
	}
	conditions, err := d.UtxoFilter.ToDbConditions()
	if err != nil {
		return nil, err
	}

	applyIfNotNil(conditions, "xpub_id", d.XpubID)

	return conditions, nil
}
