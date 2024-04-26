package utxos

import (
	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// SearchUtxos is a model for handling searching with filters and metadata
type SearchUtxos = common.SearchModel[filter.UtxoFilter]

// CountUtxos is a model for handling counting filtered UTXOs
type CountUtxos = common.ConditionsModel[filter.UtxoFilter]
