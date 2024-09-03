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
		err := checkSignature(c)
		if err != nil {
			spverrors.AbortWithErrorResponse(c, err, reqctx.Logger(c))
		} else {
			c.Next()
		}
	}
}

func checkSignature(c *gin.Context) error {
	appConfig := reqctx.AppConfig(c)
	userContext := reqctx.GetUserContext(c)

	requireSigning := userContext.IsAuthorizedByAccessKey() || appConfig.Authentication.RequireSigning
	if requireSigning {
		if err := verify(c, userContext); err != nil {
			return err
		}
	}
	return nil
}

type payload struct {
	AuthHash     string `json:"auth_hash"`
	AuthNonce    string `json:"auth_nonce"`
	BodyContents string `json:"body_contents"`
	Signature    string `json:"signature"`
	xPub         string
	accessKey    string
	AuthTime     int64 `json:"auth_time"`
}

// verify check the signature for the provided auth payload
func verify(c *gin.Context, userContext *reqctx.UserContext) error {
	bodyContent, err := readBodyContents(c)
	if err != nil {
		return err
	}
	authTime, _ := strconv.Atoi(c.GetHeader(models.AuthHeaderTime))
	xPub, authAccessKey := userContext.GetValuesForCheckSignature()
	sigData := &payload{
		AuthHash:     c.GetHeader(models.AuthHeaderHash),
		AuthNonce:    c.GetHeader(models.AuthHeaderNonce),
		AuthTime:     int64(authTime),
		BodyContents: bodyContent,
		Signature:    c.GetHeader(models.AuthSignature),
		xPub:         xPub,
		accessKey:    authAccessKey,
	}

	if err := checkSignatureRequirements(sigData); err != nil {
		return err
	}

	if sigData.xPub != "" {
		return verifyKeyXPub(sigData.xPub, sigData)
	}
	return verifyMessageAndSignature(sigData.accessKey, sigData)
}

func readBodyContents(c *gin.Context) (string, error) {
	if c.Request.Body == nil {
		return "", spverrors.ErrMissingBody
	}
	defer func() {
		_ = c.Request.Body.Close()
	}()
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", spverrors.ErrMissingBody
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(b))
	return string(b), nil
}

// checkSignatureRequirements will check the payload for basic signature requirements
func checkSignatureRequirements(auth *payload) error {
	if auth == nil || auth.Signature == "" {
		return spverrors.ErrMissingSignature
	}

	bodyHash := createBodyHash(auth.BodyContents)
	if auth.AuthHash != bodyHash {
		return spverrors.ErrHashesDoNotMatch
	}

	if time.Now().UTC().After(time.UnixMilli(auth.AuthTime).Add(models.AuthSignatureTTL)) {
		return spverrors.ErrSignatureExpired
	}
	return nil
}

// createBodyHash will create the hash of the body, removing any carriage returns
func createBodyHash(bodyContents string) string {
	return utils.Hash(strings.TrimSuffix(bodyContents, "\n"))
}

// verifyKeyXPub will verify the xPub key and the signature payload
func verifyKeyXPub(xPub string, auth *payload) error {
	if _, err := utils.ValidateXPub(xPub); err != nil {
		return spverrors.ErrValidateXPub
	}

	if auth == nil {
		return spverrors.ErrMissingSignature
	}

	key, err := bitcoin.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		return spverrors.ErrGettingHdKeyFromXpub
	}

	if key, err = utils.DeriveChildKeyFromHex(key, auth.AuthNonce); err != nil {
		return spverrors.ErrDeriveChildKey
	}

	var address *bscript.Address
	if address, err = bitcoin.GetAddressFromHDKey(key); err != nil {
		return spverrors.ErrGettingAddressFromHdKey
	}

	message := getSigningMessage(xPub, auth)
	if err = bitcoin.VerifyMessage(
		address.AddressString,
		auth.Signature,
		message,
	); err != nil {
		return spverrors.ErrInvalidSignature
	}
	return nil
}

// verifyMessageAndSignature will verify the access key and the signature payload
func verifyMessageAndSignature(key string, auth *payload) error {
	address, err := bitcoin.GetAddressFromPubKeyString(
		key, true,
	)
	if err != nil {
		return spverrors.ErrGettingAddressFromPublicKey
	}

	if err := bitcoin.VerifyMessage(
		address.AddressString,
		auth.Signature,
		getSigningMessage(key, auth),
	); err != nil {
		return spverrors.ErrInvalidSignature
	}
	return nil
}

// getSigningMessage will build the signing message string
func getSigningMessage(xPub string, auth *payload) string {
	return fmt.Sprintf("%s%s%s%d", xPub, auth.AuthHash, auth.AuthNonce, auth.AuthTime)
}
