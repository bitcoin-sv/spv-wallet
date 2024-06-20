package engine

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/spverrors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/mrz1836/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testExternalAddress = "1CfaQw9udYNPccssFJFZ94DN8MqNZm9nGt"
	testDraftID         = "z50bb0d4eda0636aae1709e7e7080485a4d00af3ca2962c6e677cf5b53dgab9l"
	testReferenceID     = "example-reference-id"
	testXPriv           = "xprv9s21ZrQH143K3N6qVJQAu4EP51qMcyrKYJLkLgmYXgz58xmVxVLSsbx2DfJUtjcnXK8NdvkHMKfmmg5AJT2nqqRWUrjSHX29qEJwBgBPkJQ"
	testXPub            = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"
	testXPubID          = "1a0b10d4eda0636aae1709e7e7080485a4d99af3ca2962c6e677cf5b53d8ab8c"
)

// TestXpub_newXpub will test the method newXpub()
func TestXpub_newXpub(t *testing.T) {
	t.Parallel()

	t.Run("init xpub", func(t *testing.T) {
		xPub := newXpub(testXPub, New())
		assert.IsType(t, Xpub{}, *xPub)
		assert.Equal(t, testXPubID, xPub.ID)
		assert.Equal(t, testXPubID, xPub.GetID())
		assert.Equal(t, testXPub, xPub.rawXpubKey)
		assert.Equal(t, "xpub", xPub.GetModelName())
	})
}

// TestXpub_getXpub will test the method getXpub()
func TestXpub_getXpub(t *testing.T) {
	t.Run("get xpub - does not exist", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		xPub, err := getXpub(ctx, testXPub, client.DefaultModelOptions()...)
		assert.NoError(t, err)
		assert.Nil(t, xPub)
	})

	t.Run("get xpub", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		xPub := newXpub(testXPub, client.DefaultModelOptions()...)
		err := xPub.Save(ctx)
		assert.NoError(t, err)

		gXPub, gErr := getXpub(ctx, testXPub, client.DefaultModelOptions()...)
		assert.NoError(t, gErr)
		assert.IsType(t, Xpub{}, *gXPub)
	})
}

// TestXpub_GetModelName will test the method GetModelName()
func TestXpub_GetModelName(t *testing.T) {
	t.Parallel()

	xPub := newXpub(testXPub, New())
	assert.Equal(t, ModelXPub.String(), xPub.GetModelName())
}

// TestXpub_GetID will test the method GetID()
func TestXpub_GetID(t *testing.T) {
	t.Parallel()

	xPub := newXpub(testXPub, New())
	assert.Equal(t, testXPubID, xPub.GetID())
}

// TestXpub_getNewDestination will test the method GetNewDestination()
func TestXpub_getNewDestination(t *testing.T) {
	t.Run("err destination", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub("test", client.DefaultModelOptions()...)
		err := xPub.Save(ctx)
		assert.NoError(t, err)

		metaData := map[string]interface{}{
			"test-key": "test-value",
		}
		_, err = xPub.getNewDestination(ctx, utils.ChainInternal, utils.ScriptTypePubKeyHash, append(client.DefaultModelOptions(), WithMetadatas(metaData))...)
		assert.ErrorIs(t, spverrors.ErrXpubInvalidLength, err)
	})

	t.Run("new internal destination", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, client.DefaultModelOptions()...)
		err := xPub.Save(ctx)
		assert.NoError(t, err)

		metaData := map[string]interface{}{
			"test-key": "test-value",
		}
		var destination *Destination
		destination, err = xPub.getNewDestination(ctx, utils.ChainInternal, utils.ScriptTypePubKeyHash, append(client.DefaultModelOptions(), WithMetadatas(metaData))...)
		assert.NoError(t, err)
		assert.Equal(t, "ac18a89055c9269622d9a00ce89047b10aab03cae39feb32cde1be1f1b9bc222", destination.ID)
		assert.Equal(t, xPub.ID, destination.XpubID)
		assert.Equal(t, "76a914296e4f4c6bf609b62b44f2d7c7c4bd5794235ead88ac", destination.LockingScript)
		assert.Equal(t, utils.ScriptTypePubKeyHash, destination.Type)
		assert.Equal(t, utils.ChainInternal, destination.Chain)
		assert.Equal(t, uint32(0), destination.Num)
		assert.Equal(t, "14n4rKed7f5vkPfV7Yj8N3E8Pxa35Rytp9", destination.Address)
		assert.Equal(t, "test-value", destination.Metadata["test-key"])
	})

	t.Run("new external destination", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, client.DefaultModelOptions()...)
		err := xPub.Save(ctx)
		assert.NoError(t, err)

		metaData := map[string]interface{}{
			"test-key": "test-value",
		}
		var destination *Destination
		destination, err = xPub.getNewDestination(ctx, utils.ChainExternal, utils.ScriptTypePubKeyHash, append(client.DefaultModelOptions(), WithMetadatas(metaData))...)
		assert.NoError(t, err)
		assert.Equal(t, "fc1e635d98151c6008f29908ee2928c60c745266f9853e945c917b1baa05973e", destination.ID)
		assert.Equal(t, xPub.ID, destination.XpubID)
		assert.Equal(t, "76a9147ff514e6ae3deb46e6644caac5cdd0bf2388906588ac", destination.LockingScript)
		assert.Equal(t, utils.ScriptTypePubKeyHash, destination.Type)
		assert.Equal(t, utils.ChainExternal, destination.Chain)
		assert.Equal(t, uint32(0), destination.Num)
		assert.Equal(t, "1CfaQw9udYNPccssFJFZ94DN8MqNZm9nGt", destination.Address)
		assert.Equal(t, "test-value", destination.Metadata["test-key"])
	})
}

