package manualtests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefreshState_SoItWillAddNewFieldsIfAdded_SoYouCanJustSetupThoseMissingValues(t *testing.T) {
	state := NewState()
	err := state.Load()
	require.NoError(t, err)
	err = state.Save()
	require.NoError(t, err)
}

func TestSwitchCurrentUserInState_ToUserWithID(t *testing.T) {
	state := NewState()
	err := state.Load()
	require.NoError(t, err)

	err = state.UseUserWithID("174DcxCYRWySRtUWSPcKkV7wtzTENjDFzf")
	require.NoError(t, err)

	err = state.Save()
	require.NoError(t, err)
}
