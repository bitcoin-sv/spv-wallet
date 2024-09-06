package reqctx

import (
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	appConfigKey = "appconfig"
	appEngineKey = "appengine"
	appLoggerKey = "applogger"
)

// AppConfig returns the app config from the request context
func AppConfig(c *gin.Context) *config.AppConfig {
	value := c.MustGet(appConfigKey)
	return value.(*config.AppConfig)
}

// Engine returns the app engine from the request context
func Engine(c *gin.Context) engine.ClientInterface {
	value := c.MustGet(appEngineKey)
	return value.(engine.ClientInterface)
}

// Logger returns the app logger from the request context
func Logger(c *gin.Context) *zerolog.Logger {
	value := c.MustGet(appLoggerKey)
	return value.(*zerolog.Logger)
}

// SetAppConfig sets the app config in the request context
func SetAppConfig(c *gin.Context, appConfig *config.AppConfig) {
	c.Set(appConfigKey, appConfig)
}

// SetEngine sets the app engine in the request context
func SetEngine(c *gin.Context, engine engine.ClientInterface) {
	c.Set(appEngineKey, engine)
}

// SetLogger sets the app logger in the request context
func SetLogger(c *gin.Context, logger *zerolog.Logger) {
	c.Set(appLoggerKey, logger)
}
