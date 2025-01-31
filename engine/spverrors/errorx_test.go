package spverrors

import (
	"fmt"
	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWrapping(t *testing.T) {
	// repository-package-level error
	dbError := errorx.InternalError.NewSubtype("database_error")

	// 3rd package return error about database
	some3rdPartyError := fmt.Errorf("3rd party error")
	fromRepo := dbError.Wrap(some3rdPartyError, "failed to connect to database")

	// domain-package-level error
	fromDomain := errorx.Decorate(fromRepo, "failed to get data")

	// request-handler-level error
	var statusCode int
	if errorx.IsOfType(fromDomain, errorx.InternalError) {
		statusCode = 500
	} else {
		statusCode = 400
	}

	require.Equal(t, 500, statusCode)

	require.True(t, errorx.IsOfType(fromDomain, errorx.InternalError))
	require.True(t, errorx.IsOfType(fromDomain, dbError))

	// full error message
	fmt.Printf("%+v", fromDomain)
	/*
		failed to get data, cause: common.internal_error.database_error: failed to connect to database, cause: 3rd party error
		 at github.com/bitcoin-sv/spv-wallet/engine/spverrors.TestWrapping()
			E:/Data/Source/4chain/spv-wallet/engine/spverrors/errorx_test.go:16
		 at testing.tRunner()
			C:/Users/Tomec/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.3.windows-amd64/src/testing/testing.go:1690
		 at runtime.goexit()
			C:/Users/Tomec/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.23.3.windows-amd64/src/runtime/asm_amd64.s:1700
	*/

	// short error message
	fmt.Printf("%v", fromDomain)
	/*
		failed to get data, cause: common.internal_error.database_error: failed to connect to database, cause: 3rd party error
	*/
}

func TestWrapping2(t *testing.T) {
	// repository-package-level error
	some3rdPartyError := fmt.Errorf("3rd party error")
	fromRepo := errorx.InternalError.Wrap(some3rdPartyError, "failed to connect to database")

	// domain-package-level error
	fromDomain := errorx.Decorate(fromRepo, "failed to get data from")

	require.True(t, errorx.IsOfType(fromDomain, errorx.InternalError))
}

func TestWrapping3(t *testing.T) {
	// repository-package-level error
	some3rdPartyError := fmt.Errorf("3rd party error")
	fromRepo := errorx.InternalError.Wrap(some3rdPartyError, "failed to connect to database")

	// domain-package-level error
	fromDomain := errorx.Decorate(fromRepo, "failed to get data from")

	require.True(t, errorx.IsOfType(fromDomain, errorx.InternalError))
}

func TestWrapping4(t *testing.T) {
	safeToExposeProperty := errorx.RegisterPrintableProperty("safe_to_expose")
	invalidArgumentProperty := errorx.RegisterPrintableProperty("invalid_argument")

	// some service
	domainValidationErrorType := errorx.IllegalArgument.NewSubtype("domain_validation_error")
	someWrongPaymail := "example.com"
	fromSomeService := domainValidationErrorType.
		New("failed to validate paymail %v", someWrongPaymail).
		WithProperty(safeToExposeProperty, "provided paymail is invalid").
		WithProperty(invalidArgumentProperty, someWrongPaymail)

	// upper layer
	fromDomain := errorx.Decorate(fromSomeService, "failed to validate paymail")

	// request-handler-level error
	var statusCode int
	if errorx.IsOfType(fromDomain, errorx.IllegalArgument) {
		statusCode = 400
	} else {
		statusCode = 500
	}

	safeToExposeMessage := ""
	if msgProperty, ok := fromDomain.Property(safeToExposeProperty); ok {
		safeToExposeMessage = msgProperty.(string)
	}

	require.Equal(t, 400, statusCode)
	require.Equal(t, "provided paymail is invalid", safeToExposeMessage)

	require.Equal(t, "failed to validate paymail, cause: common.illegal_argument.domain_validation_error: failed to validate paymail example.com {invalid_argument: example.com, safe_to_expose: provided paymail is invalid}", fromDomain.Error())
}

func TestWrapping5(t *testing.T) {
	// repository-package-level error
	dbTimeoutError := errorx.InternalError.NewSubtype("database_error", errorx.Timeout())

	// 3rd package return error about database
	some3rdPartyError := fmt.Errorf("3rd party error")
	fromRepo := dbTimeoutError.Wrap(some3rdPartyError, "failed to connect to database")

	// domain-package-level error
	fromDomain := errorx.Decorate(fromRepo, "failed to get data")

	// request-handler-level error
	var statusCode int
	if errorx.IsOfType(fromDomain, errorx.InternalError) {
		statusCode = 500
	} else {
		statusCode = 400
	}

	var errorCode string
	if errorx.HasTrait(fromDomain, errorx.Timeout()) {
		errorCode = "timeout"
	} else {
		errorCode = "unknown"
	}

	require.Equal(t, 500, statusCode)
	require.Equal(t, "timeout", errorCode)
}
