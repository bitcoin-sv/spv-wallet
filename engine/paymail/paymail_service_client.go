package paymail

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/mrz1836/go-cachestore"
)

const cacheKeyCapabilities = "paymail-capabilities-"
const cacheTTLCapabilities = 60 * time.Minute

type service struct {
	cache         cachestore.ClientInterface
	paymailClient paymail.ClientInterface
}

// NewServiceClient creates a new paymail service client
func NewServiceClient(cache cachestore.ClientInterface, paymailClient paymail.ClientInterface) ServiceClient {
	if paymailClient == nil {
		panic(spverrors.Newf("paymail client is required to create a new paymail service"))
	}
	return &service{
		cache:         cache,
		paymailClient: paymailClient,
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
	// Attempt to get from cachestore
	// todo: allow user to configure the time that they want to cache the capabilities (if they want to cache or not)
	capabilities := new(paymail.CapabilitiesPayload)
	if s.cache != nil {
		if err := s.cache.GetModel(
			ctx, cacheKeyCapabilities+domain, capabilities,
		); err != nil && !errors.Is(err, cachestore.ErrKeyNotFound) {
			return nil, spverrors.Wrapf(err, "failed to get capabilities from cachestore")
		} else if len(capabilities.Capabilities) > 0 {
			return capabilities, nil
		}
	}

	// Get SRV record (domain can be different!)
	var response *paymail.CapabilitiesResponse
	srv, err := s.paymailClient.GetSRVRecord(
		paymail.DefaultServiceName, paymail.DefaultProtocol, domain,
	)
	if err != nil {
		// Error returned was a real error
		if !strings.Contains(err.Error(), "zero SRV records found") { // This error is from no SRV record being found
			return nil, err //nolint:wrapcheck // we have handler for paymail errors
		}

		// Try to get capabilities without the SRV record
		// 'Should no record be returned, a paymail client should assume a host of <domain>.<tld> and a port of 443.'
		// http://bsvalias.org/02-01-host-discovery.html

		// Get the capabilities via target
		if response, err = s.paymailClient.GetCapabilities(
			domain, paymail.DefaultPort,
		); err != nil {
			return nil, err //nolint:wrapcheck // we have handler for paymail errors
		}
	} else {
		// Get the capabilities via SRV record
		if response, err = s.paymailClient.GetCapabilities(
			srv.Target, int(srv.Port),
		); err != nil {
			return nil, err //nolint:wrapcheck // we have handler for paymail errors
		}
	}

	// Save to cachestore
	if s.cache != nil && !s.cache.Engine().IsEmpty() {
		_ = s.cache.SetModel(
			context.Background(), cacheKeyCapabilities+domain,
			&response.CapabilitiesPayload, cacheTTLCapabilities,
		)
	}

	return &response.CapabilitiesPayload, nil
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
