package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *APIAdminContacts) RejectInvitation(c *gin.Context, id int) {
	_, err := s.engine.ContactService().RejectContactByID(c.Request.Context(), uint(id))
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.Status(http.StatusOK)
}
