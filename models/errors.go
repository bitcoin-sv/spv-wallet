package models

// SPVError is extended error which holds information about http status and code that should be returned
type SPVError struct {
	Code       string
	Message    string
	StatusCode int
}

// ResponseError is an error which will be returned in HTTP response
type ResponseError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

const UnknownErrorCode = "error-unknown"

func (e SPVError) Error() string {
	return e.Message
}
