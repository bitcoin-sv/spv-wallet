package testabilities

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type faucetFixture struct {
	engine  engine.ClientInterface
	user    fixtures.User
	t       testing.TB
	assert  *assert.Assertions
	require *require.Assertions
	arc     ARCFixture
	bhs     BlockHeadersServiceFixture
}

func (f *faucetFixture) TopUp(satoshis bsv.Satoshis) fixtures.GivenTXSpec {
	f.t.Helper()

	txSpec := fixtures.GivenTX(f.t).
		WithSender(fixtures.ExternalFaucet).
		WithInput(uint64(satoshis + 1)).
		WithRecipient(f.user).
		WithP2PKHOutput(uint64(satoshis))

	transaction := &txmodels.NewTransaction{
		ID:       txSpec.ID(),
		TxStatus: txmodels.TxStatusMined,
		Outputs: []txmodels.NewOutput{
			txmodels.NewOutputForP2PKH(
				bsv.Outpoint{TxID: txSpec.ID(), Vout: 0},
				f.user.ID(),
				satoshis,
				bsv.CustomInstructions{
					{
						Type:        "sign",
						Instruction: "P2PKH",
					},
				},
			),
		},
	}

	beefHex, err := txSpec.TX().BEEFHex()
	f.require.NoError(err)
	f.require.NotEmpty(beefHex)

	transaction.SetBEEFHex(beefHex)

	operation := txmodels.NewOperation{
		UserID: f.user.ID(),

		Type:  "incoming",
		Value: int64(satoshis),

		Transaction: transaction,
	}

	err = f.engine.Repositories().Operations.SaveAll(context.Background(), func(yield func(*txmodels.NewOperation) bool) {
		yield(&operation)
	})
	f.assert.NoError(err)

	// Additional check - assertion if the top-up operation was saved correctly
	balance, err := f.engine.Repositories().Users.GetBalance(context.Background(), f.user.ID(), bucket.BSV)
	f.assert.NoError(err)
	f.assert.GreaterOrEqual(balance, satoshis)

	return txSpec
}

func (f *faucetFixture) StoreData(data string) (fixtures.GivenTXSpec, string) {
	f.t.Helper()

	txSpec := fixtures.GivenTX(f.t).
		WithSender(fixtures.ExternalFaucet).
		WithInput(uint64(1000)).
		WithOPReturn(data)

	outpoint := bsv.Outpoint{TxID: txSpec.ID(), Vout: 0}

	operation := txmodels.NewOperation{
		UserID: f.user.ID(),

		Type:  "data",
		Value: 0,

		Transaction: &txmodels.NewTransaction{
			ID:       txSpec.ID(),
			TxStatus: txmodels.TxStatusMined,
			Outputs: []txmodels.NewOutput{
				txmodels.NewOutputForData(
					outpoint,
					f.user.ID(),
					[]byte(data),
				),
			},
		},
	}

	err := f.engine.Repositories().Operations.SaveAll(context.Background(), func(yield func(*txmodels.NewOperation) bool) {
		yield(&operation)
	})
	f.assert.NoError(err)

	return txSpec, outpoint.String()
}
