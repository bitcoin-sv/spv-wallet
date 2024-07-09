package notifications

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type notifierWithCtx struct {
	notifier   *WebhookNotifier
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// WebhookManager is a manager for webhooks. It is responsible for creating, updating and removing webhooks.
type WebhookManager struct {
	repository       WebhooksRepository
	rootContext      context.Context
	cancelAllFunc    context.CancelFunc
	webhookNotifiers *sync.Map // [string, *notifierWithCtx]
	ticker           *time.Ticker
	updateMsg        chan bool
	banMsg           chan string // url
	notifications    *Notifications
	logger           *zerolog.Logger
}

// NewWebhookManager creates a new WebhookManager. It starts a goroutine which checks for webhook updates.
func NewWebhookManager(ctx context.Context, logger *zerolog.Logger, notifications *Notifications, repository WebhooksRepository) *WebhookManager {
	rootContext, cancelAllFunc := context.WithCancel(ctx)
	manager := WebhookManager{
		repository:       repository,
		rootContext:      rootContext,
		cancelAllFunc:    cancelAllFunc,
		webhookNotifiers: &sync.Map{},
		ticker:           time.NewTicker(5 * time.Second),
		notifications:    notifications,
		updateMsg:        make(chan bool),
		banMsg:           make(chan string),
		logger:           logger,
	}

	go manager.checkForUpdates()

	return &manager
}

// Stop stops the WebhookManager.
func (w *WebhookManager) Stop() {
	w.cancelAllFunc()
}

// Subscribe subscribes to a webhook. It adds the webhook to the database and starts a notifier for it.
func (w *WebhookManager) Subscribe(ctx context.Context, url, tokenHeader, tokenValue string) error {
	err := w.repository.CreateWebhook(ctx, url, tokenHeader, tokenValue)
	if err == nil {
		w.updateMsg <- true
	}
	return errors.Wrap(err, "failed to create webhook")
}

// Unsubscribe unsubscribes from a webhook. It removes the webhook from the database and stops the notifier for it.
func (w *WebhookManager) Unsubscribe(ctx context.Context, url string) error {
	err := w.repository.RemoveWebhook(ctx, url)
	if err == nil {
		w.updateMsg <- true
	}
	return errors.Wrap(err, "failed to remove webhook")
}

func (w *WebhookManager) checkForUpdates() {
	defer func() {
		w.logger.Info().Msg("WebhookManager stopped")
		if err := recover(); err != nil {
			w.logger.Warn().Msgf("WebhookManager failed: %v", err)
		}
	}()

	w.logger.Info().Msg("WebhookManager started")
	w.update()

	for {
		select {
		case <-w.ticker.C:
			w.update()
		case <-w.updateMsg:
			w.update()
		case url := <-w.banMsg:
			err := w.repository.BanWebhook(w.rootContext, url, time.Now().Add(banTime))
			if err != nil {
				w.logger.Warn().Msgf("failed to mark a webhook as banned: %v", err)
			}
			w.removeNotifier(url)
		case <-w.rootContext.Done():
			return
		}
	}
}

func (w *WebhookManager) update() {
	defer func() {
		if err := recover(); err != nil {
			w.logger.Warn().Msgf("WebhookManager update failed: %v", err)
		}
	}()
	dbWebhooks, err := w.repository.GetWebhooks(w.rootContext)
	if err != nil {
		w.logger.Warn().Msgf("failed to get webhooks: %v", err)
		return
	}

	// filter out banned webhooks
	var filteredWebhooks []*WebhookModel
	for _, webhook := range dbWebhooks {
		if !webhook.Banned() {
			filteredWebhooks = append(filteredWebhooks, webhook)
		}
	}

	// add notifiers which are not in the map
	for _, model := range filteredWebhooks {
		if _, ok := w.webhookNotifiers.Load(model.URL); !ok {
			w.addNotifier(model)
		}
	}

	// remove notifiers which are not in the database
	w.webhookNotifiers.Range(func(key, _ any) bool {
		url := key.(string)
		if !containsWebhook(filteredWebhooks, url) {
			w.removeNotifier(url)
		}
		return true
	})

	// update definition of remained webhooks
	for _, model := range filteredWebhooks {
		if item, ok := w.webhookNotifiers.Load(model.URL); ok {
			item.(*notifierWithCtx).notifier.Update(*model)
		}
	}
}

func (w *WebhookManager) addNotifier(model *WebhookModel) {
	w.logger.Info().Msgf("Add a webhook notifier. URL: %s", model.URL)
	ctx, cancel := context.WithCancel(w.rootContext)
	notifier := NewWebhookNotifier(ctx, w.logger, *model, w.banMsg)
	w.webhookNotifiers.Store(model.URL, &notifierWithCtx{notifier: notifier, ctx: ctx, cancelFunc: cancel})
	w.notifications.AddNotifier(model.URL, notifier.Channel)
}

func (w *WebhookManager) removeNotifier(url string) {
	if item, ok := w.webhookNotifiers.Load(url); ok {
		w.logger.Info().Msgf("Remove a webhook notifier. URL: %s", url)
		item := item.(*notifierWithCtx)
		item.cancelFunc()
		w.webhookNotifiers.Delete(url)
		w.notifications.RemoveNotifier(url)
	}
}

func containsWebhook(webhooks []*WebhookModel, url string) bool {
	for _, webhook := range webhooks {
		if webhook.URL == url {
			return true
		}
	}
	return false
}
