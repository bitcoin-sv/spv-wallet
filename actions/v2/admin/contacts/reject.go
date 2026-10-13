package contacts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bsv-blockchain/spv-wallet/errdef/clienterr"
)

// AdminRejectInvitation rejects an invitation from a contact.
func (s *APIAdminContacts) AdminRejectInvitation(c *gin.Context, id uint) {
	_, err := s.contactsService.RejectContactByID(c.Request.Context(), id)
	if err != nil {
		clienterr.Response(c, err, s.logger)
		return
	}

	c.Status(http.StatusOK)
}
