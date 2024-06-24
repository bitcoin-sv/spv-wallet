package admin

import (
	spverrors2 "github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/gin-gonic/gin"
)

// create will make a new model using the services defined in the action object
// Create xPub godoc
// @Summary		Create xPub
// @Description	Create xPub
// @Tags		Admin
// @Produce		json
// @Param		CreateXpub body CreateXpub true " "
// @Success		201 {object} models.Xpub "Created Xpub"
// @Failure		400	"Bad request - Error while parsing CreateXpub from request body"
// @Failure 	500	"Internal server error - Error while creating xpub"
// @Router		/v1/admin/xpub [post]
// @Security	x-auth-xpub
func (a *Action) xpubsCreate(c *gin.Context) {
	var requestBody CreateXpub
	if err := c.Bind(&requestBody); err != nil {
		spverrors2.ErrorResponse(c, spverrors2.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	xPub, err := a.Services.SpvWalletEngine.NewXpub(
		c.Request.Context(), requestBody.Key,
		engine.WithMetadatas(requestBody.Metadata),
	)
	if err != nil {
		spverrors2.ErrorResponse(c, err, a.Services.Logger)
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
// @Param		SearchXpubs body filter.SearchXpubs false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Xpub "List of xpubs"
// @Failure		400	"Bad request - Error while parsing SearchXpubs from request body"
// @Failure 	500	"Internal server error - Error while searching for xpubs"
// @Router		/v1/admin/xpubs/search [post]
// @Security	x-auth-xpub
func (a *Action) xpubsSearch(c *gin.Context) {
	var reqParams filter.SearchXpubs
	if err := c.Bind(&reqParams); err != nil {
		spverrors2.ErrorResponse(c, spverrors2.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	xpubs, err := a.Services.SpvWalletEngine.GetXPubs(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors2.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	xpubContracts := make([]*models.Xpub, 0)
	for _, xpub := range xpubs {
		xpubContracts = append(xpubContracts, mappings.MapToXpubContract(xpub))
	}

	c.JSON(http.StatusOK, xpubContracts)
}

// xpubsCount will count all xpubs filtered by metadata
// Count xpubs filtering by metadata godoc
// @Summary		Count xpubs
// @Description	Count xpubs
// @Tags		Admin
// @Produce		json
// @Param		CountXpubs body filter.CountXpubs false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of access keys"
// @Failure		400	"Bad request - Error while parsing CountXpubs from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of xpubs"
// @Router		/v1/admin/xpubs/count [post]
// @Security	x-auth-xpub
func (a *Action) xpubsCount(c *gin.Context) {
	var reqParams filter.CountXpubs
	if err := c.Bind(&reqParams); err != nil {
		spverrors2.ErrorResponse(c, spverrors2.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	count, err := a.Services.SpvWalletEngine.GetXPubsCount(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		spverrors2.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	c.JSON(http.StatusOK, count)
}
