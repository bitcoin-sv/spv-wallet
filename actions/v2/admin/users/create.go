package users

import (
	"net/http"
	"slices"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	adminerrors "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/errors"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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

	newUser := &database.User{
		PubKey: requestBody.PublicKey,
	}
	if requestBody.PaymailDefined() {
		alias, domain, err := parsePaymail(requestBody.Paymail)
		if err != nil {
			spverrors.ErrorResponse(c, err, logger)
			return
		}

		config := reqctx.AppConfig(c)
		if config.Paymail.DomainValidationEnabled {
			if !slices.Contains(config.Paymail.Domains, domain) {
				spverrors.ErrorResponse(c, spverrors.ErrInvalidDomain, logger)
				return
			}
		}

		newUser.Paymails = append(newUser.Paymails, &database.Paymail{
			Alias:  alias,
			Domain: domain,

			PublicName: requestBody.Paymail.PublicName,
			Avatar:     requestBody.Paymail.Avatar,
		})
	}
	if err := reqctx.Engine(c).Repositories().Users.Save(c, newUser); err != nil {
		spverrors.ErrorResponse(c, adminerrors.ErrCreatingUser.Wrap(err), logger)
		return
	}

	c.JSON(http.StatusCreated, mapping.UserResponse(newUser))
}

func invalidPubKey(pubKey string) bool {
	_, err := primitives.PublicKeyFromString(pubKey)
	return err != nil
}
