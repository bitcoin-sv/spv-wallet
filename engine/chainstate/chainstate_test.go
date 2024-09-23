package chainstate

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// NewTestClient returns a test client
func NewTestClient(ctx context.Context, t *testing.T, opts ...ClientOps) ClientInterface {
	logger := zerolog.Nop()
	c, err := NewClient(
		ctx, append(opts, WithDebugging(), WithLogger(&logger))...,
	)
	require.NoError(t, err)
	require.NotNil(t, c)
	return c
}
