package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/gin-gonic/gin"
)

// create will make a new model using the services defined in the action object
// Create xPub godoc
// @Summary		Create xPub
// @Description	Create xPub
// @Tags		xPub
// @Produce		json
// @Param		key query string true "key"
// @Param		metadata query string false "metadata"
// @Success		201
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
		c.JSON(http.StatusUnprocessableEntity, err.Error())
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
// @Param		page query int false "page"
// @Param		page_size query int false "page_size"
// @Param		order_by_field query string false "order_by_field"
// @Param		sort_direction query string false "sort_direction"
// @Param		metadata query string false "Metadata filter"
// @Param		conditions query string false "Conditions filter"
// @Success		200
// @Router		/v1/admin/xpubs/search [post]
// @Security	x-auth-xpub
func (a *Action) xpubsSearch(c *gin.Context) {
	queryParams, metadata, conditions, err := actions.GetQueryParameters(c)
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

	c.JSON(http.StatusOK, xpubs)
}

// xpubsCount will count all xpubs filtered by metadata
// Count xpubs filtering by metadata godoc
// @Summary		Count xpubs
// @Description	Count xpubs
// @Tags		Admin
// @Produce		json
// @Param		metadata query string false "Metadata filter"
// @Param		conditions query string false "Conditions filter"
// @Success		200
// @Router		/v1/admin/xpubs/count [post]
// @Security	x-auth-xpub
func (a *Action) xpubsCount(c *gin.Context) {
	_, metadata, conditions, err := actions.GetQueryParameters(c)
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
