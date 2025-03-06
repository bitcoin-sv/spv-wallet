package contacts

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/gin-gonic/gin"
)

// AdminGetContacts returns a list of contacts for the admin.
func (s *APIAdminContacts) AdminGetContacts(c *gin.Context, params api.AdminGetContactsParams) {
	page := mapContactsParamToFilterPage(params)
	conditions := mapping.MapToDBConditions(params)

	pagedResult, err := s.engine.ContactService().PaginatedForAdmin(c, page, conditions)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(http.StatusOK, mapping.ContactsPagedResponse(pagedResult))
}

func mapContactsParamToFilterPage(params api.AdminGetContactsParams) filter.Page {
	return filter.Page{
		Number: mapping.GetPointerValue(params.Page),
		Size:   mapping.GetPointerValue(params.Size),
		Sort:   mapping.GetPointerValue(params.Sort),
		SortBy: mapping.GetPointerValue(params.SortBy),
	}
}
