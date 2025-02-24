package userapi

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/samber/lo"
)

func TestUserAPIRead(t *testing.T) {
	var calls = map[string]struct {
		call manualtests.CallWithState
	}{
		"currentUser": {
			call: func(_ manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
				return c.CurrentUserWithResponse(context.Background())
			},
		},
		"searchOperations": {
			call: func(_ manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
				return c.SearchOperationsWithResponse(context.Background(), nil)
			},
		},
		"searchOperationsWithQueryParams": {
			call: func(_ manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
				return c.SearchOperationsWithResponse(context.Background(), &client.SearchOperationsParams{
					Page:   lo.ToPtr(1),
					Size:   lo.ToPtr(10),
					Sort:   lo.ToPtr("asc"),
					SortBy: lo.ToPtr("tx_id"),
				})
			},
		},
		"dataById": {
			call: func(state manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
				if state.DataID == "" {
					state.T.Skip("no data id")
				}
				return c.DataByIdWithResponse(context.Background(), state.DataID)
			},
		},
	}
	for name, endpoint := range calls {
		t.Run(name, func(t *testing.T) {
			manualtests.APICallForUser(t).CallWithState(endpoint.call).RequireSuccess()
		})
	}
}
