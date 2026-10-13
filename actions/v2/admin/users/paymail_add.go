package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bsv-blockchain/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bsv-blockchain/spv-wallet/api"
	configerrors "github.com/bsv-blockchain/spv-wallet/config/errors"
	"github.com/bsv-blockchain/spv-wallet/engine/v2/paymails/paymailerrors"
	"github.com/bsv-blockchain/spv-wallet/errdef/clienterr"
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

	createdPaymail, err := s.engine.PaymailsService().Create(c.Request.Context(), newPaymail)
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
