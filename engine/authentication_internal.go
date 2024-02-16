package engine

import (
	"context"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/bip32"
)

const (
	// AuthHeader is the header to use for authentication (raw xPub)
	AuthHeader = "x-auth-xpub"

	// AuthAccessKey is the header to use for access key authentication (access public key)
	AuthAccessKey = "x-auth-key"

	// AuthSignature is the given signature (body + timestamp)
	AuthSignature = "x-auth-signature"

	// AuthHeaderHash hash of the body coming from the request
	AuthHeaderHash = "x-auth-hash"

	// AuthHeaderNonce random nonce for the request
	AuthHeaderNonce = "x-auth-nonce"

	// AuthHeaderTime the time of the request, only valid for 30 seconds
	AuthHeaderTime = "x-auth-time"

	// AuthSignatureTTL is the max TTL for a signature to be valid
	AuthSignatureTTL = 20 * time.Second
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

// ParamRequestKey for context key
type ParamRequestKey string

const (
	// ParamXPubKey the request parameter for the xpub string
	ParamXPubKey ParamRequestKey = "xpub"

	// ParamXPubHashKey the request parameter for the xpub ID
	ParamXPubHashKey ParamRequestKey = "xpub_hash"

	// ParamAdminRequest the request parameter whether this is an admin request
	ParamAdminRequest ParamRequestKey = "auth_admin"

	// ParamAuthSigned the request parameter that says whether the request was signed
	ParamAuthSigned ParamRequestKey = "auth_signed"
)

// createBodyHash will create the hash of the body, removing any carriage returns
func createBodyHash(bodyContents string) string {
	return utils.Hash(strings.TrimSuffix(bodyContents, "\n"))
}

// createSignature will create a signature for the given key & body contents
func createSignature(xPriv *bip32.ExtendedKey, bodyString string) (payload *AuthPayload, err error) {
	// No key?
	if xPriv == nil {
		err = ErrMissingXPriv
		return
	}

	// Get the xPub
	payload = new(AuthPayload)
	if payload.xPub, err = bitcoin.GetExtendedPublicKey(
		xPriv,
	); err != nil { // Should never error if key is correct
		return
	}

	// auth_nonce is a random unique string to seed the signing message
	// this can be checked server side to make sure the request is not being replayed
	if payload.AuthNonce, err = utils.RandomHex(32); err != nil { // Should never error if key is correct
		return
	}

	// Derive the address for signing
	var key *bip32.ExtendedKey
	if key, err = utils.DeriveChildKeyFromHex(
		xPriv, payload.AuthNonce,
	); err != nil {
		return
	}

	var privateKey *bec.PrivateKey
	if privateKey, err = bitcoin.GetPrivateKeyFromHDKey(key); err != nil {
		return // Should never error if key is correct
	}

	return createSignatureCommon(payload, bodyString, privateKey)
}

// createSignatureAccessKey will create a signature for the given access key & body contents
func createSignatureAccessKey(privateKeyHex, bodyString string) (payload *AuthPayload, err error) {
	// No key?
	if privateKeyHex == "" {
		err = ErrMissingAccessKey
		return
	}

	var privateKey *bec.PrivateKey
	if privateKey, err = bitcoin.PrivateKeyFromString(
		privateKeyHex,
	); err != nil {
		return
	}
	publicKey := privateKey.PubKey()

	// Get the xPub
	payload = new(AuthPayload)
	payload.accessKey = hex.EncodeToString(publicKey.SerialiseCompressed())

	// auth_nonce is a random unique string to seed the signing message
	// this can be checked server side to make sure the request is not being replayed
	payload.AuthNonce, err = utils.RandomHex(32)
	if err != nil {
		return nil, err
	}

	return createSignatureCommon(payload, bodyString, privateKey)
}

// createSignatureCommon will create a signature
func createSignatureCommon(payload *AuthPayload, bodyString string, privateKey *bec.PrivateKey) (*AuthPayload, error) {
	// Create the auth header hash
	payload.AuthHash = utils.Hash(bodyString)

	// auth_time is the current time and makes sure a request can not be sent after 30 secs
	payload.AuthTime = time.Now().UnixMilli()

	key := payload.xPub
	if key == "" && payload.accessKey != "" {
		key = payload.accessKey
	}

	// Signature, using bitcoin signMessage
	var err error
	if payload.Signature, err = bitcoin.SignMessage(
		hex.EncodeToString(privateKey.Serialise()),
		getSigningMessage(key, payload),
		true,
	); err != nil {
		return nil, err
	}

	return payload, nil
}

// setOnRequest will set the value on the request with the given key
func setOnRequest(req *http.Request, keyName ParamRequestKey, value interface{}) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), keyName, value))
}

// getFromRequest gets the stored value from the request if found
func getFromRequest(req *http.Request, key ParamRequestKey) (v string, ok bool) {
	v, ok = req.Context().Value(key).(string)
	return
}

// getBoolFromRequest gets the stored bool value from the request if found
func getBoolFromRequest(req *http.Request, key ParamRequestKey) (v bool, ok bool) {
	v, ok = req.Context().Value(key).(bool)
	return
}
