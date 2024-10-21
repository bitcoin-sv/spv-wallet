package spverrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	baseInternalErr := NewError("test error")

	tests := map[string]struct {
		err error

		expectHTTPStatus  int
		expectHTTPMessage string
		expectDetails     string
		expectCode        string
	}{
		"simple internal": {
			err: baseInternalErr,

			expectHTTPStatus:  500,
			expectHTTPMessage: internalErrMessage,
			expectDetails:     "test error",
			expectCode:        genericErrorCode,
		},
		"with details": {
			err: baseInternalErr.WithDetails("details"),

			expectHTTPStatus:  500,
			expectHTTPMessage: internalErrMessage,
			expectDetails:     "test error; details",
			expectCode:        genericErrorCode,
		},
		"with formatted details": {
			err: baseInternalErr.WithDetails("details %d", 1),

			expectHTTPStatus:  500,
			expectHTTPMessage: internalErrMessage,
			expectDetails:     "test error; details 1",
			expectCode:        genericErrorCode,
		},
		"with http status": {
			err: NewError("user did sth wrong").WithHTTPStatus(404, "not found"),

			expectHTTPStatus:  404,
			expectHTTPMessage: "not found",
			expectDetails:     "user did sth wrong",
			expectCode:        genericErrorCode,
		},
		"with code": {
			err: baseInternalErr.WithCode("test-code"),

			expectHTTPStatus:  500,
			expectHTTPMessage: internalErrMessage,
			expectDetails:     "test error",
			expectCode:        "test-code",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var e Error
			if !errors.As(test.err, &e) {
				t.Fatalf("error is not of type Error")
			}

			require.Equal(t, test.expectHTTPStatus, e.GetStatusCode())
			require.Equal(t, test.expectHTTPMessage, e.GetMessage())
			require.Equal(t, test.expectDetails, e.Error())
			require.Equal(t, test.expectCode, e.GetCode())
		})
	}
}

func TestErrorIs(t *testing.T) {
	tests := map[string]struct {
		left  error
		right error

		expect bool
	}{
		"spverror and std error": {
			left:  NewError("test error"),
			right: errors.New("test error"),

			expect: false,
		},
		"two spverrors with no code specified": {
			left:  NewError("test error1"),
			right: NewError("test error2"),

			expect: true,
		},
		"two spverrors with custom code": {
			left:  NewError("test error1").WithCode("test-code"),
			right: NewError("test error2"),

			expect: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, test.expect, errors.Is(test.left, test.right))
		})
	}
}
