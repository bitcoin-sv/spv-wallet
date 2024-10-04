package config

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/rs/zerolog"
)

const pastMerkleRootsJSON = `[
	{
		"merkleRoot": "c6867259e635c18f5fbea7b76aba799a43ae43d6daa4095002b2a5ec2cd656fe",
		"blockHeight": 826599
	}
]`

const (
	pleaseCheck                          = "Please check Block Headers Service configuration and service status."
	appWillContinue                      = "Application will continue to operate but cannot receive transactions until Block Headers Service is online."
	blockHeadersServiceIsOfflineWarning  = "Unable to connect to Block Headers Service service at startup. " + appWillContinue + " " + pleaseCheck
	unexpectedResponse                   = "Unexpected response from Block Headers Service service. " + pleaseCheck
	blockHeadersServiceIsNotReadyWarning = "Block Headers Service is responding but is not ready to verify transactions. " + appWillContinue
)

// CheckBlockHeadersService tries to make a request to the Block Headers Service to check if it is online and ready to verify transactions.
// AppConfig should be validated before calling this method.
// This method returns nothing, instead it logs either an error or a warning based on the state of the Block Headers Service.
func (c *AppConfig) CheckBlockHeadersService(ctx context.Context, logger *zerolog.Logger) {
	if !c.BlockHeadersServiceEnabled() {
		// this method works only with Beef/Block Headers Service enabled
		return
	}
	b := c.Paymail.Beef

	logger.Info().Msg("checking Block Headers Service")

	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(timedCtx, "POST", b.BlockHeadersServiceHeaderValidationURL, bytes.NewBufferString(pastMerkleRootsJSON))
	if err != nil {
		logger.Error().Err(err).Msg("error checking Block Headers Service - failed to create request")
		return
	}

	if b.BlockHeadersServiceAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+b.BlockHeadersServiceAuthToken)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if res != nil {
		defer func() {
			_ = res.Body.Close()
		}()
	}

	if err != nil {
		logger.Error().Err(err).Msg(blockHeadersServiceIsOfflineWarning)
		return
	}
	if res.StatusCode != http.StatusOK {
		logger.Error().Msgf("%s Response statusCode: %d", blockHeadersServiceIsOfflineWarning, res.StatusCode)
		return
	}

	var responseModel chainstate.MerkleRootsConfirmationsResponse
	err = json.NewDecoder(res.Body).Decode(&responseModel)
	if err != nil {
		logger.Error().Err(err).Msg(unexpectedResponse)
		return
	}

	if responseModel.ConfirmationState != chainstate.Confirmed {
		logger.Error().Msg(blockHeadersServiceIsNotReadyWarning)
		return
	}

	logger.Info().Msg("Block Headers Service is ready to verify transactions.")
}

// BlockHeadersServiceEnabled returns true if the Block Headers Service is enabled in the AppConfig
func (c *AppConfig) BlockHeadersServiceEnabled() bool {
	return c.Paymail != nil && c.Paymail.Beef.enabled()
}
