package tester

import "errors"

// ErrAppNameRequired is when the app name is required
var ErrAppNameRequired = errors.New("app name is required")

// ErrFailedLoadingPostgresql is when loading postgresql failed
var ErrFailedLoadingPostgresql = errors.New("failed loading postgresql server")
