package engine

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClient_NewDestination will test the method NewDestination()
func (ts *EmbeddedDBTestSuite) TestClient_NewDestination() {
	for _, testCase := range dbTestCases {

		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			ctx := context.Background()

			_, err := tc.client.NewXpub(ctx, testXPub, tc.client.DefaultModelOptions()...)
			assert.NoError(t, err)

			metadata := map[string]interface{}{
				"test-key": "test-value",
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			var destination *Destination
			destination, err = tc.client.NewDestination(
				ctx, testXPub, utils.ChainExternal, utils.ScriptTypePubKeyHash, opts...,
			)
			assert.NoError(t, err)
			assert.Equal(t, "fc1e635d98151c6008f29908ee2928c60c745266f9853e945c917b1baa05973e", destination.ID)
			assert.Equal(t, testXPubID, destination.XpubID)
			assert.Equal(t, utils.ScriptTypePubKeyHash, destination.Type)
			assert.Equal(t, utils.ChainExternal, destination.Chain)
			assert.Equal(t, uint32(0), destination.Num)
			assert.Equal(t, testExternalAddress, destination.Address)
			assert.Equal(t, "test-value", destination.Metadata["test-key"])

			destination2, err2 := tc.client.NewDestination(
				ctx, testXPub, utils.ChainExternal, utils.ScriptTypePubKeyHash, opts...,
			)
			assert.NoError(t, err2)
			assert.Equal(t, testXPubID, destination2.XpubID)
			// assert.Equal(t, "1234567", destination2.Metadata[ReferenceIDField])
			assert.Equal(t, utils.ScriptTypePubKeyHash, destination2.Type)
			assert.Equal(t, utils.ChainExternal, destination2.Chain)
			assert.Equal(t, uint32(1), destination2.Num)
			assert.NotEqual(t, testExternalAddress, destination2.Address)
			assert.Equal(t, "test-value", destination2.Metadata["test-key"])
		})

		ts.T().Run(testCase.name+" - error - missing xpub", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			metadata := map[string]interface{}{
				"test-key": "test-value",
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			destination, err := tc.client.NewDestination(
				context.Background(), testXPub, utils.ChainExternal,
				utils.ScriptTypePubKeyHash, opts...,
			)
			require.Error(t, err)
			require.Nil(t, destination)
			assert.ErrorIs(t, err, ErrMissingXpub)
		})
	}
}

// TestClient_NewDestinationForLockingScript will test the method NewDestinationForLockingScript()
func (ts *EmbeddedDBTestSuite) TestClient_NewDestinationForLockingScript() {
	for _, testCase := range dbTestCases {

		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, err := tc.client.NewXpub(tc.ctx, testXPub, tc.client.DefaultModelOptions()...)
			assert.NoError(t, err)

			lockingScript := "14c91e5cc393bb9d6da3040a7c72b4b569b237e450517901687f517f7c76767601ff9c636d75587f7c6701fe" +
				"9c636d547f7c6701fd9c6375527f7c67686868817f7b6d517f7c7f77605b955f937f517f787f517f787f567f01147f527f7577" +
				"7e777e7b7c7e7b7c7ea77b885279887601447f01207f75776baa517f7c818b7c7e263044022079be667ef9dcbbac55a06295ce" +
				"870b07029bfcdb2dce28d959f2815b16f8179802207c7e01417e2102b405d7f0322a89d0f9f3a98e6f938fdc1c969a8d1382a2" +
				"bf66a71ae74a1e83b0ad046d6574612102b8e6b4441609460d1605ce328d7a39e7216050e105738725b05b7b542dcf1f51205f" +
				"a8b5671a8b577a44ea2d1e70ca9c291145d3da3a7c649fc4e9ea389a8053646c886d76a9146e12c6d84b06757bd4316c33cac4" +
				"4e1e5965589088ac6a0b706172656e74206e6f6465"

			metadata := map[string]interface{}{"test_key": "test_value"}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			var destination *Destination
			destination, err = tc.client.NewDestinationForLockingScript(
				tc.ctx, testXPubID, lockingScript, opts...,
			)
			assert.NoError(t, err)
			assert.Equal(t, "a64c7aca7110c7cde92245252a58bb18a4317381fc31fc293f6aafa3fcc7019f", destination.ID)
			assert.Equal(t, testXPubID, destination.XpubID)
			assert.Equal(t, utils.ScriptTypeNonStandard, destination.Type)
			assert.Equal(t, utils.ChainExternal, destination.Chain)
			assert.Equal(t, uint32(0), destination.Num)
			assert.Equal(t, "test_value", destination.Metadata["test_key"])
		})

		ts.T().Run(testCase.name+" - error - missing locking script", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			metadata := map[string]interface{}{"test_key": "test_value"}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			destination, err := tc.client.NewDestinationForLockingScript(
				tc.ctx, testXPubID, "",
				opts...,
			)
			require.Error(t, err)
			require.Nil(t, destination)
			assert.ErrorIs(t, err, ErrMissingLockingScript)
		})
	}
}

