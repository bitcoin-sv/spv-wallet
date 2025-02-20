//go:build errorx
// +build errorx

package domainerr

import "github.com/joomcode/errorx"

var NotFound = errorx.DataUnavailable.NewSubtype("not_found")
var QueryFailed = errorx.InternalError.NewSubtype("repository.query_failed")
var SaveFailed = errorx.InternalError.NewSubtype("repository.write_failed")
