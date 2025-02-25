package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *APIAdminContacts) ConfirmContact(c *gin.Context) {
	var reqParams *api.RequestsAdminConfirmContact
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.WithTrace(err), s.logger)
		return
	}

	if err := s.engine.ContactService().AdminConfirmContacts(c.Request.Context(), reqParams.PaymailA, reqParams.PaymailB); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrConfirmContact.WithTrace(err), s.logger)
		return
	}

	c.Status(http.StatusOK)

}
