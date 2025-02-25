package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *APIAdminContacts) GetContacts(c *gin.Context, params api.GetContactsParams) {
	page := mapping.MapContactsParamToFilterPage(params)
	conditions := mapping.MapToDBConditions(params)

	pagedResult, err := s.engine.ContactService().PaginatedForAdmin(c.Request.Context(), page, conditions)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(http.StatusOK, mapping.ContactsPagedResponse(pagedResult))
}
