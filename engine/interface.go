package engine

import (
	"context"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/metrics"
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
	paymailclient "github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/addresses"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/data"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database/repository"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/operations"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/record"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txsync"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
)

// AccessKeyService is the access key actions
type AccessKeyService interface {
	GetAccessKey(ctx context.Context, xPubID, pubAccessKey string) (*AccessKey, error)
	GetAccessKeys(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*AccessKey, error)
	GetAccessKeysCount(ctx context.Context, metadata *Metadata,
		conditions map[string]interface{}, opts ...ModelOps) (int64, error)
	GetAccessKeysByXPubID(ctx context.Context, xPubID string, metadata *Metadata, conditions map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*AccessKey, error)
	GetAccessKeysByXPubIDCount(ctx context.Context, xPubID string, metadata *Metadata,
		conditions map[string]interface{}, opts ...ModelOps) (int64, error)
	NewAccessKey(ctx context.Context, rawXpubKey string, opts ...ModelOps) (*AccessKey, error)
	RevokeAccessKey(ctx context.Context, rawXpubKey, id string, opts ...ModelOps) (*AccessKey, error)
}

// AdminService is the SPV Wallet Engine admin service interface comprised of all services available for admins
type AdminService interface {
	GetStats(ctx context.Context, opts ...ModelOps) (*AdminStats, error)
	GetPaymailAddresses(ctx context.Context, metadataConditions *Metadata, conditions map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*PaymailAddress, error)
	GetPaymailAddressesCount(ctx context.Context, metadataConditions *Metadata,
		conditions map[string]interface{}, opts ...ModelOps) (int64, error)
	GetXPubs(ctx context.Context, metadataConditions *Metadata,
		conditions map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Xpub, error)
	GetXPubsCount(ctx context.Context, metadataConditions *Metadata,
		conditions map[string]interface{}, opts ...ModelOps) (int64, error)
}

// ClientService is the client related services
type ClientService interface {
	Cachestore() cachestore.ClientInterface
	Cluster() cluster.ClientInterface
	Datastore() datastore.ClientInterface
	Logger() *zerolog.Logger
	Notifications() *notifications.Notifications
	PaymailClient() paymail.ClientInterface
	PaymailService() paymailclient.ServiceClient
	TransactionOutlinesService() outlines.Service
	TransactionRecordService() *record.Service
	Taskmanager() taskmanager.TaskEngine
}

// ContactService is the service for managing contacts
type ContactService interface {
	UpsertContact(ctx context.Context, fullName, paymailAdress, requesterXPubID, requesterPaymail string, opts ...ModelOps) (*Contact, error)
	AddContactRequest(ctx context.Context, fullName, paymailAdress, requesterXPubID string, opts ...ModelOps) (*Contact, error)

	AdminChangeContactStatus(ctx context.Context, id string, status ContactStatus) (*Contact, error)
	AdminCreateContact(ctx context.Context, contactPaymail, creatorPaymail, fullName string, metadata *Metadata) (*Contact, error)
	AdminConfirmContacts(ctx context.Context, paymailA string, paymailB string) error
	UpdateContact(ctx context.Context, id, fullName string, metadata *Metadata) (*Contact, error)
	DeleteContactByID(ctx context.Context, id string) error
	AdminUnconfirmContact(ctx context.Context, id string) error

	DeleteContact(ctx context.Context, xPubID, paymail string) error
	AcceptContact(ctx context.Context, xPubID, paymail string) error
	RejectContact(ctx context.Context, xPubID, paymail string) error
	ConfirmContact(ctx context.Context, xPubID, paymail string) error
	UnconfirmContact(ctx context.Context, xPubID, paymail string) error

	GetContacts(ctx context.Context, metadata *Metadata, conditions map[string]interface{}, queryParams *datastore.QueryParams) ([]*Contact, error)
	GetContactsByXpubID(ctx context.Context, xPubID string, metadata *Metadata, conditions map[string]interface{}, queryParams *datastore.QueryParams) ([]*Contact, error)
	GetContactsByXPubIDCount(ctx context.Context, xPubID string, metadata *Metadata, conditions map[string]interface{}, opts ...ModelOps) (int64, error)
	GetContactsCount(ctx context.Context, metadata *Metadata, conditions map[string]interface{}, opts ...ModelOps) (int64, error)
}

// HTTPInterface is the HTTP client interface
type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// ModelService is the "model" related services
type ModelService interface {
	DefaultModelOptions(opts ...ModelOps) []ModelOps
}

