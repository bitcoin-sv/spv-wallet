package engine

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

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
