package engine

import (
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// todo: finish unit tests!

var (
	testLockingScript = "76a9147ff514e6ae3deb46e6644caac5cdd0bf2388906588ac"
	testAddressID     = "fc1e635d98151c6008f29908ee2928c60c745266f9853e945c917b1baa05973e"
	testDestinationID = "c775e7b757ede630cd0aa1113bd102661ab38829ca52a6422ab782862f268646"
)

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
		assert.Equal(t, script.ScriptTypeNonStandard, destination.Type)
		assert.Equal(t, testDestinationID, destination.GetID())
	})
}

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
		assert.Equal(t, script.ScriptTypePubKeyHash, address.Type)
		assert.Equal(t, testAddressID, address.GetID())
	})
}

func TestDestination_GetModelName(t *testing.T) {
	t.Parallel()

	t.Run("model name", func(t *testing.T) {
		address, err := newAddress(testXPub, 0, 0, New())
		require.NotNil(t, address)
		require.NoError(t, err)

		assert.Equal(t, ModelDestination.String(), address.GetModelName())
	})
}

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
