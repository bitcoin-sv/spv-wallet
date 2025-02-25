package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeleteContact deletes a contact
func (s *APIAdminContacts) DeleteContact(c *gin.Context, id int) {
	err := s.engine.ContactService().RemoveContactByID(c.Request.Context(), uint(id))
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrDeleteContact.WithTrace(err), s.logger)
		return
	}

	c.Status(http.StatusNoContent)
}
