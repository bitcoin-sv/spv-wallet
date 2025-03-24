package dictionary

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestGetInternalMessage(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase        string
		expectedMessage string
		code            ErrorCode
	}{
		{
			"ErrorBadErrorCode",
			"failed to find internal error message from error code",
			ErrorBadErrorCode,
		},
		{
			"ErrorMissingEnv",
			"missing required environment variable: %s",
			ErrorMissingEnv,
		},
		{
			"ErrorInvalidEnv",
			"invalid environment variable value: %s",
			ErrorInvalidEnv,
		},
		{
			"ErrorReadingConfig",
			"error reading environment configuration: %s",
			ErrorReadingConfig,
		},
		{
			"ErrorViper",
			"error in viper unmarshal into config.Values: %s",
			ErrorViper,
		},
		{
			"ErrorConfigValidation",
			"error in environment configuration validation: %s",
			ErrorConfigValidation,
		},
		{
			"ErrorDecryptEnv",
			"error in decrypting %s: %s",
			ErrorDecryptEnv,
		},
		{
			"ErrorLoadingConfig",
			"fatal error loading configuration: %s",
			ErrorLoadingConfig,
		},
		{
			"ErrorLoadingCache",
			"failed to enable cache: %s - cache is disabled",
			ErrorLoadingCache,
		},
		{
			"unknown code",
			"failed to find internal error message from error code",
			9999,
		},
	}

	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			message := GetInternalMessage(test.code)
			assert.Equal(t, test.expectedMessage, message)
		})
	}
}

func TestGetPublicMessage(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase        string
		expectedMessage string
		code            ErrorCode
	}{
		{
			"ErrorBadErrorCode",
			"error not found",
			ErrorBadErrorCode,
		},
		{
			"ErrorMissingEnv",
			"missing required environment variable",
			ErrorMissingEnv,
		},
		{
			"ErrorInvalidEnv",
			"invalid environment variable value",
			ErrorInvalidEnv,
		},
		{
			"ErrorReadingConfig",
			"error reading environment configuration",
			ErrorReadingConfig,
		},
		{
			"ErrorViper",
			"error in loading configuration",
			ErrorViper,
		},
		{
			"ErrorConfigValidation",
			"error in environment configuration validation",
			ErrorConfigValidation,
		},
		{
			"ErrorDecryptEnv",
			"error decrypting an encrypted environment variable",
			ErrorDecryptEnv,
		},
		{
			"ErrorLoadingConfig",
			"error loading configuration",
			ErrorLoadingConfig,
		},
		{
			"ErrorLoadingCache",
			"failed to enable cache",
			ErrorLoadingCache,
		},
		{
			"unknown error",
			"error not found",
			9999,
		},
	}

	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			message := GetPublicMessage(test.code)
			assert.Equal(t, test.expectedMessage, message)
		})
	}
}

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
