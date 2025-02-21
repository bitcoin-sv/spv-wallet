package mockmiddleware

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/clienterr"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/errdef"
)

var WrongAuth = clienterr.RegisterSubtype(
	clienterr.Unauthorized,
	"wrong_auth",
	"Client provided wrong API key",
)

var MissingAuth = clienterr.RegisterSubtype(
	clienterr.Unauthorized,
	"missing_auth",
	"Auth headers missing",
)

var MiddlewareErr = errdef.Internal.NewSubNamespace("middleware", errdef.TraitAuth)

var SomeAlgoFailed = MiddlewareErr.NewType("some_auth_algo_failed")

func Auth(fail *api.ModelsFailingPoint) error {
	if fail == nil {
		return nil
	}

	switch *fail {
	case api.ClientAuthWrong:
		return WrongAuth.New().
			WithDetailf("Some details xyz %d", 123).
			Err()
	case api.ClientAuthMissing:
		return MissingAuth.New().Err()
	case api.InternalDuringAuth:
		return SomeAlgoFailed.New("Some internal algorithm failed xyz %d", 123)
	default:
		return nil
	}
}
