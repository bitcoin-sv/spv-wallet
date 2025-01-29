package operations

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/operations/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func search(c *gin.Context, userContext *reqctx.UserContext) {
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, reqctx.Logger(c))
		return
	}

	logger := reqctx.Logger(c)

	searchParams, err := query.ParseSearchParams[struct{}](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams.WithTrace(err), logger)
		return
	}

	pagedResult, err := reqctx.Engine(c).OperationsService().PaginatedForUser(c.Request.Context(), userID, searchParams.Page)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, mapping.OperationsPagedResponse(pagedResult))
}
