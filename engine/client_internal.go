package engine

import (
	"context"

	paymailclient "github.com/bitcoin-sv/go-paymail"
	paymailserver "github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/addresses"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/data"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database/repository"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/operations"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails"
	paymailprovider "github.com/bitcoin-sv/spv-wallet/engine/v2/paymailserver"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/beef"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/utxo"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/record"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txsync"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users"
	"github.com/mrz1836/go-cachestore"
)

// loadCache will load caching configuration and start the Cachestore client
func (c *Client) loadCache(ctx context.Context) (err error) {
	// Load if a custom interface was NOT provided
	if c.options.cacheStore.ClientInterface == nil {
		c.options.cacheStore.ClientInterface, err = cachestore.NewClient(ctx, c.options.cacheStore.options...)
	}
	return
}

// loadCluster will load the cluster coordinator
func (c *Client) loadCluster(ctx context.Context) (err error) {
	// Load if a custom interface was NOT provided
	if c.options.cluster.ClientInterface == nil {
		c.options.cluster.ClientInterface, err = cluster.NewClient(ctx, c.options.cluster.options...)
	}

	return
}

// loadDatastore will load the Datastore and start the Datastore client
//
// NOTE: this WON't run database migrations
func (c *Client) loadDatastore() (err error) {
	if c.options.dataStore.ClientInterface != nil {
		return
	}

	c.options.dataStore.options = append(
		c.options.dataStore.options,
		datastore.WithCustomFields(
			[]string{ // Array fields
				"xpub_in_ids",
				"xpub_out_ids",
			}, []string{ // Object fields
				"xpub_metadata",
				"xpub_output_value",
			},
		))

	c.options.dataStore.ClientInterface, err = datastore.NewClient(c.options.dataStore.options...)
	return
}

func (c *Client) autoMigrate(ctx context.Context) error {
	if c.Datastore() == nil {
		return spverrors.Newf("datastore is not loaded")
	}

	db := c.Datastore().DB().WithContext(ctx)
	models := AllDBModels(c.options.paymail.serverConfig.ExperimentalProvider)

	if err := db.AutoMigrate(models...); err != nil {
		return spverrors.Wrapf(err, "failed to auto-migrate models")
	}

	// Legacy code compatibility:
	// Some models implement post-migration logic to e.g. manually add some indexes
	// NOTE: In the future, we should remove this and stick to GORM features
	for _, model := range models {
		if migrator, ok := model.(interface {
			PostMigrate(client datastore.ClientInterface) error
		}); ok {
			if err := migrator.PostMigrate(c.Datastore()); err != nil {
				return spverrors.Wrapf(err, "failed to post-migrate model %+v", model)
			}
		}
	}
	return nil
}

// loadNotificationClient will load the notifications client
func (c *Client) loadNotificationClient(ctx context.Context) (err error) {
	if c.options.notifications == nil || !c.options.notifications.enabled {
		return
	}
	logger := c.Logger().With().Str("subservice", "notification").Logger()
	notificationService := notifications.NewNotifications(ctx, &logger)
	c.options.notifications.client = notificationService
	c.options.notifications.webhookManager = notifications.NewWebhookManager(ctx, &logger, notificationService, &WebhooksRepository{client: c})
	return
}

// SubscribeWebhook adds URL to the list of subscribed webhooks
func (c *Client) SubscribeWebhook(ctx context.Context, url, tokenHeader, token string) error {
	if c.options.notifications == nil || c.options.notifications.webhookManager == nil {
		return spverrors.ErrNotificationsDisabled
	}

	err := c.options.notifications.webhookManager.Subscribe(ctx, url, tokenHeader, token)
	if err != nil {
		return spverrors.ErrWebhookSubscriptionFailed
	}
	return nil
}

// UnsubscribeWebhook removes URL from the list of subscribed webhooks
func (c *Client) UnsubscribeWebhook(ctx context.Context, url string) error {
	if c.options.notifications == nil || c.options.notifications.webhookManager == nil {
		return spverrors.ErrNotificationsDisabled
	}

	//nolint:wrapcheck //we're returning our custom errors
	return c.options.notifications.webhookManager.Unsubscribe(ctx, url)
}

