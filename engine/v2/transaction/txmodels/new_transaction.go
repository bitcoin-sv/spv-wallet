package txmodels

import "github.com/samber/lo"

// NewTransaction is a data for creating a new transaction.
type NewTransaction struct {
	ID       string
	TxStatus TxStatus

	Inputs  []TrackedOutput
	Outputs []NewOutput

	transactionInputSources []string
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

// SetRawHex sets the raw hexadecimal representation of the transaction and source transaction IDs.
func (t *NewTransaction) SetRawHex(hex string, sourceTXIDs ...string) {
	t.rawHex = hex
	t.transactionInputSources = lo.Map(t.Outputs, func(item NewOutput, index int) string {
		return item.TxID
	})
}

// SetBEEFHex sets the BEEF-encoded hexadecimal representation of the transaction.
func (t *NewTransaction) SetBEEFHex(hex string) { t.beefHex = hex }

// TransactionInputSources returns the list of input sources IDs associated with the transaction.
func (t *NewTransaction) TransactionInputSources() []string {
	return t.transactionInputSources
}

// AddInputs adds outpoints to spend in the transaction.
func (t *NewTransaction) AddInputs(tracked ...TrackedOutput) {
	t.Inputs = append(t.Inputs, tracked...)
}
