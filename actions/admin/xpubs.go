package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
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
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	xPub, err := a.Services.SpvWalletEngine.NewXpub(
		c.Request.Context(), requestBody.Key,
		engine.WithMetadatas(requestBody.Metadata),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
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
// @Param		SearchRequestParameters body actions.SearchRequestParameters false "Supports targeted resource searches with filters for metadata and custom conditions, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Xpub "List of xpubs"
// @Failure		400	"Bad request - Error while parsing SearchRequestParameters from request body"
// @Failure 	500	"Internal server error - Error while searching for xpubs"
// @Router		/v1/admin/xpubs/search [post]
// @Security	x-auth-xpub
func (a *Action) xpubsSearch(c *gin.Context) {
	queryParams, metadata, conditions, err := actions.GetSearchQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var xpubs []*engine.Xpub
	if xpubs, err = a.Services.SpvWalletEngine.GetXPubs(
		c.Request.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
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
// @Param		CountRequestParameters body actions.CountRequestParameters false "Enables precise filtering of resource counts using custom conditions or metadata, catering to specific business or analysis needs"
// @Success		200	{number} int64 "Count of access keys"
// @Failure		400	"Bad request - Error while parsing CountRequestParameters from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of xpubs"
// @Router		/v1/admin/xpubs/count [post]
// @Security	x-auth-xpub
func (a *Action) xpubsCount(c *gin.Context) {
	metadata, conditions, err := actions.GetCountQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.SpvWalletEngine.GetXPubsCount(
		c.Request.Context(),
		metadata,
		conditions,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
