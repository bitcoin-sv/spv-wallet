package engine

import "github.com/bitcoin-sv/spv-wallet/engine/spverrors"

// ErrDatastoreNotSupported is when a Datastore cannot be 100% tested
var ErrDatastoreNotSupported = spverrors.Newf("this Datastore is not supported for testing at this time")

// ErrLoadServerFirst is used when calling a database before the server was loaded first
var ErrLoadServerFirst = spverrors.Newf("the embedded database server must be loaded first")
