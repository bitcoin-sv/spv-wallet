package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	maxBatchSize           = 100
	mexRetries             = 2
	retriesDelay           = 1 * time.Second
	banTime                = 60 * time.Minute
	lengthOfWebhookChannel = 100
)

// WebhookNotifier - notifier for sending events to webhook
type WebhookNotifier struct {
	Channel       chan *RawEvent
	banMsg        chan string
	httpClient    *http.Client
	definition    WebhookModel
	definitionMtx sync.Mutex
}

// NewWebhookNotifier - creates a new instance of WebhookNotifier
func NewWebhookNotifier(ctx context.Context, model WebhookModel, banMsg chan string) *WebhookNotifier {
	notifier := &WebhookNotifier{
		Channel:    make(chan *RawEvent, lengthOfWebhookChannel),
		definition: model,
		banMsg:     banMsg,
		httpClient: &http.Client{},
	}

	go notifier.consumer(ctx)

	return notifier
}

// Update - updates the webhook model
func (w *WebhookNotifier) Update(model WebhookModel) {
	w.definitionMtx.Lock()
	defer w.definitionMtx.Unlock()

	w.definition = model
}

func (w *WebhookNotifier) currentDefinition() WebhookModel {
	w.definitionMtx.Lock()
	defer w.definitionMtx.Unlock()

	return w.definition
}

// consumer - consumer for webhook notifier
// It accumulates events (produced during http call) and sends them to webhook
// If sending fails, it retries several times
// If sending fails after several retries, it bans notifier for some time
func (w *WebhookNotifier) consumer(ctx context.Context) {
	for {
		select {
		case event := <-w.Channel:
			events, done := w.accumulateEvents(ctx, event)
			if done {
				return
			}
			var err error
			for i := 0; i < mexRetries; i++ {
				err = w.sendEventsToWebhook(ctx, events)
				if err == nil {
					break
				}
				select {
				case <-ctx.Done():
					return
				case <-time.After(retriesDelay):
				}
			}

			if err != nil {
				w.banMsg <- w.currentDefinition().URL
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *WebhookNotifier) accumulateEvents(ctx context.Context, event *RawEvent) (events []*RawEvent, done bool) {
	events = append(events, event)
loop:
	for i := 0; i < maxBatchSize; i++ {
		select {
		case event := <-w.Channel:
			events = append(events, event)
		case <-ctx.Done():
			return nil, true
		default:
			break loop
		}
	}
	return events, false
}

func (w *WebhookNotifier) sendEventsToWebhook(ctx context.Context, events []*RawEvent) error {
	definition := w.currentDefinition()
	data, err := json.Marshal(events)
	if err != nil {
		return errors.Wrap(err, "failed to marshal events")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", definition.URL, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	tokenHeader, tokenValue := definition.TokenHeader, definition.TokenValue
	if tokenHeader != "" {
		req.Header.Set(tokenHeader, tokenValue)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return nil
}
