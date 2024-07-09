package notifications

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type WebhookRepository interface {
	CreateWebhook(webhook *WebhookModel) error
	RemoveWebhook(url string) error
	GetWebhooks() ([]*WebhookModel, error)
}

type notifierWithCtx struct {
	notifier   *WebhookNotifier
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type WebhookManager struct {
	repository       WebhookRepository
	rootContext      context.Context
	cancelAllFunc    context.CancelFunc
	webhookNotifiers *sync.Map // [string, *notifierWithCtx]
	ticker           *time.Ticker
	updateMsg        chan bool
	notifications    *Notifications
}

func NewWebhookManager(ctx context.Context, notifications *Notifications, repository WebhookRepository) *WebhookManager {
	rootContext, cancelAllFunc := context.WithCancel(ctx)
	manager := WebhookManager{
		repository:       repository,
		rootContext:      rootContext,
		cancelAllFunc:    cancelAllFunc,
		webhookNotifiers: &sync.Map{},
		ticker:           time.NewTicker(5 * time.Second),
		notifications:    notifications,
		updateMsg:        make(chan bool),
	}

	go manager.checkForUpdates()

	return &manager
}

func (w *WebhookManager) Stop() {
	w.cancelAllFunc()
}

func (w *WebhookManager) Subscribe(webhookModel *WebhookModel) error {
	err := w.repository.CreateWebhook(webhookModel)
	if err == nil {
		w.updateMsg <- true
	}
	return errors.Wrap(err, "failed to create webhook")
}

func (w *WebhookManager) Unsubscribe(url string) error {
	err := w.repository.RemoveWebhook(url)
	if err != nil {
		w.updateMsg <- true
	}
	return errors.Wrap(err, "failed to remove webhook")
}

func (w *WebhookManager) checkForUpdates() {
	w.update()

	for {
		select {
		case <-w.ticker.C:
			w.update()
		case <-w.updateMsg:
			w.update()
		case <-w.rootContext.Done():
			return
		}
	}
}

func (w *WebhookManager) update() {
	dbWebhooks, err := w.repository.GetWebhooks()
	if err != nil {
		// log error
		return
	}

	// add notifiers which are not in the map
	for _, model := range dbWebhooks {
		if _, ok := w.webhookNotifiers.Load(model.URL); !ok {
			w.addNotifier(model)
		}
	}

	// remove notifiers which are not in the database
	w.webhookNotifiers.Range(func(key, _ any) bool {
		url := key.(string)
		if !containsWebhook(dbWebhooks, url) {
			w.removeNotifier(url)
		}
		return true
	})
}

func (w *WebhookManager) addNotifier(model *WebhookModel) {
	ctx, cancel := context.WithCancel(w.rootContext)
	notifier := NewWebhookNotifier(ctx, model)
	w.webhookNotifiers.Store(model.URL, &notifierWithCtx{notifier: notifier, ctx: ctx, cancelFunc: cancel})
	w.notifications.AddNotifier(model.URL, notifier.Channel)
}

func (w *WebhookManager) removeNotifier(url string) {
	if item, ok := w.webhookNotifiers.Load(url); ok {
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
