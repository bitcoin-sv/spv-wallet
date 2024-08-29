package middleware

import (
	"strings"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware will check the request for the xPub or AccessKey header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		xPub := strings.TrimSpace(c.GetHeader(models.AuthHeader))
		authAccessKey := strings.TrimSpace(c.GetHeader(models.AuthAccessKey))

		userContext, err := tryAuth(c, xPub, authAccessKey)

		if err == nil {
			c.Set(reqctx.UserContextKey, userContext)
			c.Next()
		} else {
			spverrors.AbortWithErrorResponse(c, err, reqctx.Logger(c))
		}
	}
}

func tryAuth(c *gin.Context, xPub, authAccessKey string) (*reqctx.UserContext, error) {
	if xPub != "" {
		return authByXPub(c, xPub)
	}
	if authAccessKey != "" {
		return authByAccessKey(c, authAccessKey)
	}
	return nil, spverrors.ErrMissingAuthHeader
}

func authByXPub(c *gin.Context, xPub string) (*reqctx.UserContext, error) {
	if _, err := utils.ValidateXPub(xPub); err != nil {
		return nil, spverrors.ErrAuthorization
	}
	config := reqctx.AppConfig(c)
	if xPub == config.Authentication.AdminKey {
		return reqctx.NewUserContextAsAdmin(xPub), nil
	}

	xPubID := utils.Hash(xPub)
	xPubObj, err := getXPubByID(c, xPubID)
	if err != nil {
		return nil, err
	}

	return reqctx.NewUserContextWithXPub(xPub, xPubID, xPubObj), nil
}

func authByAccessKey(c *gin.Context, authAccessKey string) (*reqctx.UserContext, error) {
	accessKey, err := reqctx.Engine(c).AuthenticateAccessKey(c, utils.Hash(authAccessKey))
	if err != nil || accessKey == nil {
		return nil, spverrors.ErrAuthorization
	}

	xPubObj, err := getXPubByID(c, accessKey.XpubID)
	if err != nil {
		return nil, err
	}

	return reqctx.NewUserContextWithAccessKey(authAccessKey, accessKey, xPubObj), nil
}

func getXPubByID(c *gin.Context, xPubID string) (*engine.Xpub, error) {
	xPubObj, err := reqctx.Engine(c).GetXpubByID(c, xPubID)
	if err != nil || xPubObj == nil {
		if err != nil {
			reqctx.Logger(c).Warn().Msgf("Could not get XPub by ID: %v", err)
		} else if xPubObj == nil {
			reqctx.Logger(c).Debug().Msgf("Provided XPubID (%v) doesn't exist", xPubID)
		}

		return nil, spverrors.ErrAuthorization
	}
	return xPubObj, nil
}
