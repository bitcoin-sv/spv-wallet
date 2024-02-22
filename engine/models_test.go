package engine

import (
	"testing"
	"time"

	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/mrz1836/go-datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testMetadataKey   = "test_key"
	testMetadataValue = "test_value"
)

// TestModelName_String will test the method String()
func TestModelName_String(t *testing.T) {
	t.Parallel()

	t.Run("all model names", func(t *testing.T) {
		assert.Equal(t, "destination", ModelDestination.String())
		assert.Equal(t, "empty", ModelNameEmpty.String())
		assert.Equal(t, "metadata", ModelMetadata.String())
		assert.Equal(t, "paymail_address", ModelPaymailAddress.String())
		assert.Equal(t, "paymail_address", ModelPaymailAddress.String())
		assert.Equal(t, "sync_transaction", ModelSyncTransaction.String())
		assert.Equal(t, "transaction", ModelTransaction.String())
		assert.Equal(t, "utxo", ModelUtxo.String())
		assert.Equal(t, "xpub", ModelXPub.String())
		assert.Len(t, AllModelNames, 9)
	})
}

// TestModelName_IsEmpty will test the method IsEmpty()
func TestModelName_IsEmpty(t *testing.T) {
	t.Parallel()

	t.Run("empty model", func(t *testing.T) {
		assert.Equal(t, true, ModelNameEmpty.IsEmpty())
		assert.Equal(t, false, ModelUtxo.IsEmpty())
	})
}

// TestModel_GetModelName will test the GetModelName function
func TestModel_GetModelName(t *testing.T) {
	t.Parallel()

	t.Run("empty model", func(t *testing.T) {
		assert.Nil(t, datastore.GetModelName(nil))
	})

	t.Run("base model names", func(t *testing.T) {
		xPub := Xpub{}
		assert.Equal(t, ModelXPub.String(), *datastore.GetModelName(xPub))

		destination := Destination{}
		assert.Equal(t, ModelDestination.String(), *datastore.GetModelName(destination))

		utxo := Utxo{}
		assert.Equal(t, ModelUtxo.String(), *datastore.GetModelName(utxo))

		transaction := Transaction{}
		assert.Equal(t, ModelTransaction.String(), *datastore.GetModelName(transaction))

		accessKey := AccessKey{}
		assert.Equal(t, ModelAccessKey.String(), *datastore.GetModelName(accessKey))

		draftTx := DraftTransaction{}
		assert.Equal(t, ModelDraftTransaction.String(), *datastore.GetModelName(draftTx))

		paymailAddress := PaymailAddress{}
		assert.Equal(t, ModelPaymailAddress.String(), *datastore.GetModelName(paymailAddress))

		syncTx := SyncTransaction{}
		assert.Equal(t, ModelSyncTransaction.String(), *datastore.GetModelName(syncTx))
	})
}

// TestModel_GetModelTableName will test the GetModelTableName function
func TestModel_GetModelTableName(t *testing.T) {
	t.Parallel()

	t.Run("empty model", func(t *testing.T) {
		assert.Nil(t, datastore.GetModelTableName(nil))
	})

	t.Run("get model table names", func(t *testing.T) {
		xPub := Xpub{}
		assert.Equal(t, tableXPubs, *datastore.GetModelTableName(xPub))

		destination := Destination{}
		assert.Equal(t, tableDestinations, *datastore.GetModelTableName(destination))

		utxo := Utxo{}
		assert.Equal(t, tableUTXOs, *datastore.GetModelTableName(utxo))

		transaction := Transaction{}
		assert.Equal(t, tableTransactions, *datastore.GetModelTableName(transaction))

		accessKey := AccessKey{}
		assert.Equal(t, tableAccessKeys, *datastore.GetModelTableName(accessKey))

		draftTx := DraftTransaction{}
		assert.Equal(t, tableDraftTransactions, *datastore.GetModelTableName(draftTx))

		paymailAddress := PaymailAddress{}
		assert.Equal(t, tablePaymailAddresses, *datastore.GetModelTableName(paymailAddress))

		syncTx := SyncTransaction{}
		assert.Equal(t, tableSyncTransactions, *datastore.GetModelTableName(syncTx))
	})
}

func (ts *EmbeddedDBTestSuite) createXpubModels(tc *TestingClient, t *testing.T, number int) {
	for i := 0; i < number; i++ {
		_, xPublicKey, err := bitcoin.GenerateHDKeyPair(bitcoin.SecureSeedLength)
		require.NoError(t, err)
		xPub := newXpub(xPublicKey, append(tc.client.DefaultModelOptions(), New())...)
		xPub.CurrentBalance = 125000
		xPub.NextExternalNum = 12
		xPub.NextInternalNum = 37
		err = xPub.Save(tc.ctx)
		require.NoError(t, err)
	}
}

type xPubFieldsTest struct {
	CurrentBalance uint64 `json:"current_balance" toml:"current_balance" yaml:"current_balance" bson:"current_balance"`
}

// TestModels_GetModels will test the method GetModels()
func (ts *EmbeddedDBTestSuite) TestModels_GetModels() {
	numberOfModels := 10
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - GetModels", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)
			ts.createXpubModels(tc, t, numberOfModels)

			queryParams := &datastore.QueryParams{Page: 0, PageSize: 10}
			var models []*Xpub
			err := tc.client.Datastore().GetModels(
				tc.ctx,
				&models,
				nil,
				queryParams,
				nil,
				30*time.Second,
			)
			require.NoError(t, err)
			require.Len(t, models, numberOfModels)
			assert.Equal(t, uint64(125000), models[0].CurrentBalance) // should be set
			assert.Equal(t, uint32(12), models[0].NextExternalNum)    // should be set
			assert.Equal(t, uint32(37), models[0].NextInternalNum)    // should be set
		})

		ts.T().Run(testCase.name+" - GetModels with projection", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)
			ts.createXpubModels(tc, t, numberOfModels)

			queryParams := &datastore.QueryParams{Page: 0, PageSize: 10}
			var models []*Xpub
			var results []*xPubFieldsTest
			err := tc.client.Datastore().GetModels(
				tc.ctx,
				&models,
				nil,
				queryParams,
				&results,
				30*time.Second,
			)
			require.NoError(t, err)
			require.Len(t, results, numberOfModels)
			assert.Equal(t, uint64(125000), results[0].CurrentBalance) // should be set
		})
	}
}
