package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// status will return the status of the admin login
// Get status godoc
// @Summary		Get status
// @Description	Get status
// @Tags		Admin
// @Produce		json
// @Success		200
// @Router		/v1/admin/status [get]
// @Security	x-auth-xpub
func (a *Action) status(c *gin.Context) {
	c.JSON(http.StatusOK, true)
}
