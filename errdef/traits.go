//nolint:revive // Traits should be self-explanatory
package errdef

import "github.com/joomcode/errorx"

var TraitConfig = errorx.RegisterTrait("config")
var TraitIllegalArgument = errorx.RegisterTrait("illegal_argument")
var TraitNotFound = errorx.RegisterTrait("not_found")
var TraitAuth = errorx.RegisterTrait("auth")
var TraitARC = errorx.RegisterTrait("arc")
var TraitShouldNeverHappen = errorx.RegisterTrait("should_never_happen")
var TraitUnsupported = errorx.RegisterTrait("unsupported")
