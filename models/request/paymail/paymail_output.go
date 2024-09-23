package paymail

import (
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/optional"
)

type Output struct {
	To       string                 `json:"to"`
	Satoshis bsv.Satoshis           `json:"satoshis"`
	From     optional.Param[string] `json:"from,omitempty"`
}

func (o Output) GetType() string {
	return "paymail"
}