// PaymailService is the paymail actions & services
type PaymailService interface {
	DeletePaymailAddress(ctx context.Context, address string, opts ...ModelOps) error
	DeletePaymailAddressByID(ctx context.Context, id string, opts ...ModelOps) error
	GetPaymailConfig() *PaymailServerOptions
	GetPaymailAddress(ctx context.Context, address string, opts ...ModelOps) (*PaymailAddress, error)
	GetPaymailAddressByID(ctx context.Context, id string, opts ...ModelOps) (*PaymailAddress, error)
	GetPaymailAddressesByXPubID(ctx context.Context, xPubID string, metadataConditions *Metadata,
		conditions map[string]interface{}, queryParams *datastore.QueryParams) ([]*PaymailAddress, error)
	NewPaymailAddress(ctx context.Context, key, address, publicName,
		avatar string, opts ...ModelOps) (*PaymailAddress, error)
}

// TransactionService is the transaction actions
type TransactionService interface {
	GetTransaction(ctx context.Context, xPubID, txID string) (*Transaction, error)
	GetAdminTransaction(ctx context.Context, txID string) (*Transaction, error)
	GetTransactionsByIDs(ctx context.Context, txIDs []string) ([]*Transaction, error)
	GetTransactions(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Transaction, error)
	GetTransactionsCount(ctx context.Context, metadata *Metadata,
		conditions map[string]interface{}, opts ...ModelOps) (int64, error)
	GetTransactionsByXpubID(ctx context.Context, xPubID string, metadata *Metadata, conditions map[string]interface{},
		queryParams *datastore.QueryParams) ([]*Transaction, error)
	GetTransactionsByXpubIDCount(ctx context.Context, xPubID string, metadata *Metadata,
		conditions map[string]interface{}) (int64, error)
	NewTransaction(ctx context.Context, rawXpubKey string, config *TransactionConfig,
		opts ...ModelOps) (*DraftTransaction, error)
	RecordTransaction(ctx context.Context, xPubKey, txHex, draftID string,
		opts ...ModelOps) (*Transaction, error)
	HandleTxCallback(ctx context.Context, callbackResp *chainmodels.TXInfo) error
	UpdateTransactionMetadata(ctx context.Context, xPubID, id string, metadata Metadata) (*Transaction, error)
	RevertTransaction(ctx context.Context, id string) error
}

// UTXOService is the utxo actions
type UTXOService interface {
	GetUtxo(ctx context.Context, xPubKey, txID string, outputIndex uint32) (*Utxo, error)
	GetUtxoByTransactionID(ctx context.Context, txID string, outputIndex uint32) (*Utxo, error)
	GetUtxos(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
		queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Utxo, error)
	GetUtxosCount(ctx context.Context, metadata *Metadata,
		conditions map[string]interface{}, opts ...ModelOps) (int64, error)
	GetUtxosByXpubID(ctx context.Context, xPubID string, metadata *Metadata, conditions map[string]interface{},
		queryParams *datastore.QueryParams) ([]*Utxo, error)
	GetUtxosByXpubIDCount(ctx context.Context, xPubID string, metadata *Metadata,
		conditions map[string]interface{}) (int64, error)
}

// XPubService is the xPub actions
type XPubService interface {
	GetXpub(ctx context.Context, xPubKey string) (*Xpub, error)
	GetXpubByID(ctx context.Context, xPubID string) (*Xpub, error)
	NewXpub(ctx context.Context, xPubKey string, opts ...ModelOps) (*Xpub, error)
	UpdateXpubMetadata(ctx context.Context, xPubID string, metadata Metadata) (*Xpub, error)
}

// V2 contains services for version 2
type V2 interface {
	Repositories() *repository.All
	UsersService() *users.Service
	PaymailsService() *paymails.Service
	AddressesService() *addresses.Service
	DataService() *data.Service
	OperationsService() *operations.Service
	ContactService() *contacts.Service
	TxSyncService() *txsync.Service
}

// ClientInterface is the client (spv wallet engine) interface comprised of all services/actions
type ClientInterface interface {
	AccessKeyService
	AdminService
	ClientService
	ModelService
	PaymailService
	TransactionService
	UTXOService
	XPubService
	ContactService
	AuthenticateAccessKey(ctx context.Context, pubAccessKey string) (*AccessKey, error)
	Close(ctx context.Context) error
	Debug(on bool)
	IsDebug() bool
	IsEncryptionKeySet() bool
	IsIUCEnabled() bool
	UserAgent() string
	Version() string
	Metrics() (metrics *metrics.Metrics, enabled bool)
	SubscribeWebhook(ctx context.Context, url, tokenHeader, token string) error
	UnsubscribeWebhook(ctx context.Context, url string) error
	GetWebhooks(ctx context.Context) ([]notifications.ModelWebhook, error)
	Chain() chain.Service
	LogBHSReadiness(ctx context.Context)
	FeeUnit() bsv.FeeUnit
	V2
}
