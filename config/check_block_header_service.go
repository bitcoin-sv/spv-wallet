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
	pleaseCheck            = "Please check Block Headers Service configuration and service status."
	appWillContinue        = "Application will continue to operate but cannot receive transactions until Block Headers Service is online."
	blockHeaderServiceIsOfflineWarning  = "Unable to connect to Block Headers Service service at startup. " + appWillContinue + " " + pleaseCheck
	unexpectedResponse     = "Unexpected response from Block Headers Service service. " + pleaseCheck
	blockHeaderServiceIsNotReadyWarning = "Block Headers Service is responding but is not ready to verify transactions. " + appWillContinue
)

// CheckBlockHeaderService tries to make a request to the Block Headers Service to check if it is online and ready to verify transactions.
// AppConfig should be validated before calling this method.
// This method returns nothing, instead it logs either an error or a warning based on the state of the Block Headers Service.
func (config *AppConfig) CheckBlockHeaderService(ctx context.Context, logger *zerolog.Logger) {
	if !config.BlockHeaderServiceEnabled() {
		// this method works only with Beef/Block Headers Service enabled
		return
	}
	b := config.Paymail.Beef

	logger.Info().Msg("checking Block Headers Service")

	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(timedCtx, "POST", b.BlockHeaderServiceHeaderValidationURL, bytes.NewBufferString(pastMerkleRootsJSON))
	if err != nil {
		logger.Error().Err(err).Msg("error checking Block Headers Service - failed to create request")
		return
	}

	if b.BlockHeaderServiceAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+b.BlockHeaderServiceAuthToken)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if res != nil {
		defer func() {
			_ = res.Body.Close()
		}()
	}

	if err != nil {
		logger.Error().Err(err).Msg(blockHeaderServiceIsOfflineWarning)
		return
	}
	if res.StatusCode != http.StatusOK {
		logger.Error().Msgf("%s Response statusCode: %d", blockHeaderServiceIsOfflineWarning, res.StatusCode)
		return
	}

	var responseModel chainstate.MerkleRootsConfirmationsResponse
	err = json.NewDecoder(res.Body).Decode(&responseModel)
	if err != nil {
		logger.Error().Err(err).Msg(unexpectedResponse)
		return
	}

	if responseModel.ConfirmationState != chainstate.Confirmed {
		logger.Error().Msg(blockHeaderServiceIsNotReadyWarning)
		return
	}

	logger.Info().Msg("Block Headers Service is ready to verify transactions.")
}

// BlockHeaderServiceEnabled returns true if the Block Headers Service is enabled in the AppConfig
func (config *AppConfig) BlockHeaderServiceEnabled() bool {
	return config.Paymail != nil && config.Paymail.Beef.enabled()
}
