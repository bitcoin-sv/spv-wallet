package engine

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	bscript2 "github.com/libsv/go-bt/v2/bscript"
	"github.com/mrz1836/go-cache"
	"github.com/mrz1836/go-datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// todo: finish unit tests!

var (
	testLockingScript = "76a9147ff514e6ae3deb46e6644caac5cdd0bf2388906588ac"
	testAddressID     = "fc1e635d98151c6008f29908ee2928c60c745266f9853e945c917b1baa05973e"
	testDestinationID = "c775e7b757ede630cd0aa1113bd102661ab38829ca52a6422ab782862f268646"
	stasHex           = "76a9146d3562a8ec96bcb3b2253fd34f38a556fb66733d88ac6976aa607f5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7c5f7f7c5e7f7c5d7f7c5c7f7c5b7f7c5a7f7c597f7c587f7c577f7c567f7c557f7c547f7c537f7c527f7c517f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01007e818b21414136d08c5ed2bf3ba048afe6dcaebafeffffffffffffffffffffffffffffff007d976e7c5296a06394677768827601249301307c7e23022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798027e7c7e7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c8276638c687f7c7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e7e01417e21038ff83d8cf12121491609c4939dc11c4aa35503508fe432dc5a5c1905608b9218ad547f7701207f01207f7701247f517f7801007e8102fd00a063546752687f7801007e817f727e7b01177f777b557a766471567a577a786354807e7e676d68aa880067765158a569765187645294567a5379587a7e7e78637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6878637c8c7c53797e577a7e6867567a6876aa587a7d54807e577a597a5a7a786354807e6f7e7eaa727c7e676d6e7eaa7c687b7eaa587a7d877663516752687c72879b69537a647500687c7b547f77517f7853a0916901247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77788c6301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f777852946301247f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e816854937f77686877517f7c52797d8b9f7c53a09b91697c76638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6876638c7c587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f777c6863587f77517f7c01007e817602fc00a06302fd00a063546752687f7c01007e81687f7768587f517f7801007e817602fc00a06302fd00a063546752687f7801007e81727e7b7b687f75537f7c0376a9148801147f775379645579887567726881766968789263556753687a76026c057f7701147f8263517f7c766301007e817f7c6775006877686b537992635379528763547a6b547a6b677c6b567a6b537a7c717c71716868547a587f7c81547a557964936755795187637c686b687c547f7701207f75748c7a7669765880748c7a76567a876457790376a9147e7c7e557967041976a9147c7e0288ac687e7e5579636c766976748c7a9d58807e6c0376a9147e748c7a7e6c7e7e676c766b8263828c007c80517e846864745aa0637c748c7a76697d937b7b58807e56790376a9147e748c7a7e55797e7e6868686c567a5187637500678263828c007c80517e846868647459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e687459a0637c748c7a76697d937b7b58807e55790376a9147e748c7a7e55797e7e68687c537a9d547963557958807e041976a91455797e0288ac7e7e68aa87726d77776a14f566909f378788e61108d619e40df2757455d14c010005546f6b656e"
)

// TestDestination_newDestination will test the method newDestination()
func TestDestination_newDestination(t *testing.T) {
	t.Parallel()

	t.Run("New empty destination model", func(t *testing.T) {
		destination := newDestination("", "", New())
		require.NotNil(t, destination)
		assert.IsType(t, Destination{}, *destination)
		assert.Equal(t, ModelDestination.String(), destination.GetModelName())
		assert.Equal(t, true, destination.IsNew())
		assert.Equal(t, "", destination.LockingScript)
		assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", destination.GetID())
	})

	t.Run("New destination model", func(t *testing.T) {
		testScript := "1234567890"
		xPubID := "xpub123456789"
		destination := newDestination(xPubID, testScript, New())
		require.NotNil(t, destination)
		assert.IsType(t, Destination{}, *destination)
		assert.Equal(t, ModelDestination.String(), destination.GetModelName())
		assert.Equal(t, true, destination.IsNew())
		assert.Equal(t, testScript, destination.LockingScript)
		assert.Equal(t, xPubID, destination.XpubID)
		assert.Equal(t, bscript2.ScriptTypeNonStandard, destination.Type)
		assert.Equal(t, testDestinationID, destination.GetID())
	})
}

