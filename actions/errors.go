package actions

import "errors"

// ErrXpubNotFound is when the xpub is not found (in Auth Header)
var ErrXpubNotFound = errors.New("xpub not found")

// ErrTxConfigNotFound is when the transaction config is not found in request body
var ErrTxConfigNotFound = errors.New("transaction config not found")

// ErrBadTxConfig is when the transaction config specified is not valid
var ErrBadTxConfig = errors.New("bad transaction config supplied")
