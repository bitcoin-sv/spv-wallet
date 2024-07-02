package engine

import (
	"context"
	"errors"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	customTypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
)

// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type Webhook struct {
	// Base model
	Model `bson:",inline"`

	URL         string               `json:"url" toml:"url" yaml:"url" gorm:"<-create;primaryKey;comment:This is the url on which notifications will be sent" bson:"url"`
	TokenHeader string               `json:"token_header" toml:"token_header" yaml:"token_header" gorm:"<-create;comment:This is optioal token header to be sent" bson:"token_header"`
	Token       string               `json:"token" toml:"token" yaml:"token" gorm:"<-create;comment:This is optional token to be sent" bson:"token"`
	BannedTo    customTypes.NullTime `json:"banned_to" toml:"banned_to" yaml:"banned_to" gorm:"comment:The time until the webhook will be banned" bson:"banned_to"`
}

func newWebhook(url, tokenHeader, token string, opts ...ModelOps) *Webhook {
	return &Webhook{
		Model:       *NewBaseModel(ModelAccessKey, opts...),
		URL:         url,
		TokenHeader: tokenHeader,
		Token:       token,
	}
}

func getWebhooks(ctx context.Context, conditions map[string]any, opts ...ModelOps) ([]*Webhook, error) {
	modelItems := make([]*Webhook, 0)
	if err := getModelsByConditions(ctx, ModelAccessKey, &modelItems, nil, conditions, nil, opts...); err != nil {
		return nil, err
	}

	return modelItems, nil
}

// GetModelName will get the name of the current model
func (m *Webhook) GetModelName() string {
	return ModelWebhook.String()
}

// GetModelTableName will get the db table name of the current model
func (m *Webhook) GetModelTableName() string {
	return tableWebhooks
}

// Save will save the model into the Datastore
func (m *Webhook) Save(ctx context.Context) error {
	return Save(ctx, m)
}

// GetID will get the ID
func (m *Webhook) GetID() string {
	return m.URL
}

// BeforeCreating will fire before the model is being inserted into the Datastore
func (m *Webhook) BeforeCreating(_ context.Context) error {
	return nil
}

// Migrate model specific migration on startup
func (m *Webhook) Migrate(client datastore.ClientInterface) error {
	return client.IndexMetadata(client.GetTableName(tableAccessKeys), metadataField)
}

// GetWebhookModel will return the webhook model
func (m *Webhook) GetWebhookModel() *notifications.WebhookModel {
	model := &notifications.WebhookModel{
		URL:         m.URL,
		TokenHeader: m.TokenHeader,
		TokenValue:  m.Token,
	}

	if m.BannedTo.Valid {
		model.BannedTo = &m.BannedTo.Time
	}
	return model
}

// WebhooksRepository is the repository for webhooks. It implements the WebhooksRepository interface
type WebhooksRepository struct {
	client *Client
}

// CreateWebhook stores a new webhook or updates an existing one in the database
func (wr *WebhooksRepository) CreateWebhook(ctx context.Context, url, tokenHeader, tokenValue string) error {
	existed, err := wr.getByURL(ctx, url)
	if err != nil {
		return err
	}
	if existed != nil {
		existed.DeletedAt.Valid = false
		existed.BannedTo.Valid = false
		existed.TokenHeader = tokenHeader
		existed.Token = tokenValue
		return existed.Save(ctx)
	}
	opts := append(wr.client.DefaultModelOptions(), New())
	model := newWebhook(url, tokenHeader, tokenValue, opts...)
	return model.Save(ctx)
}

// getByURL gets a webhook by its URL. If the webhook does not exist, it returns a nil pointer and no error
func (wr *WebhooksRepository) getByURL(ctx context.Context, url string) (*Webhook, error) {
	conditions := map[string]any{
		"url": url,
	}

	webhook := &Webhook{}
	webhook.enrich(ModelContact, wr.client.DefaultModelOptions()...)

	if err := Get(ctx, webhook, conditions, false, defaultDatabaseReadTimeout, false); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return webhook, nil
}

// RemoveWebhook removes a webhook from the database
func (wr *WebhooksRepository) RemoveWebhook(ctx context.Context, url string) error {
	webhook, err := wr.getByURL(ctx, url)
	if err != nil {
		return err
	}

	webhook.DeletedAt.Valid = true
	webhook.DeletedAt.Time = time.Now()

	return Save(ctx, webhook)
}

// BanWebhook bans a webhook until the specified time
func (wr *WebhooksRepository) BanWebhook(ctx context.Context, url string, bannedTo time.Time) error {
	webhook, err := wr.getByURL(ctx, url)
	if err != nil {
		return err
	}

	webhook.BannedTo.Valid = true
	webhook.BannedTo.Time = bannedTo

	return Save(ctx, webhook)
}

// GetWebhooks gets all webhooks from the database
func (wr *WebhooksRepository) GetWebhooks(ctx context.Context) ([]*notifications.WebhookModel, error) {
	conditions := map[string]any{
		deletedAtField: nil,
	}
	list, err := getWebhooks(ctx, conditions, wr.client.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}
	// map to slice of WebhookInterface
	res := make([]*notifications.WebhookModel, len(list))
	for i, elem := range list {
		res[i] = elem.GetWebhookModel()
	}
	return res, nil
}
