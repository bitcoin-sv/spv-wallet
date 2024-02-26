package routes

import (
	"github.com/gin-gonic/gin"
)

// AdminEndpointsFunc wrapping type for function to mark it as implementation of AdminEndpoints.
type AdminEndpointsFunc func(router *gin.RouterGroup)

// AdminEndpoints registrar which will register routes with admin auth middleware.
type AdminEndpoints interface {
	// RegisterAdminEndpoints register root endpoints.
	RegisterAdminEndpoints(router *gin.RouterGroup)
}

// APIEndpointsFunc wrapping type for function to mark it as implementation of ApiEndpoints.
type APIEndpointsFunc func(router *gin.RouterGroup)

// APIEndpoints registrar which will register routes in ADMIN routes group.
type APIEndpoints interface {
	// RegisterAPIEndpoints register ADMIN endpoints.
	RegisterAPIEndpoints(router *gin.RouterGroup)
}

// BasicEndpointsFunc wrapping type for function to mark it as implementation of BasicEndpoints.
type BasicEndpointsFunc func(router *gin.RouterGroup)

// BasicEndpoints registrar which will register routes in BASIC routes group.
type BasicEndpoints interface {
	// RegisterBasicEndpoints register BASIC endpoints.
	RegisterBasicEndpoints(router *gin.RouterGroup)
}

// BaseEndpointsFunc wrapping type for function to mark it as implementation of BaseEndpoints.
type BaseEndpointsFunc func(router *gin.RouterGroup)

// BaseEndpoints registrar which will register routes in BASE routes group.
type BaseEndpoints interface {
	// RegisterBaseEndpoints register BASE endpoints.
	RegisterBaseEndpoints(router *gin.RouterGroup)
}

// RegisterAdminEndpoints register root endpoints by registrar AdminEndpointsFunc.
func (f AdminEndpointsFunc) RegisterAdminEndpoints(router *gin.RouterGroup) {
	f(router)
}

// RegisterAPIEndpoints register API endpoints by registrar ApiEndpointsFunc.
func (f APIEndpointsFunc) RegisterAPIEndpoints(router *gin.RouterGroup) {
	f(router)
}

// RegisterBasicEndpoints register API endpoints by registrar BasicEndpointsFunc.
func (f BasicEndpointsFunc) RegisterBasicEndpoints(router *gin.RouterGroup) {
	f(router)
}

// RegisterBaseEndpoints register API endpoints by registrar BaseEndpointsFunc.
func (f BaseEndpointsFunc) RegisterBaseEndpoints(router *gin.RouterGroup) {
	f(router)
}

// APIMiddleware middleware that should handle API requests.
type APIMiddleware interface {
	//ApplyToAPI handle API request by middleware.
	ApplyToAPI(c *gin.Context)
}

// APIMiddlewareFunc wrapping type for function to mark it as implementation of ApiMiddleware.
type APIMiddlewareFunc func(c *gin.Context)

// ApplyToAPI handle API request by middleware function.
func (f APIMiddlewareFunc) ApplyToAPI(c *gin.Context) {
	f(c)
}
