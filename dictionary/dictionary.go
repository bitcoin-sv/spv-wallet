// Package dictionary is for our standardized error responses
package dictionary

import (
	"fmt"
	"net/http"
)

// ErrorMessage is used for standardizing the messages/codes for errors
type ErrorMessage struct {
	Code            ErrorCode `json:"code"`
	StatusCode      int       `json:"status_code"`
	InternalMessage string    `json:"internal_message"`
	PublicMessage   string    `json:"public_message"`
}

// errorMessages is the public
var errorMessages map[ErrorCode]ErrorMessage

// GetError will return the error based on the given error code
// Set replacements for Sprintf() on the internal message if needed
func GetError(code ErrorCode, internalMessageReplacements ...interface{}) ErrorMessage {
	if val, ok := errorMessages[code]; ok {
		if len(internalMessageReplacements) > 0 { // strings.Contains(val.InternalMessage,"%")
			val.InternalMessage = fmt.Sprintf(val.InternalMessage, internalMessageReplacements...)
		}
		return val
	}
	return errorMessages[ErrorBadErrorCode]
}

// GetInternalMessage will return the internal message based on the given error code
func GetInternalMessage(code ErrorCode) string {
	if val, ok := errorMessages[code]; ok {
		return val.InternalMessage
	}
	return errorMessages[ErrorBadErrorCode].InternalMessage
}

// GetPublicMessage will return the public message based on the given error code
func GetPublicMessage(code ErrorCode) string {
	if val, ok := errorMessages[code]; ok {
		return val.PublicMessage
	}
	return errorMessages[ErrorBadErrorCode].PublicMessage
}

// GetStatusCode will return the http status code on the given error code
func GetStatusCode(code ErrorCode) int {
	if val, ok := errorMessages[code]; ok {
		return val.StatusCode
	}
	return errorMessages[ErrorBadErrorCode].StatusCode
}

// ErrorCode is a unique code for a specific error message
type ErrorCode uint16

// ErrorCode constants for the different error responses
// Leave the start and last constants in place
const (
	_                           ErrorCode = iota
	ErrorBadErrorCode                     = 1
	ErrorConfigValidation                 = 2
	ErrorDatabaseOpen                     = 3
	ErrorDatabasePing                     = 4
	ErrorDatabaseUnknownDriver            = 5
	ErrorDecryptEnv                       = 6
	ErrorInvalidEnv                       = 7
	ErrorLoadingCache                     = 8
	ErrorLoadingConfig                    = 9
	ErrorMethodNotAllowed                 = 10
	ErrorMissingEnv                       = 11
	ErrorReadingConfig                    = 12
	ErrorRequestNotFound                  = 13
	ErrorViper                            = 14
	ErrorLoadingService                   = 15
	ErrorAuthenticationError              = 16
	ErrorAuthenticationScheme             = 17
	ErrorAuthenticationNotAdmin           = 18
	ErrorAuthenticationCallback           = 19

	errorCodeLast = iota
)

// IsValid tests if the error code is valid or not
func (e ErrorCode) IsValid() bool {
	return e > 0 && e < errorCodeLast
}

