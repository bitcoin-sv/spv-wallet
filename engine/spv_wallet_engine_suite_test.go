package engine

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	embeddedPostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tryvium-travels/memongo"
)

const (
	defaultNewRelicTx        = "testing-transaction"
	defaultNewRelicApp       = "testing-app"
	postgresqlTestHost       = "localhost"
	postgresqlTestName       = "postgres"
	postgresqlTestPort       = uint32(61333)
	postgresqlTestUser       = "postgres"
	postgresTestPassword     = "postgres"
	testIdleTimeout          = 240 * time.Second
	testMaxActiveConnections = 0
	testMaxConnLifetime      = 60 * time.Second
	testMaxIdleConnections   = 10
	testQueueName            = "test_queue"
)

// dbTestCase is a database test case
type dbTestCase struct {
	name     string
	database datastore.Engine
}

// dbTestCases is the list of supported databases
var dbTestCases = []dbTestCase{
	{name: "[mongo] [in-memory]", database: datastore.MongoDB},
	{name: "[mysql] [in-memory]", database: datastore.MySQL},
	{name: "[postgresql] [in-memory]", database: datastore.PostgreSQL},
	{name: "[sqlite] [in-memory]", database: datastore.SQLite},
}

// EmbeddedDBTestSuite is for testing the entire package using real/mocked services
type EmbeddedDBTestSuite struct {
	suite.Suite
	MongoServer      *memongo.Server                    // In-memory  Mongo server 	// In-memory MySQL server
	PostgresqlServer *embeddedPostgres.EmbeddedPostgres // In-memory Postgresql server
}

// SetupSuite runs at the start of the suite
func (ts *EmbeddedDBTestSuite) SetupSuite() {
	var err error

	// Create the Mongo server
	// if ts.MongoServer, err = tester.CreateMongoServer(mongoTestVersion); err != nil {
	// 	require.NoError(ts.T(), err)
	// }

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
	// Stop the Mongo server
	if ts.MongoServer != nil {
		ts.MongoServer.Stop()
	}

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
	tablePrefix string, mockDB, mockRedis bool, opts ...ClientOps,
) (*TestingClient, error) {
	var err error

	// Start the suite
	tc := &TestingClient{
		ctx:         ctx,
		database:    database,
		mocking:     mockDB,
		tablePrefix: tablePrefix,
	}

	// Are we mocking SQL?
	if mockDB {

		// Create new SQL mocked connection
		if tc.SQLConn, tc.MockSQLDB, err = sqlmock.New(
			sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
		); err != nil {
			return nil, err
		}

		// Switch on database types
		if database == datastore.SQLite {
			opts = append(opts, WithSQLite(&datastore.SQLiteConfig{
				CommonConfig: datastore.CommonConfig{
					MaxConnectionIdleTime: 0,
					MaxConnectionTime:     0,
					MaxIdleConnections:    1,
					MaxOpenConnections:    1,
					TablePrefix:           tablePrefix,
				},
				ExistingConnection: tc.SQLConn,
			}))
		} else if database == datastore.MySQL {
			opts = append(opts, WithSQLConnection(datastore.MySQL, tc.SQLConn, tablePrefix))
		} else if database == datastore.PostgreSQL {
			opts = append(opts, WithSQLConnection(datastore.PostgreSQL, tc.SQLConn, tablePrefix))
		} else { // todo: finish more Datastore support (missing: Mongo)
			// "https://medium.com/@victor.neuret/mocking-the-official-mongo-golang-driver-5aad5b226a78"
			return nil, ErrDatastoreNotSupported
		}

	} else {
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
		} else if database == datastore.MongoDB {

			// Sanity check
			if ts.MongoServer == nil {
				return nil, ErrLoadServerFirst
			}

			// Add the new Mongo connection
			opts = append(opts, WithMongoDB(&datastore.MongoDBConfig{
				CommonConfig: datastore.CommonConfig{
					MaxIdleConnections: 1,
					MaxOpenConnections: 1,
					TablePrefix:        tablePrefix,
				},
				URI:          ts.MongoServer.URIWithRandomDB(),
				DatabaseName: memongo.RandomDatabase(),
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
				Host:                      postgresqlTestHost,
				Name:                      postgresqlTestName,
				User:                      postgresqlTestUser,
				Password:                  postgresTestPassword,
				Port:                      fmt.Sprintf("%d", postgresqlTestPort),
				SkipInitializeWithVersion: true,
			}))

		} else {
			return nil, ErrDatastoreNotSupported
		}
	}

	// Custom for SQLite and Mocking (cannot ignore the version check that GORM does)
	if mockDB && database == datastore.SQLite {
		tc.MockSQLDB.ExpectQuery(
			"select sqlite_version()",
		).WillReturnRows(tc.MockSQLDB.NewRows([]string{"version"}).FromCSVString(sqliteTestVersion))
	}

	// Are we mocking redis?
	if mockRedis {
		tc.redisClient, tc.redisConn = tester.LoadMockRedis(
			testIdleTimeout,
			testMaxConnLifetime,
			testMaxActiveConnections,
			testMaxIdleConnections,
		)
		opts = append(opts, WithRedisConnection(tc.redisClient))
	}

	// Add a custom user agent (future: make this passed into the function via opts)
	opts = append(opts, WithUserAgent("spv wallet engine test suite"))
	opts = append(opts, WithMinercraft(&chainstate.MinerCraftBase{}))

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
		WithChainstateOptions(false, false, false, false),
		WithAutoMigrate(BaseModels...),
		WithAutoMigrate(&PaymailAddress{}),
	)
	if taskManagerEnabled {
		opts = append(opts, WithTaskqConfig(taskmanager.DefaultTaskQConfig(prefix+"_queue")))
	} else {
		opts = append(opts, withTaskManagerMockup())
	}

	tc, err := ts.createTestClient(
		tester.GetNewRelicCtx(
			t, defaultNewRelicApp, defaultNewRelicTx,
		),
		database, prefix,
		false, false,
		opts...,
	)
	require.NoError(t, err)
	require.NotNil(t, tc)
	return tc
}

// genericMockedDBClient is a helpful wrapper for getting the same type of client
//
// NOTE: you need to close the client: ts.Close()
func (ts *EmbeddedDBTestSuite) genericMockedDBClient(t *testing.T, database datastore.Engine) *TestingClient {
	prefix := tester.RandomTablePrefix()
	tc, err := ts.createTestClient(
		tester.GetNewRelicCtx(
			t, defaultNewRelicApp, defaultNewRelicTx,
		),
		database, prefix,
		true, true, WithDebugging(),
		withTaskManagerMockup(),
	)
	require.NoError(t, err)
	require.NotNil(t, tc)
	return tc
}
