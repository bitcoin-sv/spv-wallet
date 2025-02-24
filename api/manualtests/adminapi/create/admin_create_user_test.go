//nolint:nolintlint,revive
package create_test // This should make it easier to run mutations separately from the queries

import (
	"context"
	"testing"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/stretchr/testify/require"
)

func TestCreateUser1(t *testing.T) {
	t.Skip("Don't run it yet")

	xpriv, xpub, err := bip32.GenerateHDKeyPair(bip32.RecommendedSeedLength)
	require.NoError(t, err)

	state := manualtests.NewState()
	err = state.Load()
	require.NoError(t, err)

	user, err := state.NewUser(xpriv, xpub)
	require.NoError(t, err)

	c, err := state.AdminClient()
	require.NoError(t, err)

	res, err := c.CreateUserWithResponse(context.Background(), client.CreateUserJSONRequestBody{
		Paymail: &client.RequestsAddPaymail{
			Address:    user.PaymailAddress(),
			AvatarURL:  user.AvatarURL(),
			PublicName: user.PublicName(),
		},
		PublicKey: user.PublicKey,
	})
	require.NoError(t, err)

	manualtests.Print(res)

	err = state.SaveOnSuccess(res)
	require.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	t.Skip("Don't run it yet")

	manualtests.APICallForAdmin(t).
		CallWithUpdateState(func(state manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
			xpriv, xpub, err := bip32.GenerateHDKeyPair(bip32.RecommendedSeedLength)
			require.NoError(t, err)

			user, err := state.NewUser(xpriv, xpub)
			require.NoError(t, err)

			return c.CreateUserWithResponse(context.Background(), client.CreateUserJSONRequestBody{
				Paymail: &client.RequestsAddPaymail{
					Address:    user.PaymailAddress(),
					AvatarURL:  user.AvatarURL(),
					PublicName: user.PublicName(),
				},
				PublicKey: user.PublicKey,
			})
		})
}
