package reqctx

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/gin-gonic/gin"
)

const userContextKey = "usercontext"

type AuthType = int

const (
	AuthTypeXPub AuthType = iota
	AuthTypeAccessKey
	AuthTypeAdmin
)

// UserContext is the context for the user
type UserContext struct {
	xPub     string
	xPubID   string
	xPubObj  *engine.Xpub
	AuthType AuthType
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

// GetXPubID returns the xPubID from the user context
func (ctx *UserContext) GetXPubID() string {
	return ctx.xPubID
}

// GetXPubObj returns an object of engine.Xpub
func (ctx *UserContext) GetXPubObj() *engine.Xpub {
	return ctx.xPubObj
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

// EnsureXPubIsSet returns the xPub if authorization was made by "regular user" with xPub (not accessKey)
// It panics on fail, so use with caution.
// This function should not be called in actions.
func EnsureXPubIsSet(ctx *UserContext) string {
	if ctx.AuthType != AuthTypeXPub {
		panic("The xPub is not available when the user is authorized by access key or is an admin")
	}
	if ctx.xPub == "" {
		panic("The xPub is not available")
	}

	return ctx.xPub
}
