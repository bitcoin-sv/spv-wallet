package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
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
	Channel    chan Event
	definition *WebhookModel
	banTime    *time.Time
	banMsg     chan bool
	httpClient *http.Client
}

// NewWebhookNotifier - creates a new instance of WebhookNotifier
func NewWebhookNotifier(ctx context.Context, hook *WebhookModel) *WebhookNotifier {
	notifier := &WebhookNotifier{
		Channel:    make(chan Event, lengthOfWebhookChannel),
		definition: hook,
		banMsg:     make(chan bool),
		httpClient: &http.Client{},
	}

	go notifier.consumer(ctx)

	return notifier
}

// consumer - consumer for webhook notifier
// It accumulates events (produced during http call) and sends them to webhook
// If sending fails, it retries several times
// If sending fails after several retries, it bans notifier for some time
func (w *WebhookNotifier) consumer(ctx context.Context) {
	for {
		select {
		case event := <-w.Channel:
			if w.banned() {
				continue // discard events
			}
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
				time.Sleep(retriesDelay)
			}

			if err != nil {
				w.ban()
			}
		case <-w.banMsg:
			w.ban()
		case <-ctx.Done():
			return
		}
	}
}

// Ban - ban notifier for some time
func (w *WebhookNotifier) Ban() {
	w.banMsg <- true
}

// internal ban function; not thread safe
func (w *WebhookNotifier) ban() {
	now := time.Now()
	w.banTime = &now
}

// banned - check if notifier is banned; this is not thread safe, for internal use only
func (w *WebhookNotifier) banned() bool {
	return w.banTime != nil && time.Now().Before(w.banTime.Add(banTime))
}

func (w *WebhookNotifier) accumulateEvents(ctx context.Context, event Event) (events []Event, done bool) {
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

func (w *WebhookNotifier) sendEventsToWebhook(ctx context.Context, events []Event) error {
	data, err := json.Marshal(events)
	if err != nil {
		return errors.Wrap(err, "failed to marshal events")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", w.definition.URL, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	if w.definition.TokenHeader != "" {
		req.Header.Set(w.definition.TokenHeader, w.definition.Token)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// TODO: Handle response

	return nil
}
