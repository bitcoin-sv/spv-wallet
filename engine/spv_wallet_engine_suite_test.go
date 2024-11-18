package engine

import (
	"context"
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	embeddedPostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	postgresqlTestHost   = "localhost"
	postgresqlTestName   = "postgres"
	postgresqlTestPort   = uint32(61333)
	postgresqlTestUser   = "postgres"
	postgresTestPassword = "postgres"
	testQueueName        = "test_queue"
)

var mockFeeUnit = bsv.FeeUnit{Satoshis: 1, Bytes: 1000}

// dbTestCase is a database test case
type dbTestCase struct {
	name     string
	database datastore.Engine
}

// dbTestCases is the list of supported databases
var dbTestCases = []dbTestCase{
	{name: "[postgresql] [in-memory]", database: datastore.PostgreSQL},
	{name: "[sqlite] [in-memory]", database: datastore.SQLite},
}

// EmbeddedDBTestSuite is for testing the entire package using real/mocked services
type EmbeddedDBTestSuite struct {
	suite.Suite
	PostgresqlServer *embeddedPostgres.EmbeddedPostgres // In-memory Postgresql server
}

// SetupSuite runs at the start of the suite
func (ts *EmbeddedDBTestSuite) SetupSuite() {
	var err error
	// Create the Postgresql server
	if ts.PostgresqlServer, err = tester.CreatePostgresServer(postgresqlTestPort); err != nil {
		require.NoError(ts.T(), err)
	}

	// Fail-safe! If a test completes or fails, this is triggered
	// Embedded servers are still running on the ports given, and causes a conflict re-running tests
	ts.T().Cleanup(func() {
		ts.TearDownSuite()
	})
}

// TearDownSuite runs after the suite finishes
func (ts *EmbeddedDBTestSuite) TearDownSuite() {
	// Stop the postgresql server
	if ts.PostgresqlServer != nil {
		_ = ts.PostgresqlServer.Stop()
	}
}

// SetupTest runs before each test
func (ts *EmbeddedDBTestSuite) SetupTest() {
	// Nothing needed here (yet)
}

// TearDownTest runs after each test
func (ts *EmbeddedDBTestSuite) TearDownTest() {
	// Nothing needed here (yet)
}

// createTestClient will make a new test client
//
// NOTE: you need to close the client: ts.Close()
func (ts *EmbeddedDBTestSuite) createTestClient(ctx context.Context, database datastore.Engine,
	tablePrefix string, opts ...ClientOps,
) (*TestingClient, error) {
	var err error

	// Start the suite
	tc := &TestingClient{
		ctx:         ctx,
		database:    database,
		tablePrefix: tablePrefix,
	}

	// Load the in-memory version of the database
	if database == datastore.SQLite {
		opts = append(opts, WithSQLite(&datastore.SQLiteConfig{
			CommonConfig: datastore.CommonConfig{
				MaxIdleConnections: 1,
				MaxOpenConnections: 1,
				TablePrefix:        tablePrefix,
			},
			Shared: true, // mrz: TestTransaction_Save requires this to be true for some reason
			// I get the error: no such table: _17a1f3e22f2eec56_utxos
		}))
	} else if database == datastore.PostgreSQL {

		// Sanity check
		if ts.PostgresqlServer == nil {
			return nil, ErrLoadServerFirst
		}

		// Add the new Postgresql connection
		opts = append(opts, WithSQL(datastore.PostgreSQL, &datastore.SQLConfig{
			CommonConfig: datastore.CommonConfig{
				MaxIdleConnections: 1,
				MaxOpenConnections: 1,
				TablePrefix:        tablePrefix,
			},
			Host:     postgresqlTestHost,
			Name:     postgresqlTestName,
			User:     postgresqlTestUser,
			Password: postgresTestPassword,
			Port:     fmt.Sprintf("%d", postgresqlTestPort),
		}))

	} else {
		return nil, ErrDatastoreNotSupported
	}

	// Add a custom user agent (future: make this passed into the function via opts)
	opts = append(opts, WithUserAgent("spv wallet engine test suite"))

	// Create the client
	testLogger := zerolog.Nop()
	opts = append(opts, WithLogger(&testLogger))

	if tc.client, err = NewClient(ctx, opts...); err != nil {
		return nil, err
	}

	// Return the suite
	return tc, nil
}

// genericDBClient is a helpful wrapper for getting the same type of client
//
// NOTE: you need to close the client: ts.Close()
//
//nolint:nolintlint,unparam,gci // opts is the way, but not yet being used
func (ts *EmbeddedDBTestSuite) genericDBClient(t *testing.T, database datastore.Engine, taskManagerEnabled bool, opts ...ClientOps) *TestingClient {
	prefix := tester.RandomTablePrefix()

	if opts == nil {
		opts = []ClientOps{}
	}
	opts = append(opts,
		WithDebugging(),
		WithAutoMigrate(BaseModels...),
		WithAutoMigrate(&PaymailAddress{}),
		WithCustomFeeUnit(mockFeeUnit),
	)
	if taskManagerEnabled {
		opts = append(opts, WithTaskqConfig(taskmanager.DefaultTaskQConfig(prefix+"_queue")))
	} else {
		opts = append(opts, withTaskManagerMockup())
	}

	tc, err := ts.createTestClient(
		context.Background(),
		database, prefix,
		opts...,
	)
	require.NoError(t, err)
	require.NotNil(t, tc)
	return tc
}
