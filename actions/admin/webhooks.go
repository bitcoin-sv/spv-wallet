package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO add swagger description
func (a *Action) subscribeWebhook(c *gin.Context) {
	requestModel := struct {
		URL         string `json:"url"`
		TokenHeader string `json:"tokenHeader"`
		TokenValue  string `json:"tokenValue"`
	}{}
	if err := c.Bind(&requestModel); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	a.Services.SpvWalletEngine.SubscribeWebhook(c.Request.Context(), requestModel.URL, requestModel.TokenHeader, requestModel.TokenValue)

	c.JSON(http.StatusOK, true)
}

// TODO add swagger description
func (a *Action) unsubscribeWebhook(c *gin.Context) {
	requestModel := struct {
		URL string `json:"url"`
	}{}
	if err := c.Bind(&requestModel); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	a.Services.SpvWalletEngine.UnsubscribeWebhook(c.Request.Context(), requestModel.URL)

	c.JSON(http.StatusOK, true)
}
