package users

import (
	"github.com/bitcoin-sv/spv-wallet/errdef/clienterr"
	"github.com/joomcode/errorx"
	"net/http"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	configerrors "github.com/bitcoin-sv/spv-wallet/config/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailerrors"
	"github.com/gin-gonic/gin"
)

// CreateUser creates a new user
func (s *APIAdminUsers) CreateUser(c *gin.Context) {
	var request api.RequestsCreateUser

	if err := c.Bind(&request); err != nil {
		// TODO: Bind does AbortWithError internally, so we should not call Response, I guess
		clienterr.UnprocessableEntity.
			Wrap(err, "cannot bind request").
			Response(c, s.logger)
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
		if errorx.IsOfType(err, configerrors.UnsupportedDomain) {
			clienterr.BadRequest.
				Wrap(err, "Unsupported domain").
				Response(c, s.logger)
		} else if errorx.IsOfType(err, paymailerrors.InvalidAvatarURL) {
			clienterr.UnprocessableEntity.
				Wrap(err, "Invalid avatar url").
				Response(c, s.logger)
		} else {
			clienterr.Response(c, err, s.logger)
		}
		return
	}

	c.JSON(http.StatusCreated, mapping.UserToResponse(createdUser))
}

func validatePubKey(pubKey string) error {
	_, err := primitives.PublicKeyFromString(pubKey)
	if err != nil {
		return clienterr.BadRequest.
			Wrap(err, "Cannot parse public key: '%s'", pubKey).Err()
	}
	return nil
}
