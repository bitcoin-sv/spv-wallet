package engine

import "errors"

// ErrDatastoreNotSupported is when a Datastore cannot be 100% tested
var ErrDatastoreNotSupported = errors.New("this Datastore is not supported for testing at this time")

// ErrLoadServerFirst is used when calling a database before the server was loaded first
var ErrLoadServerFirst = errors.New("the embedded database server must be loaded first")
