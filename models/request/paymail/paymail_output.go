package paymail

import (
	"github.com/bitcoin-sv/spv-wallet/models/request/optional"
)

type Output struct {
	To       string                 `json:"to"`
	Satoshis uint                   `json:"satoshis"`
	From     optional.Param[string] `json:"from,omitempty"`
}

func (o Output) GetType() string {
	return "paymail"
}
