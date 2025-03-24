package engine

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// todo: finish unit tests!

var (
	utxoID       = "e6d250a2dc725ccd237ff8edec0da58537c198960cc2c9f231972464c73ca2ef"
	testDraftID2 = "test-draft-id2"
	testDraftID3 = "test-draft-id3"
)

func createTestUtxos(ctx context.Context, client ClientInterface) error {
	opts := append(client.DefaultModelOptions(), New())

	_utxo := newUtxo(testXPubID, testTxID, testLockingScript, 12, 1225, opts...)
	err := _utxo.Save(ctx)
	if err != nil {
		return err
	}

	_utxo1 := newUtxo(testXPubID, testTxID, testLockingScript, 13, 1225, opts...)
	err = _utxo1.Save(ctx)
	if err != nil {
		return err
	}

	_utxo2 := newUtxo(testXPubID, testTxID, testLockingScript, 14, 1225, opts...)
	err = _utxo2.Save(ctx)
	if err != nil {
		return err
	}

	_utxo3 := newUtxo(testXPubID, testTxID, testLockingScript, 15, 1225, opts...)
	err = _utxo3.Save(ctx)
	if err != nil {
		return err
	}

	_utxo4 := newUtxo(testXPubID, testTxID, testLockingScript, 16, 1225, opts...)
	err = _utxo4.Save(ctx)
	return err
}

func TestUtxo_newUtxo(t *testing.T) {
	t.Parallel()

	t.Run("newUtxo", func(t *testing.T) {
		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 12, 1200, New())
		assert.IsType(t, Utxo{}, *utxo)
		assert.Equal(t, testTxID, utxo.TransactionID)
		assert.Equal(t, testXPubID, utxo.XpubID)
		assert.Equal(t, "", utxo.ID)
		assert.Equal(t, utxoID, utxo.GetID())
		assert.Equal(t, utxoID, utxo.ID)
		assert.Equal(t, uint32(12), utxo.OutputIndex)
		assert.Equal(t, uint64(1200), utxo.Satoshis)
		assert.Equal(t, testLockingScript, utxo.ScriptPubKey)
		assert.Equal(t, "", utxo.Type)
		assert.Equal(t, ModelUtxo.String(), utxo.GetModelName())
	})
}

func TestUtxo_newUtxoFromTxID(t *testing.T) {
	t.Run("newUtxo", func(t *testing.T) {
		utxo := newUtxoFromTxID(testTxID, 12, New())
		assert.IsType(t, Utxo{}, *utxo)
		assert.Equal(t, testTxID, utxo.TransactionID)
		assert.Equal(t, uint32(12), utxo.OutputIndex)
		assert.Equal(t, ModelUtxo.String(), utxo.GetModelName())
	})
}

func TestUtxo_getUtxo(t *testing.T) {
	// t.Parallel()

	t.Run("getUtxo empty", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		utxo, err := getUtxo(ctx, testTxID, 12, client.DefaultModelOptions()...)
		assert.NoError(t, err)
		assert.Nil(t, utxo)
	})

	t.Run("getUtxo", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		_utxo := newUtxo(testXPubID, testTxID, testLockingScript, 12, 1225, append(client.DefaultModelOptions(), New())...)
		_ = _utxo.Save(ctx)

		utxo, err := getUtxo(ctx, testTxID, 12, client.DefaultModelOptions()...)
		assert.NoError(t, err)
		checkUtxoValues(t, utxo, uint32(12), uint64(1225))
	})
}

func TestUtxo_getUtxosByXpubID(t *testing.T) {
	t.Run("getUtxos empty", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		utxos, err := getUtxosByXpubID(
			ctx, testXPubID,
			nil,
			nil,
			nil,
			client.DefaultModelOptions()...,
		)
		assert.NoError(t, err)
		assert.Nil(t, utxos)
	})

	t.Run("getUtxos", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		err := createTestUtxos(ctx, client)
		require.NoError(t, err)

		var utxos []*Utxo
		utxos, err = getUtxosByXpubID(
			ctx, testXPubID,
			nil,
			nil,
			nil,
			client.DefaultModelOptions()...,
		)
		assert.NoError(t, err)
		assert.Len(t, utxos, 5)
	})
}

func TestUtxo_GetModelName(t *testing.T) {
	t.Parallel()

	utxo := newUtxoFromTxID("", 0, New())
	assert.Equal(t, ModelUtxo.String(), utxo.GetModelName())
}

