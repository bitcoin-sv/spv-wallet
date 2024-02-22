package auth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/gin-gonic/gin"
	"github.com/libsv/go-bt/v2/bscript"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// ParamXPubKey the request parameter for the xpub string
	ParamXPubKey = "xpub"

	// ParamXPubHashKey the request parameter for the xpub ID
	ParamXPubHashKey = "xpub_hash"

	// ParamAccessKey the request parameter for the xpub ID
	ParamAccessKey = "access_key"

	// ParamAdminRequest the request parameter whether this is an admin request
	ParamAdminRequest = "auth_admin"

	// ParamAuthSigned the request parameter that says whether the request was signed
	ParamAuthSigned = "auth_signed"
)

// AuthPayload is the authentication payload for checking or creating a signature
type AuthPayload struct {
	AuthHash     string `json:"auth_hash"`
	AuthNonce    string `json:"auth_nonce"`
	AuthTime     int64  `json:"auth_time"`
	BodyContents string `json:"body_contents"`
	Signature    string `json:"signature"`
	xPub         string
	accessKey    string
}

// CorsMiddleware is a middleware that handles CORS.
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Next()
	}
}

// AuthMiddleware will check the request for the xPub or AccessKey header
func AuthMiddleware(engine engine.ClientInterface, appConfig *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		xPub := strings.TrimSpace(c.GetHeader(models.AuthHeader))
		authAccessKey := strings.TrimSpace(c.GetHeader(models.AuthAccessKey))
		if len(xPub) == 0 && len(authAccessKey) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authentication header"})
			return
		}

		fmt.Println("xPub: ", xPub)
		fmt.Println("authAccessKey: ", authAccessKey)

		xPubID := utils.Hash(xPub)
		xPubOrAccessKey := xPub

		if xPub != "" {
			// Validate that the xPub is an HD key (length, validation)
			if _, err := utils.ValidateXPub(xPubOrAccessKey); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
				return
			}

			if xPub == appConfig.Authentication.AdminKey {
				c.Set(ParamAdminRequest, true)
			}

			c.Set(ParamXPubKey, xPub)

		} else if authAccessKey != "" {
			xPubOrAccessKey = authAccessKey
			accessKey, err := engine.AuthenticateAccessKey(context.Background(), utils.Hash(authAccessKey))
			if err != nil || accessKey == nil {
				fmt.Println("Error: ", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
				return
			}

			xPubID = accessKey.XpubID

			c.Set(ParamAccessKey, authAccessKey)
		}

		c.Set(ParamXPubHashKey, xPubID)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(ParamAdminRequest) {
			c.Next()
			return
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "xpub provided is not an admin key")
		}
	}
}

func SignatureMiddleware(appConfig *config.AppConfig, requireSigning, adminRequired bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "missing body")
		}
		defer func() {
			_ = c.Request.Body.Close()
		}()
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(b))

		authTime, _ := strconv.Atoi(c.GetHeader(models.AuthHeaderTime))
		authData := &AuthPayload{
			AuthHash:     c.GetHeader(models.AuthHeaderHash),
			AuthNonce:    c.GetHeader(models.AuthHeaderNonce),
			AuthTime:     int64(authTime),
			BodyContents: string(b),
			Signature:    c.GetHeader(models.AuthSignature),
			xPub:         c.GetString(ParamXPubKey),
			accessKey:    c.GetString(ParamAccessKey),
		}

		// adminRequired will always force checking of a signature
		if (requireSigning || adminRequired) && !appConfig.Authentication.SigningDisabled {
			if err = checkSignature(authData); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			}
			c.Set(ParamAuthSigned, true)
		} else {
			// check the signature and add to request, but do not fail if incorrect
			err = checkSignature(authData)
			c.Set(ParamAuthSigned, err == nil)

			// NOTE: you can not use an access key if signing is invalid - ever
			if authData.accessKey != "" && err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			}
		}
		c.Next()
	}
}

// checkSignature check the signature for the provided auth payload
func checkSignature(auth *AuthPayload) error {
	if err := checkSignatureRequirements(auth); err != nil {
		return err
	}

	if auth.xPub != "" {
		return verifyKeyXPub(auth.xPub, auth)
	}
	return verifyMessageAndSignature(auth.accessKey, auth)
}

// checkSignatureRequirements will check the payload for basic signature requirements
func checkSignatureRequirements(auth *AuthPayload) error {
	if auth == nil || auth.Signature == "" {
		return errors.New("missing signature")
	}

	bodyHash := createBodyHash(auth.BodyContents)
	if auth.AuthHash != bodyHash {
		return errors.New("auth hash and body hash do not match")
	}

	if time.Now().UTC().After(time.UnixMilli(auth.AuthTime).Add(models.AuthSignatureTTL)) {
		return errors.New("signature has expired")
	}
	return nil
}

// createBodyHash will create the hash of the body, removing any carriage returns
func createBodyHash(bodyContents string) string {
	return utils.Hash(strings.TrimSuffix(bodyContents, "\n"))
}

// verifyKeyXPub will verify the xPub key and the signature payload
func verifyKeyXPub(xPub string, auth *AuthPayload) error {
	if _, err := utils.ValidateXPub(xPub); err != nil {
		return err
	}

	if auth == nil {
		return errors.New("missing signature")
	}

	key, err := bitcoin.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		return err
	}

	if key, err = utils.DeriveChildKeyFromHex(key, auth.AuthNonce); err != nil {
		return err
	}

	var address *bscript.Address
	if address, err = bitcoin.GetAddressFromHDKey(key); err != nil {
		return err // Should never error
	}

	message := getSigningMessage(xPub, auth)
	if err = bitcoin.VerifyMessage(
		address.AddressString,
		auth.Signature,
		message,
	); err != nil {
		return errors.New("signature invalid")
	}
	return nil
}

// verifyMessageAndSignature will verify the access key and the signature payload
func verifyMessageAndSignature(key string, auth *AuthPayload) error {
	address, err := bitcoin.GetAddressFromPubKeyString(
		key, true,
	)
	if err != nil {
		return err
	}

	if err := bitcoin.VerifyMessage(
		address.AddressString,
		auth.Signature,
		getSigningMessage(key, auth),
	); err != nil {
		return errors.New("signature invalid")
	}
	return nil
}

// getSigningMessage will build the signing message string
func getSigningMessage(xPub string, auth *AuthPayload) string {
	return fmt.Sprintf("%s%s%s%d", xPub, auth.AuthHash, auth.AuthNonce, auth.AuthTime)
}
