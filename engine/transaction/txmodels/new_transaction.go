package txmodels

import (
	"iter"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// NewTransaction is a data for creating a new transaction.
type NewTransaction struct {
	ID       string
	TxStatus TxStatus

	OutpointsToSpend []bsv.Outpoint
	Outputs          []NewOutput
}

// AddOutputs adds list of NewOutput types to the transaction.
func (t *NewTransaction) AddOutputs(output ...NewOutput) {
	t.Outputs = append(t.Outputs, output...)
}

// AddInputs adds outpoints to spend in the transaction.
func (t *NewTransaction) AddInputs(outpoints iter.Seq[bsv.Outpoint]) {
	t.OutpointsToSpend = slices.AppendSeq(t.OutpointsToSpend, outpoints)
}