func TestUtxo_UnReserveUtxos(t *testing.T) {
	t.Run("un-reserve 2000", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		err := createTestUtxos(ctx, client)
		require.NoError(t, err)

		var utxos []*Utxo
		utxos, err = reserveUtxos(ctx, testXPubID, testDraftID2, 2000, 0.5, nil, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.Len(t, utxos, 2)
		for _, utxo := range utxos {
			assert.True(t, utxo.DraftID.Valid)
			assert.True(t, utxo.ReservedAt.Valid)
		}

		err = unReserveUtxos(ctx, testXPubID, testDraftID2, client.DefaultModelOptions()...)
		require.NoError(t, err)
		for _, utxo := range utxos {
			var u *Utxo
			u, err = getUtxo(ctx, utxo.TransactionID, utxo.OutputIndex, client.DefaultModelOptions()...)
			require.NoError(t, err)
			assert.Equal(t, utxo.TransactionID, u.TransactionID)
			assert.Equal(t, utxo.OutputIndex, u.OutputIndex)
			assert.False(t, u.DraftID.Valid)
			assert.False(t, u.ReservedAt.Valid)
		}
	})
}

func TestUtxo_ReserveUtxos(t *testing.T) {
	t.Run("reserve 1000", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		err := createTestUtxos(ctx, client)
		require.NoError(t, err)

		var utxos []*Utxo
		utxos, err = reserveUtxos(ctx, testXPubID, testDraftID2, 1000, 0.5, nil, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.Len(t, utxos, 1)
		assert.Equal(t, testDraftID2, utxos[0].DraftID.String)
		assert.True(t, utxos[0].ReservedAt.Valid)
	})

	t.Run("reserve 2000", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		err := createTestUtxos(ctx, client)
		require.NoError(t, err)

		var utxos []*Utxo
		utxos, err = reserveUtxos(ctx, testXPubID, testDraftID2, 2000, 0.5, nil, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.Len(t, utxos, 2)
		assert.Equal(t, testDraftID2, utxos[0].DraftID.String)
		assert.True(t, utxos[0].ReservedAt.Valid)
		assert.Equal(t, testDraftID2, utxos[1].DraftID.String)
		assert.True(t, utxos[1].ReservedAt.Valid)
	})

	t.Run("reserve 20000", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		err := createTestUtxos(ctx, client)
		require.NoError(t, err)

		_, err = reserveUtxos(ctx, testXPubID, testDraftID2, 20000, 0.5, nil, client.DefaultModelOptions()...)
		require.Error(t, err, spverrors.ErrNotEnoughUtxos)
	})

	t.Run("reserve fromUtxos", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		err := createTestUtxos(ctx, client)
		require.NoError(t, err)

		fromUtxos := []*UtxoPointer{{
			TransactionID: testTxID,
			OutputIndex:   16,
		}}

		var utxos []*Utxo
		utxos, err = reserveUtxos(ctx, testXPubID, testDraftID2, 1000, 0.5, fromUtxos, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.Len(t, utxos, 1)
		assert.Equal(t, testDraftID2, utxos[0].DraftID.String)
		assert.True(t, utxos[0].ReservedAt.Valid)
		assert.Equal(t, testTxID, utxos[0].TransactionID)
		assert.Equal(t, uint32(16), utxos[0].OutputIndex)
	})

	t.Run("reserve fromUtxos 2", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		err := createTestUtxos(ctx, client)
		require.NoError(t, err)

		fromUtxos := []*UtxoPointer{{
			TransactionID: testTxID,
			OutputIndex:   15,
		}, {
			TransactionID: testTxID,
			OutputIndex:   16,
		}}

		var utxos []*Utxo
		utxos, err = reserveUtxos(ctx, testXPubID, testDraftID2, 2000, 0.5, fromUtxos, client.DefaultModelOptions()...)
		require.NoError(t, err)
		assert.Len(t, utxos, 2)
		assert.Equal(t, testDraftID2, utxos[0].DraftID.String)
		assert.True(t, utxos[0].ReservedAt.Valid)
		assert.Equal(t, testTxID, utxos[0].TransactionID)
		assert.Equal(t, uint32(15), utxos[0].OutputIndex)
		assert.Equal(t, testDraftID2, utxos[1].DraftID.String)
		assert.True(t, utxos[1].ReservedAt.Valid)
		assert.Equal(t, testTxID, utxos[1].TransactionID)
		assert.Equal(t, uint32(16), utxos[1].OutputIndex)
	})

	t.Run("reserve fromUtxos err", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		err := createTestUtxos(ctx, client)
		require.NoError(t, err)

		fromUtxos := []*UtxoPointer{{
			TransactionID: testTxID,
			OutputIndex:   16,
		}}
		_, err = reserveUtxos(ctx, testXPubID, testDraftID2, 2000, 0.5, fromUtxos, client.DefaultModelOptions()...)
		require.Error(t, err, spverrors.ErrNotEnoughUtxos)
	})

	t.Run("reserve utxos paginated", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		err := createTestUtxos(ctx, client)
		require.NoError(t, err)

		var utxos []*Utxo
		utxos, err = reserveUtxos(ctx, testXPubID, testDraftID2, 4000, 0.5, nil, client.DefaultModelOptions(WithPageSize(2))...)
		require.NoError(t, err)
		assert.Len(t, utxos, 4)
	})

	t.Run("duplicate inputs", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		opts := append(client.DefaultModelOptions(), New())
		utxo := newUtxo(testXPubID, testTxID, testLockingScript, 12, 1225, opts...)
		err := utxo.Save(ctx)
		require.NoError(t, err)

		fromUtxos := []*UtxoPointer{{
			TransactionID: utxo.TransactionID,
			OutputIndex:   utxo.OutputIndex,
		}, {
			TransactionID: utxo.TransactionID,
			OutputIndex:   utxo.OutputIndex,
		}}

		_, err = reserveUtxos(ctx, testXPubID, testDraftID2, 2200, 0.05, fromUtxos, client.DefaultModelOptions()...)
		require.ErrorIs(t, err, spverrors.ErrDuplicateUTXOs)
	})
}

