package spverrors

import (
	"github.com/rs/zerolog"
)

// SetupGlobalZerologErrorHandler setup the ErrorMarshalFunc to print detailed error info
func SetupGlobalZerologErrorHandler() {
	zerolog.ErrorMarshalFunc = func(err error) any {
		return UnfoldError(err)
	}
}
