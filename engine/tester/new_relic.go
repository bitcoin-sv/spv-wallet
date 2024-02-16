package tester

import (
	"context"
	"testing"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/require"
)

// GetNewRelicApp will return a dummy new relic app
func GetNewRelicApp(appName string) (*newrelic.Application, error) {
	if len(appName) == 0 {
		return nil, ErrAppNameRequired
	}
	return newrelic.NewApplication(
		func(config *newrelic.Config) {
			config.AppName = appName
			config.DistributedTracer.Enabled = true
			config.Enabled = false
		},
	)
}

// GetNewRelicCtx will return a dummy ctx
func GetNewRelicCtx(t *testing.T, appName, txnName string) context.Context {

	// Load new relic (dummy)
	newRelic, err := GetNewRelicApp(appName)
	require.NoError(t, err)
	require.NotNil(t, newRelic)

	// Create new relic tx
	return newrelic.NewContext(
		context.Background(), newRelic.StartTransaction(txnName),
	)
}
