package errdef

import "github.com/joomcode/errorx"

// INTERNAL ERRORS

var PropSpecificProblemOccurrence = errorx.RegisterPrintableProperty("instance")

var ServerNamespace = errorx.NewNamespace("spv-wallet")

var UnsupportedOperation = ServerNamespace.NewType("unsupported_operation")
var NotImplementedYet = UnsupportedOperation.NewSubtype("not_implemented_yet")
