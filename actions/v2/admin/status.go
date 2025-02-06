package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server represents server with API endpoints
type Server struct{}

// GetApiV2AdminStatus return the status of the server only after admin authentication
func (s *Server) GetApiV2AdminStatus(c *gin.Context) {
	c.Status(http.StatusOK)
}
