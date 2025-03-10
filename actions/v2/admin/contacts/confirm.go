package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
)

// AdminConfirmContact confirms a contact between two users.
func (s *APIAdminContacts) AdminConfirmContact(c *gin.Context) {
	var reqParams *api.RequestsAdminConfirmContact
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.WithTrace(err), s.logger)
		return
	}

	if err := s.contactsService.AdminConfirmContacts(c, reqParams.PaymailA, reqParams.PaymailB); err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.Status(http.StatusOK)

}
