package engine

import (
	"context"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/metrics"
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
	"github.com/rs/zerolog"
)

// AccessKeyService is the access key actions
type AccessKeyService interface {
	GetAccessKey(ctx context.Context, xPubID, pubAccessKey string) (*AccessKey, error)
	GetAccessKeys(ctx context.Context, metadata *Metadata, conditions *map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*AccessKey, error)
	GetAccessKeysCount(ctx context.Context, metadata *Metadata,
		conditions *map[string]interface{}, opts ...ModelOps) (int64, error)
	GetAccessKeysByXPubID(ctx context.Context, xPubID string, metadata *Metadata, conditions *map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*AccessKey, error)
	GetAccessKeysByXPubIDCount(ctx context.Context, xPubID string, metadata *Metadata,
		conditions *map[string]interface{}, opts ...ModelOps) (int64, error)
	NewAccessKey(ctx context.Context, rawXpubKey string, opts ...ModelOps) (*AccessKey, error)
	RevokeAccessKey(ctx context.Context, rawXpubKey, id string, opts ...ModelOps) (*AccessKey, error)
}

// AdminService is the SPV Wallet Engine admin service interface comprised of all services available for admins
type AdminService interface {
	GetStats(ctx context.Context, opts ...ModelOps) (*AdminStats, error)
	GetPaymailAddresses(ctx context.Context, metadataConditions *Metadata, conditions *map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*PaymailAddress, error)
	GetPaymailAddressesCount(ctx context.Context, metadataConditions *Metadata,
		conditions *map[string]interface{}, opts ...ModelOps) (int64, error)
	GetXPubs(ctx context.Context, metadataConditions *Metadata,
		conditions *map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Xpub, error)
	GetXPubsCount(ctx context.Context, metadataConditions *Metadata,
		conditions *map[string]interface{}, opts ...ModelOps) (int64, error)
}

// ClientService is the client related services
type ClientService interface {
	Cachestore() cachestore.ClientInterface
	Cluster() cluster.ClientInterface
	Chainstate() chainstate.ClientInterface
	Datastore() datastore.ClientInterface
	HTTPClient() HTTPInterface
	Logger() *zerolog.Logger
	Notifications() notifications.ClientInterface
	PaymailClient() paymail.ClientInterface
	Taskmanager() taskmanager.TaskEngine
}

// DestinationService is the destination actions
type DestinationService interface {
	GetDestinationByID(ctx context.Context, xPubID, id string) (*Destination, error)
	GetDestinationByAddress(ctx context.Context, xPubID, address string) (*Destination, error)
	GetDestinationByLockingScript(ctx context.Context, xPubID, lockingScript string) (*Destination, error)
	GetDestinations(ctx context.Context, metadata *Metadata, conditions *map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Destination, error)
	GetDestinationsCount(ctx context.Context, metadata *Metadata,
		conditions *map[string]interface{}, opts ...ModelOps) (int64, error)
	GetDestinationsByXpubID(ctx context.Context, xPubID string, usingMetadata *Metadata, conditions *map[string]interface{},
		queryParams *datastore.QueryParams) ([]*Destination, error)
	GetDestinationsByXpubIDCount(ctx context.Context, xPubID string, usingMetadata *Metadata,
		conditions *map[string]interface{}) (int64, error)
	NewDestination(ctx context.Context, xPubKey string, chain uint32, destinationType string,
		opts ...ModelOps) (*Destination, error)
	NewDestinationForLockingScript(ctx context.Context, xPubID, lockingScript string,
		opts ...ModelOps) (*Destination, error)
	UpdateDestinationMetadataByID(ctx context.Context, xPubID, id string, metadata Metadata) (*Destination, error)
	UpdateDestinationMetadataByLockingScript(ctx context.Context, xPubID,
		lockingScript string, metadata Metadata) (*Destination, error)
	UpdateDestinationMetadataByAddress(ctx context.Context, xPubID, address string,
		metadata Metadata) (*Destination, error)
}

// DraftTransactionService is the draft transactions actions
type DraftTransactionService interface {
	GetDraftTransactions(ctx context.Context, metadata *Metadata, conditions *map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*DraftTransaction, error)
	GetDraftTransactionsCount(ctx context.Context, metadata *Metadata,
		conditions *map[string]interface{}, opts ...ModelOps) (int64, error)
}

// HTTPInterface is the HTTP client interface
type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// ModelService is the "model" related services
type ModelService interface {
	AddModels(ctx context.Context, autoMigrate bool, models ...interface{}) error
	DefaultModelOptions(opts ...ModelOps) []ModelOps
	GetModelNames() []string
}

