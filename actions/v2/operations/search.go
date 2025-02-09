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

	// TODO: some mapping is missing here
	//searchParams, err := query.ParseSearchParams[struct{}](c)
	//if err != nil {
	//	spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams.WithTrace(err), logger)
	//	return
	//}

	page := filter.Page{
		Number: *params.Page,
		Size:   *params.Size,
		Sort:   *params.Sort,
		SortBy: *params.SortBy,
	}

	pagedResult, err := reqctx.Engine(c).OperationsService().PaginatedForUser(c.Request.Context(), userID, page)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, mapping.OperationsPagedResponse(pagedResult))
}
