package reqctx

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
)

const userContextKey = "usercontext"

// AuthType is the type of authentication
type AuthType = int

const (
	// AuthTypeXPub is when user provides xPub
	AuthTypeXPub AuthType = iota

	// AuthTypeAccessKey is when user provides access key
	AuthTypeAccessKey

	// AuthTypeAdmin is when provided xpub matches the admin key
	AuthTypeAdmin
)

// UserContext is the context for the user
type UserContext struct {
	xPub     string
	xPubID   string
	xPubObj  *engine.Xpub
	AuthType AuthType

	// v2
	userID    string
	publicKey string
}

// NewUserContextWithXPub creates a new UserContext based on xpub authorization
func NewUserContextWithXPub(xpub, xpubID string, xPubObj *engine.Xpub) *UserContext {
	return &UserContext{
		xPub:     xpub,
		xPubID:   xpubID,
		xPubObj:  xPubObj,
		AuthType: AuthTypeXPub,
	}
}

// NewUserContextWithAccessKey creates a new UserContext based on accessKey authorization
func NewUserContextWithAccessKey(xpubID string, xPubObj *engine.Xpub) *UserContext {
	return &UserContext{
		xPubID:   xpubID,
		xPubObj:  xPubObj,
		AuthType: AuthTypeAccessKey,
	}
}

// NewUserContextWithPublicKeys creates a new UserContext based on public keys
// Note: This is used for API v2 authentication only
func NewUserContextWithPublicKeys(xpub, xpubID, publicKey, userID string) *UserContext {
	return &UserContext{
		xPub:      xpub,
		xPubID:    xpubID,
		userID:    userID,
		publicKey: publicKey,
		AuthType:  AuthTypeXPub,
	}
}

// NewUserContextAsAdmin creates a new UserContext as an admin
func NewUserContextAsAdmin() *UserContext {
	return &UserContext{
		AuthType: AuthTypeAdmin,
	}
}

// GetAuthType returns the authentication type from the user context
func (ctx *UserContext) GetAuthType() AuthType {
	return ctx.AuthType
}

// ShouldGetXPub returns the xPub from the user context
// If the authentication type is not xPub, it will return an error
func (ctx *UserContext) ShouldGetXPub() (string, error) {
	if ctx.AuthType != AuthTypeXPub {
		return "", spverrors.ErrXPubAuthRequired
	}
	if ctx.xPub == "" {
		// if AuthType is XPub, xPub should not be empty (by design)
		// if it is empty, it is a bug
		return "", spverrors.ErrInternal
	}

	return ctx.xPub, nil
}

// GetXPubID returns the xPubID from the user context
func (ctx *UserContext) GetXPubID() string {
	return ctx.xPubID
}

// GetXPubObj returns an object of engine.Xpub
func (ctx *UserContext) GetXPubObj() *engine.Xpub {
	return ctx.xPubObj
}

// ShouldGetUserID returns userID for NEW DB SCHEMA
// Warning: Don't use it for old DB schema
func (ctx *UserContext) ShouldGetUserID() (string, error) {
	if ctx.AuthType != AuthTypeXPub {
		return "", spverrors.ErrXPubAuthRequired
	}
	if ctx.userID == "" {
		return "", spverrors.ErrInternal
	}

	return ctx.userID, nil
}

// GetUserContext returns the user context from the request context
func GetUserContext(c *gin.Context) *UserContext {
	value := c.MustGet(userContextKey)
	return value.(*UserContext)
}

// SetUserContext sets the user context in the request context
func SetUserContext(c *gin.Context, userContext *UserContext) {
	c.Set(userContextKey, userContext)
}
