package spverrors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/models"
)

type joinedErrors interface {
	Unwrap() []error
}

// UnfoldError unfolds the error chain into a single string
// example:
//
//	[errType1] outer error message -> [errType2] inner error message -> [errTypeN] innermost error message
//
// if error is a joined error, it will be unfolded as:
//
//	[type] ([errType1] joined error1 AND [errType2] joined error2)
//
// for SPVError, it will print the status code in parentheses:
//
//	[models.SPVError(404)] error message
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

		if joined, ok := current.(joinedErrors); ok {
			unfoldJoinedErrors(joined, &result)
			break // joined errors cannot be unfolded to keep this as chain of errors instead of a tree
		}

		result.WriteRune(' ')
		result.WriteString(msg)
		printDetailsForInternalError(current, &result)
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

func printDetailsForInternalError(err error, builder *strings.Builder) {
	//nolint:errorlint //errors.Is/As would check also the wrapped error but here only the current one should be concerned
	if internalErr, ok := err.(InternalError); ok {
		builder.WriteString(fmt.Sprintf(" (%s)", internalErr.details))
	}
}

func unfoldJoinedErrors(joined joinedErrors, builder *strings.Builder) {
	errList := joined.Unwrap()
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