// TestClient_GetDestinations will test the method GetDestinationsByXpubID()
func (ts *EmbeddedDBTestSuite) TestClient_GetDestinations() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			var getDestinations []*Destination
			getDestinations, err = tc.client.GetDestinationsByXpubID(
				tc.ctx, xPubID, nil, nil, nil,
			)
			require.NoError(t, err)
			require.NotNil(t, getDestinations)
			assert.Equal(t, 1, len(getDestinations))
			assert.Equal(t, destination.Address, getDestinations[0].Address)
			assert.Equal(t, testReferenceID, getDestinations[0].Metadata[ReferenceIDField])
			assert.Equal(t, destination.XpubID, getDestinations[0].XpubID)
		})

		ts.T().Run(testCase.name+" - no destinations found", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)

			metadata := map[string]interface{}{testMetadataKey: testMetadataValue}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			// use the wrong xpub
			var getDestinations []*Destination
			getDestinations, err = tc.client.GetDestinationsByXpubID(
				tc.ctx, testXPubID, nil, nil, nil,
			)
			require.NoError(t, err)
			assert.Equal(t, 0, len(getDestinations))
		})

		ts.T().Run(testCase.name+" with locking_script filter", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			conditions := map[string]interface{}{
				"locking_script": destination.LockingScript,
			}

			var getDestinations []*Destination
			getDestinations, err = tc.client.GetDestinationsByXpubID(
				tc.ctx, xPubID, nil, conditions, nil,
			)
			fmt.Printf("Destinatosn %+v", getDestinations[0].LockingScript)
			require.NoError(t, err)
			require.NotNil(t, getDestinations)
			assert.Equal(t, 1, len(getDestinations))
			assert.Equal(t, destination.XpubID, getDestinations[0].XpubID)
			assert.Equal(t, destination.LockingScript, getDestinations[0].LockingScript)
		})

		ts.T().Run(testCase.name+" with address filter", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			conditions := map[string]interface{}{
				"address": destination.Address,
			}

			var getDestinations []*Destination
			getDestinations, err = tc.client.GetDestinationsByXpubID(
				tc.ctx, xPubID, nil, conditions, nil,
			)
			require.NoError(t, err)
			require.NotNil(t, getDestinations)
			assert.Equal(t, 1, len(getDestinations))
			assert.Equal(t, destination.XpubID, getDestinations[0].XpubID)
			assert.Equal(t, destination.Address, getDestinations[0].Address)
		})

		ts.T().Run(testCase.name+" with draft_id filter", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			// conditions := models.DestinationFilters{DraftID: &destination.DraftID}

			//var getDestinations []*Destination
			//getDestinations, err = tc.client.GetDestinationsByXpubID(
			//	tc.ctx, xPubID, nil, &conditions, nil,
			//)
			dests, err := tc.client.GetDestinationsByXpubID(
				tc.ctx, xPubID, nil, nil, nil,
			)
			fmt.Printf("bbb %+v", &dests)
			require.NoError(t, err)
			require.NotNil(t, dests)
			assert.Equal(t, 1, len(dests))
			assert.Equal(t, destination.XpubID, dests[0].XpubID)
			assert.Equal(t, destination.DraftID, dests[0].DraftID)
		})

		ts.T().Run(testCase.name+" with include_deleted true filter", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			// deleted items should be present by default (empty conditions)
			conditions := make(map[string]interface{})

			var getDestinations []*Destination
			getDestinations, err = tc.client.GetDestinationsByXpubID(
				tc.ctx, xPubID, nil, conditions, nil,
			)
			require.NoError(t, err)
			require.NotNil(t, getDestinations)
			assert.Equal(t, 1, len(getDestinations))
			assert.Equal(t, destination.XpubID, getDestinations[0].XpubID)
		})

		ts.T().Run(testCase.name+" with include_deleted false filter", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			// when deleted_at is NULL id db - we treat it as not deleted
			conditions := map[string]interface{}{
				"deleted_at": nil,
			}

			var getDestinations []*Destination
			getDestinations, err = tc.client.GetDestinationsByXpubID(
				tc.ctx, xPubID, nil, conditions, nil,
			)
			require.NoError(t, err)
			require.NotNil(t, getDestinations)
			assert.Equal(t, 1, len(getDestinations))
			assert.Equal(t, destination.XpubID, getDestinations[0].XpubID)
		})

		ts.T().Run(testCase.name+" with created_range filter valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			fmt.Println("New dest")
			require.NoError(t, err)
			require.NotNil(t, destination)
			fmt.Println("test", destination.CreatedAt)

			fromTime, _ := time.Parse(time.RFC3339Nano, "2020-02-26T11:01:28.069911Z")
			toTime, _ := time.Parse(time.RFC3339Nano, "2035-02-26T11:01:28.069911Z")

			conditions := map[string]interface{}{
				"created_at": map[string]interface{}{
					"$gte": fromTime,
					"$lte": toTime,
				},
			}

			dests, err := tc.client.GetDestinationsByXpubID(
				tc.ctx, xPubID, nil, conditions, nil,
			)
			require.NoError(t, err)
			fmt.Println("Updated", dests)
			require.NotNil(t, dests)
			assert.Equal(t, 1, len(dests))
			assert.Equal(t, destination.XpubID, dests[0].XpubID)
		})

		ts.T().Run(testCase.name+" with created_range filter invalid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			conditions := map[string]interface{}{
				"created_at": 123,
			}

			var getDestinations []*Destination
			getDestinations, err = tc.client.GetDestinationsByXpubID(
				tc.ctx, xPubID, nil, conditions, nil,
			)
			require.NoError(t, err)
			require.Nil(t, getDestinations)
		})

		ts.T().Run(testCase.name+" with created_range filter valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			fromTime, _ := time.Parse(time.RFC3339Nano, "2020-02-26T11:01:28.069911Z")
			toTime, _ := time.Parse(time.RFC3339Nano, "2030-02-26T11:01:28.069911Z")

			conditions := map[string]interface{}{
				"updated_at": map[string]interface{}{
					"$gte": fromTime,
					"$lte": toTime,
				},
			}

			dests, err := tc.client.GetDestinationsByXpubID(
				tc.ctx, xPubID, nil, conditions, nil,
			)
			require.NoError(t, err)
			require.NotNil(t, dests)
			assert.Equal(t, 1, len(dests))
			assert.Equal(t, destination.XpubID, dests[0].XpubID)
		})

		ts.T().Run(testCase.name+" with updated_range filter invalid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			conditions := map[string]interface{}{
				"updated_at": 123,
			}

			// var getDestinations []*Destination
			dests, err := tc.client.GetDestinationsByXpubID(
				tc.ctx, xPubID, nil, conditions, nil,
			)
			require.NoError(t, err)
			require.NotNil(t, dests)
			assert.Equal(t, 1, len(dests))
			// assert.Equal(t, destination.XpubID, getDestinations[0].XpubID)
		})

	}
}

