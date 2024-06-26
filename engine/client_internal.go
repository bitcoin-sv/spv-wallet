package engine

import (
	"context"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"

	// "github.com/bitcoin-sv/spv-wallet/engine/notifications"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
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
	notificationService := notifications.NewNotifications(ctx)
	c.options.notifications.client = notificationService
	c.options.notifications.webhookManager = notifications.NewWebhookManager(ctx, notificationService, &WebhooksRepository{client: c})

	// for development purposes only
	i := 0
	ticker := time.NewTicker(100 * time.Microsecond)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				notificationService.Notify(i)
				i++
			}
		}
	}()
	return
}

func (c *Client) SubscribeWebhook(ctx context.Context, url, tokenHeader, token string) error {
	if c.options.notifications == nil || c.options.notifications.webhookManager == nil {
		return nil
	}

	return c.options.notifications.webhookManager.Subscribe(ctx, url, tokenHeader, token)
}

// loadPaymailClient will load the Paymail client
func (c *Client) loadPaymailClient() (err error) {
	// Only load if it's not set (the client can be overloaded)
	if c.options.paymail.client == nil {
		c.options.paymail.client, err = paymail.NewClient()
	}
	return
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

	return c.Taskmanager().CronJobsInit(cronJobs)
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
