package models

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestErrorIsForWrappedExtendedErrors(t *testing.T) {
	testedError := SPVError{
		Code:       "some-code1",
		Message:    "some-message1",
		StatusCode: 1,
	}

	otherErr1 := SPVError{
		Code:       "some-code2",
		Message:    "some-message2",
		StatusCode: 2,
	}

	otherErr2 := SPVError{
		Code:       "some-code3",
		Message:    "some-message3",
		StatusCode: 3,
	}

	tests := map[string]struct {
		err error
	}{
		"single": {
			err: testedError,
		},
		"wrapped": {
			err: errors.Wrap(testedError, "wrapped"),
		},
		"double wrapped": {
			err: errors.Wrap(errors.Wrap(testedError, "wrapped"), "double wrapped"),
		},
		"middle wrapped": {
			err: errors.Wrapf(testedError.Wrap(errors.New("source")), "middle wrapped"),
		},
		"as cause": {
			err: otherErr1.Wrap(testedError),
		},
		"as double cause": {
			err: otherErr1.Wrap(otherErr2).Wrap(testedError),
		},
		"as middle cause": {
			err: otherErr1.Wrap(testedError.Wrap(otherErr2)),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			is := errors.Is(test.err, testedError)
			require.True(t, is)
		})
	}
}
