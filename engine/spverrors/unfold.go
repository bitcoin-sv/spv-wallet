package spverrors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/models"
)

type errListUnwrapper interface {
	Unwrap() []error
}

// UnfoldError unfolds the error chain into a single string.
func UnfoldError(err error) string {
	if err == nil {
		return ""
	}

	prevMsg := ""
	result := strings.Builder{}
	for current := err; current != nil; current = errors.Unwrap(current) {
		msg := current.Error()

		if prevMsg != "" {
			result.WriteString(" -> ")
		}

		printTypename(current, &result)

		if errListUnwrapper, ok := current.(errListUnwrapper); ok {
			unfoldJoinedErrors(errListUnwrapper, &result)
			break // joned errors cannot be unfolded to keep this as chain of errors instead of a tree
		}

		if prevMsg == "" || !strings.Contains(prevMsg, msg) {
			result.WriteRune(' ')
			result.WriteString(msg)
		}
		prevMsg = msg
	}
	return result.String()
}

func printTypename(err error, builder *strings.Builder) {
	typename := reflect.TypeOf(err).String()
	builder.WriteRune('[')
	builder.WriteString(typename)
	printStatusCodeForSPVError(err, builder)
	builder.WriteRune(']')
}

func printStatusCodeForSPVError(err error, builder *strings.Builder) {
	//nolint:errorlint //errors.Is/As would check also the wrapped error but here only the current one should be concerned
	if spvErr, ok := err.(models.SPVError); ok {
		builder.WriteString(fmt.Sprintf("(%d)", spvErr.GetStatusCode()))
	}
}

func unfoldJoinedErrors(unwrapper errListUnwrapper, builder *strings.Builder) {
	errList := unwrapper.Unwrap()
	if len(errList) == 0 {
		return
	}
	builder.WriteString(" (")
	for i, jErr := range errList {
		if i > 0 {
			builder.WriteString(" AND ")
		}
		printTypename(jErr, builder)
		builder.WriteRune(' ')
		builder.WriteString(jErr.Error())
	}
	builder.WriteRune(')')
}
