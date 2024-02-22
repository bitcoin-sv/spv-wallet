package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// index basic request to /
func index(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{"message": "Welcome to the SPV Wallet ✌(◕‿-)✌"})
}
