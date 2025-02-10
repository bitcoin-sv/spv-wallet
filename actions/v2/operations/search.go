package operations

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/operations/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// GetApiV2OperationsSearch return operations based on given filter parameters
func (s *APIOperations) GetApiV2OperationsSearch(c *gin.Context, params api.GetApiV2OperationsSearchParams) {
	logger := reqctx.Logger(c)

	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.AbortWithErrorResponse(c, err, logger)
		return
	}

	page := mapToFilter(params)
	pagedResult, err := reqctx.Engine(c).OperationsService().PaginatedForUser(c.Request.Context(), userID, page)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, mapping.OperationsPagedResponse(pagedResult))
}

func mapToFilter(params api.GetApiV2OperationsSearchParams) filter.Page {
	page := filter.Page{}

	if params.Page != nil {
		page.Number = *params.Page
	}
	if params.Size != nil {
		page.Size = *params.Size
	}
	if params.Sort != nil {
		page.Sort = *params.Sort
	}
	if params.SortBy != nil {
		page.SortBy = *params.SortBy
	}

	return page
}
