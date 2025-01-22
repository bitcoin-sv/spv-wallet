package operations

import (
	"net/http"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/models/response"
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

	pagedResult, err := reqctx.Engine(c).Repositories().Operations.PaginatedForUser(c.Request.Context(), userID, searchParams.Page)
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, response.PageModel[response.Operation]{
		Page: pagedResult.PageDescription,
		Content: slices.AppendSeq(
			make([]*response.Operation, 0, len(pagedResult.Content)),
			func(yield func(operation *response.Operation) bool) {
				for _, operation := range pagedResult.Content {
					yield(&response.Operation{
						CreatedAt:    operation.CreatedAt,
						Value:        operation.Value,
						TxID:         operation.TxID,
						Type:         operation.Type,
						Counterparty: operation.Counterparty,
					})
				}
			}),
	})
}
