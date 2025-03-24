package manualtests

import (
	"testing"

	"github.com/joomcode/errorx"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestRefreshState_SoItWillAddNewFieldsIfAdded_SoYouCanJustSetupThoseMissingValues(t *testing.T) {
	t.Skip("Don't run it yet")

	state := NewState()
	err := state.Load()
	require.NoError(t, err)
	err = state.Save()
	require.NoError(t, err)
}

func TestSwitchCurrentUserInState_ToUserWithID(t *testing.T) {
	t.Skip("Don't run it yet")

	state := NewState()
	err := state.Load()
	require.NoError(t, err)

	err = state.UseUserWithID("174DcxCYRWySRtUWSPcKkV7wtzTENjDFzf")
	require.NoError(t, err)

	err = state.Save()
	require.NoError(t, err)
}

func TestCleanupOldUsersByTag(t *testing.T) {
	t.Skip("Don't run it yet")

	tag := "deleted"

	state := NewState()
	err := state.Load()
	require.NoError(t, err)

	state.CleanupOldUsersByTag(tag)

	if lo.Contains(state.CurrentUser().Tags, tag) {
		logger := Logger()
		logger.Warn().Msg("Current user has also the tag that should be removed. Will remove it also.")
		userYaml, err := yaml.Marshal(state.User)
		if err != nil {
			logger.Err(err).Msg("Failed to marshal user state into yaml, so it want be recoverable :( ")
		}
		logger.Warn().Msgf("In case it wasn't intended, put back him manually: \n%s", userYaml)

		old, err := state.GetLastOldUserFromState()
		if err == nil {
			err = state.UseUserWithID(old.ID)
			require.NoError(t, err)
			logger.Warn().Msgf("Switching current user to the last old user %s with note '%s'.", old.ID, old.Note)
			state.CleanupOldUsersByTag(tag)
		} else if errorx.IsOfType(err, NotFound) {
			state.User = User{}
			logger.Warn().Msg("No old user found, so the current user will be removed.")
		} else {
			require.NoError(t, err)
		}
	}

	err = state.Save()
	require.NoError(t, err)
}
