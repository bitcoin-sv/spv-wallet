package txmodels

// NewTransaction is a data for creating a new transaction.
type NewTransaction struct {
	ID       string
	TxStatus TxStatus

	Inputs  []TrackedOutput
	Outputs []NewOutput
}

// AddOutputs adds list of NewOutput types to the transaction.
func (t *NewTransaction) AddOutputs(output ...NewOutput) {
	t.Outputs = append(t.Outputs, output...)
}

// AddInputs adds outpoints to spend in the transaction.
func (t *NewTransaction) AddInputs(tracked ...TrackedOutput) {
	t.Inputs = append(t.Inputs, tracked...)
}
