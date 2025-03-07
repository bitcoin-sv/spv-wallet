package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/gin-gonic/gin"
)

// AdminCreateContact creates a new contact for a user.
func (s *APIAdminContacts) AdminCreateContact(c *gin.Context, paymail string) {
	var req api.RequestsAdminCreateContact
	if err := c.Bind(&req); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.WithTrace(err), s.logger)
		return
	}

	newContact := contactsmodels.NewContact{
		FullName:          req.FullName,
		NewContactPaymail: paymail,
		RequesterPaymail:  req.CreatorPaymail,
		Status:            contactsmodels.ContactNotConfirmed,
	}

	contact, err := s.contactsService.AdminCreateContact(c, newContact)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(http.StatusOK, mapping.MapToContactContract(contact))
}
