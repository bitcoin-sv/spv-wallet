package paymailmock

import (
	"encoding/json"
	"net/http"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/jarcoal/httpmock"
)

// MockedRecordBEEFResponse is a mocked response for the record tx endpoint
type MockedRecordBEEFResponse struct{}

// Responder returns a httpmock responder for the mocked P2P destinations response
func (m *MockedRecordBEEFResponse) Responder() httpmock.Responder {
	return func(request *http.Request) (*http.Response, error) {
		var payload struct {
			BEEF string `json:"beef"`
		}
		err := json.NewDecoder(request.Body).Decode(&payload)
		if err != nil {
			return httpmock.NewStringResponse(http.StatusBadRequest, "invalid json"), nil
		}

		tx, err := trx.NewTransactionFromBEEFHex(payload.BEEF)
		if err != nil {
			return httpmock.NewStringResponse(http.StatusBadRequest, "invalid transaction BEEF"), nil
		}

		r, err := httpmock.NewJsonResponse(http.StatusOK, map[string]any{
			"note": "",
			"txid": tx.TxID().String(),
		})
		if err != nil {
			panic(spverrors.Wrapf(err, "cannot create mocked responder for record tx response"))
		}

		return r, nil
	}
}

// RecordBEEFResponse returns a new mocked response for the record tx endpoint
func RecordBEEFResponse() *MockedRecordBEEFResponse {
	return &MockedRecordBEEFResponse{}
}
