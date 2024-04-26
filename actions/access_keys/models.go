package accesskeys

import (
	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// CreateAccessKey is the model for creating an access key
type CreateAccessKey struct {
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}

// SearchAccessKeys is a model for handling searching with filters and metadata
type SearchAccessKeys = common.SearchModel[filter.AccessKeyFilter]

// CountAccessKeys is a model for handling counting filtered transactions
type CountAccessKeys = common.ConditionsModel[filter.AccessKeyFilter]