// PaymailService is the paymail actions & services
type PaymailService interface {
	DeletePaymailAddress(ctx context.Context, address string, opts ...ModelOps) error
	GetPaymailConfig() *PaymailServerOptions
	GetPaymailAddress(ctx context.Context, address string, opts ...ModelOps) (*PaymailAddress, error)
	GetPaymailAddressesByXPubID(ctx context.Context, xPubID string, metadataConditions *Metadata,
		conditions *map[string]interface{}, queryParams *datastore.QueryParams) ([]*PaymailAddress, error)
	NewPaymailAddress(ctx context.Context, key, address, publicName,
		avatar string, opts ...ModelOps) (*PaymailAddress, error)
	UpdatePaymailAddress(ctx context.Context, address, publicName,
		avatar string, opts ...ModelOps) (*PaymailAddress, error)
	UpdatePaymailAddressMetadata(ctx context.Context, address string,
		metadata Metadata, opts ...ModelOps) (*PaymailAddress, error)
}

// TransactionService is the transaction actions
type TransactionService interface {
	GetTransaction(ctx context.Context, xPubID, txID string) (*Transaction, error)
	GetTransactionsByIDs(ctx context.Context, txIDs []string) ([]*Transaction, error)
	GetTransactionByHex(ctx context.Context, hex string) (*Transaction, error)
	GetTransactions(ctx context.Context, metadata *Metadata, conditions *map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Transaction, error)
	GetTransactionsCount(ctx context.Context, metadata *Metadata,
		conditions *map[string]interface{}, opts ...ModelOps) (int64, error)
	GetTransactionsByXpubID(ctx context.Context, xPubID string, metadata *Metadata, conditions *map[string]interface{},
		queryParams *datastore.QueryParams) ([]*Transaction, error)
	GetTransactionsByXpubIDCount(ctx context.Context, xPubID string, metadata *Metadata,
		conditions *map[string]interface{}) (int64, error)
	NewTransaction(ctx context.Context, rawXpubKey string, config *TransactionConfig,
		opts ...ModelOps) (*DraftTransaction, error)
	RecordTransaction(ctx context.Context, xPubKey, txHex, draftID string,
		opts ...ModelOps) (*Transaction, error)
	RecordRawTransaction(ctx context.Context, txHex string, opts ...ModelOps) (*Transaction, error)
	UpdateTransaction(ctx context.Context, txInfo *broadcast.SubmittedTx) error
	UpdateTransactionMetadata(ctx context.Context, xPubID, id string, metadata Metadata) (*Transaction, error)
	RevertTransaction(ctx context.Context, id string) error
}

// UTXOService is the utxo actions
type UTXOService interface {
	GetUtxo(ctx context.Context, xPubKey, txID string, outputIndex uint32) (*Utxo, error)
	GetUtxoByTransactionID(ctx context.Context, txID string, outputIndex uint32) (*Utxo, error)
	GetUtxos(ctx context.Context, metadata *Metadata, conditions *map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Utxo, error)
	GetUtxosCount(ctx context.Context, metadata *Metadata,
		conditions *map[string]interface{}, opts ...ModelOps) (int64, error)
	GetUtxosByXpubID(ctx context.Context, xPubID string, metadata *Metadata, conditions *map[string]interface{},
		queryParams *datastore.QueryParams) ([]*Utxo, error)
	UnReserveUtxos(ctx context.Context, xPubID, draftID string) error
}

// XPubService is the xPub actions
type XPubService interface {
	GetXpub(ctx context.Context, xPubKey string) (*Xpub, error)
	GetXpubByID(ctx context.Context, xPubID string) (*Xpub, error)
	NewXpub(ctx context.Context, xPubKey string, opts ...ModelOps) (*Xpub, error)
	UpdateXpubMetadata(ctx context.Context, xPubID string, metadata Metadata) (*Xpub, error)
}

// ClientInterface is the client (spv wallet engine) interface comprised of all services/actions
type ClientInterface interface {
	AccessKeyService
	AdminService
	ClientService
	DestinationService
	DraftTransactionService
	ModelService
	PaymailService
	TransactionService
	UTXOService
	XPubService
	AuthenticateAccessKey(ctx context.Context, pubAccessKey string) (*AccessKey, error)
	Close(ctx context.Context) error
	Debug(on bool)
	DefaultSyncConfig() *SyncConfig
	EnableNewRelic()
	GetOrStartTxn(ctx context.Context, name string) context.Context
	IsDebug() bool
	IsEncryptionKeySet() bool
	IsIUCEnabled() bool
	IsMigrationEnabled() bool
	IsNewRelicEnabled() bool
	SetNotificationsClient(notifications.ClientInterface)
	UserAgent() string
	Version() string
	Metrics() (metrics *metrics.Metrics, enabled bool)
}
