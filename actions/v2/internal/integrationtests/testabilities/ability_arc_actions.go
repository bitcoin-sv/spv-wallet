package testabilities

import (
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/stretchr/testify/require"
	"testing"
)

const ARCCallbackToken = "arc-test-token"

type arcActions struct {
	t       testing.TB
	fixture *fixture
}

func (a *arcActions) Callbacks(txInfo chainmodels.TXInfo) {
	client := a.fixture.HttpClient().ForAnonymous()

	res, _ := client.R().
		SetBody(txInfo).
		SetAuthToken(ARCCallbackToken).
		Post("/transactions/transaction/broadcast/callback")

	require.Equal(a.t, 200, res.StatusCode())
}