// Load all error messages on init
func init() {
	// Create all the standard error messages
	errorMessages = make(map[ErrorCode]ErrorMessage, errorCodeLast)

	// First error - could not find corresponding error via the code given
	errorMessages[ErrorBadErrorCode] = ErrorMessage{Code: ErrorBadErrorCode, InternalMessage: "failed to find internal error message from error code", PublicMessage: "error not found", StatusCode: http.StatusExpectationFailed}

	// Loading environment variables
	errorMessages[ErrorMissingEnv] = ErrorMessage{Code: ErrorMissingEnv, InternalMessage: "missing required environment variable: %s", PublicMessage: "missing required environment variable", StatusCode: http.StatusBadRequest}
	errorMessages[ErrorInvalidEnv] = ErrorMessage{Code: ErrorInvalidEnv, InternalMessage: "invalid environment variable value: %s", PublicMessage: "invalid environment variable value", StatusCode: http.StatusBadRequest}
	errorMessages[ErrorReadingConfig] = ErrorMessage{Code: ErrorReadingConfig, InternalMessage: "error reading environment configuration: %s", PublicMessage: "error reading environment configuration", StatusCode: http.StatusExpectationFailed}
	errorMessages[ErrorViper] = ErrorMessage{Code: ErrorViper, InternalMessage: "error in viper unmarshal into config.Values: %s", PublicMessage: "error in loading configuration", StatusCode: http.StatusExpectationFailed}
	errorMessages[ErrorConfigValidation] = ErrorMessage{Code: ErrorConfigValidation, InternalMessage: "error in environment configuration validation: %s", PublicMessage: "error in environment configuration validation", StatusCode: http.StatusExpectationFailed}
	errorMessages[ErrorDecryptEnv] = ErrorMessage{Code: ErrorDecryptEnv, InternalMessage: "error in decrypting %s: %s", PublicMessage: "error decrypting an encrypted environment variable", StatusCode: http.StatusExpectationFailed}

	// Loading the application
	errorMessages[ErrorLoadingConfig] = ErrorMessage{Code: ErrorLoadingConfig, InternalMessage: "fatal error loading configuration: %s", PublicMessage: "error loading configuration", StatusCode: http.StatusExpectationFailed}
	errorMessages[ErrorLoadingCache] = ErrorMessage{Code: ErrorLoadingCache, InternalMessage: "failed to enable cache: %s - cache is disabled", PublicMessage: "failed to enable cache", StatusCode: http.StatusExpectationFailed}
	errorMessages[ErrorLoadingService] = ErrorMessage{Code: ErrorLoadingService, InternalMessage: "fatal error loading service: %s: %s", PublicMessage: "error loading service", StatusCode: http.StatusExpectationFailed}

	// Database
	errorMessages[ErrorDatabaseUnknownDriver] = ErrorMessage{Code: ErrorDatabaseUnknownDriver, InternalMessage: "unknown database driver specified: %s", PublicMessage: "failed to connect to database", StatusCode: http.StatusExpectationFailed}
	errorMessages[ErrorDatabaseOpen] = ErrorMessage{Code: ErrorDatabaseOpen, InternalMessage: "database open error: %s at address: %s error: %s", PublicMessage: "failed to connect to database", StatusCode: http.StatusExpectationFailed}
	errorMessages[ErrorDatabasePing] = ErrorMessage{Code: ErrorDatabasePing, InternalMessage: "database ping error: %s at address: %s error: %s", PublicMessage: "failed to connect to database", StatusCode: http.StatusExpectationFailed}

	// Router - Basics
	errorMessages[ErrorMethodNotAllowed] = ErrorMessage{Code: ErrorMethodNotAllowed, InternalMessage: "method: %s is not allowed, request: %s", PublicMessage: "method not allowed", StatusCode: http.StatusMethodNotAllowed}
	errorMessages[ErrorRequestNotFound] = ErrorMessage{Code: ErrorRequestNotFound, InternalMessage: "request: %s was not found", PublicMessage: "request was not found", StatusCode: http.StatusNotFound}

	// Authentication
	errorMessages[ErrorAuthenticationError] = ErrorMessage{Code: ErrorAuthenticationError, InternalMessage: "authentication error: %s", PublicMessage: "authentication failed", StatusCode: http.StatusUnauthorized}
	errorMessages[ErrorAuthenticationScheme] = ErrorMessage{Code: ErrorAuthenticationScheme, InternalMessage: "authentication scheme unknown: %s", PublicMessage: "authentication failed", StatusCode: http.StatusUnauthorized}
	errorMessages[ErrorAuthenticationNotAdmin] = ErrorMessage{Code: ErrorAuthenticationNotAdmin, InternalMessage: "xpub provided is not an admin key: %s", PublicMessage: "authentication failed", StatusCode: http.StatusUnauthorized}
	errorMessages[ErrorAuthenticationCallback] = ErrorMessage{Code: ErrorAuthenticationCallback, InternalMessage: "callback authentication failed: %s", PublicMessage: "authentication failed", StatusCode: http.StatusUnauthorized}
}
