package paymail

import (
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/optional"
)

// Output represents a paymail output.
type Output struct {
	To       string                 `json:"to"`
	Satoshis bsv.Satoshis           `json:"satoshis"`
	From     optional.Param[string] `json:"from,omitempty"`
}

// GetType returns a string typename of the output.
func (o Output) GetType() string {
	return "paymail"
}
