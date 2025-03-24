package users

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	configerrors "github.com/bitcoin-sv/spv-wallet/config/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailerrors"
	"github.com/bitcoin-sv/spv-wallet/errdef/clienterr"
	"github.com/gin-gonic/gin"
)

// AddPaymailToUser add paymails to the user
func (s *APIAdminUsers) AddPaymailToUser(c *gin.Context, id string) {
	var request api.RequestsAddPaymail
	if err := c.Bind(&request); err != nil {
		clienterr.UnprocessableEntity.New().Wrap(err).Response(c, s.logger)
		return
	}

	newPaymail, err := mapping.RequestAddPaymailToNewPaymailModel(&request, id)
	if err != nil {
		clienterr.Response(c, err, s.logger)
		return
	}

	createdPaymail, err := s.engine.PaymailsService().Create(c, newPaymail)

	if err != nil {
		clienterr.Map(err).
			IfOfType(configerrors.UnsupportedDomain).
			Then(
				clienterr.BadRequest.Detailed("unsupported_domain", "Unsupported domain: '%s'", newPaymail.Domain),
			).
			IfOfType(paymailerrors.InvalidAvatarURL).
			Then(
				clienterr.UnprocessableEntity.Detailed("invalid_avatar_url", "Invalid avatar URL: '%s'", newPaymail.Avatar),
			).
			Response(c, s.logger)
		return
	}

	c.JSON(http.StatusCreated, mapping.PaymailToAdminResponse(createdPaymail))
}
