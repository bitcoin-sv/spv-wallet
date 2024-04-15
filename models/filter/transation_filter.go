package filter

// TransactionFilter is a struct for handling request parameters for destination search requests
type TransactionFilter struct {
	ModelFilter          `json:",inline"`
	Hex                  *string `json:"hex,omitempty"`
	BlockHash            *string `json:"block_hash,omitempty"`
	BlockHeight          *uint64 `json:"block_height,omitempty"`
	Fee                  *uint64 `json:"fee,omitempty"`
	NumberOfInputs       *uint32 `json:"number_of_inputs,omitempty"`
	NumberOfOutputs      *uint32 `json:"number_of_outputs,omitempty"`
	DraftID              *string `json:"draft_id,omitempty"`
	TotalValue           *uint64 `json:"total_value,omitempty"`
	Status               *string `json:"status,omitempty" enums:"UNKNOWN,QUEUED,RECEIVED,STORED,ANNOUNCED_TO_NETWORK,REQUESTED_BY_NETWORK,SENT_TO_NETWORK,ACCEPTED_BY_NETWORK,SEEN_ON_NETWORK,MINED,SEEN_IN_ORPHAN_MEMPOOL,CONFIRMED,REJECTED"`
	TransactionDirection *string `json:"direction,omitempty" enums:"incoming,outgoing,reconcile"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *TransactionFilter) ToDbConditions() (map[string]interface{}, error) {
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

	// NOTE that the "direction" is not a column in the database
	// this field is transformed into final form in the processDBConditions function /engine/tx_repository.go
	direction, err := strOption(d.TransactionDirection, "incoming", "outgoing", "reconcile")
	if err != nil {
		return nil, err
	}
	applyIfNotNil(conditions, "direction", direction)

	return conditions, nil
}
