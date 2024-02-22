// Package notifications is a basic internal notifications module
package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// GetWebhookEndpoint will get the configured webhook endpoint
func (c *Client) GetWebhookEndpoint() string {
	return c.options.config.webhookEndpoint
}

// Notify will create a new notification event
func (c *Client) Notify(ctx context.Context, modelType string, eventType EventType,
	model interface{}, id string) error {

	if len(c.options.config.webhookEndpoint) == 0 {
		if c.IsDebug() {
			c.Logger().Info().Msgf("NOTIFY %s: %s - %v", eventType, id, model)
		}
	} else {
		jsonData, err := json.Marshal(map[string]interface{}{
			"event_type": eventType,
			"id":         id,
			"model":      model,
			"model_type": modelType,
		})
		if err != nil {
			return err
		}

		var req *http.Request
		if req, err = http.NewRequestWithContext(ctx,
			http.MethodPost,
			c.options.config.webhookEndpoint,
			bytes.NewBuffer(jsonData),
		); err != nil {
			return err
		}

		var response *http.Response
		if response, err = c.options.httpClient.Do(req); err != nil {
			return err
		}
		defer func() {
			_ = response.Body.Close()
		}()

		if response.StatusCode != http.StatusOK {
			// todo queue notification for another try ...
			c.Logger().Error().Msgf("received invalid response from notification endpoint: %d",
				response.StatusCode)
		}
	}

	return nil
}
