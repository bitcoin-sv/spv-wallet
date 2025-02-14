package txmodels

// TransactionInputSource represents a link between a transaction and its source transaction.
// It is used to track which transaction inputs originate from which previous transactions.
type TransactionInputSource struct {
	// SourceTxID is the ID of the transaction that provided the input.
	SourceTxID string

	// TxID is the ID of the transaction that consumes the input.
	TxID string
}

// NewTransaction is a data for creating a new transaction.
type NewTransaction struct {
	ID       string
	TxStatus TxStatus

	Inputs  []TrackedOutput
	Outputs []NewOutput

	transactionInputSources []TransactionInputSource
	beefHex                 string
	rawHex                  string
}

// AddOutputs adds list of NewOutput types to the transaction.
func (t *NewTransaction) AddOutputs(output ...NewOutput) {
	t.Outputs = append(t.Outputs, output...)
}

// BEEFHex returns the BEEF-encoded hexadecimal representation of the transaction.
func (t *NewTransaction) BEEFHex() string { return t.beefHex }

// RawHex returns the raw hexadecimal representation of the transaction.
func (t *NewTransaction) RawHex() string { return t.rawHex }

// SetRawHex sets the raw hexadecimal representation of the transaction.
func (t *NewTransaction) SetRawHex(hex string) {
	t.rawHex = hex
	t.transactionInputSources = make([]TransactionInputSource, 0, len(t.Outputs))

	for _, o := range t.Outputs {
		t.transactionInputSources = append(t.transactionInputSources, TransactionInputSource{TxID: t.ID, SourceTxID: o.TxID})
	}
}

// SetBEEFHex sets the BEEF-encoded hexadecimal representation of the transaction.
func (t *NewTransaction) SetBEEFHex(hex string) { t.beefHex = hex }

// TransactionInputSources returns the list of input sources associated with the transaction.
func (t *NewTransaction) TransactionInputSources() []TransactionInputSource {
	return t.transactionInputSources
}

// AddInputs adds outpoints to spend in the transaction.
func (t *NewTransaction) AddInputs(tracked ...TrackedOutput) {
	t.Inputs = append(t.Inputs, tracked...)
}
