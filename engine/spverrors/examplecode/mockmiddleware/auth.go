package mockmiddleware

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/clienterr"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/errdef"
	"github.com/joomcode/errorx"
)

var MiddlewareErr = errorx.NewNamespace("middleware", errdef.TraitAuth)
var SomeAlgoFailed = MiddlewareErr.NewType("some_auth_algo_failed")

func Auth(fail *api.ModelsFailingPoint) error {
	if fail == nil {
		return nil
	}

	switch *fail {
	case api.ClientAuthWrong:
		return clienterr.Unauthorized.New().
			WithDetailf("Some details xyz %d", 123).
			WithInstance("wrong_auth").
			Err()
	case api.ClientAuthMissing:
		return clienterr.Unauthorized.New().
			WithDetailf("Some details xyz %d", 123).
			WithInstance("missing_auth").
			Err()
	case api.InternalDuringAuth:
		return SomeAlgoFailed.New("Some internal algorithm failed xyz %d", 123)
	default:
		return nil
	}
}
