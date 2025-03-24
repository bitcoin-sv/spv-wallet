package engine

import (
	"context"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/logging"
	"github.com/bitcoin-sv/spv-wallet/engine/metrics"
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
	paymailclient "github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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
	"github.com/go-resty/resty/v2"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
)

type (

	// Client is the SPV Wallet Engine client & options
	Client struct {
		options *clientOptions
	}

	// clientOptions holds all the configuration for the client
	clientOptions struct {
		cacheStore                 *cacheStoreOptions    // Configuration options for Cachestore (ristretto, redis, etc.)
		cluster                    *clusterOptions       // Configuration options for the cluster coordinator
		dataStore                  *dataStoreOptions     // Configuration options for the DataStore (PostgreSQL, etc.)
		debug                      bool                  // If the client is in debug mode
		encryptionKey              string                // Encryption key for encrypting sensitive information (IE: paymail xPub) (hex encoded key)
		httpClient                 *resty.Client         // HTTP client to use for http calls
		iuc                        bool                  // (Input UTXO Check) True will check input utxos when saving transactions
		logger                     *zerolog.Logger       // Internal logging
		metrics                    *metrics.Metrics      // Metrics with a collector interface
		notifications              *notificationsOptions // Configuration options for Notifications
		paymail                    *paymailOptions       // Paymail options & client
		transactionOutlinesService outlines.Service      // Service for transaction outlines
		transactionRecordService   *record.Service       // Service for recording transactions
		taskManager                *taskManagerOptions   // Configuration options for the TaskManager (TaskQ, etc.)
		userAgent                  string                // User agent for all outgoing requests
		chainService               chain.Service         // Chain service
		arcConfig                  chainmodels.ARCConfig // Configuration for ARC
		bhsConfig                  chainmodels.BHSConfig // Configuration for BHS
		feeUnit                    *bsv.FeeUnit          // Fee unit for transactions

		// v2
		repositories *repository.All   // Repositories for all db models
		users        *users.Service    // User domain service
		paymails     *paymails.Service // Paymail domain service
		addresses    *addresses.Service
		operations   *operations.Service
		txSync       *txsync.Service
		data         *data.Service
		config       *config.AppConfig
		contacts     *contacts.Service
	}

	// cacheStoreOptions holds the cache configuration and client
	cacheStoreOptions struct {
		cachestore.ClientInterface                        // Client for Cachestore
		options                    []cachestore.ClientOps // List of options
	}

	// clusterOptions holds the cluster configuration for SPV Wallet Engine clusters
	// at the moment we only support redis as the cluster coordinator
	clusterOptions struct {
		cluster.ClientInterface
		options []cluster.ClientOps // List of options
	}

	// dataStoreOptions holds the data storage configuration and client
	dataStoreOptions struct {
		datastore.ClientInterface                       // Client for Datastore
		options                   []datastore.ClientOps // List of options
	}

	// notificationsOptions holds the configuration for notifications
	notificationsOptions struct {
		enabled        bool
		client         *notifications.Notifications
		webhookManager *notifications.WebhookManager
	}

	// paymailOptions holds the configuration for Paymail
	paymailOptions struct {
		client       paymail.ClientInterface     // Paymail client for communicating with Paymail providers
		service      paymailclient.ServiceClient // Paymail service for handling Paymail requests
		serverConfig *PaymailServerOptions       // Server configuration if Paymail is enabled
	}

	// PaymailServerOptions is the options for the Paymail server
	PaymailServerOptions struct {
		*server.Configuration                    // Server configuration if Paymail is enabled
		options               []server.ConfigOps // Options for the paymail server
		DefaultFromPaymail    string             // IE: from@domain.com
		ExperimentalProvider  bool
	}

	// taskManagerOptions holds the configuration for taskmanager
	taskManagerOptions struct {
		taskmanager.TaskEngine                          // Client for TaskManager
		options                []taskmanager.Options    // List of options
		cronCustomPeriods      map[string]time.Duration // will override the default period of cronJob
	}
)

// NewClient creates a new client for all SPV Wallet Engine functionality
//
// If no options are given, it will use the defaultClientOptions()
func NewClient(ctx context.Context, opts ...ClientOps) (ClientInterface, error) {
	// Create a new client with defaults
	client := &Client{options: defaultClientOptions()}

	// Overwrite defaults with any custom options provided by the user
	for _, opt := range opts {
		opt(client.options)
	}

	// Set the logger (if no custom logger was detected)
	if client.options.logger == nil {
		client.options.logger = logging.GetDefaultLogger()
	}

	// Load the Cachestore client
	var err error
	if err = client.loadCache(ctx); err != nil {
		return nil, err
	}

	// Load the cluster coordinator
	if err = client.loadCluster(ctx); err != nil {
		return nil, err
	}

	// Load the Datastore (automatically migrate models)
	if err = client.loadDatastore(); err != nil {
		return nil, err
	}

	if err = client.autoMigrate(ctx); err != nil {
		return nil, err
	}

	client.loadRepositories()

	client.loadUsersService()
	client.loadPaymailsService()
	client.loadAddressesService()
	client.loadDataService()
	client.loadOperationsService()

	// Load the Paymail client and service (if does not exist)
	if err = client.loadPaymailComponents(); err != nil {
		return nil, err
	}

	// Load the Contacts service after the Paymail service
	client.loadContactsService()

	// Load the Notification client (if client does not exist)
	if err = client.loadNotificationClient(ctx); err != nil {
		return nil, err
	}

	// Load the Taskmanager (automatically start consumers and tasks)
	if err = client.loadTaskmanager(ctx); err != nil {
		return nil, err
	}

	client.loadChainService()
	client.loadTxSyncService()

	if err = client.loadTransactionRecordService(); err != nil {
		return nil, err
	}

	// Register all cron jobs
	if err = client.registerCronJobs(); err != nil {
		return nil, err
	}

	if err = client.loadPaymailServer(); err != nil {
		return nil, err
	}

	if client.options.feeUnit == nil {
		if err = client.askForFeeUnit(ctx); err != nil {
			return nil, err
		}
	}

	if err = client.loadTransactionOutlinesService(); err != nil {
		return nil, err
	}

	// Return the client
	return client, nil
}