// GetWebhooks returns all the webhooks stored in database
func (c *Client) GetWebhooks(ctx context.Context) ([]notifications.ModelWebhook, error) {
	if c.options.notifications == nil || c.options.notifications.webhookManager == nil {
		return nil, spverrors.ErrNotificationsDisabled
	}

	//nolint:wrapcheck //we're returning our custom errors
	return c.options.notifications.webhookManager.GetAll(ctx)
}

// loadPaymailComponents will load the Paymail client
func (c *Client) loadPaymailComponents() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = spverrors.Wrapf(e, "error when creating paymail components")
			} else {
				err = spverrors.Newf("error when creating paymail components: %v", r)
			}
		}
	}()

	// Only load if it's not set (the client can be overloaded)
	if c.options.paymail.client == nil {
		c.options.paymail.client, err = paymailclient.NewClient()
		if err != nil {
			return
		}
		c.options.paymail.client.WithCustomHTTPClient(c.options.httpClient)
	}

	if c.options.paymail.service == nil {
		logger := c.Logger().With().Str("subservice", "paymail").Logger()
		c.options.paymail.service = paymail.NewServiceClient(c.Cachestore(), c.options.paymail.client, logger)
	}
	return
}

func (c *Client) loadTransactionOutlinesService() error {
	if c.options.transactionOutlinesService == nil {
		logger := c.Logger().With().Str("subservice", "transactionOutlines").Logger()
		utxoSelector := utxo.NewSelector(c.Datastore().DB(), c.FeeUnit())
		beefService := beef.NewService(c.Repositories().Transactions)

		c.options.transactionOutlinesService = outlines.NewService(c.PaymailService(), c.options.paymails, beefService, utxoSelector, c.FeeUnit(), logger, c.UsersService())
	}
	return nil
}

func (c *Client) loadTransactionRecordService() error {
	if c.options.transactionRecordService == nil {
		logger := c.Logger().With().Str("subservice", "transactionRecord").Logger()
		c.options.transactionRecordService = record.NewService(
			logger,
			c.AddressesService(),
			c.UsersService(),
			c.Repositories().Outputs,
			c.Repositories().Operations,
			c.Repositories().Transactions,
			c.Chain(),
			c.PaymailService(),
		)
	}
	return nil
}

func (c *Client) loadRepositories() {
	if c.options.repositories == nil {
		c.options.repositories = repository.NewRepositories(c.Datastore().DB())
	}
}

func (c *Client) loadUsersService() {
	if c.options.users == nil {
		c.options.users = users.NewService(c.Repositories().Users, c.options.config)
	}
}

func (c *Client) loadPaymailsService() {
	if c.options.paymails == nil {
		c.options.paymails = paymails.NewService(c.Repositories().Paymails, c.UsersService(), c.options.config)
	}
}

func (c *Client) loadAddressesService() {
	if c.options.addresses == nil {
		c.options.addresses = addresses.NewService(c.Repositories().Addresses)
	}
}

func (c *Client) loadDataService() {
	if c.options.data == nil {
		c.options.data = data.NewService(c.Repositories().Data)
	}
}

func (c *Client) loadOperationsService() {
	if c.options.operations == nil {
		c.options.operations = operations.NewService(c.Repositories().Operations)
	}
}

func (c *Client) loadContactsService() {
	if c.options.contacts == nil {
		logger := c.Logger().With().Str("subservice", "contacts").Logger()
		c.options.contacts = contacts.NewService(c.Repositories().Contacts, c.PaymailsService(), c.PaymailService(), logger)
	}
}

func (c *Client) loadChainService() {
	if c.options.chainService == nil {
		logger := c.Logger().With().Str("subservice", "chain").Logger()
		c.options.arcConfig.TxsGetter = newSDKTxGetter(c)
		c.options.chainService = chain.NewChainService(logger, c.options.httpClient, c.options.arcConfig, c.options.bhsConfig)
	}
}

