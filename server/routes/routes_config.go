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

// ApiEndpointsFunc wrapping type for function to mark it as implementation of ApiEndpoints.
type ApiEndpointsFunc func(router *gin.RouterGroup)

// ApiEndpoints registrar which will register routes in ADMIN routes group.
type ApiEndpoints interface {
	// RegisterApiEndpoints register ADMIN endpoints.
	RegisterApiEndpoints(router *gin.RouterGroup)
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

// RegisterApiEndpoints register API endpoints by registrar ApiEndpointsFunc.
func (f ApiEndpointsFunc) RegisterApiEndpoints(router *gin.RouterGroup) {
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

// ApiMiddleware middleware that should handle API requests.
type ApiMiddleware interface {
	//ApplyToApi handle API request by middleware.
	ApplyToApi(c *gin.Context)
}

// ApiMiddlewareFunc wrapping type for function to mark it as implementation of ApiMiddleware.
type ApiMiddlewareFunc func(c *gin.Context)

// ApplyToApi handle API request by middleware function.
func (f ApiMiddlewareFunc) ApplyToApi(c *gin.Context) {
	f(c)
}
