package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminServer represents server with API endpoints
type AdminServer struct{}

// GetApiV2AdminStatus return the status of the server only after admin authentication
func (s *AdminServer) GetApiV2AdminStatus(c *gin.Context) {
	c.Status(http.StatusOK)
}
