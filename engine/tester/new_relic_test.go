package tester

import (
	"testing"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetNewRelicApp will test the method GetNewRelicApp()
func TestGetNewRelicApp(t *testing.T) {
	t.Parallel()

	t.Run("create an app", func(t *testing.T) {
		app, err := GetNewRelicApp("test-app")
		require.NoError(t, err)
		require.NotNil(t, app)
		txn := app.StartTransaction("test-transaction")
		require.NotNil(t, txn)
	})

	t.Run("missing an app", func(t *testing.T) {
		app, err := GetNewRelicApp("")
		require.Error(t, err)
		require.Nil(t, app)
		assert.ErrorIs(t, err, ErrAppNameRequired)
	})
}

// TestGetNewRelicCtx will test the method GetNewRelicCtx()
func TestGetNewRelicCtx(t *testing.T) {
	t.Parallel()

	t.Run("load the ctx", func(t *testing.T) {
		ctx := GetNewRelicCtx(t, "test-app", "test-transaction")
		require.NotNil(t, ctx)

		txn := newrelic.FromContext(ctx)
		require.NotNil(t, txn)
		app := txn.Application()
		require.NotNil(t, app)
	})
}
