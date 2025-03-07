package engine

import (
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/utils/must"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

func setupPaymailClient(overridesToApply *overrides, httpClient *resty.Client) paymail.ClientInterface {
	var paymailClient paymail.ClientInterface
	if overridesToApply.paymailClient != nil {
		paymailClient = overridesToApply.paymailClient
	} else {
		var err error
		paymailClient, err = paymail.NewClient()
		must.HaveNoErrorf(err, "failed to setup paymail client")
		paymailClient.WithCustomHTTPClient(httpClient)
	}
	return paymailClient
}

func setupPaymailServer(cfg *config.AppConfig, logger zerolog.Logger, serviceProvider server.PaymailServiceProvider) *server.Configuration {
	pmCfg := cfg.Paymail

	logger = logger.With().Str("service", "paymail-server").Logger()

	if !pmCfg.Beef.Enabled() {
		logger.Warn().Msg("In V2, BEEF capability cannot be disabled")
	}
	options := []server.ConfigOps{
		server.WithP2PCapabilities(),
		server.WithBeefCapabilities(),
	}

	for _, domain := range pmCfg.Domains {
		options = append(options, server.WithDomain(domain))
	}

	if pmCfg.SenderValidationEnabled {
		options = append(options, server.WithSenderValidation())
	}

	if !pmCfg.DomainValidationEnabled {
		options = append(options, server.WithDomainValidationDisabled())
	}

	if cfg.ExperimentalFeatures.PikeContactsEnabled {
		logger.Warn().Msg("In V2, Pike Payment is not yet supported")
	}

	if cfg.ExperimentalFeatures.PikePaymentEnabled {
		logger.Warn().Msg("In V2, Pike Payment is not yet supported")
	}

	paymailLogger := logger.With().Str("subservice", "go-paymail").Logger()
	options = append(options, server.WithLogger(&paymailLogger))

	paymailLocator := &server.PaymailServiceLocator{}

	paymailLocator.RegisterPaymailService(serviceProvider)

	configuration, err := server.NewConfig(
		paymailLocator,
		options...,
	)
	must.HaveNoErrorf(err, "failed to setup paymail server")

	return configuration
}
