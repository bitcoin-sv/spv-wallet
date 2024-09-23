package models

import (
	"errors"
	"fmt"
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
		require.Equal(t, "'test error' <of type [*errors.errorString]>", unfolded.ToString())
		require.JSONEq(t, `{"msg":"test error","type":"*errors.errorString"}`, string(unfolded.ToJSON()))
	})

	t.Run("unfold single SPVError", func(t *testing.T) {
		err := SPVError{Code: "test-err", Message: "test error", StatusCode: 500}
		unfolded := UnfoldError(err)
		require.NotNil(t, unfolded)
		require.Equal(t, err, unfolded.Err)
		require.Equal(t, "test error", unfolded.Msg)
		require.Equal(t, "models.SPVError", unfolded.Type)
		require.Equal(t, "'test error' <of type [models.SPVError]>", unfolded.ToString())
		require.JSONEq(t, `{"msg":"test error","type":"models.SPVError"}`, string(unfolded.ToJSON()))
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
		require.JSONEq(t, `{
			"msg": "test error1",
			"type": "models.SPVError",
			"causes": [
				{
				"msg": "test error2",
				"type": "models.SPVError"
				}
			]
		}`, string(unfolded.ToJSON()))
	})
	t.Run("unfold SPVError wrapped by pkg.Wrap and then wrapped into SPVError", func(t *testing.T) {
		err1 := SPVError{Code: "test-err1", Message: "test error1", StatusCode: 500}
		err3 := SPVError{Code: "test-err2", Message: "test error3", StatusCode: 500}
		err := err1.Wrap(pkgerrors.Wrap(err3, "test error2"))
		unfolded := UnfoldError(err)
		require.NotNil(t, unfolded)
		require.Equal(t, err3, unfolded.InitialCause().Err)
		require.Equal(t, "'test error1' <of type [models.SPVError]> was caused by { 'test error3' <of type [models.SPVError]> annotated with 'test error2' }", unfolded.ToString())
		require.JSONEq(t, `{
			"msg": "test error1",
			"type": "models.SPVError",
			"causes": [
				{
				"msg": "test error2: test error3",
				"type": "*errors.withStack",
				"causes": [
					{
					"msg": "test error2: test error3",
					"type": "*errors.withMessage",
					"causes": [
						{
						"msg": "test error3",
						"type": "models.SPVError"
						}
					]
					}
				]
				}
			]
		}`, string(unfolded.ToJSON()))
	})
	t.Run("wrapping with fmt.Errorf", func(t *testing.T) {
		err1 := SPVError{Code: "test-err1", Message: "test error1", StatusCode: 500}
		err := fmt.Errorf("Some error is here %w - in the middle of a string", err1)
		unfolded := UnfoldError(err)
		require.NotNil(t, unfolded)
		require.Equal(t, err1, unfolded.InitialCause().Err)
		require.Equal(t, "'Some error is here test error1 - in the middle of a string' <of type [*fmt.wrapError]> was caused by { 'test error1' <of type [models.SPVError]> }", unfolded.ToString())
		require.JSONEq(t, `{
			"msg": "Some error is here test error1 - in the middle of a string",
			"type": "*fmt.wrapError",
			"causes": [
				{
				"msg": "test error1",
				"type": "models.SPVError"
				}
			]
		}`, string(unfolded.ToJSON()))
	})

	t.Run("joining errors", func(t *testing.T) {
		err1 := SPVError{Code: "test-err1", Message: "test error1", StatusCode: 500}
		err2 := SPVError{Code: "test-err2", Message: "test error2", StatusCode: 500}
		err3 := SPVError{Code: "test-err3", Message: "test error3", StatusCode: 500}
		err := err1.Wrap(errors.Join(err2, err3))
		unfolded := UnfoldError(err)
		require.NotNil(t, unfolded)
		require.Equal(t, err2, unfolded.InitialCause().Err)
		require.Equal(t, "'test error1' <of type [models.SPVError]> was caused by { 'test error2' <of type [models.SPVError]> AND 'test error3' <of type [models.SPVError]> }", unfolded.ToString())
		require.JSONEq(t, `{
			"msg": "test error1",
			"type": "models.SPVError",
			"causes": [
				{
				"msg": "test error2\ntest error3",
				"type": "*errors.joinError",
				"causes": [
					{
					"msg": "test error2",
					"type": "models.SPVError"
					},
					{
					"msg": "test error3",
					"type": "models.SPVError"
					}
				]
				}
			]
		}`, string(unfolded.ToJSON()))
	})
}
