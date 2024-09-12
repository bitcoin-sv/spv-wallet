// Package tests provides the base test suite for the entire package
package tests

import (
	"context"
	"os"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/logging"
	"github.com/bitcoin-sv/spv-wallet/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TestSuite is for testing the entire package using real/mocked services
type TestSuite struct {
	AppConfig   *config.AppConfig   // App config
	Router      *gin.Engine         // Gin router with handlers
	Services    *config.AppServices // Services
	suite.Suite                     // Extends the suite.Suite package
}

// BaseSetupSuite runs at the start of the suite
func (ts *TestSuite) BaseSetupSuite() {
	ts.AppConfig = config.LoadForTest()
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
	nop := zerolog.Nop()

	ts.Services, err = ts.AppConfig.LoadTestServices(context.Background())

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(logging.GinMiddleware(&nop), gin.Recovery())
	engine.Use(middleware.AppContextMiddleware(ts.AppConfig, ts.Services.SpvWalletEngine, ts.Services.Logger))
	engine.Use(middleware.CorsMiddleware())

	ts.Router = engine
	require.NotNil(ts.T(), ts.Router)

	require.NoError(ts.T(), err)
}

// BaseTearDownTest runs after each test
func (ts *TestSuite) BaseTearDownTest() {
	if ts.Services != nil {
		ts.Services.CloseAll(context.Background())
		ts.Services = nil
	}
}
