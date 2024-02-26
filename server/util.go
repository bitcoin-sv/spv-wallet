package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
