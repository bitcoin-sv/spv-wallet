package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AdminAcceptInvitation accepts an invitation from a contact.
func (s *APIAdminContacts) AdminAcceptInvitation(c *gin.Context, id int) {
	contact, err := s.engine.ContactService().AcceptContactByID(c.Request.Context(), uint(id))
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	res := mapping.MapToContactContract(contact)

	c.JSON(http.StatusOK, res)
}
