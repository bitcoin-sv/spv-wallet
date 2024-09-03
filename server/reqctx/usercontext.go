package reqctx

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/gin-gonic/gin"
)

const userContextKey = "usercontext"

// UserContext is the context for the user
type UserContext struct {
	xPub          string
	xPubID        string
	xPubObj       *engine.Xpub
	authAccessKey string
	isAdmin       bool
}

// NewUserContextWithXPub creates a new UserContext based on xpub authorization
func NewUserContextWithXPub(xpub, xpubID string, xPubObj *engine.Xpub) *UserContext {
	return &UserContext{
		xPub:    xpub,
		xPubID:  xpubID,
		xPubObj: xPubObj,
	}
}

// NewUserContextWithAccessKey creates a new UserContext based on accessKey authorization
func NewUserContextWithAccessKey(authAccessKey string, accessKey *engine.AccessKey, xPubObj *engine.Xpub) *UserContext {
	return &UserContext{
		authAccessKey: authAccessKey,
		xPubID:        accessKey.XpubID,
		xPubObj:       xPubObj,
	}
}

// NewUserContextAsAdmin creates a new UserContext as an admin
func NewUserContextAsAdmin(adminXPub string) *UserContext {
	return &UserContext{
		xPub:    adminXPub,
		isAdmin: true,
	}
}

// GetXPub returns the xPub from the user context
func (ctx *UserContext) GetXPub() string {
	if ctx.IsAdmin() {
		panic("You should not get the admin xPub using this GetXPub method")
	}
	// if authentication was made using accessKey there is not xPub; only xPubID can be used
	return ctx.xPub
}

// GetXPubID returns the xPubID from the user context
func (ctx *UserContext) GetXPubID() string {
	return ctx.xPubID
}

// GetXPubObj returns an object of engine.Xpub
func (ctx *UserContext) GetXPubObj() *engine.Xpub {
	// if authentication was made using accessKey there is not xPub; only xPubID can be used
	return ctx.xPubObj
}

// IsAdmin checks if the user is an admin
func (ctx *UserContext) IsAdmin() bool {
	return ctx.isAdmin
}

// IsAuthorizedByAccessKey checks if the user is authorized by access key
func (ctx *UserContext) IsAuthorizedByAccessKey() bool {
	return ctx.authAccessKey != ""
}

// GetValuesForCheckSignature returns values needed to check signature
func (ctx *UserContext) GetValuesForCheckSignature() (xpub, authAccessKey string) {
	xpub = ctx.xPub
	authAccessKey = ctx.authAccessKey
	return
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
