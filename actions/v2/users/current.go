package users

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// GetApiV2UsersCurrent returns current user information
func (s *APIUsers) GetApiV2UsersCurrent(c *gin.Context) {
	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	satoshis, err := reqctx.Engine(c).Repositories().Users.GetBalance(c.Request.Context(), userID, "bsv")
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, &response.UserInfo{
		CurrentBalance: satoshis,
	})
}
