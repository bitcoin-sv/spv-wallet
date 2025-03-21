package userapi

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/samber/lo"
)

func TestCurrentUser(t *testing.T) {
	manualtests.APICallForCurrentUser(t).Call(func(c *client.ClientWithResponses) (manualtests.Result, error) {
		return c.CurrentUserWithResponse(context.Background())
	}).RequireSuccess()
}

func TestSearchOperations(t *testing.T) {
	manualtests.APICallForCurrentUser(t).Call(func(c *client.ClientWithResponses) (manualtests.Result, error) {
		return c.SearchOperationsWithResponse(context.Background(), nil)
	}).RequireSuccess()
}

func TestSearchOperationsWithQueryParams(t *testing.T) {
	manualtests.APICallForCurrentUser(t).Call(func(c *client.ClientWithResponses) (manualtests.Result, error) {
		return c.SearchOperationsWithResponse(context.Background(), &client.SearchOperationsParams{
			Page:   lo.ToPtr(1),
			Size:   lo.ToPtr(10),
			Sort:   lo.ToPtr("asc"),
			SortBy: lo.ToPtr("tx_id"),
		})
	}).RequireSuccess()
}

func TestDataById(t *testing.T) {
	manualtests.APICallForCurrentUser(t).CallWithState(func(state manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
		dataId := state.LatestDataID()
		if dataId == "" {
			state.T.Skip("no data id")
		}
		return c.DataByIdWithResponse(context.Background(), dataId)
	}).RequireSuccess()
}

func TestMerkleRootsWithoutParams(t *testing.T) {
	manualtests.APICallForCurrentUser(t).CallForSuccess(func(c *client.ClientWithResponses) (manualtests.Result, error) {
		return c.MerkleRootsWithResponse(context.Background(), &client.MerkleRootsParams{
			BatchSize:        lo.ToPtr(5),
			LastEvaluatedKey: lo.ToPtr("d727e7764d395318ad17527346343532664483d9428c7afbf313a7f414bc76e2"),
		})
	})
}
