package config

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/BuxOrg/bux/chainstate"
	"github.com/rs/zerolog"
)

const pastMerkleRootsJSON = `[
	{
		"merkleRoot": "c6867259e635c18f5fbea7b76aba799a43ae43d6daa4095002b2a5ec2cd656fe",
		"blockHeight": 826599
	}
]`

const (
	pleaseCheck            = "Please check Pulse configuration and service status."
	appWillContinue        = "Application will continue to operate but cannot receive transactions until Pulse is online."
	pulseIsOfflineWarning  = "Unable to connect to Pulse service at startup. " + appWillContinue + " " + pleaseCheck
	unexpectedResponse     = "Unexpected response from Pulse service. " + pleaseCheck
	pulseIsNotReadyWarning = "Pulse is responding but is not ready to verify transactions. " + appWillContinue
)

// CheckPulse tries to make a request to the Pulse service to check if it is online and ready to verify transactions.
// AppConfig should be validated before calling this method.
// This method returns nothing, instead it logs either an error or a warning based on the state of the Pulse service.
func (config *AppConfig) CheckPulse(ctx context.Context, logger *zerolog.Logger) {
	if !config.PulseEnabled() {
		// this method works only with Beef/Pulse enabled
		return
	}
	b := config.Paymail.Beef

	logger.Info().Msg("checking Pulse")

	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(timedCtx, "POST", b.PulseHeaderValidationURL, bytes.NewBufferString(pastMerkleRootsJSON))
	if err != nil {
		logger.Error().Err(err).Msg("error checking Pulse - failed to create request")
		return
	}

	if b.PulseAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+b.PulseAuthToken)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		logger.Error().Err(err).Msg(pulseIsOfflineWarning)
		return
	}
	if res.StatusCode != http.StatusOK {
		logger.Error().Msgf("%s Response statusCode: %d", pulseIsOfflineWarning, res.StatusCode)
		return
	}

	var responseModel chainstate.MerkleRootsConfirmationsResponse
	err = json.NewDecoder(res.Body).Decode(&responseModel)
	if err != nil {
		logger.Error().Err(err).Msg(unexpectedResponse)
		return
	}

	if responseModel.ConfirmationState != chainstate.Confirmed {
		logger.Error().Msg(pulseIsNotReadyWarning)
		return
	}

	logger.Info().Msg("Pulse is ready to verify transactions.")
}

// PulseEnabled returns true if the Pulse service is enabled in the AppConfig
func (config *AppConfig) PulseEnabled() bool {
	return config.Paymail != nil && config.Paymail.Beef.enabled()
}
