package models

type ExtendedError interface {
	error
	GetCode() string
	GetMessage() string
	GetStatusCode() int
}

// SPVError is extended error which holds information about http status and code that should be returned
type SPVError struct {
	Code       string
	Message    string
	StatusCode int
}

// ResponseError is an error which will be returned in HTTP response
type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const UnknownErrorCode = "error-unknown"

// Error returns the error message string for SPVError, satisfying the error interface
func (e SPVError) Error() string {
	return e.Message
}

// GetCode returns the error code string for SPVError
func (e SPVError) GetCode() string {
	return e.Code
}

// GetMessage returns the error message string for SPVError
func (e SPVError) GetMessage() string {
	return e.Message
}

// GetStatusCode returns the error status code for SPVError
func (e SPVError) GetStatusCode() int {
	return e.StatusCode
}
