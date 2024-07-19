package tester

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// ErrAppNameRequired is when the app name is required
var ErrAppNameRequired = spverrors.Newf("app name is required")

// ErrFailedLoadingPostgresql is when loading postgresql failed
var ErrFailedLoadingPostgresql = spverrors.Newf("failed loading postgresql server")
