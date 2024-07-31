package datastore

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// ErrUnsupportedEngine is used when the engine given is not a known datastore engine
var ErrUnsupportedEngine = spverrors.Newf("unsupported datastore engine")

// ErrNoResults error when no results are found
var ErrNoResults = spverrors.Newf("no results found")

// ErrDuplicateKey error when a record is inserted and conflicts with an existing record
var ErrDuplicateKey = spverrors.Newf("duplicate key")

// ErrUnknownCollection is thrown when the collection can not be found using the model/name
var ErrUnknownCollection = spverrors.Newf("could not determine collection name from model")

// ErrUnsupportedDriver is when the given SQL driver is not determined to be known or supported
var ErrUnsupportedDriver = spverrors.Newf("sql driver unsupported")

// ErrNoSourceFound is when no source database is found in all given configurations
var ErrNoSourceFound = spverrors.Newf("no source database found in all given configurations")

// ErrUnknownSQL is an error when using a SQL engine that is not known for indexes and migrations
var ErrUnknownSQL = spverrors.Newf("unknown sql implementation")

// ErrNotImplemented is an error when a method is not implemented
var ErrNotImplemented = spverrors.Newf("not implemented")
