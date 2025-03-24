package engine

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (ts *EmbeddedDBTestSuite) TestClient_NewXpub() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, tc.client.DefaultModelOptions()...)
			require.NoError(t, err)
			assert.Equal(t, testXPubID, xPub.ID)

			xPub2, err2 := tc.client.GetXpub(tc.ctx, testXPub)
			require.NoError(t, err2)
			assert.Equal(t, testXPubID, xPub2.ID)
		})

		ts.T().Run(testCase.name+" - valid with metadata (key->val)", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			opts := append(tc.client.DefaultModelOptions(), WithMetadata(testMetadataKey, testMetadataValue))

			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, opts...)
			require.NoError(t, err)
			assert.Equal(t, testXPubID, xPub.ID)
			assert.Equal(t, Metadata{testMetadataKey: testMetadataValue}, xPub.Metadata)

			xPub2, err2 := tc.client.GetXpub(tc.ctx, testXPub)
			require.NoError(t, err2)
			assert.Equal(t, testXPubID, xPub2.ID)
			assert.Equal(t, Metadata{testMetadataKey: testMetadataValue}, xPub2.Metadata)
		})

		ts.T().Run(testCase.name+" - valid with metadatas", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			opts := append(
				tc.client.DefaultModelOptions(),
				WithMetadatas(map[string]interface{}{
					testMetadataKey: testMetadataValue,
				}),
			)

			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, opts...)
			require.NoError(t, err)
			assert.Equal(t, testXPubID, xPub.ID)
			assert.Equal(t, Metadata{testMetadataKey: testMetadataValue}, xPub.Metadata)

			xPub2, err2 := tc.client.GetXpub(tc.ctx, testXPub)
			require.NoError(t, err2)
			assert.Equal(t, testXPubID, xPub2.ID)
			assert.Equal(t, Metadata{testMetadataKey: testMetadataValue}, xPub2.Metadata)
		})

		ts.T().Run(testCase.name+" - errors", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, err := tc.client.NewXpub(tc.ctx, "test", tc.client.DefaultModelOptions()...)
			assert.ErrorIs(t, err, spverrors.ErrXpubInvalidLength)

			_, err = tc.client.NewXpub(tc.ctx, "", tc.client.DefaultModelOptions()...)
			assert.ErrorIs(t, err, spverrors.ErrMissingFieldXpubID)
		})

		ts.T().Run(testCase.name+" - duplicate xPub", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, tc.client.DefaultModelOptions()...)
			require.NoError(t, err)
			assert.Equal(t, testXPubID, xPub.ID)

			_, err2 := tc.client.NewXpub(tc.ctx, testXPub, tc.client.DefaultModelOptions()...)
			require.Error(t, err2)
		})
	}
}

func (ts *EmbeddedDBTestSuite) TestClient_GetXpub() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, tc.client.DefaultModelOptions()...)
			require.NoError(t, err)
			assert.Equal(t, testXPubID, xPub.ID)

			xPub2, err2 := tc.client.GetXpub(tc.ctx, testXPub)
			require.NoError(t, err2)
			assert.Equal(t, testXPubID, xPub2.ID)

			xPub3, err3 := tc.client.GetXpubByID(tc.ctx, xPub2.ID)
			require.NoError(t, err3)
			assert.Equal(t, testXPubID, xPub3.ID)
		})

		ts.T().Run(testCase.name+" - error - invalid xpub", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			xPub, err := tc.client.GetXpub(tc.ctx, "test")
			require.Error(t, err)
			require.Nil(t, xPub)
			assert.ErrorIs(t, err, spverrors.ErrCouldNotFindXpub)
		})

		ts.T().Run(testCase.name+" - error - missing xpub", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			xPub, err := tc.client.GetXpub(tc.ctx, testXPub)
			require.Error(t, err)
			require.Nil(t, xPub)
			assert.ErrorIs(t, err, spverrors.ErrCouldNotFindXpub)
		})
	}
}

func (ts *EmbeddedDBTestSuite) TestClient_GetXpubByID() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, tc.client.DefaultModelOptions()...)
			require.NoError(t, err)
			assert.Equal(t, testXPubID, xPub.ID)

			xPub2, err2 := tc.client.GetXpubByID(tc.ctx, xPub.ID)
			require.NoError(t, err2)
			assert.Equal(t, testXPubID, xPub2.ID)
		})

		ts.T().Run(testCase.name+" - error - invalid xpub", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			xPub, err := tc.client.GetXpubByID(tc.ctx, "test")
			require.Error(t, err)
			require.Nil(t, xPub)
		})

		ts.T().Run(testCase.name+" - error - missing xpub", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			xPub, err := tc.client.GetXpubByID(tc.ctx, testXPub)
			require.Error(t, err)
			require.Nil(t, xPub)
			assert.ErrorIs(t, err, spverrors.ErrCouldNotFindXpub)
		})
	}
}

func (ts *EmbeddedDBTestSuite) TestClient_UpdateXpubMetadata() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - valid", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			metadata := Metadata{
				"test-key-1": "test-value-1",
				"test-key-2": "test-value-2",
				"test-key-3": "test-value-3",
			}
			opts := tc.client.DefaultModelOptions()
			opts = append(opts, WithMetadatas(metadata))

			xPub, err := tc.client.NewXpub(tc.ctx, testXPub, opts...)
			require.NoError(t, err)
			assert.Equal(t, testXPubID, xPub.ID)
			assert.Equal(t, metadata, xPub.Metadata)

			xPub, err = tc.client.UpdateXpubMetadata(tc.ctx, xPub.ID, Metadata{"test-key-new": "new-value"})
			require.NoError(t, err)
			assert.Len(t, xPub.Metadata, 1)
			assert.Equal(t, "new-value", xPub.Metadata["test-key-new"])

			xPub, err = tc.client.UpdateXpubMetadata(tc.ctx, xPub.ID, Metadata{
				"test-key-new-2": "new-value-2",
				"test-key-1":     nil,
				"test-key-2":     nil,
				"test-key-3":     nil,
			})
			require.NoError(t, err)
			assert.Len(t, xPub.Metadata, 4)
			assert.Equal(t, nil, xPub.Metadata["test-key-1"])
			assert.Equal(t, nil, xPub.Metadata["test-key-2"])
			assert.Equal(t, nil, xPub.Metadata["test-key-3"])
			assert.Equal(t, "new-value-2", xPub.Metadata["test-key-new-2"])

			err = xPub.Save(tc.ctx)
			require.NoError(t, err)

			// make sure it was saved
			xPub2, err2 := tc.client.GetXpubByID(tc.ctx, xPub.ID)
			require.NoError(t, err2)
			assert.Len(t, xPub2.Metadata, 4)
			assert.Equal(t, nil, xPub.Metadata["test-key-1"])
			assert.Equal(t, nil, xPub.Metadata["test-key-2"])
			assert.Equal(t, nil, xPub.Metadata["test-key-3"])
			assert.Equal(t, "new-value-2", xPub2.Metadata["test-key-new-2"])
		})
	}
}
