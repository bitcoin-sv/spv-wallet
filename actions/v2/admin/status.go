package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *AdminServer) GetApiV2AdminStatus(c *gin.Context) {
	c.Status(http.StatusOK)
}
