package errdef

import "github.com/joomcode/errorx"

var PropSpecificProblemOccurrence = errorx.RegisterPrintableProperty("instance")
var PropPublicHint = errorx.RegisterPrintableProperty("public_hint")

var RootNamespace = errorx.NewNamespace("spv-wallet")
var Internal = RootNamespace.NewSubNamespace("5xx")
var External = RootNamespace.NewSubNamespace("4xx")

var UnsupportedOperation = Internal.NewType("unsupported_operation")
var NotImplementedYet = UnsupportedOperation.NewSubtype("not_implemented_yet")
