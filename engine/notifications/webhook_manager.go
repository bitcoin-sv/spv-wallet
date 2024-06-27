package notifications

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type notifierWithCtx struct {
	notifier   *WebhookNotifier
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type WebhookManager struct {
	repository       WebhooksRepository
	rootContext      context.Context
	cancelAllFunc    context.CancelFunc
	webhookNotifiers *sync.Map // [string, *notifierWithCtx]
	ticker           *time.Ticker
	updateMsg        chan bool
	banMsg           chan string // url
	notifications    *Notifications
}

func NewWebhookManager(ctx context.Context, notifications *Notifications, repository WebhooksRepository) *WebhookManager {
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
	}

	go manager.checkForUpdates()

	return &manager
}

func (w *WebhookManager) Stop() {
	w.cancelAllFunc()
}

func (w *WebhookManager) Subscribe(ctx context.Context, url, tokenHeader, tokenValue string) error {
	err := w.repository.CreateWebhook(ctx, url, tokenHeader, tokenValue)
	if err == nil {
		w.updateMsg <- true
	}
	return errors.Wrap(err, "failed to create webhook")
}

func (w *WebhookManager) Unsubscribe(ctx context.Context, url string) error {
	err := w.repository.RemoveWebhook(ctx, url)
	if err == nil {
		w.updateMsg <- true
	}
	return errors.Wrap(err, "failed to remove webhook")
}

func (w *WebhookManager) checkForUpdates() {
	defer func() {
		fmt.Printf("WebhookManager stopped\n")
		if err := recover(); err != nil {
			fmt.Printf("WebhookManager stopped with error: %v\n", err)
		}
	}()

	w.update()

	for {
		select {
		case <-w.ticker.C:
			w.update()
		case <-w.updateMsg:
			w.update()
		case url := <-w.banMsg:
			w.repository.BanWebhook(w.rootContext, url, time.Now().Add(banTime)) // TODO log error from this method
			w.removeNotifier(url)
		case <-w.rootContext.Done():
			return
		}
	}
}

func (w *WebhookManager) update() {
	dbWebhooks, err := w.repository.GetWebhooks(w.rootContext)
	if err != nil {
		// log error
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
	ctx, cancel := context.WithCancel(w.rootContext)
	notifier := NewWebhookNotifier(ctx, *model, w.banMsg)
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
