package contacts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bsv-blockchain/spv-wallet/errdef/clienterr"
)

// AdminDeleteContact deletes a contact
func (s *APIAdminContacts) AdminDeleteContact(c *gin.Context, id uint) {
	err := s.contactsService.RemoveContactByID(c.Request.Context(), id)
	if err != nil {
		clienterr.Response(c, err, s.logger)
		return
	}

	c.Status(http.StatusOK)
}
