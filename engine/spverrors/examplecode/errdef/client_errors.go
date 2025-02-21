package errdef

import "github.com/joomcode/errorx"

type ClientError struct {
	title    string
	httpCode int
	code     string
}

var propClientError = errorx.RegisterProperty("client_error")

var ClientUnauthorized = ClientError{
	title:    "Unauthorized",
	httpCode: 401,
	code:     "unauthorized",
}

var ClientNotFound = ClientError{
	title:    "Not Found",
	httpCode: 404,
	code:     "not_found",
}

var ClientBadRequest = ClientError{
	title:    "Bad Request",
	httpCode: 400,
	code:     "bad_request",
}

var ClientUnprocessableEntity = ClientError{
	title:    "Unprocessable Entity",
	httpCode: 422,
	code:     "unprocessable_entity",
}

func AsClientError(err error, details ClientError) error {
	if e := errorx.Cast(err); e != nil {
		return e.WithProperty(propClientError, details)
	}
	return err
}
