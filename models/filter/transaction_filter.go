package filter

// TransactionFilter is a struct for handling request parameters for transactions search requests
type TransactionFilter struct {
	// ModelFilter is a struct for handling typical request parameters for search requests
	//lint:ignore SA5008 We want to reuse json tags also to mapstructure.
	ModelFilter     `json:",inline,squash"`
	Hex             *string `json:"hex,omitempty"`
	BlockHash       *string `json:"blockHash,omitempty" example:"0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"`
	BlockHeight     *uint64 `json:"blockHeight,omitempty" example:"839376"`
	Fee             *uint64 `json:"fee,omitempty" example:"1"`
	NumberOfInputs  *uint32 `json:"numberOfInputs,omitempty" example:"1"`
	NumberOfOutputs *uint32 `json:"numberOfOutputs,omitempty" example:"2"`
	DraftID         *string `json:"draftId,omitempty" example:"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"`
	TotalValue      *uint64 `json:"totalValue,omitempty" example:"100000000"`
	Status          *string `json:"status,omitempty" enums:"UNKNOWN,QUEUED,RECEIVED,STORED,ANNOUNCED_TO_NETWORK,REQUESTED_BY_NETWORK,SENT_TO_NETWORK,ACCEPTED_BY_NETWORK,SEEN_ON_NETWORK,MINED,SEEN_IN_ORPHAN_MEMPOOL,CONFIRMED,REJECTED"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *TransactionFilter) ToDbConditions() map[string]interface{} {
	if d == nil {
		return nil
	}
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_transactions.go
	applyIfNotNil(conditions, "hex", d.Hex)
	applyIfNotNil(conditions, "block_hash", d.BlockHash)
	applyIfNotNil(conditions, "block_height", d.BlockHeight)
	applyIfNotNil(conditions, "fee", d.Fee)
	applyIfNotNil(conditions, "number_of_inputs", d.NumberOfInputs)
	applyIfNotNil(conditions, "number_of_outputs", d.NumberOfOutputs)
	applyIfNotNil(conditions, "draft_id", d.DraftID)
	applyIfNotNil(conditions, "total_value", d.TotalValue)
	applyIfNotNil(conditions, "tx_status", d.Status) // be aware that the name of db the dbcolumn is "tx_status" not "status"

	return conditions
}