func (c *Client) loadTxSyncService() {
	if c.options.txSync == nil {
		logger := c.Logger().With().Str("subservice", "tx_sync").Logger()
		c.options.txSync = txsync.NewService(logger, c.Repositories().Transactions)
	}
}

// loadTaskmanager will load the TaskManager and start the TaskManager client
func (c *Client) loadTaskmanager(ctx context.Context) (err error) {
	// Load if a custom interface was NOT provided
	if c.options.taskManager.TaskEngine == nil {
		c.options.taskManager.TaskEngine, err = taskmanager.NewTaskManager(
			ctx, c.options.taskManager.options...,
		)
	}
	return
}

func (c *Client) registerCronJobs() error {
	cronJobs := c.cronJobs()

	if c.options.taskManager.cronCustomPeriods != nil {
		// override the default periods
		for name, job := range cronJobs {
			if custom, ok := c.options.taskManager.cronCustomPeriods[name]; ok {
				job.Period = custom
				cronJobs[name] = job
			}
		}
	}

	err := c.Taskmanager().CronJobsInit(cronJobs)
	return spverrors.Wrapf(err, "failed to init cron jobs")
}

// loadPaymailServer will load the default paymail server configuration
func (c *Client) loadPaymailServer() (err error) {
	// Default FROM paymail
	if len(c.options.paymail.serverConfig.DefaultFromPaymail) == 0 {
		c.options.paymail.serverConfig.DefaultFromPaymail = defaultSenderPaymail
	}

	// Set default options if none are found
	if len(c.options.paymail.serverConfig.options) == 0 {
		c.options.paymail.serverConfig.options = append(c.options.paymail.serverConfig.options,
			paymailserver.WithP2PCapabilities(),
			paymailserver.WithDomainValidationDisabled(),
		)
	}

	paymailLogger := c.Logger().With().Str("subservice", "go-paymail").Logger()
	c.options.paymail.serverConfig.options = append(c.options.paymail.serverConfig.options, paymailserver.WithLogger(&paymailLogger))

	// Create the paymail configuration using the client and default service provider
	paymailLocator := &paymailserver.PaymailServiceLocator{}

	var serviceProvider paymailprovider.ServiceProvider
	var pikeContactProvider paymailserver.PikeContactServiceProvider
	if c.options.paymail.serverConfig.ExperimentalProvider {
		paymailServiceLogger := c.Logger().With().Str("subservice", "paymail-service-provider").Logger()
		serviceProvider = paymailprovider.NewServiceProvider(
			&paymailServiceLogger,
			c.PaymailsService(),
			c.UsersService(),
			c.AddressesService(),
			c.ContactService(),
			c.Chain(),
			c.TransactionRecordService(),
		)

		pikeContactProvider = serviceProvider
	} else {
		serviceProvider = &PaymailDefaultServiceProvider{client: c}
		pikeContactProvider = &PikeContactServiceProvider{client: c}
	}

	paymailLocator.RegisterPaymailService(serviceProvider)
	paymailLocator.RegisterPikeContactService(pikeContactProvider)
	paymailLocator.RegisterPikePaymentService(&PikePaymentServiceProvider{client: c})

	c.options.paymail.serverConfig.Configuration, err = paymailserver.NewConfig(
		paymailLocator,
		c.options.paymail.serverConfig.options...,
	)
	return

}

func (c *Client) askForFeeUnit(ctx context.Context) error {
	feeUnit, err := c.Chain().GetFeeUnit(ctx)
	if err != nil {
		return spverrors.ErrAskingForFeeUnit.Wrap(err)
	}
	c.options.feeUnit = feeUnit
	c.Logger().Info().Msgf("Fee unit set by ARC policy: %d satoshis per %d bytes", feeUnit.Satoshis, feeUnit.Bytes)
	return nil
}
