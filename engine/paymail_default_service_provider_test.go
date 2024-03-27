package engine

import (
	"testing"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PaymailDefaultServiceProvider(t *testing.T) {
	t.Run("PaymailDefaultServiceProvider.GetPaymailByAlias", test_GetPaymailByAlias)
	t.Run("PaymailDefaultServiceProvider.GetPaymailByAlias - multiple call", test_GetPaymailByAlias_MultipleRequest_ShouldReturnStablePubKey)
	t.Run("PaymailDefaultServiceProvider.CreateAddressResolutionResponse - multiple call", test_CreateAddressResolutionResponse_ShouldReturnDifferentResponses)
	t.Run("PaymailDefaultServiceProvider.CreateP2PDestinationResponse - multiple call", test_CreateP2PDestinationResponse_ShouldReturnDifferentResponses)

}

func test_GetPaymailByAlias(t *testing.T) {
	// given
	ctx, c, deferMe := CreateTestSQLiteClient(t, false, false, WithAutoMigrate(&PaymailAddress{}), WithFreeCache())
	defer deferMe()

	pm := newPaymail("paymail@domain.sc", 0, WithClient(c), WithXPub(testXPub))
	err := pm.Save(ctx)
	require.NoError(t, err)

	sut := &PaymailDefaultServiceProvider{client: c}

	// when
	res, err := sut.GetPaymailByAlias(ctx, pm.Alias, pm.Domain, nil)
	require.NoError(t, err)

	// then

	assert.Equal(t, pm.ID, res.ID)
	assert.Equal(t, pm.Alias, res.Alias)
	assert.Equal(t, pm.Avatar, res.Avatar)
	assert.Equal(t, pm.Domain, res.Domain)
	assert.Equal(t, pm.PublicName, res.Name)

	expectedPk, err := pm.GetPubKey()
	require.NoError(t, err)
	assert.Equal(t, expectedPk, res.PubKey)

}

func test_GetPaymailByAlias_MultipleRequest_ShouldReturnStablePubKey(t *testing.T) {
	// given
	ctx, c, deferMe := CreateTestSQLiteClient(t, false, false, WithAutoMigrate(&PaymailAddress{}), WithFreeCache())
	defer deferMe()

	pm := newPaymail("paymail@domain.sc", 0, WithClient(c), WithXPub(testXPub))
	err := pm.Save(ctx)
	require.NoError(t, err)

	sut := &PaymailDefaultServiceProvider{client: c}

	// when
	expectedRes, err := sut.GetPaymailByAlias(ctx, pm.Alias, pm.Domain, nil)
	require.NoError(t, err)

	for i := 0; i < 100; i++ {
		res, err := sut.GetPaymailByAlias(ctx, pm.Alias, pm.Domain, nil)
		require.NoErrorf(t, err, "error in %d iteration", i)

		// then
		require.Equalf(t, expectedRes.PubKey, res.PubKey, "different pub key return in %d iteration")
	}

}

func test_CreateAddressResolutionResponse_ShouldReturnDifferentResponses(t *testing.T) {
	// given
	ctx, c, deferMe := CreateTestSQLiteClient(t, false, false, WithAutoMigrate(&PaymailAddress{}), WithFreeCache())
	defer deferMe()

	pm := newPaymail("paymail@domain.sc", 0, WithClient(c), WithXPub(testXPub))
	err := pm.Save(ctx)
	require.NoError(t, err)

	sut := &PaymailDefaultServiceProvider{client: c}

	// when
	results := make([]*paymail.ResolutionPayload, 0)

	for i := 0; i < 100; i++ {
		res, err := sut.CreateAddressResolutionResponse(ctx, pm.Alias, pm.Domain, false, nil)
		require.NoErrorf(t, err, "error in %d iteration", i)

		results = append(results, res)
	}

	// then
	seen := make([]*paymail.ResolutionPayload, 0)
	for _, res := range results {

		for _, seenRes := range seen {
			require.NotEqual(t, res.Address, seenRes.Address)
			require.NotEqual(t, res.Output, seenRes.Output)
		}

		seen = append(seen, res)
	}
}

func test_CreateP2PDestinationResponse_ShouldReturnDifferentResponses(t *testing.T) {
	// given
	ctx, c, deferMe := CreateTestSQLiteClient(t, false, false, WithAutoMigrate(&PaymailAddress{}), WithFreeCache())
	defer deferMe()

	pm := newPaymail("paymail@domain.sc", 0, WithClient(c), WithXPub(testXPub))
	err := pm.Save(ctx)
	require.NoError(t, err)

	sut := &PaymailDefaultServiceProvider{client: c}

	// when
	results := make([]*paymail.PaymentDestinationPayload, 0)

	for i := 0; i < 100; i++ {
		res, err := sut.CreateP2PDestinationResponse(ctx, pm.Alias, pm.Domain, uint64(100), nil)
		require.NoErrorf(t, err, "error in %d iteration", i)

		results = append(results, res)
	}

	// then
	seen := make([]*paymail.PaymentDestinationPayload, 0)
	for _, res := range results {
		for _, out := range res.Outputs {

			for _, seenRes := range seen {

				for _, seenOut := range seenRes.Outputs {

					require.NotEqual(t, out.Address, seenOut.Address)
					require.NotEqual(t, out.Script, seenOut.Script)
				}
			}
		}
		seen = append(seen, res)
	}
}
