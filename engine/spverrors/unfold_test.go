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

	t.Run("unfold single string error", func(t *testing.T) {
		err := errors.New("test error")
		unfolded := UnfoldError(err)
		require.Equal(t, "[*errors.errorString] test error", unfolded)
	})

	t.Run("unfold single SPVError", func(t *testing.T) {
		unfolded := UnfoldError(err1)
		require.Equal(t, "[models.SPVError(404)] test error1", unfolded)
	})

	t.Run("unfold SPVError wrapped in SPVError", func(t *testing.T) {
		err := err1.Wrap(err2)
		unfolded := UnfoldError(err)
		require.Equal(t, "[models.SPVError(404)] test error1 -> [models.SPVError(400)] test error2", unfolded)
	})
	t.Run("unfold SPVError wrapped by pkg.Wrap and then wrapped into SPVError", func(t *testing.T) {
		err := err1.Wrap(pkgerrors.Wrap(err3, "test error2"))
		unfolded := UnfoldError(err)
		require.Equal(t, "[models.SPVError(404)] test error1 -> [*errors.withStack] test error2: test error3 -> [*errors.withMessage] test error2: test error3 -> [models.SPVError(500)] test error3", unfolded)
	})
	t.Run("wrapping with fmt.Errorf", func(t *testing.T) {
		err := fmt.Errorf("Some error is here %w - in the middle of a string", err1)
		unfolded := UnfoldError(err)
		require.Equal(t, "[*fmt.wrapError] Some error is here test error1 - in the middle of a string -> [models.SPVError(404)] test error1", unfolded)
	})

	t.Run("joining errors", func(t *testing.T) {
		err := err1.Wrap(errors.Join(err2, fmt.Errorf("test error3")))
		unfolded := UnfoldError(err)
		require.Equal(t, "[models.SPVError(404)] test error1 -> [*errors.joinError] ([models.SPVError(400)] test error2 AND [*errors.errorString] test error3)", unfolded)
	})
}
