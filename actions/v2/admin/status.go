package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetApiV2AdminStatus return the status of the server only after admin authentication
func (s *APIAdmin) GetApiV2AdminStatus(c *gin.Context) {
	c.Status(http.StatusOK)
}
