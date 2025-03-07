package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// AdminUpdateContact updates a contact.
func (s *APIAdminContacts) AdminUpdateContact(c *gin.Context, id uint) {
	var requestBody api.RequestsUpdateContact
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
		return
	}

	contact, err := s.contactsService.UpdateFullNameByID(c, id, requestBody.FullName)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(http.StatusOK, mapping.MapToContactContract(contact))
}