func TestUtxo_GetSpendableUtxos(t *testing.T) {
	t.Run("spendable", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		err := createTestUtxos(ctx, client)
		require.NoError(t, err)

		opts := client.DefaultModelOptions()

		var utxos []*Utxo
		utxos, err = getSpendableUtxos(ctx, testXPubID, utils.ScriptTypePubKeyHash, nil, nil, opts...)
		require.NoError(t, err)
		assert.Len(t, utxos, 5)

		_, err = reserveUtxos(ctx, testXPubID, testDraftID2, 2000, 0.5, nil, opts...)
		require.NoError(t, err)

		utxos, err = getSpendableUtxos(ctx, testXPubID, utils.ScriptTypePubKeyHash, nil, nil, opts...)
		require.NoError(t, err)
		assert.Len(t, utxos, 3)

		_, err = reserveUtxos(ctx, testXPubID, testDraftID3, 1000, 0.5, nil, opts...)
		require.NoError(t, err)

		utxos, err = getSpendableUtxos(ctx, testXPubID, utils.ScriptTypePubKeyHash, nil, nil, opts...)
		require.NoError(t, err)
		assert.Len(t, utxos, 2)

		err = unReserveUtxos(ctx, testXPubID, testDraftID2, opts...)
		require.NoError(t, err)

		utxos, err = getSpendableUtxos(ctx, testXPubID, utils.ScriptTypePubKeyHash, nil, nil, opts...)
		require.NoError(t, err)
		assert.Len(t, utxos, 4)
	})

	t.Run("paginated spendable", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		err := createTestUtxos(ctx, client)
		require.NoError(t, err)

		opts := client.DefaultModelOptions()

		queryParams := &datastore.QueryParams{Page: 1, PageSize: 2}

		var utxos []*Utxo
		utxos, err = getSpendableUtxos(ctx, testXPubID, utils.ScriptTypePubKeyHash, queryParams, nil, opts...)
		require.NoError(t, err)
		assert.Len(t, utxos, 2)

		queryParams = &datastore.QueryParams{Page: 2, PageSize: 2}
		utxos, err = getSpendableUtxos(ctx, testXPubID, utils.ScriptTypePubKeyHash, queryParams, nil, opts...)
		require.NoError(t, err)
		assert.Len(t, utxos, 2)

		queryParams = &datastore.QueryParams{Page: 3, PageSize: 2}
		utxos, err = getSpendableUtxos(ctx, testXPubID, utils.ScriptTypePubKeyHash, queryParams, nil, opts...)
		require.NoError(t, err)
		assert.Len(t, utxos, 1)

		queryParams = &datastore.QueryParams{Page: 4, PageSize: 2}
		utxos, err = getSpendableUtxos(ctx, testXPubID, utils.ScriptTypePubKeyHash, queryParams, nil, opts...)
		require.NoError(t, err)
		assert.Len(t, utxos, 0)
	})
}

func TestUtxo_Save(t *testing.T) {
	// t.Parallel()

	t.Run("Save empty", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		_utxo := newUtxo("", "", "", 0, 0, append(client.DefaultModelOptions(), New())...)
		err := _utxo.Save(ctx)
		assert.ErrorIs(t, err, spverrors.ErrMissingFieldScriptPubKey)
	})

	t.Run("Save", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		index := uint32(12)
		satoshis := uint64(1225)
		utxo := newUtxo(testXPubID, testTxID, testLockingScript, index, satoshis, append(client.DefaultModelOptions(), New())...)
		err := utxo.Save(ctx)
		assert.NoError(t, err)
		checkUtxoValues(t, utxo, index, satoshis)
	})
}

func checkUtxoValues(t *testing.T, utxo *Utxo, index uint32, satoshis uint64) {
	assert.Equal(t, testTxID, utxo.TransactionID)
	assert.Equal(t, testXPubID, utxo.XpubID)
	assert.Equal(t, utxoID, utxo.ID)
	assert.Equal(t, index, utxo.OutputIndex)
	assert.Equal(t, satoshis, utxo.Satoshis)
	assert.Equal(t, testLockingScript, utxo.ScriptPubKey)
	assert.Equal(t, utils.ScriptTypePubKeyHash, utxo.Type)
	assert.Equal(t, ModelUtxo.String(), utxo.GetModelName())
}
