package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// count will fetch a count of contacts filtered on conditions and metadata
// Count of transactions godoc
// @Summary		Count of contacts
// @Description	Count of contacts
// @Tags		Contacts
// @Produce		json
// @Param		CountContacts body CountContacts false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of contacts"
// @Failure		400	"Bad request - Error while parsing CountRequestParameters from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of contacts"
// @Router		/v1/contact/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams CountContacts
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var count int64
	count, err := a.Services.SpvWalletEngine.GetContactsByXPubIDCount(
		c.Request.Context(),
		reqXPubID,
		reqParams.Metadata,
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
