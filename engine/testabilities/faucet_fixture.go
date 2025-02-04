package testabilities

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/assert"
)

type faucetFixture struct {
	engine engine.ClientInterface
	user   fixtures.User
	t      testing.TB
	assert *assert.Assertions
	arc    ARCFixture
	bhs    BlockHeadersServiceFixture
}

func (f *faucetFixture) TopUp(satoshis bsv.Satoshis) fixtures.GivenTXSpec {
	f.t.Helper()

	txSpec := fixtures.GivenTX(f.t).
		WithSender(fixtures.ExternalFaucet).
		WithInput(uint64(satoshis + 1)).
		WithRecipient(f.user).
		WithP2PKHOutput(uint64(satoshis))

	operation := txmodels.NewOperation{
		UserID: f.user.ID(),

		Type:  "incoming",
		Value: int64(satoshis), //nolint:gosec // This is a test fixture, values won't exceed int64

		Transaction: &txmodels.NewTransaction{
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
		},
	}

	err := f.engine.Repositories().Operations.SaveAll(context.Background(), func(yield func(*txmodels.NewOperation) bool) {
		yield(&operation)
	})
	f.assert.NoError(err)

	// Additional check - assertion if the top-up operation was saved correctly
	balance, err := f.engine.Repositories().Users.GetBalance(context.Background(), f.user.ID(), "bsv")
	f.assert.NoError(err)
	f.assert.GreaterOrEqual(balance, satoshis)

	return txSpec
}
