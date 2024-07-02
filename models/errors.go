package models

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
