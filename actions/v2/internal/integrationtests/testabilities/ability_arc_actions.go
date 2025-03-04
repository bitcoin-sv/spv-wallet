package testabilities

import (
	"testing"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/stretchr/testify/require"
)

type arcActions struct {
	t       testing.TB
	fixture *fixture
}

func (a *arcActions) SendsCallback(txInfo chainmodels.TXInfo) {
	client := a.fixture.HttpClient().ForAnonymous()
	token := a.fixture.Config().ARC.Callback.Token

	res, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(txInfo).
		SetAuthToken(token).
		Post("/transaction/broadcast/callback")

	require.Equal(a.t, 200, res.StatusCode())
}