// TestXpub_childModels will test the method ChildModels()
func TestXpub_childModels(t *testing.T) {
	t.Run("with 1 child model", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, client.DefaultModelOptions()...)
		err := xPub.Save(ctx)
		assert.NoError(t, err)

		_, err = xPub.getNewDestination(ctx, utils.ChainExternal, utils.ScriptTypePubKeyHash, client.DefaultModelOptions()...)
		assert.NoError(t, err)

		childModels := xPub.ChildModels()
		assert.Len(t, childModels, 1)
		assert.IsType(t, &Destination{}, childModels[0])
	})

	t.Run("with 2 child model", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()
		xPub := newXpub(testXPub, client.DefaultModelOptions()...)
		err := xPub.Save(ctx)
		assert.NoError(t, err)

		_, err = xPub.getNewDestination(ctx, utils.ChainExternal, utils.ScriptTypePubKeyHash, client.DefaultModelOptions()...)
		assert.NoError(t, err)
		_, err = xPub.getNewDestination(ctx, utils.ChainExternal, utils.ScriptTypePubKeyHash, client.DefaultModelOptions()...)
		assert.NoError(t, err)

		childModels := xPub.ChildModels()
		assert.Len(t, childModels, 2)
		assert.IsType(t, &Destination{}, childModels[0])
		assert.IsType(t, &Destination{}, childModels[1])
	})
}

// TestXpub_BeforeCreating will test the method BeforeCreating()
func TestXpub_BeforeCreating(t *testing.T) {
	// t.Parallel()

	t.Run("valid xpub", func(t *testing.T) {
		xPub := newXpub(testXPub, New())
		require.NotNil(t, xPub)

		opts := DefaultClientOpts(false, false)
		client, _ := NewClient(context.Background(), opts...)
		xPub.client = client

		err := xPub.BeforeCreating(context.Background())
		require.NoError(t, err)
		require.NotNil(t, xPub)
	})

	t.Run("valid random xpub", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()

		_, xPub, _ := CreateNewXPub(ctx, t, client)

		err := xPub.BeforeCreating(ctx)
		require.NoError(t, err)
		require.NotNil(t, xPub)
	})

	t.Run("incorrect xpub", func(t *testing.T) {
		xPub := newXpub("test", New())
		require.NotNil(t, xPub)

		opts := DefaultClientOpts(false, false)
		client, _ := NewClient(context.Background(), opts...)
		xPub.client = client

		err := xPub.BeforeCreating(context.Background())
		assert.Error(t, err)
		assert.EqualError(t, err, "xpub is an invalid length")
	})
}

// TestXpub_AfterCreated will test the method AfterCreated()
func TestXpub_AfterCreated(t *testing.T) {
	// t.Parallel()

	t.Run("no cache store", func(t *testing.T) {
		xPub := newXpub(testXPub, New())
		require.NotNil(t, xPub)

		opts := DefaultClientOpts(false, false)
		client, _ := NewClient(context.Background(), opts...)
		xPub.client = client

		err := xPub.BeforeCreating(context.Background())
		require.NoError(t, err)
		require.NotNil(t, xPub)

		err = xPub.AfterCreated(context.Background())
		require.NoError(t, err)
	})
}

// TestXpub_AfterUpdated will test the method AfterUpdated()
func TestXpub_AfterUpdated(t *testing.T) {
	// t.Parallel()

	t.Run("no cache store", func(t *testing.T) {
		xPub := newXpub(testXPub, New())
		require.NotNil(t, xPub)

		opts := DefaultClientOpts(false, false)
		client, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		xPub.client = client

		err = xPub.BeforeUpdating(context.Background())
		require.NoError(t, err)
		require.NotNil(t, xPub)

		err = xPub.AfterUpdated(context.Background())
		require.NoError(t, err)
	})
}

// TestXpub_RemovePrivateData will test the method RemovePrivateData()
func TestXpub_RemovePrivateData(t *testing.T) {
	t.Run("remove private data", func(t *testing.T) {
		xPub := newXpub(testXPub, New())
		require.NotNil(t, xPub)

		xPub.Metadata = Metadata{
			"test-key": "test-value",
		}
		xPub.NextInternalNum = uint32(123)
		xPub.NextExternalNum = uint32(321)

		assert.NotNil(t, xPub.Metadata)
		assert.Equal(t, "test-value", xPub.Metadata["test-key"])
		assert.Equal(t, uint32(123), xPub.NextInternalNum)
		assert.Equal(t, uint32(321), xPub.NextExternalNum)

		xPub.RemovePrivateData()
		assert.Nil(t, xPub.Metadata)
		assert.Equal(t, uint32(0), xPub.NextInternalNum)
		assert.Equal(t, uint32(0), xPub.NextExternalNum)
	})
}

