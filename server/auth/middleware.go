package auth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/spverrors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/gin-gonic/gin"
	"github.com/libsv/go-bt/v2/bscript"
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

// Payload is the authentication payload for checking or creating a signature
type Payload struct {
	AuthHash     string `json:"auth_hash"`
	AuthNonce    string `json:"auth_nonce"`
	BodyContents string `json:"body_contents"`
	Signature    string `json:"signature"`
	xPub         string
	accessKey    string
	AuthTime     int64 `json:"auth_time"`
}

// CorsMiddleware is a middleware that handles CORS.
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		corsAllowedHeaders := []string{
			"Content-Type",
			"Cache-Control",
			models.AuthHeader,
			models.AuthAccessKey,
			models.AuthSignature,
			models.AuthHeaderHash,
			models.AuthHeaderNonce,
			models.AuthHeaderTime,
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(corsAllowedHeaders, ","))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// BasicMiddleware will check the request for the xPub or AccessKey header
func BasicMiddleware(engine engine.ClientInterface, appConfig *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		xPub := strings.TrimSpace(c.GetHeader(models.AuthHeader))
		authAccessKey := strings.TrimSpace(c.GetHeader(models.AuthAccessKey))
		if len(xPub) == 0 && len(authAccessKey) == 0 {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrMissingAuthHeader, nil)
			return
		}

		xPubID := utils.Hash(xPub)

		if xPub != "" {
			// Validate that the xPub is an HD key (length, validation)
			if _, err := utils.ValidateXPub(xPub); err != nil {
				spverrors.AbortWithErrorResponse(c, spverrors.ErrAuthorization, nil)
				return
			}

			if xPub == appConfig.Authentication.AdminKey {
				c.Set(ParamAdminRequest, true)
			}

			c.Set(ParamXPubKey, xPub)

		} else if authAccessKey != "" {
			accessKey, err := engine.AuthenticateAccessKey(context.Background(), utils.Hash(authAccessKey))
			if err != nil || accessKey == nil {
				spverrors.AbortWithErrorResponse(c, spverrors.ErrAuthorization, nil)
				return
			}

			xPubID = accessKey.XpubID

			c.Set(ParamAccessKey, authAccessKey)
		}

		c.Set(ParamXPubHashKey, xPubID)
		c.Next()
	}
}

// AdminMiddleware will check if the request is authorized with admin xpub
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(ParamAdminRequest) {
			c.Next()
		} else {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrNotAnAdminKey, nil)
		}
	}
}

// SignatureMiddleware will check the request for a signature
func SignatureMiddleware(appConfig *config.AppConfig, requireSigning, adminRequired bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrMissingBody, nil)
			return
		}
		defer func() {
			_ = c.Request.Body.Close()
		}()
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrAuthorization, nil)
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(b))

		authTime, _ := strconv.Atoi(c.GetHeader(models.AuthHeaderTime))
		authData := &Payload{
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
				spverrors.AbortWithErrorResponse(c, spverrors.ErrCheckSignature, nil)
			}
			c.Set(ParamAuthSigned, true)
		} else {
			// check the signature and add to request, but do not fail if incorrect
			err = checkSignature(authData)
			c.Set(ParamAuthSigned, err == nil)

			// NOTE: you can not use an access key if signing is invalid - ever
			if authData.accessKey != "" && err != nil {
				spverrors.AbortWithErrorResponse(c, spverrors.ErrCheckSignature, nil)
			}
		}
		c.Next()
	}
}

// CallbackTokenMiddleware verifies the callback token - if it's valid and matches the Bearer scheme.
func CallbackTokenMiddleware(appConfig *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrMissingAuthHeader, nil)
		}

		if !strings.HasPrefix(authHeader, BearerSchema) || len(authHeader) <= len(BearerSchema) {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrInvalidOrMissingToken, nil)
		}

		providedToken := authHeader[len(BearerSchema):]
		if providedToken != appConfig.Nodes.Callback.Token {
			spverrors.AbortWithErrorResponse(c, spverrors.ErrInvalidToken, nil)
		}

		c.Next()
	}
}

// checkSignature check the signature for the provided auth payload
func checkSignature(auth *Payload) error {
	if err := checkSignatureRequirements(auth); err != nil {
		return err
	}

	if auth.xPub != "" {
		return verifyKeyXPub(auth.xPub, auth)
	}
	return verifyMessageAndSignature(auth.accessKey, auth)
}

// checkSignatureRequirements will check the payload for basic signature requirements
func checkSignatureRequirements(auth *Payload) error {
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
func verifyKeyXPub(xPub string, auth *Payload) error {
	if _, err := utils.ValidateXPub(xPub); err != nil {
		err := fmt.Errorf("error occurred while validating xPub key: %w", err)
		return err
	}

	if auth == nil {
		return errors.New("missing signature")
	}

	key, err := bitcoin.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		err = fmt.Errorf("error occurred while getting HD key from xPub: %w", err)
		return err
	}

	if key, err = utils.DeriveChildKeyFromHex(key, auth.AuthNonce); err != nil {
		err = fmt.Errorf("error occurred while deriving child key: %w", err)
		return err
	}

	var address *bscript.Address
	if address, err = bitcoin.GetAddressFromHDKey(key); err != nil {
		err = fmt.Errorf("error occurred while getting address from HD key: %w", err)
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
func verifyMessageAndSignature(key string, auth *Payload) error {
	address, err := bitcoin.GetAddressFromPubKeyString(
		key, true,
	)
	if err != nil {
		err = fmt.Errorf("error occurred while getting address from public key: %w", err)
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
func getSigningMessage(xPub string, auth *Payload) string {
	return fmt.Sprintf("%s%s%s%d", xPub, auth.AuthHash, auth.AuthNonce, auth.AuthTime)
}