// TestDestination_newAddress will test the method newAddress()
func TestDestination_newAddress(t *testing.T) {
	t.Parallel()

	t.Run("New empty address model", func(t *testing.T) {
		address, err := newAddress("", 0, 0, New())
		assert.Nil(t, address)
		assert.Error(t, err)
	})

	t.Run("invalid xPub", func(t *testing.T) {
		address, err := newAddress("test", 0, 0, New())
		assert.Nil(t, address)
		assert.Error(t, err)
	})

	t.Run("valid xPub", func(t *testing.T) {
		address, err := newAddress(testXPub, 0, 0, New())
		require.NotNil(t, address)
		require.NoError(t, err)

		// Default values
		assert.IsType(t, Destination{}, *address)
		assert.Equal(t, ModelDestination.String(), address.GetModelName())
		assert.Equal(t, true, address.IsNew())

		// Check set address
		assert.Equal(t, testXPubID, address.XpubID)
		assert.Equal(t, testExternalAddress, address.Address)

		// Check set locking script
		assert.Equal(t, testLockingScript, address.LockingScript)
		assert.Equal(t, bscript2.ScriptTypePubKeyHash, address.Type)
		assert.Equal(t, testAddressID, address.GetID())
	})
}

// TestDestination_GetModelName will test the method GetModelName()
func TestDestination_GetModelName(t *testing.T) {
	t.Parallel()

	t.Run("model name", func(t *testing.T) {
		address, err := newAddress(testXPub, 0, 0, New())
		require.NotNil(t, address)
		require.NoError(t, err)

		assert.Equal(t, ModelDestination.String(), address.GetModelName())
	})
}

// TestDestination_GetID will test the method GetID()
func TestDestination_GetID(t *testing.T) {
	t.Parallel()

	t.Run("valid id - address", func(t *testing.T) {
		address, err := newAddress(testXPub, 0, 0, New())
		require.NotNil(t, address)
		require.NoError(t, err)

		assert.Equal(t, testAddressID, address.GetID())
	})

	t.Run("valid id - destination", func(t *testing.T) {
		testScript := "1234567890"
		xPubID := "xpub123456789"
		destination := newDestination(xPubID, testScript, New())
		require.NotNil(t, destination)

		assert.Equal(t, testDestinationID, destination.GetID())
	})
}

// TestDestination_setAddress will test the method setAddress()
func TestDestination_setAddress(t *testing.T) {
	t.Run("internal 1", func(t *testing.T) {
		destination := newDestination(testXPubID, testLockingScript)
		destination.Chain = utils.ChainInternal
		destination.Num = 1
		err := destination.setAddress(testXPub)
		require.NoError(t, err)
		assert.Equal(t, "1PQW54xMn5KA6uK7wgfzN4y7ZXMi6o7Qtm", destination.Address)
	})

	t.Run("external 1", func(t *testing.T) {
		destination := newDestination(testXPubID, testLockingScript)
		destination.Chain = utils.ChainExternal
		destination.Num = 1
		err := destination.setAddress(testXPub)
		require.NoError(t, err)
		assert.Equal(t, "16fq7PmmXXbFUG5maT5Xvr2zDBUgN1xdMF", destination.Address)
	})

	t.Run("internal 2", func(t *testing.T) {
		destination := newDestination(testXPubID, testLockingScript)
		destination.Chain = utils.ChainInternal
		destination.Num = 2
		err := destination.setAddress(testXPub)
		require.NoError(t, err)
		assert.Equal(t, "13St2SHkw1b8ZuaExyMf6ZMEzNjYbWRqL4", destination.Address)
	})

	t.Run("external 2", func(t *testing.T) {
		destination := newDestination(testXPubID, testLockingScript)
		destination.Chain = utils.ChainExternal
		destination.Num = 2
		err := destination.setAddress(testXPub)
		require.NoError(t, err)
		assert.Equal(t, "19jswATg9vBFta1aRnEjPHa2KMwafkmANj", destination.Address)
	})
}

