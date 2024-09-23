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
		require.Equal(t, unfolded.Err, err)
		require.Equal(t, unfolded.Msg, "test error")
		require.Equal(t, unfolded.Type, "*errors.errorString")
	})

	t.Run("unfold single SPVError", func(t *testing.T) {
		err := SPVError{Code: "test-err", Message: "test error", StatusCode: 500}
		unfolded := UnfoldError(err)
		require.NotNil(t, unfolded)
		require.Equal(t, unfolded.Err, err)
		require.Equal(t, unfolded.Msg, "test error")
		require.Equal(t, unfolded.Type, "models.SPVError")
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
		require.Equal(t, unfolded.Causes[0].Msg, "test error2")
		require.Equal(t, unfolded.Causes[0].Type, "models.SPVError")
		require.Equal(t, unfolded.InitialCause().Err, err2)
		require.Equal(t, unfolded.ToString(), "'test error1' <of type [models.SPVError]> was caused by { 'test error2' <of type [models.SPVError]> }")
	})
	t.Run("unfold SPVError wrapped in string error wrapped by pkg/errors", func(t *testing.T) {
		err1 := SPVError{Code: "test-err1", Message: "test error1", StatusCode: 500}
		err3 := SPVError{Code: "test-err2", Message: "test error3", StatusCode: 500}
		err := err1.Wrap(pkgerrors.Wrap(err3, "test error2"))
		unfolded := UnfoldError(err)
		require.NotNil(t, unfolded)
		require.Equal(t, unfolded.Msg, "test error1")
		require.Equal(t, unfolded.Type, "models.SPVError")
		require.Len(t, unfolded.Causes, 1)
		require.Equal(t, unfolded.Causes[0].Msg, "test error2")
		require.Equal(t, unfolded.Causes[0].Type, "*errors.fundamental")
		require.Equal(t, unfolded.InitialCause().Err, err3)
		require.Equal(t, unfolded.ToString(), "'test error1' <of type [models.SPVError]> was caused by { 'test error2' <of type [*errors.fundamental]> was caused by { 'test error3' <of type [models.SPVError]> } }")
	})
}
