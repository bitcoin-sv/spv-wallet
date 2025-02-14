package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminStatus return the status of the server only after admin authentication
func (s *APIAdmin) AdminStatus(c *gin.Context) {
	c.Status(http.StatusOK)
}
