package fixtures

import "github.com/bitcoin-sv/spv-wallet/models/bsv"

// DefaultFeeUnit is the default fee unit used in the tests.
var DefaultFeeUnit = bsv.FeeUnit{
	Satoshis: 1,
	Bytes:    1000,
}

// EstimatedUnlockingScriptSizeForP2PKH is the estimated unlocking script size for a P2PKH transaction.
const EstimatedUnlockingScriptSizeForP2PKH = 106
