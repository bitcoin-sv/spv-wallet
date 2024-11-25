package datastore

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// ErrUnsupportedEngine is used when the engine given is not a known datastore engine
var ErrUnsupportedEngine = spverrors.Newf("unsupported datastore engine")

// ErrUnsupportedDriver is when the given SQL driver is not determined to be known or supported
var ErrUnsupportedDriver = spverrors.Newf("sql driver unsupported")

// ErrNoSourceFound is when no source database is found in all given configurations
var ErrNoSourceFound = spverrors.Newf("no source database found in all given configurations")

// ErrMigrationError is when something went wrong during migration
var ErrMigrationError = spverrors.Newf("error occurred during migration")
