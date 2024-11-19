package engine

import (
	"context"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	customTypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/pkg/errors"
)

// Webhook stores information about subscriptions to notifications via webhooks
//
// Gorm related models & indexes: https://gorm.io/docs/models.html - https://gorm.io/docs/indexes.html
type Webhook struct {
	// Base model
	Model

	URL         string               `json:"url" toml:"url" yaml:"url" gorm:"<-create;primaryKey;comment:This is the url on which notifications will be sent"`
	TokenHeader string               `json:"token_header" toml:"token_header" yaml:"token_header" gorm:"<-create;comment:This is optional token header to be sent"`
	Token       string               `json:"token" toml:"token" yaml:"token" gorm:"<-create;comment:This is optional token to be sent"`
	BannedTo    customTypes.NullTime `json:"banned_to" toml:"banned_to" yaml:"banned_to" gorm:"comment:The time until the webhook will be banned"`
}

func newWebhook(url, tokenHeader, token string, opts ...ModelOps) *Webhook {
	return &Webhook{
		Model:       *NewBaseModel(ModelWebhook, opts...),
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

// PostMigrate is called after the model is migrated
func (m *Webhook) PostMigrate(client datastore.ClientInterface) error {
	err := client.IndexMetadata(client.GetTableName(tableAccessKeys), metadataField)
	return spverrors.Wrapf(err, "failed to index metadata column on model %s", m.GetModelName())
}

func (m *Webhook) delete() {
	m.DeletedAt.Valid = true
	m.DeletedAt.Time = time.Now()
}

// Banned returns true if the webhook is banned right now
func (m *Webhook) Banned() bool {
	if !m.BannedTo.Valid {
		return false
	}
	ret := !time.Now().After(m.BannedTo.Time)
	return ret
}

// GetURL returns the URL of the webhook
func (m *Webhook) GetURL() string {
	return m.URL
}

// GetTokenHeader returns the token header of the webhook
func (m *Webhook) GetTokenHeader() string {
	return m.TokenHeader
}

// GetTokenValue returns the token value of the webhook
func (m *Webhook) GetTokenValue() string {
	return m.Token
}

// BanUntil sets BannedTo field to the given time
func (m *Webhook) BanUntil(bannedTo time.Time) {
	m.BannedTo.Valid = true
	m.BannedTo.Time = bannedTo
}

// Refresh sets the DeletedAt and BannedTo fields to the zero value and updates the token header and value
func (m *Webhook) Refresh(tokenHeader, tokenValue string) {
	m.DeletedAt.Valid = false
	m.BannedTo.Valid = false
	m.TokenHeader = tokenHeader
	m.Token = tokenValue
}

// Deleted returns true if the webhook is deleted
func (m *Webhook) Deleted() bool {
	return m.DeletedAt.Valid
}

// WebhooksRepository is the repository for webhooks. It implements the WebhooksRepository interface
type WebhooksRepository struct {
	client *Client
}

// Create makes a new webhook instance and saves it to the database, it will fail if the webhook already exists in the database
func (wr *WebhooksRepository) Create(ctx context.Context, url, tokenHeader, tokenValue string) error {
	opts := append(wr.client.DefaultModelOptions(), New())
	model := newWebhook(url, tokenHeader, tokenValue, opts...)
	return model.Save(ctx)
}

// Save stores a model in the database
func (wr *WebhooksRepository) Save(ctx context.Context, model notifications.ModelWebhook) error {
	webhook, ok := model.(*Webhook)
	if !ok {
		return spverrors.Newf("unknown implementation of notifications.ModelWebhook")
	}
	err := webhook.Save(ctx)
	return spverrors.Wrapf(err, "cannot save the ModelWebhook")
}

// Delete stores a model in the database
func (wr *WebhooksRepository) Delete(ctx context.Context, model notifications.ModelWebhook) error {
	webhook, ok := model.(*Webhook)
	if !ok {
		return spverrors.Newf("unknown implementation of notifications.ModelWebhook")
	}
	webhook.delete()
	err := webhook.Save(ctx)
	return spverrors.Wrapf(err, "cannot save the ModelWebhook")
}

// GetByURL gets a webhook by its URL. If the webhook does not exist, it returns a nil pointer and no error
func (wr *WebhooksRepository) GetByURL(ctx context.Context, url string) (notifications.ModelWebhook, error) {
	conditions := map[string]any{
		"url": url,
	}

	webhook := &Webhook{}
	webhook.enrich(ModelWebhook, wr.client.DefaultModelOptions()...)

	if err := Get(ctx, webhook, conditions, false, defaultDatabaseReadTimeout, false); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return webhook, nil
}

// GetAll gets all webhooks from the database
func (wr *WebhooksRepository) GetAll(ctx context.Context) ([]notifications.ModelWebhook, error) {
	conditions := map[string]any{
		deletedAtField: nil,
	}
	list, err := getWebhooks(ctx, conditions, wr.client.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}
	// map to slice of ModelWebhook
	res := make([]notifications.ModelWebhook, len(list))
	for i, elem := range list {
		res[i] = elem
	}
	return res, nil
}
