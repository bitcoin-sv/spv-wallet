package destinations

import (
	"context"
	"os"
	"path"
	"testing"

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

// SetupSuite runs at the start of the suite
func (ts *TestSuite) SetupSuite() {

	// Set the env to test
	err := os.Setenv(config.EnvironmentKey, config.EnvironmentTest)
	require.NoError(ts.T(), err)

	// Get current working directory
	var dirname string
	dirname, err = os.Getwd()
	require.NoError(ts.T(), err)

	// Go up one package
	var dir *os.File
	dir, err = os.Open(path.Join(dirname, "../../"))
	require.NoError(ts.T(), err)

	// Load the configuration
	ts.AppConfig, err = config.Load(dir.Name())
	require.NoError(ts.T(), err)
}

// TearDownSuite runs after the suite finishes
func (ts *TestSuite) TearDownSuite() {

	// Ensure all connections are closed
	if ts.Services != nil {
		ts.Services.CloseAll(context.Background())
		ts.Services = nil
	}

	ts.T().Cleanup(func() {
		_ = os.Remove("datastore.db")
	})
}

// SetupTest runs before each test
func (ts *TestSuite) SetupTest() {

	// Load the services
	var err error
	ts.Services, err = ts.AppConfig.LoadTestServices(context.Background())
	require.NoError(ts.T(), err)

	// Load the router & register routes
	ts.Router = apirouter.New()
	require.NotNil(ts.T(), ts.Router)
	RegisterRoutes(ts.Router, ts.AppConfig, ts.Services)
}

// TearDownTest runs after each test
func (ts *TestSuite) TearDownTest() {
	if ts.Services != nil {
		ts.Services.CloseAll(context.Background())
		ts.Services = nil
	}
}

// TestTestSuite kick-starts all suite tests
func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
