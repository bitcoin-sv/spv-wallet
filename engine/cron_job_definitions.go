package engine

import (
	"context"
	"errors"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"gorm.io/gorm"
)

// taskCleanupDraftTransactions will clean up all old expired draft transactions
func taskCleanupDraftTransactions(ctx context.Context, client *Client) error {
	client.Logger().Info().Msg("running cleanup draft transactions task...")

	// Construct an empty model
	var models []DraftTransaction
	conditions := map[string]interface{}{
		statusField: DraftStatusDraft,
		// todo: add DB condition for date "expires_at": map[string]interface{}{"$lte": time.Now()},
	}

	queryParams := &datastore.QueryParams{
		Page:          1,
		PageSize:      20,
		OrderByField:  idField,
		SortDirection: datastore.SortAsc,
	}

	// Get the records
	if err := getModels(
		ctx, client.Datastore(),
		&models, conditions, queryParams, defaultDatabaseReadTimeout,
	); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	// Loop and update
	var err error
	timeNow := time.Now().UTC()
	for index := range models {
		if timeNow.After(models[index].ExpiresAt) {
			models[index].enrich(ModelDraftTransaction, WithClient(client))
			models[index].Status = DraftStatusExpired
			if err = models[index].Save(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

// taskSyncTransactions will sync any transactions
func taskSyncTransactions(ctx context.Context, client *Client) error {
	logClient := client.Logger()
	logClient.Info().Msg("running sync transaction(s) task...")

	processSyncTransactions(ctx, client)
	return nil
}

func taskCalculateMetrics(ctx context.Context, client *Client) error {
	m, enabled := client.Metrics()
	if !enabled {
		return spverrors.Newf("metrics are not enabled")
	}

	modelOpts := client.DefaultModelOptions()

	if xpubsCount, err := getXPubsCount(ctx, nil, nil, modelOpts...); err != nil {
		client.options.logger.Error().Err(err).Msg("error getting xpubs count")
	} else {
		m.SetXPubCount(xpubsCount)
	}

	if utxosCount, err := getUtxosCount(ctx, nil, nil, modelOpts...); err != nil {
		client.options.logger.Error().Err(err).Msg("error getting utxos count")
	} else {
		m.SetUtxoCount(utxosCount)
	}

	if paymailsCount, err := getPaymailAddressesCount(ctx, nil, nil, modelOpts...); err != nil {
		client.options.logger.Error().Err(err).Msg("error getting paymails count")
	} else {
		m.SetPaymailCount(paymailsCount)
	}

	if destinationsCount, err := getDestinationsCount(ctx, nil, nil, modelOpts...); err != nil {
		client.options.logger.Error().Err(err).Msg("error getting destinations count")
	} else {
		m.SetDestinationCount(destinationsCount)
	}

	if accessKeysCount, err := getAccessKeysCount(ctx, nil, nil, modelOpts...); err != nil {
		client.options.logger.Error().Err(err).Msg("error getting access keys count")
	} else {
		m.SetAccessKeyCount(accessKeysCount)
	}

	return nil
}
