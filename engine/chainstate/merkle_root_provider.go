package chainstate

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/rs/zerolog"
)

// MerkleRootConfirmationState represents the state of each Merkle Root verification
// process and can be one of three values: Confirmed, Invalid and UnableToVerify.
type MerkleRootConfirmationState string

const (
	// Confirmed state occurs when Merkle Root is found in the longest chain.
	Confirmed MerkleRootConfirmationState = "CONFIRMED"
	// Invalid state occurs when Merkle Root is not found in the longest chain.
	Invalid MerkleRootConfirmationState = "INVALID"
	// UnableToVerify state occurs when Block Header Service is behind in synchronization with the longest chain.
	UnableToVerify MerkleRootConfirmationState = "UNABLE_TO_VERIFY"
)

// MerkleRootConfirmationRequestItem is a request type for verification
// of Merkle Roots inclusion in the longest chain.
type MerkleRootConfirmationRequestItem struct {
	MerkleRoot  string `json:"merkleRoot"`
	BlockHeight uint64 `json:"blockHeight"`
}

// MerkleRootConfirmation is a confirmation
// of merkle roots inclusion in the longest chain.
type MerkleRootConfirmation struct {
	Hash         string                      `json:"blockHash"`
	BlockHeight  uint64                      `json:"blockHeight"`
	MerkleRoot   string                      `json:"merkleRoot"`
	Confirmation MerkleRootConfirmationState `json:"confirmation"`
}

// MerkleRootsConfirmationsResponse is an API response for confirming
// merkle roots inclusion in the longest chain.
type MerkleRootsConfirmationsResponse struct {
	ConfirmationState MerkleRootConfirmationState `json:"confirmationState"`
	Confirmations     []MerkleRootConfirmation    `json:"confirmations"`
}

type blockHeadersServiceClientProvider struct {
	url        string
	authToken  string
	httpClient *http.Client
}

func newBlockHeaderServiceClientProvider(url, authToken string) *blockHeadersServiceClientProvider {
	return &blockHeadersServiceClientProvider{url: url, authToken: authToken, httpClient: &http.Client{}}
}

func (p *blockHeadersServiceClientProvider) verifyMerkleRoots(
	ctx context.Context,
	logger *zerolog.Logger,
	merkleRoots []MerkleRootConfirmationRequestItem,
) (*MerkleRootsConfirmationsResponse, error) {
	jsonData, err := json.Marshal(merkleRoots)
	if err != nil {
		return nil, _fmtAndLogError(err, logger, "Error occurred while marshaling merkle roots.")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, _fmtAndLogError(err, logger, "Error occurred while creating request for the Block Headers Service client.")
	}

	if p.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+p.authToken)
	}
	res, err := p.httpClient.Do(req)
	if res != nil {
		defer func() {
			_ = res.Body.Close()
		}()
	}
	if err != nil {
		return nil, _fmtAndLogError(err, logger, "Error occurred while sending request to the Block Headers Service.")
	}

	if res.StatusCode != 200 {
		return nil, _fmtAndLogError(_statusError(res.StatusCode), logger, "Received unexpected status code from Block Headers Service.")
	}

	// Parse response body.
	var merkleRootsRes MerkleRootsConfirmationsResponse
	err = json.NewDecoder(res.Body).Decode(&merkleRootsRes)
	if err != nil {
		return nil, _fmtAndLogError(err, logger, "Error occurred while parsing response from the Block Headers Service.")
	}

	return &merkleRootsRes, nil
}

// _fmtAndLogError returns brief error for http response message and logs detailed information with original error
func _fmtAndLogError(err error, logger *zerolog.Logger, message string) error {
	logger.Error().Err(err).Msg("[verifyMerkleRoots] " + message)
	return spverrors.Newf("cannot verify transaction - %s", message)
}

func _statusError(statusCode int) error {
	return spverrors.Newf("Block Headers Service client returned status code %d - check Block Headers Service configuration and status", statusCode)
}
