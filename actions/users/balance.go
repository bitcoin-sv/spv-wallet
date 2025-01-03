package users

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func balance(c *gin.Context, userContext *reqctx.UserContext) {
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, reqctx.Logger(c))
		return
	}

	satoshi, err := reqctx.Engine(c).UsersDAO().GetBalance(c.Request.Context(), userID, "bsv")
	if err != nil {
		spverrors.ErrorResponse(c, err, reqctx.Logger(c))
		return
	}

	c.JSON(http.StatusOK, &response.Balance{
		CurrentBalance: satoshi,
	})
}
