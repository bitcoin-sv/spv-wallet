package engine

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestPrepareBeefFactorsForTransactionWithMultipleInputsFromSingleTransaction(t *testing.T) {
	ctx := context.Background()
	tx := &Transaction{
		TransactionBase: TransactionBase{
			Hex: "0100000006523c54c6e2fcde4fd07e2b507f6ec47e5c9b30d58b821dadf7fda53edc60d8320000000000ffffffff523c54c6e2fcde4fd07e2b507f6ec47e5c9b30d58b821dadf7fda53edc60d8320100000000ffffffff523c54c6e2fcde4fd07e2b507f6ec47e5c9b30d58b821dadf7fda53edc60d8320200000000ffffffff523c54c6e2fcde4fd07e2b507f6ec47e5c9b30d58b821dadf7fda53edc60d8320300000000ffffffff523c54c6e2fcde4fd07e2b507f6ec47e5c9b30d58b821dadf7fda53edc60d8320400000000ffffffff523c54c6e2fcde4fd07e2b507f6ec47e5c9b30d58b821dadf7fda53edc60d8320500000000ffffffff020a000000000000001976a914dc59ed4b67ce6f0863dea4aba2b855810409edec88ac4a0c0300000000001976a914170775226cd3ac1316c3109c90a10edc9710c46488ac00000000",
		},
	}
	getter := transactionGetterFn(ReturnOnlySpecificSourceTransaction)

	_, err := prepareBEEFFactors(ctx, tx, getter)
	require.NoError(t, err)
}

type transactionGetterFn func(ctx context.Context, txIDs []string) ([]*Transaction, error)

func (f transactionGetterFn) GetTransactionsByIDs(ctx context.Context, txIDs []string) ([]*Transaction, error) {
	return f(ctx, txIDs)
}

func ReturnOnlySpecificSourceTransaction(_ context.Context, txIDs []string) ([]*Transaction, error) {
	unique := lo.Uniq(txIDs)

	filtered := lo.Filter(unique, func(item string, index int) bool {
		return item == "32d860dc3ea5fdf7ad1d828bd5309b5c7ec46e7f502b7ed04fdefce2c6543c52"
	})

	result := lo.Map(filtered, func(item string, index int) *Transaction {
		var bump BUMP
		bumpJson := []byte(`{"blockHeight":"882514","path":[[{"offset":"10","hash":"32d860dc3ea5fdf7ad1d828bd5309b5c7ec46e7f502b7ed04fdefce2c6543c52","txid":true},{"offset":"11","hash":"bf8de00d36bb634181c4630d9c51c6d927e8f3826dd581f01e67f38f7c90086d"}],[{"offset":"4","hash":"ed67e545b08e1a27bfb71bd3a47418464301efe73413f8d1d0a70d18d1a9f041"}],[{"offset":"3","hash":"497294b4e4c60b6d70a6ccf3daa42f93ad037e3e6a18b51d0fe081c970e04071"}],[{"offset":"0","hash":"e93b685d58d8a4f39698cc9bbc560792f3eb265bf61a9ead0b97e3c9d07f8b58"}],[{"offset":"1","hash":"7486f71b5b6a540e64cf5a59f89a703139366357f5884edcc6921f2bb3cf90c7"}],[{"offset":"1","hash":"a2e05afba28db0a75288489cdcad5fdf5e188add527db8c7ef733ddb0190ee96"}],[{"offset":"1","hash":"e5bd1cc5b2d367f2aa7807b9bbc805286b5e6708e89cbb1b1a930dcf37512ec6"}],[{"offset":"1","hash":"1aefdcee271cf8c3532265873156b566516a627dcf4fda0ef280a6a54a919d67"}]]}`)
		err := json.Unmarshal(bumpJson, &bump)
		if err != nil {
			panic(err)
		}

		return &Transaction{
			TransactionBase: TransactionBase{
				ID:  item,
				Hex: "0100000001801faaf6afe4e4eb20eb24ca699954e0fb2c8f97eada9bf8ad209e199e0b737c010000006a473044022006b87316deca9116dff17b143285f33bb30be0887809b7774e01574bbe9cfe9e02203fd2357bb295382fbb8ba17b99cfafb64cb126790a5ac50353fac74e32667833412102f35c2a01ccad456e095ba4aeae66e458047d02b7f181ae87a5521081d30f4a30ffffffff0601000000000000001976a914d79b1757654b2bb53470eb40a6586e7c2fa6d7e788ac01000000000000001976a914d79b1757654b2bb53470eb40a6586e7c2fa6d7e788ac01000000000000001976a914d79b1757654b2bb53470eb40a6586e7c2fa6d7e788ac01000000000000001976a914d79b1757654b2bb53470eb40a6586e7c2fa6d7e788ac01000000000000001976a914d79b1757654b2bb53470eb40a6586e7c2fa6d7e788ac500c0300000000001976a91455b8f7fa309ac27405b2e6fbc2b1089bd53b5bff88ac00000000",
			},
			BUMP: bump,
		}
	})

	return result, nil
}
