package users

import (
	"net/http"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	adminerrors "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/errors"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	configerrors "github.com/bitcoin-sv/spv-wallet/config/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/bitcoin-sv/spv-wallet/models/request/adminrequest"
	"github.com/gin-gonic/gin"
)

// PostApiV2AdminUsers creates a new user
func (s *APIAdminUsers) PostApiV2AdminUsers(c *gin.Context) {
	var requestBody adminrequest.CreateUser
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
		return
	}

	if err := validatePubKey(requestBody.PublicKey); err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	newUser := &usersmodels.NewUser{
		PublicKey: requestBody.PublicKey,
	}
	if requestBody.PaymailDefined() {
		alias, domain, err := parsePaymail(requestBody.Paymail)
		if err != nil {
			spverrors.ErrorResponse(c, err, s.logger)
			return
		}

		newUser.Paymail = &usersmodels.NewPaymail{
			Alias:  alias,
			Domain: domain,

			PublicName: requestBody.Paymail.PublicName,
			Avatar:     requestBody.Paymail.Avatar,
		}
	}

	createdUser, err := s.engine.UsersService().Create(c, newUser)
	if err != nil {
		spverrors.MapResponse(c, err, s.logger).
			If(configerrors.ErrUnsupportedDomain).Then(adminerrors.ErrInvalidDomain).
			Else(adminerrors.ErrCreatingUser)
		return
	}

	c.JSON(http.StatusCreated, mapping.UserToResponse(createdUser))
}

func validatePubKey(pubKey string) error {
	_, err := primitives.PublicKeyFromString(pubKey)
	if err != nil {
		return adminerrors.ErrInvalidPublicKey.Wrap(err)
	}
	return nil
}
