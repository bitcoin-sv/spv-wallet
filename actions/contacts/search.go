package contacts

import (
	"net/http"
	"strconv"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// Search will fetch a list of contacts
// Get contacts godoc
// @Summary		Search contacts
// @Description	Search contacts
// @Tags		Contact
// @Produce		json
// @Param		SearchContacts body SearchContacts false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Contact "List of contacts"
// @Failure		400	"Bad request - Error while parsing SearchContacts from request body"
// @Failure 	500	"Internal server error - Error while searching for contacts"
// @Router		/v1/contact/search [POST]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	addCount, err := strconv.ParseBool(c.DefaultQuery("count", "true"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var reqParams SearchContacts
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	conditions, err := reqParams.Conditions.ToDbConditions()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if addCount && reqParams.QueryParams == nil {
		reqParams.QueryParams = common.LoadDefaultQueryParams()
	}

	contacts, err := a.Services.SpvWalletEngine.GetContactsByXpubID(
		c.Request.Context(),
		reqXPubID,
		reqParams.Metadata,
		conditions,
		reqParams.QueryParams,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contracts := mappings.MapToContactContracts(contacts)

	if !addCount {
		c.JSON(http.StatusOK, common.WrapBasicSearchResponse(contracts, len(contracts)))
		return
	}

	count, err := a.Services.SpvWalletEngine.GetContactsByXPubIDCount(
		c.Request.Context(),
		reqXPubID,
		reqParams.Metadata,
		conditions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, common.WrapCountResponse(contracts, count, reqParams.QueryParams))
}
