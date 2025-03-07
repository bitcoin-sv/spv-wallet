package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
)

// AdminDeleteContact deletes a contact
func (s *APIAdminContacts) AdminDeleteContact(c *gin.Context, id uint) {
	err := s.contactsService.RemoveContactByID(c, id)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.Status(http.StatusOK)
}
