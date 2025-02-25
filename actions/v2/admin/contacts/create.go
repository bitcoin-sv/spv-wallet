package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *APIAdminContacts) CreateContact(c *gin.Context, paymail string) {
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

	contact, err := s.engine.ContactService().AdminCreateContact(c.Request.Context(), newContact)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(http.StatusOK, mapping.MapToContactContract(contact))
}
