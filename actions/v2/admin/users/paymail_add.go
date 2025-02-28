package users

import (
	"net/http"

	adminerrors "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/errors"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	configerrors "github.com/bitcoin-sv/spv-wallet/config/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailerrors"
	"github.com/gin-gonic/gin"
)

// AddPaymailToUser add paymails to the user
func (s *APIAdminUsers) AddPaymailToUser(c *gin.Context, id string) {
	var request api.RequestsAddPaymail
	if err := c.Bind(&request); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
		return
	}

	newPaymail, err := mapping.RequestAddPaymailToNewPaymailModel(&request, id)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	// map response to correct error
	createdPaymail, err := s.engine.PaymailsService().Create(c, newPaymail)
	if err != nil {
		spverrors.MapResponse(c, err, s.logger).
			If(configerrors.ErrUnsupportedDomain).Then(adminerrors.ErrInvalidDomain).
			If(paymailerrors.ErrInvalidAvatarURL).Then(adminerrors.ErrInvalidAvatarURL).
			Else(adminerrors.ErrAddingPaymail)
		return
	}

	c.JSON(http.StatusCreated, mapping.PaymailToAdminResponse(createdPaymail))
}
