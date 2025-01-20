package users

import (
	"net/http"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	adminerrors "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/errors"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/user/usermodels"
	"github.com/bitcoin-sv/spv-wallet/models/request/adminrequest"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func create(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)

	var requestBody adminrequest.CreateUser
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), logger)
		return
	}

	if invalidPubKey(requestBody.PublicKey) {
		spverrors.ErrorResponse(c, adminerrors.ErrInvalidPublicKey, logger)
		return
	}

	newUser := &usermodels.NewUser{
		PublicKey: requestBody.PublicKey,
	}
	if requestBody.PaymailDefined() {
		alias, domain, err := parsePaymail(requestBody.Paymail)
		if err != nil {
			spverrors.ErrorResponse(c, err, logger)
			return
		}

		if err = checkDomain(c, domain); err != nil {
			spverrors.ErrorResponse(c, err, logger)
			return
		}

		newUser.Paymail = &usermodels.NewPaymail{
			Alias:  alias,
			Domain: domain,

			PublicName: requestBody.Paymail.PublicName,
			Avatar:     requestBody.Paymail.Avatar,
		}
	}

	createdUser, err := reqctx.Engine(c).UserService().CreateUser(c, newUser)
	if err != nil {
		spverrors.ErrorResponse(c, adminerrors.ErrCreatingUser.Wrap(err), logger)
		return
	}

	c.JSON(http.StatusCreated, mapping.CreatedUserResponse(createdUser))
}

func invalidPubKey(pubKey string) bool {
	_, err := primitives.PublicKeyFromString(pubKey)
	return err != nil
}
