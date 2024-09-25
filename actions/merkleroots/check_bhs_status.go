package merkleroots

import (
	"context"
	"net/http"
	"time"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/rs/zerolog"
)

const (
	errCheckingBHSMsg                  = "error checking Block Header Service"
	errBlockHeadersServiceIsOfflineMsg = "Unable to connect to Block Headers Service. Please check Block Header Service configuration and status"
)

func CheckBlockHeaderServiceStatus(ctx context.Context, bhsConfig *config.BHSConfig, httpClient *http.Client, logger *zerolog.Logger) bool {
	logger.Info().Msg("checking Block Headers Service")

	if bhsConfig.URL == "" {
		logger.Error().Msgf("%s - url not configured", errCheckingBHSMsg)
	}

	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(timedCtx, "GET", bhsConfig.URL+"/status", nil)
	if err != nil {
		logger.Error().Err(err).Msgf("%s - failed to create request", errCheckingBHSMsg)
		return false
	}

	if bhsConfig.AuthToken == "" {
		logger.Warn().Msg("warning checking Block Headers Service - auth token is not set. Some requests might not work")
	}

	res, err := httpClient.Do(req)
	if res != nil {
		defer func() {
			_ = res.Body.Close()
		}()
	}
	if err != nil {
		logger.Error().Err(err).Msg(errBlockHeadersServiceIsOfflineMsg)
		return false
	}

	if res.StatusCode != http.StatusOK {
		logger.Error().Msgf("%s Response statusCode: %d", errBlockHeadersServiceIsOfflineMsg, res.StatusCode)
		return false
	}

	return true
}
