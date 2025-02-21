package logging

import (
	"errors"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog"
)

// SetupGlobalZerologErrorHandler setup the ErrorMarshalFunc to print detailed error info
func SetupGlobalZerologErrorHandler() {
	zerolog.ErrorMarshalFunc = func(err error) any {
		var ex *errorx.Error
		if errors.As(err, &ex) {
			// TODO: how to put stack trace ("%+v") into the log?
			return fmt.Sprintf("%v", ex)
		}
		return spverrors.UnfoldError(err)
	}
}
