package paymails

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bsv-blockchain/spv-wallet/actions/common"
	"github.com/bsv-blockchain/spv-wallet/engine/spverrors"
	"github.com/bsv-blockchain/spv-wallet/internal/query"
	"github.com/bsv-blockchain/spv-wallet/mappings"
	"github.com/bsv-blockchain/spv-wallet/models/filter"
	"github.com/bsv-blockchain/spv-wallet/models/response"
	"github.com/bsv-blockchain/spv-wallet/server/reqctx"
)

// paymailAddressesSearch will fetch a list of paymail addresses filtered by metadata
// Paymail addresses search by metadata
// @Summary		Paymail addresses search
// @Description	Paymail addresses search
// @Tags		Users
// @Produce		json
// @Param		SearchPaymails body filter.PaymailFilter false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []response.PaymailAddress "List of paymail addresses"
// @Failure		400	"Bad request - Error while parsing SearchPaymails from request body"
// @Failure 	500	"Internal server error - Error while searching for paymail addresses"
// @Router		/api/v1/paymails [get]
// @Security	x-auth-xpub
func paymailAddressesSearch(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)

	searchParams, err := query.ParseSearchParams[filter.PaymailFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams.WithTrace(err), logger)
		return
	}

	conditions := searchParams.Conditions.ToDbConditions()
	metadata := mappings.MapToMetadata(searchParams.Metadata)
	pageOptions := mappings.MapToDbQueryParams(&searchParams.Page)

	paymailAddresses, err := reqctx.Engine(c).GetPaymailAddressesByXPubID(
		c.Request.Context(),
		userContext.GetXPubID(),
		metadata,
		conditions,
		pageOptions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindPaymail.WithTrace(err), logger)
		return
	}

	count, err := reqctx.Engine(c).GetPaymailAddressesCount(c.Request.Context(), metadata, conditions)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindPaymail.WithTrace(err), logger)
		return
	}

	paymailAddressContracts := common.MapToTypeContracts(paymailAddresses, mappings.MapToPaymailContract)

	result := response.PageModel[response.PaymailAddress]{
		Content: paymailAddressContracts,
		Page:    common.GetPageDescriptionFromSearchParams(pageOptions, count),
	}
	c.JSON(http.StatusOK, result)
}
