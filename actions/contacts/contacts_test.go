package contacts

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	"github.com/bitcoin-sv/spv-wallet/tests"
	"github.com/stretchr/testify/suite"
)

// TestSuite is for testing the entire package using real/mocked services
type TestSuite struct {
	tests.TestSuite
}

// SetupSuite runs at the start of the suite
func (ts *TestSuite) SetupSuite() {
	ts.BaseSetupSuite()
	ts.AppConfig.ExperimentalFeatures.PikeContactsEnabled = true
}

// TearDownSuite runs after the suite finishes
func (ts *TestSuite) TearDownSuite() {
	ts.BaseTearDownSuite()
}

// SetupTest runs before each test
func (ts *TestSuite) SetupTest() {
	ts.BaseSetupTest()

	handlersManager := handlers.NewManager(ts.Router, ts.AppConfig)
	RegisterRoutes(handlersManager)
}

// TearDownTest runs after each test
func (ts *TestSuite) TearDownTest() {
	ts.BaseTearDownTest()
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
