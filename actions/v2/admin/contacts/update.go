package contacts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/bsv-blockchain/spv-wallet/actions/v2/internal/mapping"
	"github.com/bsv-blockchain/spv-wallet/api"
	"github.com/bsv-blockchain/spv-wallet/errdef/clienterr"
)

// AdminUpdateContact updates a contact.
func (s *APIAdminContacts) AdminUpdateContact(c *gin.Context, id uint) {
	var requestBody api.RequestsUpdateContact
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		clienterr.UnprocessableEntity.New().Wrap(err).Response(c, s.logger)
		return
	}

	contact, err := s.contactsService.UpdateFullNameByID(c.Request.Context(), id, requestBody.FullName)
	if err != nil {
		clienterr.Response(c, err, s.logger)
		return
	}

	c.JSON(http.StatusOK, mapping.MapToContactContract(contact))
}
