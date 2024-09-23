package models

import (
	"errors"
	"testing"

	pkgerrors "github.com/pkg/errors"

	"github.com/stretchr/testify/require"
)

func TestUnfoldError(t *testing.T) {
	t.Run("unfold single string error", func(t *testing.T) {
		err := errors.New("test error")
		unfolded := UnfoldError(err)
		require.NotNil(t, unfolded)
		require.Equal(t, err, unfolded.Err)
		require.Equal(t, "test error", unfolded.Msg)
		require.Equal(t, "*errors.errorString", unfolded.Type)
	})

	t.Run("unfold single SPVError", func(t *testing.T) {
		err := SPVError{Code: "test-err", Message: "test error", StatusCode: 500}
		unfolded := UnfoldError(err)
		require.NotNil(t, unfolded)
		require.Equal(t, err, unfolded.Err)
		require.Equal(t, "test error", unfolded.Msg)
		require.Equal(t, "models.SPVError", unfolded.Type)
	})

	t.Run("unfold SPVError wrapped in SPVError", func(t *testing.T) {
		err1 := SPVError{Code: "test-err1", Message: "test error1", StatusCode: 500}
		err2 := SPVError{Code: "test-err2", Message: "test error2", StatusCode: 500}
		err := err1.Wrap(err2)
		unfolded := UnfoldError(err)
		require.NotNil(t, unfolded)
		require.Equal(t, unfolded.Msg, "test error1")
		require.Equal(t, unfolded.Type, "models.SPVError")
		require.Len(t, unfolded.Causes, 1)
		require.Equal(t, "test error2", unfolded.Causes[0].Msg)
		require.Equal(t, "models.SPVError", unfolded.Causes[0].Type)
		require.Equal(t, err2, unfolded.InitialCause().Err)
		require.Equal(t, "'test error1' <of type [models.SPVError]> was caused by { 'test error2' <of type [models.SPVError]> }", unfolded.ToString())
	})
	t.Run("unfold SPVError wrapped in string error wrapped by pkg/errors", func(t *testing.T) {
		err1 := SPVError{Code: "test-err1", Message: "test error1", StatusCode: 500}
		err3 := SPVError{Code: "test-err2", Message: "test error3", StatusCode: 500}
		err := err1.Wrap(pkgerrors.Wrap(err3, "test error2"))
		unfolded := UnfoldError(err)
		require.NotNil(t, unfolded)
		require.Equal(t, err3, unfolded.InitialCause().Err)
		require.Equal(t, "'test error1' <of type [models.SPVError]> was caused by { 'test error3' <of type [models.SPVError]> annotated with 'test error2' }", unfolded.ToString())
	})
}
