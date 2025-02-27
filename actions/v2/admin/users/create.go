package users

import (
	"net/http"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	adminerrors "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/errors"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	configerrors "github.com/bitcoin-sv/spv-wallet/config/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/bitcoin-sv/spv-wallet/lox"
	"github.com/gin-gonic/gin"
)

// CreateUser creates a new user
func (s *APIAdminUsers) CreateUser(c *gin.Context) {
	var request api.RequestsCreateUser
	if err := c.Bind(&request); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
		return
	}

	requestBody := &createUserRequest{&request}

	if err := validatePubKey(requestBody.PublicKey); err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	newUser := &usersmodels.NewUser{
		PublicKey: requestBody.PublicKey,
	}
	if requestBody.PaymailDefined() {
		alias, domain, err := parsePaymail(requestBody.PaymailRequestBody())
		if err != nil {
			spverrors.ErrorResponse(c, err, s.logger)
			return
		}

		newUser.Paymail = &usersmodels.NewPaymail{
			Alias:  alias,
			Domain: domain,

			PublicName: lox.Unwrap(requestBody.Paymail.PublicName).Else(""),
			Avatar:     lox.Unwrap(requestBody.Paymail.AvatarURL).Else(""),
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

type createUserRequest struct {
	*api.RequestsCreateUser
}

func (r *createUserRequest) PaymailDefined() bool {
	return r.RequestsCreateUser.Paymail != nil
}

func (r *createUserRequest) PaymailRequestBody() *addPaymailRequest {
	return &addPaymailRequest{r.RequestsCreateUser.Paymail}
}
