package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// create will make a new model using the services defined in the action object
// Create xPub godoc
// @Summary		Create xPub
// @Description	Create xPub
// @Tags		Admin
// @Produce		json
// @Param		CreateXpub body CreateXpub true " "
// @Success		201 {object} response.Xpub "Created Xpub"
// @Failure		400	"Bad request - Error while parsing CreateXpub from request body"
// @Failure 	500	"Internal server error - Error while creating xpub"
// @Router		/api/v1/admin/users [post]
// @Security	x-auth-xpub
func xpubsCreate(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var requestBody CreateXpub
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.WithTrace(err), logger)
		return
	}

	xPub, err := reqctx.Engine(c).NewXpub(
		c.Request.Context(), requestBody.Key,
		engine.WithMetadatas(requestBody.Metadata),
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrInvalidRequesterXpub.WithTrace(err), logger)
		return
	}

	contract := mappings.MapToXpubContract(xPub)
	c.JSON(http.StatusCreated, contract)
}

// xpubsSearch will fetch a list of xpubs filtered by metadata
// Search for xpubs filtering by metadata godoc
// @Summary		Search for xpubs
// @Description	Search for xpubs
// @Tags		Admin
// @Produce		json
// @Param		SearchXpubs body filter.XpubFilter false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []response.Xpub "List of xpubs"
// @Failure		400	"Bad request - Error while parsing SearchXpubs from request body"
// @Failure 	500	"Internal server error - Error while searching for xpubs"
// @Router 		/api/v1/admin/users [get]
// @Security	x-auth-xpub
func xpubsSearch(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)

	searchParams, err := query.ParseSearchParams[filter.XpubFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams, logger)
		return
	}

	xpubs, err := reqctx.Engine(c).GetXPubs(
		c.Request.Context(),
		mappings.MapToMetadata(searchParams.Metadata),
		searchParams.Conditions.ToDbConditions(),
		mappings.MapToDbQueryParams(&searchParams.Page),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	xpubContracts := make([]*response.Xpub, 0)
	for _, xpub := range xpubs {
		xpubContracts = append(xpubContracts, mappings.MapToXpubContract(xpub))
	}

	c.JSON(http.StatusOK, xpubContracts)
}
