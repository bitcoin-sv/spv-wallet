package tests

import (
	"context"
	"os"
	"path"
	"runtime"

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
	// Set the env to test
	err := os.Setenv(config.EnvironmentKey, config.EnvironmentTest)
	require.NoError(ts.T(), err)

	// Get current working directory
	var dirname string
	dirname, err = os.Getwd()
	require.NoError(ts.T(), err)

	// Go up one package
	var dir *os.File
	pathPrefix := "../../"

	if runtime.GOOS == "windows" {
		pathPrefix = "../../../"
	}

	dir, err = os.Open(path.Join(dirname, pathPrefix))
	require.NoError(ts.T(), err)

	// Load the configuration
	ts.AppConfig, err = config.Load(dir.Name())
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
