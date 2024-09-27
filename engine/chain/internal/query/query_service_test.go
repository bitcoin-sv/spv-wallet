package query

import (
	"context"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestQueryService(t *testing.T) {
	logger := zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t)))
	deepSuffix, _ := uuid.NewUUID()
	deploymentID := "spv-wallet-" + deepSuffix.String()

	testCases := map[string]struct {
		txID         string
		arcToken     string
		arcURL       string
		expectErr    error
		expectNil    bool
		expectTxID   string
		expectStatus string
		applyTimeout timeoutDst
	}{
		"Query for MINED transaction": {
			txID:         minedTxID,
			arcToken:     arcToken,
			arcURL:       arcURL,
			expectErr:    nil,
			expectTxID:   minedTxID,
			expectStatus: "MINED",
		},
		"Query for unknown transaction": {
			txID:      unknownTxID,
			arcToken:  arcToken,
			arcURL:    arcURL,
			expectErr: nil,
			expectNil: true,
		},
		"Query for invalid transaction": {
			txID:      "invalid",
			arcToken:  arcToken,
			arcURL:    arcURL,
			expectErr: spverrors.ErrInvalidTransactionID,
			expectNil: true,
		},
		"Query with wrong token": {
			txID:      minedTxID,
			arcToken:  "wrong-token",
			arcURL:    arcURL,
			expectErr: spverrors.ErrARCUnauthorized,
			expectNil: true,
		},
		"Query 404 endpoint but reachable": {
			txID:      minedTxID,
			arcToken:  arcToken,
			arcURL:    arcURL + wrongButReachable,
			expectErr: spverrors.ErrARCUnreachable,
			expectNil: true,
		},
		"Query 404 endpoint with wrong arcURL": {
			txID:      minedTxID,
			arcToken:  arcToken,
			arcURL:    "wrong-url",
			expectErr: spverrors.ErrARCUnreachable,
			expectNil: true,
		},
		"Query interrupted by ctx timeout": {
			txID:         minedTxID,
			arcToken:     arcToken,
			arcURL:       arcURL,
			expectErr:    spverrors.ErrARCUnreachable,
			expectNil:    true,
			applyTimeout: applyTimeoutCtx,
		},
		"Query interrupted by resty timeout": {
			txID:         minedTxID,
			arcToken:     arcToken,
			arcURL:       arcURL,
			expectErr:    spverrors.ErrARCUnreachable,
			expectNil:    true,
			applyTimeout: applyTimeoutResty,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			httpClient, reset := arcMockActivate()
			defer reset()

			service := NewQueryService(logger, httpClient, tc.arcURL, tc.arcToken, deploymentID)

			ctx := context.Background()
			if tc.applyTimeout == applyTimeoutCtx {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 1*time.Millisecond)
				defer cancel()
			} else if tc.applyTimeout == applyTimeoutResty {
				service.httpClient.SetTimeout(1 * time.Millisecond)
			}

			txInfo, err := service.Query(ctx, tc.txID)

			if tc.expectErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectErr)
			} else {
				require.NoError(t, err)
			}

			if tc.expectNil {
				require.Nil(t, txInfo)
			} else {
				require.NotNil(t, txInfo)
				require.Equal(t, tc.expectTxID, txInfo.TxID)
				require.Equal(t, tc.expectStatus, string(txInfo.TXStatus))
			}
		})
	}
}

type timeoutDst int

const (
	applyTimeoutCtx timeoutDst = iota + 1
	applyTimeoutResty
)
