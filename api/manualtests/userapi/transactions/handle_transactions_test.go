package transactions

import (
	"context"
	hexEncoding "encoding/hex"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestTopUp(t *testing.T) {
	// t.Skip("don't run yet")

	state := manualtests.NewState()
	err := state.Load()
	require.NoError(t, err)

	err = state.Faucet.TopUp()
	require.NoError(t, err)

	api := manualtests.APICallForUser(t)

	api.CallForSuccess(RequestCurrentUser)
	api.CallForSuccess(RequestOperations)
}

func TestTransactionWithStringsData(t *testing.T) {
	t.Skip("don't run yet")

	logger := manualtests.Logger()

	api := manualtests.APICallForUser(t)

	logger.Info().Msg("step 1: Create Transaction Outline with Strings OpReturn")
	storedOutline, callForStringsOpReturnTransactionOutline := manualtests.StoreResponse(RequestStringsOpReturnTransactionOutline)
	api.CallWithStateForSuccess(callForStringsOpReturnTransactionOutline)
	outline := storedOutline.MustGetResponse().JSON200

	logger.Info().Msg("step 2: Unlock Outline")
	hex, err := UnlockOutline(outline, api.State())
	require.NoError(t, err)

	logger.Info().Msg("step 3: Record Outline")
	storedRecord, callRecordOutline := manualtests.StoreResponseOfCall(RequestRecordOutline(hex, outline))
	api.CallForSuccess(callRecordOutline)
	recorded := storedRecord.MustGetResponse().JSON201

	logger.Info().Msg("step 4: Save Data Outpoint in state")
	err = api.State().SaveUserDataOutpoint(recorded.TxID, 0)
	require.NoError(t, err)

	logger.Info().Msg("step 5: Additional Calls for checking state")
	api.CallForSuccess(RequestCurrentUser)
	api.CallForSuccess(RequestOperations)

	api.CallWithStateForSuccess(RequestData)
}

func TestTransactionWithBytesData(t *testing.T) {
	t.Skip("don't run yet")

	logger := manualtests.Logger()

	api := manualtests.APICallForUser(t)

	logger.Info().Msg("step 1: Create Transaction Outline with Strings OpReturn")
	storedOutline, callForStringsOpReturnTransactionOutline := manualtests.StoreResponse(RequestBytesOpReturnTransactionOutline)
	api.CallWithStateForSuccess(callForStringsOpReturnTransactionOutline)
	outline := storedOutline.MustGetResponse().JSON200

	logger.Info().Msg("step 2: Unlock Outline")
	hex, err := UnlockOutline(outline, api.State())
	require.NoError(t, err)

	logger.Info().Msg("step 3: Record Outline")
	storedRecord, callRecordOutline := manualtests.StoreResponseOfCall(RequestRecordOutline(hex, outline))
	api.CallForSuccess(callRecordOutline)
	recorded := storedRecord.MustGetResponse().JSON201

	logger.Info().Msg("step 4: Save Data Outpoint in state")
	err = api.State().SaveUserDataOutpoint(recorded.TxID, 0)
	require.NoError(t, err)

	logger.Info().Msg("step 5: Additional Calls for checking state")
	api.CallForSuccess(RequestCurrentUser)
	api.CallForSuccess(RequestOperations)

	api.CallWithStateForSuccess(RequestData)
}

func UnlockOutline(outline *client.ResponsesCreateTransactionOutlineSuccess, state *manualtests.State) (string, error) {
	format := string(outline.Format)

	logger := manualtests.Logger().With().Str("format", format).Logger()

	logger.Info().Bool("signed", false).Msgf("%s", outline.Hex)

	hex, err := state.UnlockOutlineHex(outline)
	if err != nil {
		return "", err
	}

	logger.Info().Str("format", format).Bool("signed", true).Msgf("%s", hex)

	return hex, nil
}

func RequestStringsOpReturnTransactionOutline(state manualtests.StateForCall, c *client.ClientWithResponses) (*client.CreateTransactionOutlineResponse, error) {
	var opReturn client.RequestsOpReturnOutputSpecification_Data
	err := opReturn.FromRequestsOpReturnStringsOutput([]string{"test", " ", time.Now().Format("2006-01-02T15:04:05")})
	require.NoError(state.T, err)

	opReturnOutput := client.RequestsOpReturnOutputSpecification{
		Data:     opReturn,
		DataType: lo.ToPtr(client.Strings),
	}

	var output client.RequestsTransactionOutlineOutputSpecification
	err = output.FromRequestsOpReturnOutputSpecification(opReturnOutput)
	require.NoError(state.T, err)

	var body client.CreateTransactionOutlineJSONRequestBody
	body.Outputs = append(body.Outputs, output)

	return c.CreateTransactionOutlineWithResponse(context.Background(), &client.CreateTransactionOutlineParams{
		Format: lo.ToPtr(client.Beef),
	}, body)
}

func RequestBytesOpReturnTransactionOutline(state manualtests.StateForCall, c *client.ClientWithResponses) (*client.CreateTransactionOutlineResponse, error) {
	var opReturn client.RequestsOpReturnOutputSpecification_Data

	dataHexPart1 := hexEncoding.EncodeToString([]byte("test bytes"))
	dataHexPart2 := hexEncoding.EncodeToString([]byte(" "))
	dataHexPart3 := hexEncoding.EncodeToString([]byte(time.Now().Format("2006-01-02T15:04:05")))

	err := opReturn.FromRequestsOpReturnHexesOutput([]string{dataHexPart1, dataHexPart2, dataHexPart3})
	require.NoError(state.T, err)

	opReturnOutput := client.RequestsOpReturnOutputSpecification{
		Data:     opReturn,
		DataType: lo.ToPtr(client.Strings),
	}

	var output client.RequestsTransactionOutlineOutputSpecification
	err = output.FromRequestsOpReturnOutputSpecification(opReturnOutput)
	require.NoError(state.T, err)

	var body client.CreateTransactionOutlineJSONRequestBody
	body.Outputs = append(body.Outputs, output)

	return c.CreateTransactionOutlineWithResponse(context.Background(), &client.CreateTransactionOutlineParams{
		Format: lo.ToPtr(client.Beef),
	}, body)
}

func RequestRecordOutline(hex string, outline *client.ResponsesCreateTransactionOutlineSuccess) manualtests.GenericCall[*client.RecordTransactionOutlineResponse] {
	return func(c *client.ClientWithResponses) (*client.RecordTransactionOutlineResponse, error) {
		body := client.RecordTransactionOutlineJSONRequestBody{
			Format: client.RequestsTransactionOutlineFormat(outline.Format),
			Hex:    hex,
			Annotations: &client.ModelsOutputsAnnotations{
				Outputs: outline.Annotations.Outputs,
			},
		}

		return c.RecordTransactionOutlineWithResponse(context.Background(), body)
	}
}

func RequestCurrentUser(c *client.ClientWithResponses) (manualtests.Result, error) {
	return c.CurrentUserWithResponse(context.Background())
}

func RequestOperations(c *client.ClientWithResponses) (manualtests.Result, error) {
	return c.SearchOperationsWithResponse(context.Background(), nil)
}

func RequestData(state manualtests.StateForCall, c *client.ClientWithResponses) (manualtests.Result, error) {
	id := state.LatestDataID()
	require.NotEmpty(state.T, id, "there should be some data id after this test")
	return c.DataByIdWithResponse(context.Background(), id)
}
