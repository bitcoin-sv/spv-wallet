package paymail

import (
	"context"
	"errors"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
)

const cacheKeyCapabilities = "paymail-capabilities-"
const cacheTTLCapabilities = 60 * time.Minute

type service struct {
	cache         cachestore.ClientInterface
	paymailClient paymail.ClientInterface
	log           zerolog.Logger
}

// NewServiceClient creates a new paymail service client
func NewServiceClient(cache cachestore.ClientInterface, paymailClient paymail.ClientInterface, log zerolog.Logger) ServiceClient {
	if paymailClient == nil {
		panic(spverrors.Newf("paymail client is required to create a new paymail service"))
	}

	if cache == nil {
		log.Info().Msg("Doesn't receive cachestore, won't use cache for capabilities")
	}

	return &service{
		cache:         cache,
		paymailClient: paymailClient,
		log:           log,
	}
}

// GetSanitizedPaymail validates and returns the sanitized version of paymail address (alias@domain.tld)
func (s *service) GetSanitizedPaymail(addr string) (*paymail.SanitisedPaymail, error) {
	if err := paymail.ValidatePaymail(addr); err != nil {
		return nil, err //nolint:wrapcheck // we have handler for paymail errors
	}

	sanitized := &paymail.SanitisedPaymail{}
	sanitized.Alias, sanitized.Domain, sanitized.Address = paymail.SanitizePaymail(addr)

	return sanitized, nil
}

// GetCapabilities is a utility function to retrieve capabilities for a Paymail provider
func (s *service) GetCapabilities(ctx context.Context, domain string) (*paymail.CapabilitiesPayload, error) {
	// todo: allow user to configure the time that they want to cache the capabilities (if they want to cache or not)

	cacheKey := cacheKeyCapabilities + domain

	capabilities, err := s.loadCapabilitiesFromCache(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if capabilities != nil {
		return capabilities, err
	}

	response, err := s.loadCapabilities(domain)
	if err != nil {
		return nil, err
	}

	s.putCapabilitiesInCache(ctx, cacheKey, response.CapabilitiesPayload)

	return &response.CapabilitiesPayload, nil
}

func (s *service) loadCapabilitiesFromCache(ctx context.Context, key string) (*paymail.CapabilitiesPayload, error) {
	if s.cache == nil {
		return nil, nil
	}

	capabilities := new(paymail.CapabilitiesPayload)
	err := s.cache.GetModel(ctx, key, capabilities)
	if errors.Is(err, cachestore.ErrKeyNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get capabilities from cachestore")
	}

	if len(capabilities.Capabilities) > 0 {
		return capabilities, nil
	}

	return nil, nil
}

func (s *service) putCapabilitiesInCache(ctx context.Context, key string, capabilities paymail.CapabilitiesPayload) {
	if s.cache == nil || s.cache.Engine().IsEmpty() {
		return
	}

	err := s.cache.SetModel(ctx, key, capabilities, cacheTTLCapabilities)
	if err != nil {
		s.log.Warn().Err(err).Msgf("failed to store capabilities for key %s in cache", key)
	}
}

func (s *service) loadCapabilities(domain string) (response *paymail.CapabilitiesResponse, err error) {
	// Get SRV record (domain can be different!)
	srv, err := s.paymailClient.GetSRVRecord(
		paymail.DefaultServiceName, paymail.DefaultProtocol, domain,
	)
	if err != nil {
		return
	}
	return s.paymailClient.GetCapabilities(srv.Target, int(srv.Port)) //nolint:wrapcheck // we have handler for paymail errors
}

// GetP2P will return the P2P urls and true if they are both found
func (s *service) GetP2P(ctx context.Context, domain string) (success bool, p2pDestinationURL, p2pSubmitTxURL string, format PayloadFormat) {
	capabilities, _ := s.GetCapabilities(ctx, domain)
	return s.extractP2P(capabilities)
}

// StartP2PTransaction will start the P2P transaction, returning the reference ID and outputs
func (s *service) StartP2PTransaction(alias, domain, p2pDestinationURL string, satoshis uint64) (*paymail.PaymentDestinationPayload, error) {
	// Start the P2P transaction request
	response, err := s.paymailClient.GetP2PPaymentDestination(
		p2pDestinationURL,
		alias, domain,
		&paymail.PaymentRequest{Satoshis: satoshis},
	)
	if err != nil {
		return nil, err //nolint:wrapcheck // we have handler for paymail errors
	}

	return &response.PaymentDestinationPayload, nil
}

// GetPkiForPaymail retrieves the PKI for a paymail address
func (s *service) GetPkiForPaymail(ctx context.Context, sPaymail *paymail.SanitisedPaymail) (*paymail.PKIResponse, error) {
	capabilities, err := s.GetCapabilities(ctx, sPaymail.Domain)
	if err != nil {
		return nil, spverrors.ErrGetCapabilities
	}

	if !capabilities.Has(paymail.BRFCPki, paymail.BRFCPkiAlternate) {
		return nil, spverrors.ErrCapabilitiesPkiUnsupported
	}

	url := capabilities.GetString(paymail.BRFCPki, paymail.BRFCPkiAlternate)
	pki, err := s.paymailClient.GetPKI(url, sPaymail.Alias, sPaymail.Domain)
	if err != nil {
		return nil, err //nolint:wrapcheck // we have handler for paymail errors
	}

	return pki, nil
}

// AddContactRequest sends a contact invitation via PIKE capability
func (s *service) AddContactRequest(ctx context.Context, receiverPaymail *paymail.SanitisedPaymail, contactData *paymail.PikeContactRequestPayload) (*paymail.PikeContactRequestResponse, error) {
	capabilities, err := s.GetCapabilities(ctx, receiverPaymail.Domain)
	if err != nil {
		return nil, spverrors.ErrGetCapabilities
	}

	if !capabilities.Has(paymail.BRFCPike, "") {
		return nil, spverrors.ErrCapabilitiesPikeUnsupported
	}

	url := capabilities.ExtractPikeInviteURL()
	response, err := s.paymailClient.AddContactRequest(url, receiverPaymail.Alias, receiverPaymail.Domain, contactData)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to send contact request")
	}

	return response, nil
}

func (s *service) extractP2P(capabilities *paymail.CapabilitiesPayload) (success bool, p2pDestinationURL, p2pSubmitTxURL string, format PayloadFormat) {
	p2pDestinationURL = capabilities.GetString(paymail.BRFCP2PPaymentDestination, "")
	p2pSubmitTxURL = capabilities.GetString(paymail.BRFCP2PTransactions, "")
	p2pBeefSubmitTxURL := capabilities.GetString(paymail.BRFCBeefTransaction, "")

	if len(p2pBeefSubmitTxURL) > 0 {
		p2pSubmitTxURL = p2pBeefSubmitTxURL
		format = BeefPaymailPayloadFormat
	}

	if len(p2pSubmitTxURL) > 0 && len(p2pDestinationURL) > 0 {
		success = true
	}
	return
}
