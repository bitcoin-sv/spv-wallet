package users

import (
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

func (s *APIUsers) GetApiV2UsersCurrent(c *gin.Context) {
	userContext := reqctx.GetUserContext(c)
	fmt.Println(userContext)
	fmt.Println("xpubID", userContext.GetXPubID())
	fmt.Println("xpub", userContext.GetXPubObj())
	fmt.Println("authType", userContext.GetAuthType())

	userID, err := userContext.ShouldGetUserID()
	fmt.Println(userID)
	fmt.Println(err)
	if err != nil {
		spverrors.AbortWithErrorResponse(c, err, reqctx.Logger(c))
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
