package merkleroots

import (
	"net/http"
	"time"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

const (
	errCheckingBHSMsg                  = "error checking Block Header Service"
	errBlockHeadersServiceIsOfflineMsg = "Unable to connect to Block Headers Service. Please check Block Header Service configuration and status"
)

func CheckBlockHeaderServiceStatus(bhsConfig *config.BHSConfig, httpClient *resty.Client, logger *zerolog.Logger) bool {
	logger.Info().Msg("checking Block Headers Service")

	if bhsConfig.URL == "" {
		logger.Error().Msgf("%s - url not configured", errCheckingBHSMsg)
	}

	if bhsConfig.AuthToken == "" {
		logger.Warn().Msg("warning checking Block Headers Service - auth token is not set. Some requests might not work")
	}

	res, err := httpClient.
		SetTimeout(5 * time.Second).
		R().
		EnableTrace().
		Get(bhsConfig.URL + "/status")

	if err != nil {
		logger.Error().Err(err).Msg(errBlockHeadersServiceIsOfflineMsg)
		return false
	}

	if res.StatusCode() != http.StatusOK {
		logger.Error().Msgf("%s Response statusCode: %d", errBlockHeadersServiceIsOfflineMsg, res.StatusCode())
		return false
	}

	return true
}
