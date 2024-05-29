package admin

import (
	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// CreatePaymail is the model for creating a paymail
type CreatePaymail struct {
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
	// The xpub with which the paymail is associated
	Key string `json:"key" example:"xpub661MyMwAqRbcGpZVrSHU..."`
	// The paymail address
	Address string `json:"address" example:"test@spv-wallet.com"`
	// The public name of the paymail
	PublicName string `json:"public_name" example:"Test"`
	// The avatar of the paymail (url address)
	Avatar string `json:"avatar" example:"https://example.com/avatar.png"`
}

// PaymailAddress is the model containing only paymail address used for getting and deleting paymail address
type PaymailAddress struct {
	// The paymail address
	Address string `json:"address" example:"test@spv-wallet.com"`
}

// RecordTransaction is the model for recording a transaction
type RecordTransaction struct {
	// The transaction hex
	Hex string `json:"hex" example:"0100000002..."`
}

// CreateXpub is the model for creating an xpub
type CreateXpub struct {
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
	// The xpub key
	Key string `json:"key" example:"xpub661MyMwAqRbcGpZVrSHU..."`
}

// UpdateContact is the model for updating a contact
type UpdateContact struct {
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
	// New name for the contact
	FullName string `json:"fullName" example:"John Doe"`
}

// SearchAccessKeys is a model for handling searching with filters and metadata
type SearchAccessKeys = common.SearchModel[filter.AdminAccessKeyFilter]

// CountAccessKeys is a model for handling counting filtered transactions
type CountAccessKeys = common.ConditionsModel[filter.AdminAccessKeyFilter]

// SearchTransactions is a model for handling searching with filters and metadata
type SearchTransactions = common.SearchModel[filter.TransactionFilter]

// CountTransactions is a model for handling counting filtered transactions
type CountTransactions = common.ConditionsModel[filter.TransactionFilter]

// SearchUtxos is a model for handling searching with filters and metadata
type SearchUtxos = common.SearchModel[filter.AdminUtxoFilter]

// CountUtxos is a model for handling counting filtered UTXOs
type CountUtxos = common.ConditionsModel[filter.AdminUtxoFilter]

// SearchPaymails is a model for handling searching with filters and metadata
type SearchPaymails = common.SearchModel[filter.AdminPaymailFilter]

// CountPaymails is a model for handling counting filtered paymails
type CountPaymails = common.ConditionsModel[filter.AdminPaymailFilter]

// SearchXpubs is a model for handling searching with filters and metadata
type SearchXpubs = common.SearchModel[filter.XpubFilter]

// CountXpubs is a model for handling counting filtered xPubs
type CountXpubs = common.ConditionsModel[filter.XpubFilter]

// CountContacts is a model for handling counting filtered contacts
type CountContacts = common.ConditionsModel[filter.ContactFilter]
