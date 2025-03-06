package errdef_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/jsonrequire"
	"github.com/bitcoin-sv/spv-wallet/logging"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
)

func TestLoggingError(t *testing.T) {
	// given:
	writer := &StringWriter{}
	logger := logging.CreateLogger(writer, "spv-wallet-default", zerolog.DebugLevel, true)

	// when:
	err := a()
	logger.Info().Stack().Err(err).Msg("test-message")

	// then:
	logMsg := writer.builder.String()

	jsonrequire.Match(t, `{
		"log.level": "info",
		"ecs.version": {{ anything }},
		"application": "spv-wallet-default",
		"error.message": "a-wrap, cause: b-wrap, cause: common.internal_error: conversion error, cause: strconv.Atoi: parsing \"a\": invalid syntax",
		"@timestamp": "{{ matchTimestamp }}",
		"log.origin": {{ anything }},
		"message": "test-message",
		"error.stack_trace": {{ anything }}
	}`, map[string]any{
		"matchStackTrace": `/.+/`, // regex to check that it's not empty
	}, logMsg)

	t.Log(logMsg)
}

func a() error {
	return errorx.Decorate(b(), "a-wrap")
}

func b() error {
	return errorx.Decorate(c(), "b-wrap")
}

func c() error {
	_, err := strconv.Atoi("a")
	return errorx.InternalError.Wrap(err, "conversion error")
}

type StringWriter struct {
	builder strings.Builder
}

func (w *StringWriter) Write(p []byte) (n int, err error) {
	return w.builder.Write(p)
}
