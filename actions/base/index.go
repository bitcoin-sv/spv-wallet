package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// index basic request to /
func index(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{"message": "Welcome to the SPV Wallet ✌(◕‿-)✌"})
}
