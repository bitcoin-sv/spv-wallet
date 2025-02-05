package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminServer struct{}

func (s *AdminServer) GetApiV2AdminStatus(c *gin.Context) {
	c.Status(http.StatusOK)
}
