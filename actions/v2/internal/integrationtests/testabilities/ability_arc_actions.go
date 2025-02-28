package testabilities

import (
	"testing"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/stretchr/testify/require"
)

const ARCCallbackToken = "arc-test-token"

type arcActions struct {
	t       testing.TB
	fixture *fixture
}

func (a *arcActions) ReceivesCallback(txInfo chainmodels.TXInfo) {
	client := a.fixture.HttpClient().ForAnonymous()

	res, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(txInfo).
		SetAuthToken(ARCCallbackToken).
		Post("/arc/broadcast/callback")

	require.Equal(a.t, 200, res.StatusCode())
}
