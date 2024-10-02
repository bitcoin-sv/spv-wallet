package chainmodels

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// FeeAmount indicating how many satoshis are paid for the given number of bytes in the transaction.
type FeeAmount struct {
	Bytes    int          `json:"bytes"`
	Satoshis bsv.Satoshis `json:"satoshis"`
}

// PolicyContent is the actual policy values returned from ARC
type PolicyContent struct {
	MaxScriptSizePolicy     int       `json:"maxscriptsizepolicy"`
	MaxTXSigOPSCountsPolicy int64     `json:"maxtxsigopscountspolicy"`
	MaxTXSizePolicy         int       `json:"maxtxsizepolicy"`
	MiningFee               FeeAmount `json:"miningFee"`
}

// Policy is the policy model returned from ARC
type Policy struct {
	Content   PolicyContent `json:"policy"`
	Timestamp time.Time     `json:"timestamp"`
}
