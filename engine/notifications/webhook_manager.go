package notifications

import (
	"context"
	"sync"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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
	found, err := w.repository.GetByURL(ctx, url)
	if err != nil {
		return spverrors.Wrapf(err, "failed to check existing webhook in database")
	}
	if found != nil {
		found.Refresh(tokenHeader, tokenValue)
		err = w.repository.Save(ctx, found)
	} else {
		err = w.repository.Create(ctx, url, tokenHeader, tokenValue)
	}

	if err != nil {
		return spverrors.Wrapf(err, "failed to store the webhook")
	}

	w.updateMsg <- true
	return nil
}

// Unsubscribe unsubscribes from a webhook. It removes the webhook from the database and stops the notifier for it.
func (w *WebhookManager) Unsubscribe(ctx context.Context, url string) error {
	model, err := w.repository.GetByURL(ctx, url)
	if err != nil || model == nil || model.Deleted() {
		return spverrors.ErrWebhookSubscriptionNotFound
	}
	err = w.repository.Delete(ctx, model)
	if err != nil {
		return spverrors.ErrWebhookUnsubscriptionFailed
	}
	w.updateMsg <- true
	return nil
}

// GetAll returns all the webhooks stored in database
func (w *WebhookManager) GetAll(ctx context.Context) ([]ModelWebhook, error) {
	webhooks, err := w.repository.GetAll(ctx)
	if err != nil {
		w.logger.Warn().Msgf("failed to get webhooks: %v", err)
		return nil, spverrors.ErrWebhookGetAll
	}
	return webhooks, nil
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
			err := w.markWebhookAsBanned(w.rootContext, url)
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
	dbWebhooks, err := w.repository.GetAll(w.rootContext)
	if err != nil {
		w.logger.Warn().Msgf("failed to get webhooks: %v", err)
		return
	}

	// filter out banned webhooks
	var filteredWebhooks []ModelWebhook
	for _, webhook := range dbWebhooks {
		if !webhook.Banned() {
			filteredWebhooks = append(filteredWebhooks, webhook)
		}
	}

	// add notifiers which are not in the map
	for _, model := range filteredWebhooks {
		if _, ok := w.webhookNotifiers.Load(model.GetURL()); !ok {
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
		if item, ok := w.webhookNotifiers.Load(model.GetURL()); ok {
			item.(*notifierWithCtx).notifier.Update(model)
		}
	}
}

func (w *WebhookManager) addNotifier(model ModelWebhook) {
	w.logger.Info().Msgf("Add a webhook notifier. URL: %s", model.GetURL())
	ctx, cancel := context.WithCancel(w.rootContext)
	notifier := NewWebhookNotifier(ctx, w.logger, model, w.banMsg)
	w.webhookNotifiers.Store(model.GetURL(), &notifierWithCtx{notifier: notifier, ctx: ctx, cancelFunc: cancel})
	w.notifications.AddNotifier(model.GetURL(), notifier.Channel)
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

func (w *WebhookManager) markWebhookAsBanned(ctx context.Context, url string) error {
	model, err := w.repository.GetByURL(ctx, url)
	if err != nil {
		return spverrors.Wrapf(err, "cannot find the webhook model")
	}
	model.BanUntil(time.Now().Add(banTime))
	err = w.repository.Save(ctx, model)
	return spverrors.Wrapf(err, "cannot update the webhook model")
}

func containsWebhook(webhooks []ModelWebhook, url string) bool {
	for _, webhook := range webhooks {
		if webhook.GetURL() == url {
			return true
		}
	}
	return false
}
