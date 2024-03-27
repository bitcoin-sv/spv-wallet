package datastore

import (
	"errors"
)

// ErrUnsupportedEngine is used when the engine given is not a known datastore engine
var ErrUnsupportedEngine = errors.New("unsupported datastore engine")

// ErrNoResults error when no results are found
var ErrNoResults = errors.New("no results found")

// ErrDuplicateKey error when a record is inserted and conflicts with an existing record
var ErrDuplicateKey = errors.New("duplicate key")

// ErrUnknownCollection is thrown when the collection can not be found using the model/name
var ErrUnknownCollection = errors.New("could not determine collection name from model")

// ErrUnsupportedDriver is when the given SQL driver is not determined to be known or supported
var ErrUnsupportedDriver = errors.New("sql driver unsupported")

// ErrNoSourceFound is when no source database is found in all given configurations
var ErrNoSourceFound = errors.New("no source database found in all given configurations")

// ErrUnknownSQL is an error when using a SQL engine that is not known for indexes and migrations
var ErrUnknownSQL = errors.New("unknown sql implementation")

// ErrNotImplemented is an error when a method is not implemented
var ErrNotImplemented = errors.New("not implemented")
