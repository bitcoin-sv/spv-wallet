package middleware

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/gin-gonic/gin"
	"github.com/libsv/go-bt/v2/bscript"
)

// CheckSignatureMiddleware is a middleware that checks the signature of the request (if required)
func CheckSignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appConfig := reqctx.AppConfig(c)
		userContext := reqctx.GetUserContext(c)

		requireSigning := userContext.GetAuthType() == reqctx.AuthTypeAccessKey || appConfig.Authentication.RequireSigning

		if requireSigning {
			if err := verifyRequest(c, userContext); err != nil {
				spverrors.AbortWithErrorResponse(c, err, reqctx.Logger(c))
				return
			}
		}

		c.Next()
	}
}

func verifyRequest(c *gin.Context, userContext *reqctx.UserContext) error {
	bodyContent, err := readBodyContents(c) // for GET methods, bodyContent is an empty string
	if err != nil {
		return err
	}
	authTime, err := strconv.Atoi(c.GetHeader(models.AuthHeaderTime))
	if err != nil {
		return spverrors.ErrInvalidSignature
	}
	validator := &sigAuth{
		AuthHash:  c.GetHeader(models.AuthHeaderHash),
		AuthNonce: c.GetHeader(models.AuthHeaderNonce),
		AuthTime:  int64(authTime),
		Signature: c.GetHeader(models.AuthSignature),
	}

	if err := validator.checkRequirements(bodyContent); err != nil {
		return err
	}

	switch userContext.GetAuthType() {
	case reqctx.AuthTypeXPub:
		return validator.verifyWithXPub(reqctx.EnsureXPubIsSet(userContext))
	case reqctx.AuthTypeAccessKey:
		return validator.verifyWithAccessKey(strings.TrimSpace(c.GetHeader(models.AuthAccessKey)))
	case reqctx.AuthTypeAdmin:
		return validator.verifyWithXPub(reqctx.AppConfig(c).Authentication.AdminKey)
	default:
		return spverrors.ErrAuthorization
	}
}

// readBodyContents reads and returns the whole body content
// To allow gin to read the body while Binding process it substitutes c.Request.Body with new io.NopCloser
// NOTE: for GET methods and other "no-body" requests this function returns empty string (with no error)
func readBodyContents(c *gin.Context) (string, error) {
	if c.Request.Body == nil {
		return "", spverrors.ErrInternal
	}
	defer func() {
		_ = c.Request.Body.Close()
	}()
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", spverrors.ErrInternal
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(b))
	return string(b), nil
}

type sigAuth struct {
	AuthHash  string
	AuthNonce string
	Signature string
	AuthTime  int64
}

func (sa *sigAuth) checkRequirements(bodyContents string) error {
	if sa.Signature == "" {
		return spverrors.ErrMissingSignature
	}

	bodyHash := utils.Hash(strings.TrimSuffix(bodyContents, "\n"))
	if sa.AuthHash != bodyHash {
		return spverrors.ErrInvalidSignature
	}

	if time.Now().UTC().After(time.UnixMilli(sa.AuthTime).Add(models.AuthSignatureTTL)) {
		return spverrors.ErrSignatureExpired
	}
	return nil
}

// verifyWithXPub will verify the xPub key and the signature payload
func (sa *sigAuth) verifyWithXPub(xPub string) error {
	key, err := bitcoin.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		return spverrors.ErrInvalidSignature
	}

	if key, err = utils.DeriveChildKeyFromHex(key, sa.AuthNonce); err != nil {
		return spverrors.ErrInvalidSignature
	}

	var address *bscript.Address
	if address, err = bitcoin.GetAddressFromHDKey(key); err != nil {
		return spverrors.ErrInvalidSignature
	}

	message := sa.getSigningMessage(xPub)
	if err = bitcoin.VerifyMessage(
		address.AddressString,
		sa.Signature,
		message,
	); err != nil {
		return spverrors.ErrInvalidSignature
	}
	return nil
}

// verifyWithAccessKey will verify the access key and the signature payload
func (sa *sigAuth) verifyWithAccessKey(accessKey string) error {
	address, err := bitcoin.GetAddressFromPubKeyString(
		accessKey, true,
	)
	if err != nil {
		return spverrors.ErrInvalidSignature
	}

	if err := bitcoin.VerifyMessage(
		address.AddressString,
		sa.Signature,
		sa.getSigningMessage(accessKey),
	); err != nil {
		return spverrors.ErrInvalidSignature
	}
	return nil
}

// getSigningMessage will build the signing message string
func (sa *sigAuth) getSigningMessage(xPub string) string {
	return fmt.Sprintf("%s%s%s%d", xPub, sa.AuthHash, sa.AuthNonce, sa.AuthTime)
}
