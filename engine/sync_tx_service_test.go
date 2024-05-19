package engine

import (
	"context"
	"errors"
	"testing"

	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func Test_processBroadcastTransactions(t *testing.T) {
	// mocked broadcast client responses
	// broadcastSuccess := broadcastClientMqResponse{
	// 	r: &broadcast_client.SubmitTxResponse{
	// 		SubmittedTx: &broadcast_client.SubmittedTx{
	// 			BaseTxResponse: broadcast_client.BaseTxResponse{},
	// 		},
	// 	},
	// }

	broadcastClientFailed := broadcastClientMqResponse{
		r: new(broadcast_client.SubmitTxResponse),
		f: broadcast_client.Failure("error", errors.New("test client error")),
	}

	broadcastTransactionDeclined := broadcastClientMqResponse{
		r: new(broadcast_client.SubmitTxResponse),
		f: broadcast_client.Failure("invalid tx", &broadcast_client.ArcError{}),
	}

	tcs := []struct {
		name                    string
		expectedBroadcastStatus SyncStatus
		broadcastResponse       broadcastClientMqResponse
	}{
		// {
		// 	name:                    "broadcast success - status is complete",
		// 	expectedBroadcastStatus: SyncStatusComplete,
		// 	broadcastResponse:       broadcastSuccess,
		// },
		{
			name:                    "broadcast failed (client error) - status is ready",
			expectedBroadcastStatus: SyncStatusReady,
			broadcastResponse:       broadcastClientFailed,
		},
		{
			name:                    "broadcast failed (arc declined tx) - status is error",
			expectedBroadcastStatus: SyncStatusError,
			broadcastResponse:       broadcastTransactionDeclined,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// given
			ctx := context.Background()

			bc := &broadcastClientMq{}
			bc.setupResponse("SubmitTransaction", tc.broadcastResponse)

			spvengine := GetEngineClient(ctx, t, WithBroadcastClient(bc))
			defer CloseClient(ctx, t, spvengine)

			tx, _ := txFromHex(testTx2Hex, WithXPub(testXPub), WithClient(spvengine))
			stx := newSyncTransaction(tx.ID, &SyncConfig{Broadcast: true}, WithClient(spvengine))

			tx.syncTransaction = stx

			err := tx.Save(ctx)
			require.NoError(t, err)

			// when
			err = processBroadcastTransactions(ctx, 1, WithClient(spvengine))
			require.NoError(t, err)

			// then
			stx, err = GetSyncTransactionByTxID(ctx, tx.ID, WithClient(spvengine))
			require.NoError(t, err)

			require.Equal(t, tc.expectedBroadcastStatus, stx.BroadcastStatus)
		})
	}

}

func GetEngineClient(ctx context.Context, t *testing.T, o ...ClientOps) ClientInterface {
	log := zerolog.Nop()
	opts := []ClientOps{
		WithLogger(&log),
		WithChainstateOptions(true, true, true, true),
		WithSQLite(tester.SQLiteTestConfig(false, false)),
		WithAutoMigrate(append(BaseModels, newPaymail("", 0))...),
	}

	opts = append(opts, o...)
	spvengine, err := NewClient(ctx, opts...)
	require.NoError(t, err)

	return spvengine
}

type broadcastClientMq struct {
	responses map[string]broadcastClientMqResponse
}

type broadcastClientMqResponse struct {
	r any
	f *broadcast_client.FailureResponse
}

func (mq *broadcastClientMq) setupResponse(method string, response broadcastClientMqResponse) {
	if mq.responses == nil {
		mq.responses = make(map[string]broadcastClientMqResponse)
	}

	mq.responses[method] = response
}

func (mq *broadcastClientMq) GetFeeQuote(ctx context.Context) ([]*broadcast_client.FeeQuote, error) {
	if r, ok := mq.responses["GetFeeQuote"]; ok {
		return r.r.([]*broadcast_client.FeeQuote), r.f
	}

	return []*broadcast_client.FeeQuote{{MiningFee: broadcast_client.MiningFeeResponse{Bytes: 1, Satoshis: 1}}}, nil
}

func (*broadcastClientMq) GetPolicyQuote(ctx context.Context) ([]*broadcast_client.PolicyQuoteResponse, error) {
	return nil, nil
}

func (*broadcastClientMq) QueryTransaction(ctx context.Context, txID string) (*broadcast_client.QueryTxResponse, error) {
	return nil, nil
}

func (*broadcastClientMq) SubmitBatchTransactions(ctx context.Context, tx []*broadcast_client.Transaction, opts ...broadcast_client.TransactionOptFunc,
) (*broadcast_client.SubmitBatchTxResponse, error) {
	return nil, nil
}

func (mq *broadcastClientMq) SubmitTransaction(ctx context.Context, tx *broadcast_client.Transaction, opts ...broadcast_client.TransactionOptFunc,
) (*broadcast_client.SubmitTxResponse, error) {
	r := mq.responses["SubmitTransaction"]
	return r.r.(*broadcast_client.SubmitTxResponse), r.f
}
