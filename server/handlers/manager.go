package handlers

import (
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/middleware"
	"github.com/gin-gonic/gin"
)

// GroupType is a type of group
type GroupType int

const (
	// GroupRoot is the root group without no prefix and no auth middleware
	GroupRoot GroupType = iota

	// GroupAPI is the group with the API prefix and auth middleware
	GroupAPI

	// GroupAPIV2 is the group with the API v2 prefix and auth middleware
	GroupAPIV2

	// GroupTransactionCallback is the group with the transaction callback prefix and callback token middleware (no auth middleware)
	GroupTransactionCallback
)

// Manager is a struct helps to group routes with proper middleware
type Manager struct {
	engine    *gin.Engine
	appConfig *config.AppConfig
	groups    map[GroupType]*gin.RouterGroup
}

// NewManager creates a new Grouper
func NewManager(engine *gin.Engine, appConfig *config.AppConfig) *Manager {
	authRouter := engine.Group("", middleware.AuthMiddleware(), middleware.CheckSignatureMiddleware())

	return &Manager{
		engine:    engine,
		appConfig: appConfig,
		groups: map[GroupType]*gin.RouterGroup{
			GroupRoot:                engine.Group(""),
			GroupAPI:                 authRouter.Group("/api" + "/" + config.APIVersion),
			GroupAPIV2:               engine.Group("/api/v2", middleware.AuthV2Middleware(), middleware.CheckSignatureMiddleware()),
			GroupTransactionCallback: engine.Group("", middleware.CallbackTokenMiddleware()),
		},
	}
}

// Group creates a new group with the given endpointType and relativePath and optional list checkers, e.g. middleware.RequireSignature
func (mg *Manager) Group(endpointType GroupType, relativePath string, middlewares ...gin.HandlerFunc) *gin.RouterGroup {
	return mg.Get(endpointType).Group(relativePath, middlewares...)
}

// Get returns the group with the given endpointType
func (mg *Manager) Get(endpointType GroupType) *gin.RouterGroup {
	return mg.groups[endpointType]
}

// GetFeatureFlags returns the experimental feature flags from app configuration
func (mg *Manager) GetFeatureFlags() *config.ExperimentalConfig {
	return mg.appConfig.ExperimentalFeatures
}

// APIVersion returns the API version from app configuration
func (mg *Manager) APIVersion() string {
	return mg.appConfig.Version
}

// GetConfig returns the app configuration
func (mg *Manager) GetConfig() *config.AppConfig {
	return mg.appConfig
}
