package dictionary

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrorCode_IsValid tests the method IsValid()
func TestErrorCode_IsValid(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase    string
		code        ErrorCode
		expectValid bool
	}{
		{
			"ErrorBadErrorCode",
			ErrorBadErrorCode,
			true,
		},
		{
			"ErrorMissingEnv",
			ErrorMissingEnv,
			true,
		},
		{
			"ErrorInvalidEnv",
			ErrorInvalidEnv,
			true,
		},
		{
			"ErrorReadingConfig",
			ErrorReadingConfig,
			true,
		},
		{
			"ErrorViper",
			ErrorViper,
			true,
		},
		{
			"ErrorConfigValidation",
			ErrorConfigValidation,
			true,
		},
		{
			"ErrorDecryptEnv",
			ErrorDecryptEnv,
			true,
		},
		{
			"ErrorLoadingConfig",
			ErrorLoadingConfig,
			true,
		},
		{
			"ErrorLoadingCache",
			ErrorLoadingCache,
			true,
		},
		{
			"unknown code",
			9999,
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			valid := test.code.IsValid()
			assert.Equal(t, test.expectValid, valid)
		})
	}
}

// TestGetInternalMessage tests the method GetInternalMessage()
func TestGetInternalMessage(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase        string
		code            ErrorCode
		expectedMessage string
	}{
		{
			"ErrorBadErrorCode",
			ErrorBadErrorCode,
			"failed to find internal error message from error code",
		},
		{
			"ErrorMissingEnv",
			ErrorMissingEnv,
			"missing required environment variable: %s",
		},
		{
			"ErrorInvalidEnv",
			ErrorInvalidEnv,
			"invalid environment variable value: %s",
		},
		{
			"ErrorReadingConfig",
			ErrorReadingConfig,
			"error reading environment configuration: %s",
		},
		{
			"ErrorViper",
			ErrorViper,
			"error in viper unmarshal into config.Values: %s",
		},
		{
			"ErrorConfigValidation",
			ErrorConfigValidation,
			"error in environment configuration validation: %s",
		},
		{
			"ErrorDecryptEnv",
			ErrorDecryptEnv,
			"error in decrypting %s: %s",
		},
		{
			"ErrorLoadingConfig",
			ErrorLoadingConfig,
			"fatal error loading configuration: %s",
		},
		{
			"ErrorLoadingCache",
			ErrorLoadingCache,
			"failed to enable cache: %s - cache is disabled",
		},
		{
			"unknown code",
			9999,
			"failed to find internal error message from error code",
		},
	}

	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			message := GetInternalMessage(test.code)
			assert.Equal(t, test.expectedMessage, message)
		})
	}
}

// TestGetPublicMessage tests the method GetPublicMessage()
func TestGetPublicMessage(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase        string
		code            ErrorCode
		expectedMessage string
	}{
		{
			"ErrorBadErrorCode",
			ErrorBadErrorCode,
			"error not found",
		},
		{
			"ErrorMissingEnv",
			ErrorMissingEnv,
			"missing required environment variable",
		},
		{
			"ErrorInvalidEnv",
			ErrorInvalidEnv,
			"invalid environment variable value",
		},
		{
			"ErrorReadingConfig",
			ErrorReadingConfig,
			"error reading environment configuration",
		},
		{
			"ErrorViper",
			ErrorViper,
			"error in loading configuration",
		},
		{
			"ErrorConfigValidation",
			ErrorConfigValidation,
			"error in environment configuration validation",
		},
		{
			"ErrorDecryptEnv",
			ErrorDecryptEnv,
			"error decrypting an encrypted environment variable",
		},
		{
			"ErrorLoadingConfig",
			ErrorLoadingConfig,
			"error loading configuration",
		},
		{
			"ErrorLoadingCache",
			ErrorLoadingCache,
			"failed to enable cache",
		},
		{
			"unknown error",
			9999,
			"error not found",
		},
	}

	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			message := GetPublicMessage(test.code)
			assert.Equal(t, test.expectedMessage, message)
		})
	}
}

// TestGetStatusCode tests the method GetStatusCode()
func TestGetStatusCode(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase       string
		code           ErrorCode
		expectedStatus int
	}{
		{
			"ErrorBadErrorCode",
			ErrorBadErrorCode,
			http.StatusExpectationFailed,
		},
		{
			"ErrorMissingEnv",
			ErrorMissingEnv,
			http.StatusBadRequest,
		},
		{

			"ErrorInvalidEnv",
			ErrorInvalidEnv,
			http.StatusBadRequest,
		},
		{
			"ErrorReadingConfig",
			ErrorReadingConfig,
			http.StatusExpectationFailed,
		},
		{
			"ErrorViper",
			ErrorViper,
			http.StatusExpectationFailed,
		},
		{
			"ErrorConfigValidation",
			ErrorConfigValidation,
			http.StatusExpectationFailed,
		},
		{
			"ErrorDecryptEnv",
			ErrorDecryptEnv,
			http.StatusExpectationFailed,
		},
		{
			"ErrorLoadingConfig",
			ErrorLoadingConfig,
			http.StatusExpectationFailed,
		},
		{
			"ErrorLoadingCache",
			ErrorLoadingCache,
			http.StatusExpectationFailed,
		},
		{
			"unknown error",
			9999,
			http.StatusExpectationFailed,
		},
	}

	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			code := GetStatusCode(test.code)
			assert.Equal(t, test.expectedStatus, code)
		})
	}
}

// TestGetError tests the method GetError()
func TestGetError(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase     string
		code         ErrorCode
		expectedCode ErrorCode
	}{
		{
			"ErrorBadErrorCode",
			ErrorBadErrorCode,
			ErrorBadErrorCode,
		},
		{
			"ErrorMissingEnv",
			ErrorMissingEnv,
			ErrorMissingEnv,
		},
		{
			"ErrorInvalidEnv",
			ErrorInvalidEnv,
			ErrorInvalidEnv,
		},
		{
			"ErrorReadingConfig",
			ErrorReadingConfig,
			ErrorReadingConfig,
		},
		{
			"ErrorViper",
			ErrorViper,
			ErrorViper,
		},
		{
			"ErrorConfigValidation",
			ErrorConfigValidation,
			ErrorConfigValidation,
		},
		{
			"ErrorDecryptEnv",
			ErrorDecryptEnv,
			ErrorDecryptEnv,
		},
		{
			"ErrorLoadingConfig",
			ErrorLoadingConfig,
			ErrorLoadingConfig,
		},
		{
			"ErrorLoadingCache",
			ErrorLoadingCache,
			ErrorLoadingCache,
		},
		{
			"unknown error",
			9999,
			ErrorBadErrorCode,
		},
	}

	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			errorObj := GetError(test.code)
			assert.Equal(t, test.expectedCode, errorObj.Code)
		})
	}
}
