//nolint:revive // Error types should be self-explanatory
package errors

import (
	"github.com/bitcoin-sv/spv-wallet/errdef"
	"github.com/joomcode/errorx"
)

var Namespace = errorx.NewNamespace("config", errdef.TraitConfig)

var UnsupportedDomain = Namespace.NewType("unsupported_domain", errdef.TraitIllegalArgument)
