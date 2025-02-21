package exampleactions

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/clienterr"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/domain/exampledomain"
	domainerr "github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/domain/exampledomain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/errdef"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/mockmiddleware"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/repos"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/joomcode/errorx"
)

var SpecificTypeUnprocessableEntity = clienterr.RegisterSubtype(
	clienterr.UnprocessableEntity,
	"cannot_parse_xyz",
	"Specific type of unprocessable entity",
)

type ExampleAPI struct{}

func (s *ExampleAPI) GetApiV2TestErrors(c *gin.Context, params api.GetApiV2TestErrorsParams) {
	log := reqctx.Logger(c)
	repo := repos.NewRepo()
	domain := exampledomain.NewService(repo)

	if err := mockmiddleware.Auth(params.Fail); err != nil {
		clienterr.Response(c, err, log)
		return
	}

	if params.Fail != nil && *params.Fail == api.ClientWrongInput {
		clienterr.UnprocessableEntity.New().
			WithDetailf("User provided wrong input %d", 123).
			Response(c, log)
		return
	}

	_, err := domain.DoSth(params.Fail)
	if err != nil {
		if errorx.HasTrait(err, errdef.TraitIllegalArgument) {
			// sometimes we want to say that the user provided wrong input
			// based on error from the domain/repository layer
			err = SpecificTypeUnprocessableEntity.
				Wrap(err, "User provided wrong data and it was checked at some lower level")
		} else if errorx.IsOfType(err, domainerr.SomeARCError) {
			err = clienterr.BadRequest.
				Wrap(err, "ARC returned an error but because of client's fault")
		}
		clienterr.Response(c, err, log)
		return
	}

	c.JSON(200, "success") // not relevant for the example
}
