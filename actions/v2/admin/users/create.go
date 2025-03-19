package users

import (
	"net/http"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	configerrors "github.com/bitcoin-sv/spv-wallet/config/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailerrors"
	"github.com/bitcoin-sv/spv-wallet/errdef/clienterr"
	"github.com/gin-gonic/gin"
)

// CreateUser creates a new user
func (s *APIAdminUsers) CreateUser(c *gin.Context) {
	var request api.RequestsCreateUser

	if err := c.Bind(&request); err != nil {
		// TODO: Bind does AbortWithError internally, so we should not call Response, I guess
		clienterr.UnprocessableEntity.New().Wrap(err).Response(c, s.logger)
		return
	}

	if err := validatePubKey(request.PublicKey); err != nil {
		clienterr.Response(c, err, s.logger)
		return
	}

	newUser, err := mapping.RequestCreateUserToNewUserModel(&request)
	if err != nil {
		clienterr.Response(c, err, s.logger)
		return
	}

	createdUser, err := s.engine.UsersService().Create(c, newUser)

	if err != nil {
		clienterr.Map(err).
			IfOfType(configerrors.UnsupportedDomain).
			Then(
				clienterr.BadRequest.Detailed("unsupported_domain", "Unsupported domain: '%s'", newUser.Paymail.Domain),
			).
			IfOfType(paymailerrors.InvalidAvatarURL).
			Then(
				clienterr.UnprocessableEntity.Detailed("invalid_avatar_url", "Invalid avatar URL: '%s'", newUser.Paymail.Avatar),
			).
			Response(c, s.logger)
		return
	}

	c.JSON(http.StatusCreated, mapping.UserToResponse(createdUser))
}

func validatePubKey(pubKey string) error {
	_, err := primitives.PublicKeyFromString(pubKey)
	if err != nil {
		return clienterr.BadRequest.
			Detailed("invalid_public_key", "Invalid public key: '%s'", pubKey).
			Wrap(err).Err()
	}
	return nil
}
