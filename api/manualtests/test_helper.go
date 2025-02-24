package manualtests

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func Logger() zerolog.Logger {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local)
	}
	logger := zerolog.New(zerolog.NewConsoleWriter(
		func(w *zerolog.ConsoleWriter) {
			w.FieldsOrder = []string{
				"method",
				"url",
				"status",
				"bodyCallback",
				"z_bodyFallback",
			}
		},
	)).With().Timestamp().Logger()
	return logger.Level(zerolog.TraceLevel)
}

func RequireSuccess(t testing.TB, result Result) {
	require.GreaterOrEqual(t, result.StatusCode(), 200, "Http Status Code is not success")
	require.Less(t, result.StatusCode(), 300, "Http Status Code is not success")
}

type StoredResponse[R Result] struct {
	Response R
}

func StoreResponse[R Result](callback func(state StateForCall, c *client.ClientWithResponses) (R, error)) (*StoredResponse[R], CallWithState) {
	var response StoredResponse[R]

	return &response, func(state StateForCall, c *client.ClientWithResponses) (Result, error) {
		r, err := callback(state, c)
		response.Response = r
		return r, err
	}
}

func StoreResponseOfCall[R Result](callback func(c *client.ClientWithResponses) (R, error)) (*StoredResponse[R], Call) {
	var response StoredResponse[R]

	return &response, func(c *client.ClientWithResponses) (Result, error) {
		r, err := callback(c)
		response.Response = r
		return r, err
	}
}

func (r *StoredResponse[R]) IsEmpty() bool {
	if r == nil {
		return false
	}

	val := reflect.ValueOf(r.Response)
	if val.Kind() == reflect.Ptr {
		return val.IsNil()
	}

	return false
}

func (r *StoredResponse[R]) IsNotEmpty() bool {
	return !r.IsEmpty()
}

func (r *StoredResponse[R]) MustGetResponse() R {
	if r.IsEmpty() {
		panic(fmt.Sprintf("stored response %T is empty", r))
	}
	return r.Response
}
