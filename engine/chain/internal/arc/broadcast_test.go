package arc_test

import (
	"context"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

const txHex = "0100000001293f17ea61f50d5ea815780c3d571f0f475533b8e812189724ab8e14b77e1616000000006a4730440220353c86782552be0c768cf675ca68c914c07cd2a35970292879353d089ca012e0022057c7a3b4b96cfedb7bb3555017044b11bd610477285b7f4e63efb5a10b906bf4412103513000984c44b7316671c1875c32eaeeacfd886f561623479794913c1cb91f73ffffffff01000000000000000038006a35323032342d31302d31312031333a31313a34322e30393234363935202b303230302043455354206d3d2b302e30313831363935303100000000"

func TestBroadcastTransaction(t *testing.T) {
	httpClient := resty.New()

	tx, err := sdk.NewTransactionFromHex(txHex)
	require.NoError(t, err)

	cfg := arcCfg(arcURL, arcToken)
	cfg.UseJunglebus = true

	service := chain.NewChainService(tester.Logger(t), httpClient, cfg, chainmodels.BHSConfig{})

	txInfo, err := service.Broadcast(context.Background(), tx)
	require.NoError(t, err)
	t.Logf("%+v", txInfo)
}
