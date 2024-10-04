package engine

import (
	"context"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
	paymailclient "github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/paymailaddress"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/go-resty/resty/v2"
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

// loadChainstate will load chainstate configuration and start the Chainstate client
func (c *Client) loadChainstate(ctx context.Context) (err error) {
	// Load chainstate if a custom interface was NOT provided
	if c.options.chainstate.ClientInterface == nil {
		c.options.chainstate.options = append(c.options.chainstate.options, chainstate.WithUserAgent(c.UserAgent()))
		c.options.chainstate.options = append(c.options.chainstate.options, chainstate.WithHTTPClient(c.HTTPClient()))
		c.options.chainstate.options = append(c.options.chainstate.options, chainstate.WithMetrics(c.options.metrics))
		c.options.chainstate.ClientInterface, err = chainstate.NewClient(ctx, c.options.chainstate.options...)
	}

	return
}

// loadDatastore will load the Datastore and start the Datastore client
//
// NOTE: this will run database migrations if the options was set
func (c *Client) loadDatastore(ctx context.Context) (err error) {
	// Add the models to migrate (after loading the client options)
	if len(c.options.models.migrateModelNames) > 0 {
		c.options.dataStore.options = append(
			c.options.dataStore.options,
			datastore.WithAutoMigrate(c.options.models.migrateModels...),
		)
	}

	// Load client (runs ALL options, IE: auto migrate models)
	if c.options.dataStore.ClientInterface == nil {

		// Add custom array and object fields
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

		// Load the datastore client
		c.options.dataStore.ClientInterface, err = datastore.NewClient(
			ctx, c.options.dataStore.options...,
		)
	}
	return
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
		c.options.paymail.client, err = paymail.NewClient()
		if err != nil {
			return
		}
	}
	if c.options.paymail.service == nil {
		logger := c.Logger().With().Str("subservice", "paymail").Logger()
		c.options.paymail.service = paymailclient.NewServiceClient(c.Cachestore(), c.options.paymail.client, logger)
	}
	return
}

func (c *Client) loadPaymailAddressService() error {
	if c.options.paymailAddressService != nil {
		return nil
	}
	c.options.paymailAddressService = paymailaddress.NewService(
		func(ctx context.Context, address string) (string, error) {
			paymailAddress, err := c.GetPaymailAddress(ctx, address)
			if err != nil {
				return "", err
			}
			return paymailAddress.XpubID, nil
		},
		func(ctx context.Context, xPubId string) ([]string, error) {
			page := &datastore.QueryParams{
				PageSize:      10,
				OrderByField:  createdAtField,
				SortDirection: datastore.SortAsc,
			}

			conditions := map[string]interface{}{
				xPubIDField: xPubId,
			}

			addresses, err := c.GetPaymailAddresses(ctx, nil, conditions, page)
			if err != nil {
				return nil, err
			}
			result := make([]string, 0, len(addresses))
			for _, address := range addresses {
				result = append(result, address.String())
			}
			return result, nil
		},
	)
	return nil
}

func (c *Client) loadTransactionDraftService() error {
	if c.options.transactionDraftService == nil {
		logger := c.Logger().With().Str("subservice", "transactionDraft").Logger()
		c.options.transactionDraftService = draft.NewDraftService(c.PaymailService(), c.options.paymailAddressService, logger)
	}
	return nil
}

func (c *Client) loadChainService() {
	if c.options.chainService == nil {
		logger := c.Logger().With().Str("subservice", "chain").Logger()
		c.options.chainService = chain.NewChainService(logger, resty.New(), c.options.arcConfig, c.options.bhsConfig)
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

// runModelMigrations will run the model Migrate() method for all models
func (c *Client) runModelMigrations(models ...interface{}) (err error) {
	// If the migrations are disabled, just return
	if c.options.dataStore.migrationDisabled {
		return nil
	}

	// Migrate the models
	d := c.Datastore()
	for _, model := range models {
		model.(ModelInterface).SetOptions(WithClient(c))
		if err = model.(ModelInterface).Migrate(d); err != nil {
			return
		}
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

// loadDefaultPaymailConfig will load the default paymail server configuration
func (c *Client) loadDefaultPaymailConfig() (err error) {
	// Default FROM paymail
	if len(c.options.paymail.serverConfig.DefaultFromPaymail) == 0 {
		c.options.paymail.serverConfig.DefaultFromPaymail = defaultSenderPaymail
	}

	// Set default options if none are found
	if len(c.options.paymail.serverConfig.options) == 0 {
		c.options.paymail.serverConfig.options = append(c.options.paymail.serverConfig.options,
			server.WithP2PCapabilities(),
			server.WithDomainValidationDisabled(),
		)
	}

	paymailLogger := c.Logger().With().Str("subservice", "go-paymail").Logger()
	c.options.paymail.serverConfig.options = append(c.options.paymail.serverConfig.options, server.WithLogger(&paymailLogger))

	// Create the paymail configuration using the client and default service provider
	paymailLocator := &server.PaymailServiceLocator{}
	paymailService := &PaymailDefaultServiceProvider{client: c}
	paymailLocator.RegisterPaymailService(paymailService)
	paymailLocator.RegisterPikeContactService(&PikeContactServiceProvider{client: c})
	paymailLocator.RegisterPikePaymentService(&PikePaymentServiceProvider{client: c})

	c.options.paymail.serverConfig.Configuration, err = server.NewConfig(
		paymailLocator,
		c.options.paymail.serverConfig.options...,
	)
	return
}