// Cachestore will return the Cachestore IF: exists and is enabled
func (c *Client) Cachestore() cachestore.ClientInterface {
	if c.options.cacheStore != nil && c.options.cacheStore.ClientInterface != nil {
		return c.options.cacheStore.ClientInterface
	}
	return nil
}

// Cluster will return the cluster coordinator client
func (c *Client) Cluster() cluster.ClientInterface {
	if c.options.cluster != nil && c.options.cluster.ClientInterface != nil {
		return c.options.cluster.ClientInterface
	}
	return nil
}

// Close will safely close any open connections (cache, datastore, etc.)
func (c *Client) Close(ctx context.Context) error {
	// Close WebhookManager
	if c.options.notifications != nil && c.options.notifications.webhookManager != nil {
		c.options.notifications.webhookManager.Stop()
	}

	// Close Datastore
	ds := c.Datastore()
	if ds != nil {
		if err := ds.Close(); err != nil {
			return spverrors.Wrapf(err, "failed to close datastore")
		}
		c.options.dataStore.ClientInterface = nil
	}

	// Close Taskmanager
	tm := c.Taskmanager()
	if tm != nil {
		if err := tm.Close(ctx); err != nil {
			return spverrors.Wrapf(err, "failed to close taskmanager")
		}
		c.options.taskManager.TaskEngine = nil
	}
	return nil
}

// Datastore will return the Datastore if it exists
func (c *Client) Datastore() datastore.ClientInterface {
	if c.options.dataStore != nil && c.options.dataStore.ClientInterface != nil {
		return c.options.dataStore.ClientInterface
	}
	return nil
}

// Debug will toggle the debug mode (for all resources)
func (c *Client) Debug(on bool) {
	// Set the flag on the current client
	c.options.debug = on

	// Set debugging on the Cachestore
	if cs := c.Cachestore(); cs != nil {
		cs.Debug(on)
	}

	// Set debugging on the Datastore
	if ds := c.Datastore(); ds != nil {
		ds.Debug(on)
	}
}

// IsDebug will return the debug flag (bool)
func (c *Client) IsDebug() bool {
	return c.options.debug
}

// IsIUCEnabled will return the flag (bool)
func (c *Client) IsIUCEnabled() bool {
	return c.options.iuc
}

// IsEncryptionKeySet will return the flag (bool) if the encryption key has been set
func (c *Client) IsEncryptionKeySet() bool {
	return len(c.options.encryptionKey) > 0
}

// Logger will return the Logger if it exists
func (c *Client) Logger() *zerolog.Logger {
	return c.options.logger
}

// Notifications will return the Notifications if it exists
func (c *Client) Notifications() *notifications.Notifications {
	if c.options.notifications == nil {
		return nil
	}
	return c.options.notifications.client
}

// Taskmanager will return the Taskmanager if it exists
func (c *Client) Taskmanager() taskmanager.TaskEngine {
	if c.options.taskManager != nil && c.options.taskManager.TaskEngine != nil {
		return c.options.taskManager.TaskEngine
	}
	return nil
}

// UserAgent will return the user agent
func (c *Client) UserAgent() string {
	return c.options.userAgent
}

// Version will return the version
func (c *Client) Version() string {
	return version
}

// Metrics will return the metrics client (if it's enabled)
func (c *Client) Metrics() (metrics *metrics.Metrics, enabled bool) {
	return c.options.metrics, c.options.metrics != nil
}

// Chain will return the chain service
func (c *Client) Chain() chain.Service {
	return c.options.chainService
}

// LogBHSReadiness tries to ping BHS server. The result is logged.
func (c *Client) LogBHSReadiness(ctx context.Context) {
	logger := c.Logger()
	shortTimeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := c.Chain().HealthcheckBHS(shortTimeoutCtx)
	if err != nil {
		logger.Warn().Err(err).Msg("Unable to connect to Block Headers Service at startup. Application will continue to operate but won't receive BEEF transactions until BHS is online.")
	} else {
		logger.Info().Msg("Block Headers Service is ready to verify transactions.")
	}
}

// FeeUnit will return the fee unit used for transactions
func (c *Client) FeeUnit() bsv.FeeUnit {
	return *c.options.feeUnit
}

// Repositories will return all the repositories
func (c *Client) Repositories() *repository.All {
	return c.options.repositories
}

// UsersService will return the user domain service
func (c *Client) UsersService() *users.Service {
	return c.options.users
}

// PaymailsService will return the paymail domain service
func (c *Client) PaymailsService() *paymails.Service {
	return c.options.paymails
}

// AddressesService will return the address domain service
func (c *Client) AddressesService() *addresses.Service {
	return c.options.addresses
}

// DataService will return the data domain service
func (c *Client) DataService() *data.Service {
	return c.options.data
}

// OperationsService will return the operations domain service
func (c *Client) OperationsService() *operations.Service {
	return c.options.operations
}

// ContactService will return the contacts domain service
func (c *Client) ContactService() *contacts.Service {
	return c.options.contacts
}

// TxSyncService will return the transaction sync service
func (c *Client) TxSyncService() *txsync.Service {
	return c.options.txSync
}
