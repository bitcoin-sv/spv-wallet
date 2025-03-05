package users

import (
	"github.com/bitcoin-sv/spv-wallet/errdef/clienterr"
	"github.com/joomcode/errorx"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	configerrors "github.com/bitcoin-sv/spv-wallet/config/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailerrors"
	"github.com/gin-gonic/gin"
)

// AddPaymailToUser add paymails to the user
func (s *APIAdminUsers) AddPaymailToUser(c *gin.Context, id string) {
	var request api.RequestsAddPaymail
	if err := c.Bind(&request); err != nil {
		clienterr.UnprocessableEntity.Wrap(err, "cannot bind request").Response(c, s.logger)
		return
	}

	newPaymail, err := mapping.RequestAddPaymailToNewPaymailModel(&request, id)
	if err != nil {
		clienterr.Response(c, err, s.logger)
		return
	}

	createdPaymail, err := s.engine.PaymailsService().Create(c, newPaymail)
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

	c.JSON(http.StatusCreated, mapping.PaymailToAdminResponse(createdPaymail))
}
