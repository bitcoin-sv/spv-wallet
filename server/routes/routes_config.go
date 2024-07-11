package routes

import (
	"github.com/gin-gonic/gin"
)

// Handler is a type that represents a handler for various types of endpoints.
type Handler struct {
	BasicEndpointsFunc
	APIEndpoints
	CallbackEndpoints
}

// AdminEndpointsFunc wrapping type for function to mark it as implementation of AdminEndpoints.
type AdminEndpointsFunc func(router *gin.RouterGroup)

// AdminEndpoints registrar which will register routes with admin auth middleware.
type AdminEndpoints interface {
	// RegisterAdminEndpoints register root endpoints.
	RegisterAdminEndpoints(router *gin.RouterGroup)
}

// OldAPIEndpointsFunc wrapping type for function to mark it as implementation of OldApiEndpoints.
type OldAPIEndpointsFunc func(router *gin.RouterGroup)

// OldAPIEndpoints registrar which will register routes in ADMIN routes group.
type OldAPIEndpoints interface {
	// RegisterAPIEndpoints register ADMIN endpoints.
	RegisterOldAPIEndpoints(router *gin.RouterGroup)
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

// CallbackEndpointsFunc wrapping type for function to mark it as implementation of BaseEndpoints.
type CallbackEndpointsFunc func(router *gin.RouterGroup)

// CallbackEndpoints registrar which will register routes in CALLBACK routes group.
type CallbackEndpoints interface {
	// RegisterCallbackEndpoints register CALLBACK endpoints.
	RegisterCallbackEndpoints(router *gin.RouterGroup)
}

// RegisterAdminEndpoints register root endpoints by registrar AdminEndpointsFunc.
func (f AdminEndpointsFunc) RegisterAdminEndpoints(router *gin.RouterGroup) {
	f(router)
}

// RegisterOldAPIEndpoints register API endpoints by registrar OldApiEndpointsFunc.
func (f OldAPIEndpointsFunc) RegisterOldAPIEndpoints(router *gin.RouterGroup) {
	f(router)
}

// RegisterAPIEndpoints register API endpoints by registrar ApiEndpointsFunc.
func (f APIEndpointsFunc) RegisterAPIEndpoints(router *gin.RouterGroup) {
	f(router)
}

// RegisterBasicEndpoints register Basic endpoints by registrar BasicEndpointsFunc.
func (f BasicEndpointsFunc) RegisterBasicEndpoints(router *gin.RouterGroup) {
	f(router)
}

// RegisterBaseEndpoints register Base endpoints by registrar BaseEndpointsFunc.
func (f BaseEndpointsFunc) RegisterBaseEndpoints(router *gin.RouterGroup) {
	f(router)
}

// RegisterCallbackEndpoints register Callback endpoints by registrar CallbackEndpointsFunc.
func (f CallbackEndpointsFunc) RegisterCallbackEndpoints(router *gin.RouterGroup) {
	f(router)
}