// TestXpub_Save will test the method Save()
func (ts *EmbeddedDBTestSuite) TestXpub_Save() {
	for _, testCase := range dbTestCases {
		ts.T().Run(testCase.name+" - valid Save (basic)", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			xPub := newXpub(testXPub, append(tc.client.DefaultModelOptions(), New())...)
			require.NotNil(t, xPub)

			err := xPub.Save(tc.ctx)
			require.NoError(t, err)

			var xPub2 *Xpub
			xPub2, err = tc.client.GetXpub(tc.ctx, testXPub)
			require.NoError(t, err)
			require.NotNil(t, xPub2)

			assert.Equal(t, xPub2.ID, testXPubID)
			require.NoError(t, err)
		})

		ts.T().Run(testCase.name+" - dynamic xPub creation", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			_, xPub, rawKey := CreateNewXPub(tc.ctx, t, tc.client)

			xPub2, err := tc.client.GetXpub(tc.ctx, rawKey)
			require.NoError(t, err)
			require.NotNil(t, xPub2)
			assert.Equal(t, xPub2.ID, xPub.ID)
		})

		ts.T().Run(testCase.name+" - error invalid xPub", func(t *testing.T) {
			tc := ts.genericDBClient(t, testCase.database, false)
			defer tc.Close(tc.ctx)

			xPub := newXpub("bad-key-val", append(tc.client.DefaultModelOptions(), New())...)
			require.NotNil(t, xPub)

			err := xPub.Save(tc.ctx)
			require.Error(t, err)
		})
	}

	ts.T().Run("[sqlite] [redis] [mocking] - create xpub", func(t *testing.T) {
		tc := ts.genericMockedDBClient(t, datastore.SQLite)
		defer tc.Close(tc.ctx)

		xPub := newXpub(testXPub, append(tc.client.DefaultModelOptions(), New())...)
		require.NotNil(t, xPub)

		// Create the expectations
		tc.MockSQLDB.ExpectBegin()

		// Create model
		tc.MockSQLDB.ExpectExec("INSERT INTO `"+tc.tablePrefix+"_"+tableXPubs+"` (`created_at`,`updated_at`,`metadata`,`deleted_at`,`id`,"+
			"`current_balance`,`next_internal_num`,`next_external_num`"+
			") VALUES (?,?,?,?,?,?,?,?)").WithArgs(
			tester.AnyTime{}, // created_at
			tester.AnyTime{}, // updated_at
			nil,              // metadata
			nil,              // deleted_at
			xPub.GetID(),     // id
			0,                // current_balance
			0,                // next_internal_num
			0,                // next_external_num
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit the TX
		tc.MockSQLDB.ExpectCommit()

		// @mrz: this is only testing a SET cmd is fired, not the data being set (that is tested elsewhere)
		setCmd := tc.redisConn.GenericCommand(cache.SetCommand).Expect("ok")

		err := xPub.Save(tc.ctx)
		require.NoError(t, err)

		err = tc.MockSQLDB.ExpectationsWereMet()
		require.NoError(t, err)

		assert.Equal(t, true, setCmd.Called)
	})

	ts.T().Run("[postgresql] [redis] [mocking] - create xpub", func(t *testing.T) {
		tc := ts.genericMockedDBClient(t, datastore.PostgreSQL)
		defer tc.Close(tc.ctx)

		xPub := newXpub(testXPub, append(tc.client.DefaultModelOptions(), New())...)
		require.NotNil(t, xPub)

		// Create the expectations
		tc.MockSQLDB.ExpectBegin()

		// Create model
		tc.MockSQLDB.ExpectExec(`INSERT INTO "`+tc.tablePrefix+`_`+tableXPubs+`" ("created_at","updated_at","metadata","deleted_at","id","current_balance","next_internal_num","next_external_num") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`).WithArgs(
			tester.AnyTime{}, // created_at
			tester.AnyTime{}, // updated_at
			nil,              // metadata
			nil,              // deleted_at
			xPub.GetID(),     // id
			0,                // current_balance
			0,                // next_internal_num
			0,                // next_external_num
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit the TX
		tc.MockSQLDB.ExpectCommit()

		// @mrz: this is only testing a SET cmd is fired, not the data being set (that is tested elsewhere)
		setCmd := tc.redisConn.GenericCommand(cache.SetCommand).Expect("ok")

		err := xPub.Save(tc.ctx)
		require.NoError(t, err)

		err = tc.MockSQLDB.ExpectationsWereMet()
		require.NoError(t, err)

		assert.Equal(t, true, setCmd.Called)
	})
}
