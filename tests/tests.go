// Package tests provides the base test suite for the entire package
package tests

import (
	"context"
	"os"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/logging"
	"github.com/bitcoin-sv/spv-wallet/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TestSuite is for testing the entire package using real/mocked services
type TestSuite struct {
	AppConfig *config.AppConfig // App config
	Router    *gin.Engine       // Gin router with handlers
	// Services    *config.AppServices // Services
	Logger          zerolog.Logger         // Logger
	SpvWalletEngine engine.ClientInterface // SPV Wallet Engine
	suite.Suite                            // Extends the suite.Suite package
}

// BaseSetupSuite runs at the start of the suite
func (ts *TestSuite) BaseSetupSuite() {
	cfg := config.GetDefaultAppConfig()
	cfg.DebugProfiling = false
	cfg.Logging.Level = zerolog.LevelDebugValue
	cfg.Logging.Format = "console"
	cfg.ARC.UseFeeQuotes = false
	cfg.ARC.FeeUnit = &config.FeeUnitConfig{
		Satoshis: 1,
		Bytes:    1000,
	}
	ts.AppConfig = cfg
}

// BaseTearDownSuite runs after the suite finishes
func (ts *TestSuite) BaseTearDownSuite() {
	ts.T().Cleanup(func() {
		_ = os.Remove("datastore.db")
	})
}

// BaseSetupTest runs before each test
func (ts *TestSuite) BaseSetupTest() {
	// Load the services
	var err error
	ts.Logger = tester.Logger(ts.T())

	opts, err := ts.AppConfig.ToEngineOptions(ts.Logger, true)
	require.NoError(ts.T(), err)

	ts.SpvWalletEngine, err = engine.NewClient(context.Background(), opts...)
	require.NoError(ts.T(), err)

	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	ginEngine.Use(logging.GinMiddleware(ts.Logger), gin.Recovery())
	ginEngine.Use(middleware.AppContextMiddleware(ts.AppConfig, ts.SpvWalletEngine, ts.Logger))
	ginEngine.Use(middleware.CorsMiddleware())

	ts.Router = ginEngine
	require.NotNil(ts.T(), ts.Router)

	require.NoError(ts.T(), err)
}

// BaseTearDownTest runs after each test
func (ts *TestSuite) BaseTearDownTest() {
	if ts.SpvWalletEngine != nil {
		err := ts.SpvWalletEngine.Close(context.Background())
		require.NoError(ts.T(), err)
	}
}