// TestClient_GetDestinationByAddress will test the method GetDestinationByAddress()
func (ts *EmbeddedDBTestSuite) TestClient_GetDestinationByAddress() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			var getDestination *Destination
			getDestination, err = tc.client.GetDestinationByAddress(
				tc.ctx, xPubID, destination.Address,
			)
			require.NoError(t, err)
			require.NotNil(t, getDestination)
			assert.Equal(t, destination.Address, getDestination.Address)
			assert.Equal(t, testReferenceID, getDestination.Metadata[ReferenceIDField])
			assert.Equal(t, destination.XpubID, getDestination.XpubID)
		})

		ts.T().Run(testCase.name+" - invalid xpub", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)

			metadata := map[string]interface{}{testMetadataKey: testMetadataValue}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			// use the wrong xpub
			var getDestination *Destination
			getDestination, err = tc.client.GetDestinationByAddress(
				tc.ctx, testXPubID, destination.Address,
			)
			require.Error(t, err)
			require.Nil(t, getDestination)
		})
	}
}

// TestClient_GetDestinationByLockingScript will test the method GetDestinationByLockingScript()
func (ts *EmbeddedDBTestSuite) TestClient_GetDestinationByLockingScript() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)
			xPubID := utils.Hash(rawKey)

			metadata := map[string]interface{}{
				ReferenceIDField: testReferenceID,
				testMetadataKey:  testMetadataValue,
			}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			var getDestination *Destination
			getDestination, err = tc.client.GetDestinationByLockingScript(
				tc.ctx, xPubID, destination.LockingScript,
			)
			require.NoError(t, err)
			require.NotNil(t, getDestination)
			assert.Equal(t, destination.Address, getDestination.Address)
			assert.Equal(t, destination.LockingScript, getDestination.LockingScript)
			assert.Equal(t, testReferenceID, getDestination.Metadata[ReferenceIDField])
			assert.Equal(t, destination.XpubID, getDestination.XpubID)
		})

		ts.T().Run(testCase.name+" - invalid xpub", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)

			metadata := map[string]interface{}{testMetadataKey: testMetadataValue}
			opts := append(tc.client.DefaultModelOptions(), WithMetadatas(metadata))

			// Create a new destination
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)

			// use the wrong xpub
			var getDestination *Destination
			getDestination, err = tc.client.GetDestinationByLockingScript(
				tc.ctx, testXPubID, destination.LockingScript,
			)
			require.Error(t, err)
			require.Nil(t, getDestination)
		})
	}
}

