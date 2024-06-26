package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO add swagger description
func (a *Action) subscribeWebhook(c *gin.Context) {
	requestModel := struct {
		URL string
	}{}
	if err := c.Bind(&requestModel); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	a.Services.SpvWalletEngine.SubscribeWebhook(c.Request.Context(), requestModel.URL, "", "")

	c.JSON(http.StatusOK, true)
}
