package merkleroots

import (
	"errors"
	"net/url"
	"syscall"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// getMerkleRootsFromBHS returns Merkle Roots from Block Header Service
func getMerkleRootsFromBHS(client *resty.Client, appConfig *config.AppConfig, logger *zerolog.Logger, query url.Values) (*any, error) {
	bhsURL, err := createBHSURL(appConfig, "/chain/merkleroot", logger)
	if err != nil {
		return nil, err
	}

	req := client.R().EnableTrace()

	if appConfig.BHS.AuthToken != "" {
		req.SetAuthToken(appConfig.BHS.AuthToken)
	}

	var response any
	res, err := req.
		SetResult(&response).
		SetQueryParamsFromValues(query).
		Get(bhsURL.String())
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) {
			return nil, ErrBHSUnreachable.Wrap(err)
		}
		return nil, spverrors.ErrInternal.Wrap(err)
	}
	if !res.IsSuccess() {
		return nil, mapBHSErrorResponseToSpverror(res)
	}

	return &response, nil
}

// createBHSURL parses Block Header Url from configuration and constructs a valid
// endpoint with provided endpointPath variable
func createBHSURL(appConfig *config.AppConfig, endpointPath string, logger *zerolog.Logger) (*url.URL, error) {
	if appConfig.BHS.URL == "" {
		logger.Error().Msgf("create Block Header Service URL - url not configured")
	}
	if appConfig.BHS.AuthToken == "" {
		logger.Warn().Msg("warning creating Block Headers Service url - auth token is not set. Some requests might not work")
	}

	url, err := url.Parse(appConfig.BHS.URL + "/api/" + appConfig.BHS.APIVersion + endpointPath)
	if err != nil {
		return nil, ErrBHSBadURL.Wrap(err)
	}

	return url, nil
}