// TestClient_UpdateDestinationMetadata will test the method UpdateDestinationMetadata()
func (ts *EmbeddedDBTestSuite) TestClient_UpdateDestinationMetadata() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, _, rawKey := CreateNewXPub(tc.ctx, t, tc.client)

			metadata := Metadata{
				"test-key-1": "test-value-1",
				"test-key-2": "test-value-2",
				"test-key-3": "test-value-3",
			}
			opts := tc.client.DefaultModelOptions()
			opts = append(opts, WithMetadatas(metadata))
			destination, err := tc.client.NewDestination(
				tc.ctx, rawKey, utils.ChainExternal, utils.ScriptTypePubKeyHash,
				opts...,
			)
			require.NoError(t, err)
			require.NotNil(t, destination)
			assert.Equal(t, metadata, destination.Metadata)

			destination, err = tc.client.UpdateDestinationMetadataByID(tc.ctx, destination.XpubID, destination.ID, Metadata{"test-key-new": "new-value"})
			require.NoError(t, err)
			assert.Len(t, destination.Metadata, 4)
			assert.Equal(t, "new-value", destination.Metadata["test-key-new"])

			destination, err = tc.client.UpdateDestinationMetadataByAddress(tc.ctx,
				destination.XpubID, destination.Address, Metadata{
					"test-key-new-2": "new-value-2",
					"test-key-1":     nil,
					"test-key-2":     nil,
					"test-key-3":     nil,
				},
			)
			require.NoError(t, err)
			assert.Len(t, destination.Metadata, 2)
			assert.Equal(t, "new-value", destination.Metadata["test-key-new"])
			assert.Equal(t, "new-value-2", destination.Metadata["test-key-new-2"])

			destination, err = tc.client.UpdateDestinationMetadataByLockingScript(tc.ctx,
				destination.XpubID, destination.LockingScript, Metadata{
					"test-key-new-5": "new-value-5",
				},
			)
			require.NoError(t, err)
			assert.Len(t, destination.Metadata, 3)
			assert.Equal(t, "new-value", destination.Metadata["test-key-new"])
			assert.Equal(t, "new-value-2", destination.Metadata["test-key-new-2"])
			assert.Equal(t, "new-value-5", destination.Metadata["test-key-new-5"])

			err = destination.Save(tc.ctx)
			require.NoError(t, err)

			// make sure it was saved
			destination2, err2 := tc.client.GetDestinationByID(tc.ctx, destination.XpubID, destination.ID)
			require.NoError(t, err2)
			assert.Len(t, destination2.Metadata, 3)
			assert.Equal(t, "new-value", destination2.Metadata["test-key-new"])
			assert.Equal(t, "new-value-2", destination2.Metadata["test-key-new-2"])
			assert.Equal(t, "new-value-5", destination2.Metadata["test-key-new-5"])
		})
	}
}
