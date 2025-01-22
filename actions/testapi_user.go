package actions

import (
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// UserServer is the implementation of the server oapi-codegen's interface
type UserServer struct{}

// GetTestapi is a handler for the GET /testapi endpoint.
func (UserServer) GetTestapi(c *gin.Context, request api.GetTestapiParams) {
	user := reqctx.GetUserContext(c)

	c.JSON(200, api.ExResponse{
		XpubID:    user.GetXPubID(),
		Something: request.Something,
	})
}

func (UserServer) PostTestpolymorphism(c *gin.Context) {
	user := reqctx.GetUserContext(c)

	processes := make([]string, 0)

	var reqBody api.TransactionSpecification

	// Bind JSON request body to struct
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if reqBody.Outputs == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Outputs are required"})
		return
	}

	// Process each Output
	for _, output := range *reqBody.Outputs {
		value, err := output.ValueByDiscriminator()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown output type: " + err.Error()})
			return
		}

		// Handle each specific type
		switch v := value.(type) {
		case api.OpReturnOutput:
			processes = append(processes, v.Data...)
		case api.PaymailOutput:
			processes = append(processes, fmt.Sprintf("%s", v.To))
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unhandled output type"})
			return
		}
	}

	c.JSON(200, api.ExResponse{
		XpubID:    user.GetXPubID(),
		Something: strings.Join(processes, ", "),
	})
}
