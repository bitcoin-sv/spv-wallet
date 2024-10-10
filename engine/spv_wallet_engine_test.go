package engine

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	broadcast_client_mock "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock"
	compat "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/mrz1836/go-cache"
	"github.com/rafaeljusto/redigomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestingClient is for testing the entire package using real/mocked services
type TestingClient struct {
	client      ClientInterface  // Local SPV Wallet Engine client for testing
	ctx         context.Context  // Current CTX
	database    datastore.Engine // Current database
	mocking     bool             // If mocking is enabled
	MockSQLDB   sqlmock.Sqlmock  // Mock Database client for SQL
	redisClient *cache.Client    // Current redis client (used for Mocking)
	redisConn   *redigomock.Conn // Current redis connection (used for Mocking)
	SQLConn     *sql.DB          // Read test client
	tablePrefix string           // Current table prefix
}

// Close will close all test services and client
func (tc *TestingClient) Close(ctx context.Context) {
	if tc.client != nil {
		/*if !tc.mocking {
			if err := dropAllTables(tc.client.Datastore(), tc.database); err != nil {
				panic(err)
			}
		}*/
		_ = tc.client.Close(ctx)
	}
	if tc.SQLConn != nil {
		_ = tc.SQLConn.Close()
	}
	if tc.redisClient != nil {
		tc.redisClient.Close()
	}
	if tc.redisConn != nil {
		_ = tc.redisConn.Close()
	}
}

// DefaultClientOpts will return a default set of client options required to load the new client
func DefaultClientOpts(debug, shared bool) []ClientOps {
	tqc := taskmanager.DefaultTaskQConfig(tester.RandomTablePrefix())
	tqc.MaxNumWorker = 2
	tqc.MaxNumFetcher = 2
	bc := broadcast_client_mock.Builder().
		WithMockArc(broadcast_client_mock.MockNilQueryTxResp).
		Build()

	opts := make([]ClientOps, 0)
	opts = append(
		opts,
		WithTaskqConfig(tqc),
		WithSQLite(tester.SQLiteTestConfig(debug, shared)),
		WithChainstateOptions(false, false, false, false),
		WithBroadcastClient(bc),
	)
	if debug {
		opts = append(opts, WithDebugging())
	}

	return opts
}

// CreateTestSQLiteClient will create a test client for SQLite
//
// NOTE: you need to close the client using the returned defer func
func CreateTestSQLiteClient(t *testing.T, debug, shared bool, clientOpts ...ClientOps) (context.Context, ClientInterface, func()) {
	ctx := context.Background()

	logger := zerolog.Nop()

	// Set the default options, add migrate models
	opts := DefaultClientOpts(debug, shared)
	opts = append(opts, WithAutoMigrate(append(BaseModels, newPaymail("", 0))...))
	opts = append(opts, WithLogger(&logger))
	opts = append(opts, clientOpts...)

	// Create the client
	client, err := NewClient(ctx, opts...)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Create a defer function
	f := func() {
		_ = client.Close(context.Background())
	}
	return ctx, client, f
}

// CreateBenchmarkSQLiteClient will create a test client for SQLite
//
// NOTE: you need to close the client using the returned defer func
func CreateBenchmarkSQLiteClient(b *testing.B, debug, shared bool, clientOpts ...ClientOps) (context.Context, ClientInterface, func()) {
	ctx := context.Background()

	logger := zerolog.Nop()

	// Set the default options, add migrate models
	opts := DefaultClientOpts(debug, shared)
	opts = append(opts, WithAutoMigrate(BaseModels...))
	opts = append(opts, WithLogger(&logger))
	opts = append(opts, clientOpts...)

	// Create the client
	client, err := NewClient(ctx, opts...)
	if err != nil {
		b.Fail()
	}

	// Create a defer function
	f := func() {
		_ = client.Close(context.Background())
	}
	return ctx, client, f
}

// CloseClient is function used in the "defer()" function
func CloseClient(ctx context.Context, t *testing.T, client ClientInterface) {
	require.NoError(t, client.Close(ctx))
}

// we need to create an interface for the unlocker
type account struct {
	PrivateKey *ec.PrivateKey
}

// Unlocker get the correct un-locker for a given locking script.
func (a *account) Unlocker(context.Context, *script.Script) (*p2pkh.P2PKH, error) {
	unlocker, err := p2pkh.Unlock(a.PrivateKey, nil)
	if err != nil {
		return nil, err
	}

	return unlocker, nil
}

// CreateFakeFundingTransaction will create a valid (fake) transaction for funding
func CreateFakeFundingTransaction(t *testing.T, masterKey *compat.ExtendedKey,
	destinations []*Destination, satoshis uint64,
) string {
	// Create new tx
	rawTx := trx.NewTransaction()
	txErr := rawTx.AddInputFrom(testTxScriptSigID, 0, testTxScriptSigOut, satoshis+354, nil)
	require.NoError(t, txErr)

	// Loop all destinations
	for _, destination := range destinations {
		s, err := script.NewFromHex(destination.LockingScript)
		require.NoError(t, err)
		require.NotNil(t, s)

		rawTx.AddOutput(&trx.TransactionOutput{
			Satoshis:      satoshis,
			LockingScript: s,
		})
	}

	// Get private key
	privateKey, err := compat.GetPrivateKeyFromHDKey(masterKey)
	require.NoError(t, err)
	require.NotNil(t, privateKey)

	err = rawTx.Sign()
	require.NoError(t, err)

	// Return the tx hex
	return rawTx.String()
}

// CreateNewXPub will create a new xPub and return all the information to use the xPub
func CreateNewXPub(ctx context.Context, t *testing.T, engineClient ClientInterface,
	opts ...ModelOps,
) (*compat.ExtendedKey, *Xpub, string) {
	// Generate a key pair
	masterKey, err := compat.GenerateHDKey(compat.SecureSeedLength)
	require.NoError(t, err)
	require.NotNil(t, masterKey)

	// Get the raw string of the xPub
	var rawXPub string
	rawXPub, err = compat.GetExtendedPublicKey(masterKey)
	require.NoError(t, err)
	require.NotNil(t, masterKey)

	// Create the new xPub
	var xPub *Xpub
	xPub, err = engineClient.NewXpub(ctx, rawXPub, opts...)
	require.NoError(t, err)
	require.NotNil(t, xPub)

	return masterKey, xPub, rawXPub
}

// GetUnlockingScript will get a locking script for valid fake transactions
func GetUnlockingScript(tx *trx.Transaction, inputIndex uint32, privateKey *ec.PrivateKey) (*p2pkh.P2PKH, error) {
	sigHashFlags := sighash.AllForkID

	sc, err := p2pkh.Unlock(privateKey, &sigHashFlags)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to create unlocking script")
	}

	return sc, nil
}
