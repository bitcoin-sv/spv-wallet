package engine

import (
	"time"
)

// Defaults for engine functionality
const (
	changeOutputSize           = uint64(35)               // Average size in bytes of a change output
	databaseLongReadTimeout    = 30 * time.Second         // For all "GET" or "SELECT" methods
	defaultBroadcastTimeout    = 25 * time.Second         // Default timeout for broadcasting
	defaultCacheLockTTL        = 20                       // in Seconds
	defaultCacheLockTTW        = 10                       // in Seconds
	defaultDatabaseReadTimeout = 20 * time.Second         // For all "GET" or "SELECT" methods
	defaultDraftTxExpiresIn    = 20 * time.Second         // Default TTL for draft transactions
	defaultOverheadSize        = uint64(8)                // 8 bytes is the default overhead in a transaction = 4 bytes version + 4 bytes nLockTime
	defaultUserAgent           = "spv-wallet: " + version // Default user agent
	dustLimit                  = uint64(1)                // Dust limit
	sqliteTestVersion          = "3.37.0"                 // SQLite Testing Version (dummy version for now)
	version                    = "v0.14.2"                // SPV Wallet Engine version
)

// All the base models
const (
	ModelAccessKey        ModelName = "access_key"
	ModelDestination      ModelName = "destination"
	ModelDraftTransaction ModelName = "draft_transaction"
	ModelMetadata         ModelName = "metadata"
	ModelNameEmpty        ModelName = "empty"
	ModelPaymailAddress   ModelName = "paymail_address"
	ModelTransaction      ModelName = "transaction"
	ModelUtxo             ModelName = "utxo"
	ModelXPub             ModelName = "xpub"
	ModelContact          ModelName = "contact"
	ModelWebhook          ModelName = "webhook"
)

// AllModelNames is a list of all models
var AllModelNames = []ModelName{
	ModelAccessKey,
	ModelDestination,
	ModelMetadata,
	ModelPaymailAddress,
	ModelPaymailAddress,
	ModelTransaction,
	ModelUtxo,
	ModelXPub,
	ModelContact,
	ModelWebhook,
}

// Internal table names
const (
	tableAccessKeys        = "access_keys"
	tableDestinations      = "destinations"
	tableDraftTransactions = "draft_transactions"
	tablePaymailAddresses  = "paymail_addresses"
	tableTransactions      = "transactions"
	tableUTXOs             = "utxos"
	tableXPubs             = "xpubs"
	tableContacts          = "contacts"
	tableWebhooks          = "webhooks"
)

const (
	// ReferenceIDField is used for Paymail
	ReferenceIDField = "reference_id"

	// Internal field names
	aliasField           = "alias"
	createdAtField       = "created_at"
	deletedAtField       = "deleted_at"
	currentBalanceField  = "current_balance"
	domainField          = "domain"
	draftIDField         = "draft_id"
	idField              = "id"
	metadataField        = "metadata"
	nextExternalNumField = "next_external_num"
	nextInternalNumField = "next_internal_num"
	satoshisField        = "satoshis"
	spendingTxIDField    = "spending_tx_id"
	statusField          = "status"
	typeField            = "type"
	xPubIDField          = "xpub_id"
	xPubMetadataField    = "xpub_metadata"
	paymailField         = "paymail"
	contactStatusField   = "status"

	// Universal statuses
	statusCanceled = "canceled"
	statusComplete = "complete"
	statusDraft    = "draft"
	statusExpired  = "expired"

	// Paymail / Handles
	defaultSenderPaymail = "example@example.com"
	handleHandcashPrefix = "$"
	handleMaxLength      = 25
	handleRelayPrefix    = "1"
	p2pMetadataField     = "p2p_tx_metadata"

	// Misc
	gormTypeText = "text"
)

// Cache keys for model caching
const (
	cacheKeyDestinationModel                = "destination-id-%s"             // model-id-<destination_id>
	cacheKeyDestinationModelByAddress       = "destination-address-%s"        // model-address-<address>
	cacheKeyDestinationModelByLockingScript = "destination-locking-script-%s" // model-locking-script-<script>
	cacheKeyXpubModel                       = "xpub-id-%s"                    // model-id-<xpub_id>
)

// AllDBModels returns all the database models, e.g. for migrations.
func AllDBModels() []any {
	return []any{
		&Xpub{},
		&AccessKey{},
		&DraftTransaction{},
		&Transaction{},
		&Destination{},
		&Utxo{},
		&Contact{},
		&Webhook{},
		&PaymailAddress{},
	}
}
