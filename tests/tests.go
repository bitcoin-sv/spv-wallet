// Package tests provides the base test suite for the entire package
package tests

import (
	"context"
	"os"

	"github.com/BuxOrg/bux-server/config"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TestSuite is for testing the entire package using real/mocked services
type TestSuite struct {
	AppConfig   *config.AppConfig   // App config
	Router      *apirouter.Router   // Router with handlers
	Services    *config.AppServices // Services
	suite.Suite                     // Extends the suite.Suite package
}

// BaseSetupSuite runs at the start of the suite
func (ts *TestSuite) BaseSetupSuite() {
	// Load the configuration
	var err error
	ts.AppConfig, err = config.Load("")
	require.NoError(ts.T(), err)
}

// BaseTearDownSuite runs after the suite finishes
func (ts *TestSuite) BaseTearDownSuite() {
	// Ensure all connections are closed
	if ts.Services != nil {
		ts.Services.CloseAll(context.Background())
		ts.Services = nil
	}

	ts.T().Cleanup(func() {
		_ = os.Remove("datastore.db")
	})
}

// BaseSetupTest runs before each test
func (ts *TestSuite) BaseSetupTest() {
	// Load the services
	var err error
	ts.Services, err = ts.AppConfig.LoadTestServices(context.Background())
	require.NoError(ts.T(), err)
}

// BaseTearDownTest runs after each test
func (ts *TestSuite) BaseTearDownTest() {
	if ts.Services != nil {
		ts.Services.CloseAll(context.Background())
		ts.Services = nil
	}
}