// TestDestination_getDestinationByID will test the method getDestinationByID()
func TestDestination_getDestinationByID(t *testing.T) {
	t.Run("does not exist", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		xPub, err := getDestinationByID(ctx, testDestinationID, client.DefaultModelOptions()...)
		assert.NoError(t, err)
		assert.Nil(t, xPub)
	})

	t.Run("get", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		destination := newDestination(testXPubID, testLockingScript, client.DefaultModelOptions()...)
		err := destination.Save(ctx)
		assert.NoError(t, err)

		gDestination, gErr := getDestinationByID(ctx, destination.ID, client.DefaultModelOptions()...)
		assert.NoError(t, gErr)
		assert.IsType(t, Destination{}, *gDestination)
		assert.Equal(t, testXPubID, gDestination.XpubID)
		assert.Equal(t, testLockingScript, gDestination.LockingScript)
		assert.Equal(t, testExternalAddress, gDestination.Address)
		assert.Equal(t, utils.ScriptTypePubKeyHash, gDestination.Type)
		assert.Equal(t, uint32(0), gDestination.Chain)
		assert.Equal(t, uint32(0), gDestination.Num)
	})
}

// TestDestination_getDestinationByAddress will test the method getDestinationByAddress()
func TestDestination_getDestinationByAddress(t *testing.T) {
	t.Run("does not exist", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		xPub, err := getDestinationByAddress(ctx, testExternalAddress, client.DefaultModelOptions()...)
		assert.NoError(t, err)
		assert.Nil(t, xPub)
	})

	t.Run("get", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		destination := newDestination(testXPubID, testLockingScript, client.DefaultModelOptions()...)
		err := destination.Save(ctx)
		assert.NoError(t, err)

		gDestination, gErr := getDestinationByAddress(ctx, testExternalAddress, client.DefaultModelOptions()...)
		assert.NoError(t, gErr)
		assert.IsType(t, Destination{}, *gDestination)
		assert.Equal(t, testXPubID, gDestination.XpubID)
		assert.Equal(t, testLockingScript, gDestination.LockingScript)
		assert.Equal(t, testExternalAddress, gDestination.Address)
		assert.Equal(t, utils.ScriptTypePubKeyHash, gDestination.Type)
		assert.Equal(t, uint32(0), gDestination.Chain)
		assert.Equal(t, uint32(0), gDestination.Num)
	})
}

// TestDestination_getDestinationByLockingScript will test the method getDestinationByLockingScript()
func TestDestination_getDestinationByLockingScript(t *testing.T) {
	t.Run("does not exist", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		xPub, err := getDestinationByLockingScript(ctx, testLockingScript, client.DefaultModelOptions()...)
		assert.NoError(t, err)
		assert.Nil(t, xPub)
	})

	t.Run("get destination", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer deferMe()

		destination := newDestination(testXPubID, testLockingScript, client.DefaultModelOptions()...)
		err := destination.Save(ctx)
		assert.NoError(t, err)

		gDestination, gErr := getDestinationByLockingScript(ctx, testLockingScript, client.DefaultModelOptions()...)
		assert.NoError(t, gErr)
		assert.IsType(t, Destination{}, *gDestination)
		assert.Equal(t, testXPubID, gDestination.XpubID)
		assert.Equal(t, testLockingScript, gDestination.LockingScript)
		assert.Equal(t, testExternalAddress, gDestination.Address)
		assert.Equal(t, utils.ScriptTypePubKeyHash, gDestination.Type)
		assert.Equal(t, uint32(0), gDestination.Chain)
		assert.Equal(t, uint32(0), gDestination.Num)
	})
}

// BenchmarkDestination_newAddress will test the method newAddress()
func BenchmarkDestination_newAddress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = newAddress(testXPub, 0, 0, New())
	}
}

