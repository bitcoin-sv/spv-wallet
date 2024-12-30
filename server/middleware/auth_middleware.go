package middleware

import (
	"strings"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
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
			reqctx.SetUserContext(c, userContext)
			c.Next()
		} else {
			spverrors.AbortWithErrorResponse(c, err, reqctx.Logger(c))
		}
	}
}

func tryAuth(c *gin.Context, xPub, authAccessKey string) (*reqctx.UserContext, error) {
	config := reqctx.AppConfig(c)
	if config.ExperimentalFeatures.NewTransactionFlowEnabled {
		if xPub != "" {
			if xPub == config.Authentication.AdminKey {
				return reqctx.NewUserContextAsAdmin(), nil
			}
			return authByXPubToPublicKey(c, xPub)
		}
		return nil, spverrors.ErrMissingAuthHeader
	}

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
		return reqctx.NewUserContextAsAdmin(), nil
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

	return reqctx.NewUserContextWithAccessKey(accessKey.XpubID, xPubObj), nil
}

func getXPubByID(c *gin.Context, xPubID string) (*engine.Xpub, error) {
	xPubObj, err := reqctx.Engine(c).GetXpubByID(c, xPubID)
	if err != nil {
		reqctx.Logger(c).Warn().Msgf("Could not get XPub by ID: %v", err)
		return nil, spverrors.ErrAuthorization
	}
	if xPubObj == nil {
		reqctx.Logger(c).Debug().Msgf("Provided XPubID (%v) doesn't exist", xPubID)
		return nil, spverrors.ErrAuthorization
	}
	return xPubObj, nil
}

func authByXPubToPublicKey(c *gin.Context, xpub string) (*reqctx.UserContext, error) {
	hdKey, err := bip32.GetHDKeyFromExtendedPublicKey(xpub)
	if err != nil {
		return nil, spverrors.ErrAuthorization.Wrap(err)
	}

	pubKey, err := bip32.GetPublicKeyFromHDKey(hdKey)
	if err != nil {
		return nil, spverrors.ErrAuthorization.Wrap(err)
	}
	pubKeyHex := pubKey.ToDERHex()

	user, err := reqctx.Engine(c).Repositories().Users.GetByPubKey(c.Request.Context(), pubKeyHex)
	if err != nil {
		return nil, spverrors.ErrAuthorization.Wrap(err)
	}

	return reqctx.NewUserContextWithPublicKeys(xpub, utils.Hash(xpub), pubKeyHex, user.ID), nil
}
