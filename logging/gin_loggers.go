package logging

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// SetGinWriters sets GIN to use zerolog logger for debugPrint, recovery messages
// and every other events when it uses fmt.Fprint(DefaultWriter/DefaultErrorWriter, ...)
// https://github.com/gin-gonic/gin/issues/1877#issuecomment-552637900
func SetGinWriters(log *zerolog.Logger) {
	gin.DefaultWriter = newGinLogsWriter(log, zerolog.DebugLevel)
	gin.DefaultErrorWriter = newGinLogsWriter(log, zerolog.ErrorLevel)
}

// GinMiddleware returns a middleware that logs requests using zerolog.
func GinMiddleware(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		params := gin.LogFormatterParams{
			Latency:      time.Since(start),
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			Path:         path,
			StatusCode:   c.Writer.Status(),
			ErrorMessage: c.Errors.ByType(gin.ErrorTypePrivate).String(),
			BodySize:     c.Writer.Size(),
		}

		if raw != "" {
			params.Path = path + "?" + raw
		}

		if params.ErrorMessage != "" {
			logWithRequestParams(log.Warn(), &params).
				Str("error_message", params.ErrorMessage).
				Msg("[GIN] Request Error")
		} else {
			logWithRequestParams(log.Info(), &params).
				Msg("[GIN] Request")
		}
	}
}

func logWithRequestParams(base *zerolog.Event, params *gin.LogFormatterParams) *zerolog.Event {
	return base.
		Str("client_ip", params.ClientIP).
		Str("method", params.Method).
		Int("status", params.StatusCode).
		Dur("latency", params.Latency).
		Int("body_size", params.BodySize).
		Str("path", params.Path)
}

type ginLogsWriter struct {
	logger *zerolog.Logger
	level  zerolog.Level
}

func newGinLogsWriter(logger *zerolog.Logger, level zerolog.Level) *ginLogsWriter {
	return &ginLogsWriter{
		logger: logger,
		level:  level,
	}
}

func (w *ginLogsWriter) Write(p []byte) (n int, err error) {
	w.logger.WithLevel(w.level).Msg(string(p))
	return len(p), nil
}
