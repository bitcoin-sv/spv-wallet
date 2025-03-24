//nolint:revive // Error types should be self-explanatory
package paymailerrors

import (
	"github.com/bitcoin-sv/spv-wallet/errdef"
	"github.com/joomcode/errorx"
)

var Namespace = errorx.NewNamespace("paymail")

var InvalidAvatarURL = Namespace.NewType("invalid_avatar_url", errdef.TraitIllegalArgument)
var InvalidPaymailAddress = Namespace.NewType("invalid_paymail_address", errdef.TraitIllegalArgument)
var UserDoesntExist = Namespace.NewType("user_doesnt_exist", errdef.TraitNotFound)
var NoDefaultPaymailAddress = Namespace.NewType("no_default_paymail_address", errdef.TraitIllegalArgument)
