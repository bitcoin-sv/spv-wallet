package filter

// TransactionFilter is a struct for handling request parameters for destination search requests
type TransactionFilter struct {
	ModelFilter          `json:",inline"`
	Hex                  *string  `json:"hex,omitempty"`
	XpubInIDs            []string `json:"xpub_in_ids,omitempty"`
	XpubOutIDs           []string `json:"xpub_out_ids,omitempty"`
	BlockHash            *string  `json:"block_hash,omitempty"`
	BlockHeight          *uint64  `json:"block_height,omitempty"`
	Fee                  *uint64  `json:"fee,omitempty"`
	NumberOfInputs       *uint32  `json:"number_of_inputs,omitempty"`
	NumberOfOutputs      *uint32  `json:"number_of_outputs,omitempty"`
	DraftID              *string  `json:"draft_id,omitempty"`
	TotalValue           *uint64  `json:"total_value,omitempty"`
	OutputValue          *uint64  `json:"output_value,omitempty"`
	Status               *string  `json:"status,omitempty"`
	TransactionDirection *string  `json:"direction,omitempty" example:"incoming|outgoing"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *TransactionFilter) ToDbConditions() map[string]interface{} {
	conditions := d.ModelFilter.ToDbConditions()

	applyIfNotNil(conditions, "hex", d.Hex)
	applyIfNotEmptySlice(conditions, "xpub_in_id", d.XpubInIDs)
	applyIfNotEmptySlice(conditions, "xpub_out_id", d.XpubOutIDs)
	applyIfNotNil(conditions, "block_hash", d.BlockHash)
	applyIfNotNil(conditions, "block_height", d.BlockHeight)
	applyIfNotNil(conditions, "fee", d.Fee)
	applyIfNotNil(conditions, "number_of_inputs", d.NumberOfInputs)
	applyIfNotNil(conditions, "number_of_outputs", d.NumberOfOutputs)
	applyIfNotNil(conditions, "draft_id", d.DraftID)
	applyIfNotNil(conditions, "total_value", d.TotalValue)
	applyIfNotNil(conditions, "output_value", d.OutputValue)

	// TODO: Should we check if the status string is valid?
	// BTW, where is the list of valid status strings?
	applyIfNotNil(conditions, "status", d.Status)

	// TODO: Should we check if the direction string is valid?
	applyIfNotNil(conditions, "direction", d.TransactionDirection)

	return conditions
}
