package middleware

import (
	"strings"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// AuthV2Middleware will check the request for the xPub and convert it to the user context.
func AuthV2Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		xPub := strings.TrimSpace(c.GetHeader(models.AuthHeader))

		userContext, err := tryAuthWithPubKey(c, xPub)

		if err == nil {
			reqctx.SetUserContext(c, userContext)
			c.Next()
		} else {
			spverrors.AbortWithErrorResponse(c, err, reqctx.Logger(c))
		}
	}
}

func tryAuthWithPubKey(c *gin.Context, xPub string) (*reqctx.UserContext, error) {
	config := reqctx.AppConfig(c)
	if xPub == "" {
		return nil, spverrors.ErrMissingAuthHeader
	}
	if xPub == config.Authentication.AdminKey {
		return reqctx.NewUserContextAsAdmin(), nil
	}

	hdKey, err := bip32.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		return nil, spverrors.ErrAuthorization.Wrap(err)
	}

	pubKey, err := bip32.GetPublicKeyFromHDKey(hdKey)
	if err != nil {
		return nil, spverrors.ErrAuthorization.Wrap(err)
	}
	pubKeyHex := pubKey.ToDERHex()

	userID, err := reqctx.Engine(c).UsersService().GetIDByPubKey(c.Request.Context(), pubKeyHex)
	if err != nil {
		return nil, spverrors.ErrAuthorization.Wrap(err)
	}

	return reqctx.NewUserContextWithPublicKeys(xPub, utils.Hash(xPub), pubKeyHex, userID), nil
}
