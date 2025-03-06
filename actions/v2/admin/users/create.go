package users

import (
	"net/http"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	adminerrors "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/errors"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	configerrors "github.com/bitcoin-sv/spv-wallet/config/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailerrors"
	"github.com/gin-gonic/gin"
)

// CreateUser creates a new user
func (s *APIAdminUsers) CreateUser(c *gin.Context) {
	var request api.RequestsCreateUser
	if err := c.Bind(&request); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
		return
	}

	if err := validatePubKey(request.PublicKey); err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	newUser, err := mapping.RequestCreateUserToNewUserModel(&request)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	createdUser, err := s.users.Create(c, newUser)
	if err != nil {
		spverrors.MapResponse(c, err, s.logger).
			If(configerrors.ErrUnsupportedDomain).Then(adminerrors.ErrInvalidDomain).
			If(paymailerrors.ErrInvalidAvatarURL).Then(adminerrors.ErrInvalidAvatarURL).
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
