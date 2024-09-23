package spverrors

import (
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/rs/zerolog"
)

type marshaller struct {
	errorNode *models.ErrorNode
}

// MarshalZerologObject converts the errorNode to zerolog Event object
func (m marshaller) MarshalZerologObject(e *zerolog.Event) {
	e.Str("error", m.errorNode.Msg)
	if initialCause := m.errorNode.InitialCause(); initialCause != m.errorNode {
		e.Str("initialCause", m.errorNode.InitialCause().Msg)
	}
}

// SetupGlobalZerologErrorHandler setup the ErrorMarshalFunc to print detailed error info, depends on log level
func SetupGlobalZerologErrorHandler(level zerolog.Level) {
	if level == zerolog.InfoLevel {
		zerolog.ErrorMarshalFunc = func(err error) any {
			return marshaller{errorNode: models.UnfoldError(err)}
		}
	} else if level == zerolog.DebugLevel {
		zerolog.ErrorMarshalFunc = func(err error) any {
			return models.UnfoldError(err).ToString()
		}
	}
	// default for other levels
}