// TestClient_NewDestination will test the method NewDestination()
func TestClient_NewDestination(t *testing.T) {
	t.Run("valid - simple destination", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()

		// Get new random key
		_, xPub, rawXPub := CreateNewXPub(ctx, t, client)
		require.NotNil(t, xPub)

		opts := append(
			client.DefaultModelOptions(),
			WithMetadatas(map[string]interface{}{
				ReferenceIDField: "some-reference-id",
				testMetadataKey:  testMetadataValue,
			}),
		)

		// Create a new destination
		destination, err := client.NewDestination(
			ctx, rawXPub, utils.ChainExternal, utils.ScriptTypePubKeyHash, opts...,
		)
		require.NoError(t, err)
		require.NotNil(t, destination)
		assert.Equal(t, "some-reference-id", destination.Metadata[ReferenceIDField])
		assert.Equal(t, 64, len(destination.ID))
		assert.Greater(t, len(destination.Address), 32)
		assert.Greater(t, len(destination.LockingScript), 32)
		assert.Equal(t, utils.ScriptTypePubKeyHash, destination.Type)
		assert.Equal(t, uint32(0), destination.Num)
		assert.Equal(t, utils.ChainExternal, destination.Chain)
		assert.Equal(t, Metadata{ReferenceIDField: "some-reference-id", testMetadataKey: testMetadataValue}, destination.Metadata)
		assert.Equal(t, xPub.ID, destination.XpubID)
	})

	t.Run("error - invalid xPub", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()

		opts := append(
			client.DefaultModelOptions(),
			WithMetadatas(map[string]interface{}{
				testMetadataKey: testMetadataValue,
			}),
		)

		// Create a new destination
		destination, err := client.NewDestination(
			ctx, "bad-value", utils.ChainExternal, utils.ScriptTypePubKeyHash,
			opts...,
		)
		require.Error(t, err)
		require.Nil(t, destination)
	})

	t.Run("error - xPub not found", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()

		opts := append(
			client.DefaultModelOptions(),
			WithMetadatas(map[string]interface{}{
				testMetadataKey: testMetadataValue,
			}),
		)

		// Create a new destination
		destination, err := client.NewDestination(
			ctx, testXPub, utils.ChainExternal, utils.ScriptTypePubKeyHash,
			opts...,
		)
		require.Error(t, err)
		require.Nil(t, destination)
		assert.ErrorIs(t, err, ErrMissingXpub)
	})

	t.Run("error - unsupported destination type", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()

		// Get new random key
		_, xPub, rawXPub := CreateNewXPub(ctx, t, client)
		require.NotNil(t, xPub)

		opts := append(
			client.DefaultModelOptions(),
			WithMetadatas(map[string]interface{}{
				testMetadataKey: testMetadataValue,
			}),
		)

		// Create a new destination
		destination, err := client.NewDestination(
			ctx, rawXPub, utils.ChainExternal, utils.ScriptTypeMultiSig,
			opts...,
		)
		require.Error(t, err)
		require.Nil(t, destination)
	})

	t.Run("stas token", func(t *testing.T) {
		ctx, client, deferMe := CreateTestSQLiteClient(t, false, true, withTaskManagerMockup())
		defer deferMe()

		// Get new random key
		_, xPub, rawXPub := CreateNewXPub(ctx, t, client)
		require.NotNil(t, xPub)

		opts := append(
			client.DefaultModelOptions(),
			WithMetadatas(map[string]interface{}{
				testMetadataKey: testMetadataValue,
			}),
		)

		// Create a new destination
		destination, err := client.NewDestinationForLockingScript(
			ctx, utils.Hash(rawXPub), stasHex,
			opts...,
		)
		require.NoError(t, err)
		require.Equal(t, utils.Hash(stasHex), destination.ID)
		require.Equal(t, utils.Hash(rawXPub), destination.XpubID)
		require.Equal(t, stasHex, destination.LockingScript)
		require.Equal(t, "1AxScC72W9tyk1Enej6dBsVZNkkgAonk4H", destination.Address)
		require.Equal(t, utils.ScriptTypeTokenStas, destination.Type)
	})
}

