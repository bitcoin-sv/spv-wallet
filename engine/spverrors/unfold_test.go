package spverrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestUnfoldError(t *testing.T) {
	err1 := models.SPVError{Code: "test-err1", Message: "test error1", StatusCode: 404}
	err2 := models.SPVError{Code: "test-err2", Message: "test error2", StatusCode: 400}
	err3 := models.SPVError{Code: "test-err3", Message: "test error3", StatusCode: 500}
	err4 := NewError("test error4")

	testCases := map[string]struct {
		input    error
		expected string
	}{
		"unfold single string error": {
			input:    errors.New("test error"),
			expected: "[*errors.errorString] test error",
		},
		"unfold single SPVError": {
			input:    err1,
			expected: "[models.SPVError(404)] test error1",
		},
		"unfold SPVError wrapped in SPVError": {
			input:    err1.Wrap(err2),
			expected: "[models.SPVError(404)] test error1 -> [models.SPVError(400)] test error2",
		},
		"unfold SPVError wrapped by pkg.Wrap and then wrapped into SPVError": {
			input:    err1.Wrap(pkgerrors.Wrap(err3, "test error2")),
			expected: "[models.SPVError(404)] test error1 -> [*errors.withStack] test error2: test error3 -> [*errors.withMessage] test error2: test error3 -> [models.SPVError(500)] test error3",
		},
		"wrapping with fmt.Errorf": {
			input:    fmt.Errorf("Some error is here %w - in the middle of a string", err1),
			expected: "[*fmt.wrapError] Some error is here test error1 - in the middle of a string -> [models.SPVError(404)] test error1",
		},
		"joining errors": {
			input:    err1.Wrap(errors.Join(err2, fmt.Errorf("test error3"))),
			expected: "[models.SPVError(404)] test error1 -> [*errors.joinError] ([models.SPVError(400)] test error2 AND [*errors.errorString] test error3)",
		},
		"internal error only": {
			input:    err4,
			expected: "[spverrors.Error(500)] test error4",
		},
		"using spverrors.Wrapf": {
			input:    Wrapf(err1, "wrapped error"),
			expected: "[spverrors.Error(500)] wrapped error -> [models.SPVError(404)] test error1",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			unfolded := UnfoldError(tc.input)
			require.Equal(t, tc.expected, unfolded)
		})
	}
}
