package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"net/http"
)

func debugWriter(logger *zerolog.Logger) io.Writer {
	w := func(p []byte) (n int, err error) {
		logger.Debug().Msg(string(p))
		return len(p), err
	}
	return WriterFunc(w)
}

// WriterFunc wrapper type for function that is implementing io.Writer interface.
type WriterFunc func(p []byte) (n int, err error)

// Write proxy to implement io.Writer interface.
func (f WriterFunc) Write(p []byte) (n int, err error) {
	return f(p)
}

func registerSwaggerEndpoints(engine *gin.Engine) {
	engine.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