// TestDestination_Save will test the method Save()
func (ts *EmbeddedDBTestSuite) TestDestination_Save() {
	ts.T().Run("[sqlite] [redis] [mocking] - create destination", func(t *testing.T) {
		tc := ts.genericMockedDBClient(t, datastore.SQLite)
		defer tc.Close(tc.ctx)

		xPub := newXpub(testXPub, append(tc.client.DefaultModelOptions(), New())...)
		require.NotNil(t, xPub)

		destination := newDestination(xPub.ID, testLockingScript, append(tc.client.DefaultModelOptions(), New())...)
		require.NotNil(t, destination)
		destination.DraftID = testDraftID

		// Create the expectations
		tc.MockSQLDB.ExpectBegin()

		// Create model
		tc.MockSQLDB.ExpectExec("INSERT INTO `"+tc.tablePrefix+"_destinations` ("+
			"`created_at`,`updated_at`,`metadata`,`deleted_at`,`id`,`xpub_id`,`locking_script`,"+
			"`type`,`chain`,`num`,`address`,`draft_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)").WithArgs(
			tester.AnyTime{},    // created_at
			tester.AnyTime{},    // updated_at
			nil,                 // metadata
			nil,                 // deleted_at
			tester.AnyGUID{},    // id
			xPub.GetID(),        // xpub_id
			testLockingScript,   // locking_script
			destination.Type,    // type
			0,                   // chain
			0,                   // num
			destination.Address, // address
			testDraftID,         // draft_id
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit the TX
		tc.MockSQLDB.ExpectCommit()

		// @mrz: this is only testing a SET cmd is fired, not the data being set (that is tested elsewhere)
		setCmd := tc.redisConn.GenericCommand(cache.SetCommand).Expect("ok")

		err := destination.Save(tc.ctx)
		require.NoError(t, err)

		err = tc.MockSQLDB.ExpectationsWereMet()
		require.NoError(t, err)
		assert.Equal(t, true, setCmd.Called)
	})

	ts.T().Run("[mysql] [redis] [mocking] - create destination", func(t *testing.T) {
		tc := ts.genericMockedDBClient(t, datastore.MySQL)
		defer tc.Close(tc.ctx)

		xPub := newXpub(testXPub, append(tc.client.DefaultModelOptions(), New())...)
		require.NotNil(t, xPub)

		destination := newDestination(xPub.ID, testLockingScript, append(tc.client.DefaultModelOptions(), New())...)
		require.NotNil(t, destination)
		destination.DraftID = testDraftID

		// Create the expectations
		tc.MockSQLDB.ExpectBegin()

		// Create model
		tc.MockSQLDB.ExpectExec("INSERT INTO `"+tc.tablePrefix+"_destinations` ("+
			"`created_at`,`updated_at`,`metadata`,`deleted_at`,`id`,`xpub_id`,`locking_script`,"+
			"`type`,`chain`,`num`,`address`,`draft_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)").WithArgs(
			tester.AnyTime{},    // created_at
			tester.AnyTime{},    // updated_at
			nil,                 // metadata
			nil,                 // deleted_at
			tester.AnyGUID{},    // id
			xPub.GetID(),        // xpub_id
			testLockingScript,   // locking_script
			destination.Type,    // type
			0,                   // chain
			0,                   // num
			destination.Address, // address
			testDraftID,         // draft_id
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit the TX
		tc.MockSQLDB.ExpectCommit()

		// @mrz: this is only testing a SET cmd is fired, not the data being set (that is tested elsewhere)
		setCmd := tc.redisConn.GenericCommand(cache.SetCommand).Expect("ok")

		err := destination.Save(tc.ctx)
		require.NoError(t, err)

		err = tc.MockSQLDB.ExpectationsWereMet()
		require.NoError(t, err)
		assert.Equal(t, true, setCmd.Called)
	})

	ts.T().Run("[postgresql] [redis] [mocking] - create destination", func(t *testing.T) {
		tc := ts.genericMockedDBClient(t, datastore.PostgreSQL)
		defer tc.Close(tc.ctx)

		xPub := newXpub(testXPub, append(tc.client.DefaultModelOptions(), New())...)
		require.NotNil(t, xPub)

		destination := newDestination(xPub.ID, testLockingScript, append(tc.client.DefaultModelOptions(), New())...)
		require.NotNil(t, destination)
		destination.DraftID = testDraftID

		// Create the expectations
		tc.MockSQLDB.ExpectBegin()

		// Create model
		tc.MockSQLDB.ExpectExec(`INSERT INTO "`+tc.tablePrefix+`_destinations" ("created_at","updated_at","metadata","deleted_at","id","xpub_id","locking_script","type","chain","num","address","draft_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`).WithArgs(
			tester.AnyTime{},    // created_at
			tester.AnyTime{},    // updated_at
			nil,                 // metadata
			nil,                 // deleted_at
			tester.AnyGUID{},    // id
			xPub.GetID(),        // xpub_id
			testLockingScript,   // locking_script
			destination.Type,    // type
			0,                   // chain
			0,                   // num
			destination.Address, // address
			testDraftID,         // draft_id
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// Commit the TX
		tc.MockSQLDB.ExpectCommit()

		// @mrz: this is only testing a SET cmd is fired, not the data being set (that is tested elsewhere)
		setCmd := tc.redisConn.GenericCommand(cache.SetCommand).Expect("ok")

		err := destination.Save(tc.ctx)
		require.NoError(t, err)

		err = tc.MockSQLDB.ExpectationsWereMet()
		require.NoError(t, err)
		assert.Equal(t, true, setCmd.Called)
	})

	ts.T().Run("[mongo] [redis] [mocking] - create destination", func(t *testing.T) {
		// todo: mocking for MongoDB
	})
}
