package logging

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/rs/zerolog"
)

// SetupGlobalZerologErrorHandler setup the ErrorMarshalFunc to print detailed error info
func SetupGlobalZerologErrorHandler() {
	zerolog.ErrorMarshalFunc = func(err error) any {
		return spverrors.UnfoldError(err)
	}
}
