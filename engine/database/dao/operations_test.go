package dao_test

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/database/dao"
	"github.com/bitcoin-sv/spv-wallet/engine/database/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/testabilities/testmode"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"slices"
	"testing"
)

func TestOperations(t *testing.T) {
	testmode.DevelopmentOnly_SetPostgresModeWithName(t, "spv-test")
	// given:
	givenDB, cleanup := testabilities.Given(t, testengine.WithNewTransactionFlowEnabled())
	defer cleanup()

	// and:
	operationsDAO := dao.NewOperationsAccessObject(givenDB.GormDB())

	// and:
	senderEntity := getUser(t, givenDB.GormDB(), fixtures.Sender)
	recipientEntity := getUser(t, givenDB.GormDB(), fixtures.RecipientInternal)

	var testState struct {
		incomingTx fixtures.GivenTXSpec
		internalTx fixtures.GivenTXSpec
	}

	t.Run("External incoming tx", func(t *testing.T) {
		// given:
		txSpec := fixtures.GivenTX(t).
			WithInput(1000).
			WithOutputScript(999, fixtures.Sender.P2PKHLockingScript())

		// and:
		transaction := &database.TrackedTransaction{
			ID:       txSpec.ID(),
			TxStatus: database.TxStatusCreated,

			Outputs: []*database.Output{
				{
					Vout: 0,
				},
			},
		}

		// and:
		operations := []*database.Operation{
			{
				User:        senderEntity,
				Transaction: transaction,
				Type:        "paymail",
				Value:       int64(1000),
			},
		}

		// when:
		err := operationsDAO.SaveOperation(context.Background(), slices.Values(operations))

		// then:
		require.NoError(t, err)

		// update:
		testState.incomingTx = txSpec
	})

	t.Run("Internal transaction from Sender to InternalRecipient", func(t *testing.T) {
		// given:
		txSpec := fixtures.GivenTX(t).
			WithInputFromUTXO(testState.incomingTx.TX(), 0).
			WithOutputScript(999, fixtures.RecipientInternal.P2PKHLockingScript())

		// and:
		transaction := &database.TrackedTransaction{
			ID:       txSpec.ID(),
			TxStatus: database.TxStatusCreated,

			Inputs: []*database.Output{
				{
					TxID:       testState.incomingTx.ID(),
					Vout:       0,
					SpendingTX: txSpec.ID(),
				},
			},

			Outputs: []*database.Output{
				{
					Vout: 0,
				},
			},
		}

		// and:
		operations := []*database.Operation{
			{
				User:        senderEntity,
				Transaction: transaction,
				Type:        "paymail",
				Value:       int64(-1000),
			},
			{
				User:        recipientEntity,
				Transaction: transaction,
				Type:        "paymail",
				Value:       int64(999),
			},
		}

		// when:
		err := operationsDAO.SaveOperation(context.Background(), slices.Values(operations))

		// then:
		require.NoError(t, err)

		// update:
		testState.internalTx = txSpec
	})

	t.Run("External outgoing tx", func(t *testing.T) {
		// given:
		txSpec := fixtures.GivenTX(t).
			WithInputFromUTXO(testState.internalTx.TX(), 0).
			WithP2PKHOutput(998)

		// and:
		transaction := &database.TrackedTransaction{
			ID:       txSpec.ID(),
			TxStatus: database.TxStatusCreated,

			Inputs: []*database.Output{
				{
					TxID:       testState.internalTx.ID(),
					Vout:       0,
					SpendingTX: txSpec.ID(),
				},
			},
		}

		// and:
		operations := []*database.Operation{
			{
				User:        recipientEntity,
				Transaction: transaction,
				Type:        "paymail",
				Value:       int64(-999),
			},
		}

		// when:
		err := operationsDAO.SaveOperation(context.Background(), slices.Values(operations))

		// then:
		require.NoError(t, err)
	})
}

func getUser(t *testing.T, db *gorm.DB, user fixtures.User) *database.User {
	t.Helper()

	var dbUser database.User
	err := db.Where("id = ?", user.Address().AddressString).First(&dbUser).Error
	assert.NoError(t, err)

	return &dbUser
}
