package paymail

import (
	"context"
	"errors"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	pmerrors "github.com/bitcoin-sv/spv-wallet/engine/paymail/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
)

const cacheKeyCapabilities = "paymail-capabilities-"
const cacheTTLCapabilities = 2 * time.Minute

// ClientInterface is an interface for the paymail client
type ClientInterface = paymail.ClientInterface

type service struct {
	cache         Cache
	paymailClient paymail.ClientInterface
	log           zerolog.Logger
}

// NewServiceClient creates a new paymail service client
func NewServiceClient(cache Cache, paymailClient ClientInterface, log zerolog.Logger) ServiceClient {
	if paymailClient == nil {
		panic(spverrors.Newf("paymail client is required to create a new paymail service"))
	}

	if cache == nil {
		panic(spverrors.Newf("cache is required to create a new paymail service"))
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
		return capabilities, nil
	}

	response, err := s.loadCapabilities(domain)
	if err != nil {
		return nil, err
	}

	s.putCapabilitiesInCache(ctx, cacheKey, response.CapabilitiesPayload)

	return &response.CapabilitiesPayload, nil
}

// GetP2PDestinations will ask a paymail host on given address for P2P destinations.
func (s *service) GetP2PDestinations(ctx context.Context, address *paymail.SanitisedPaymail, satoshis bsv.Satoshis) (*paymail.PaymentDestinationPayload, error) {
	capabilities, err := s.GetCapabilities(ctx, address.Domain)
	if err != nil {
		return nil, pmerrors.ErrPaymailHostResponseError.Wrap(err)
	}

	p2pDestinationURL := capabilities.GetString(paymail.BRFCP2PPaymentDestination, "")
	if len(p2pDestinationURL) == 0 {
		return nil, pmerrors.ErrPaymailHostNotSupportingP2P
	}

	response, err := s.paymailClient.GetP2PPaymentDestination(
		p2pDestinationURL,
		address.Alias, address.Domain,
		&paymail.PaymentRequest{Satoshis: uint64(satoshis)},
	)
	if err != nil {
		return nil, pmerrors.ErrPaymailHostResponseError.Wrap(err)
	}

	err = s.validatePaymentDestinationResponse(response, satoshis)
	if err != nil {
		return nil, err
	}

	return &response.PaymentDestinationPayload, nil
}

func (s *service) loadCapabilitiesFromCache(ctx context.Context, key string) (*paymail.CapabilitiesPayload, error) {
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
	err := s.cache.SetModel(ctx, key, capabilities, cacheTTLCapabilities)
	if err != nil {
		s.log.Warn().Err(err).Msgf("failed to store capabilities for key %s in cache", key)
	}
}

func (s *service) loadCapabilities(domain string) (*paymail.CapabilitiesResponse, error) {
	// Get SRV record (domain can be different!)
	srv, err := s.paymailClient.GetSRVRecord(
		paymail.DefaultServiceName, paymail.DefaultProtocol, domain,
	)
	if err != nil {
		return nil, err //nolint:wrapcheck // we have handler for paymail errors
	}
	return s.paymailClient.GetCapabilities(srv.Target, int(srv.Port)) //nolint:wrapcheck // we have handler for paymail errors
}

// GetP2P will return the P2P urls and true if they are both found
func (s *service) GetP2P(ctx context.Context, domain string) (success bool, p2pDestinationURL, p2pSubmitTxURL string, format PayloadFormat) {
	capabilities, err := s.GetCapabilities(ctx, domain)
	if err != nil {
		return false, "", "", BasicPaymailPayloadFormat
	}
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

func (s *service) validatePaymentDestinationResponse(response *paymail.PaymentDestinationResponse, satoshis bsv.Satoshis) error {
	var sum bsv.Satoshis
	for _, out := range response.Outputs {
		outSatoshis := bsv.Satoshis(out.Satoshis)
		sum += outSatoshis
	}
	if sum != satoshis {
		return spverrors.Wrapf(pmerrors.ErrPaymailHostInvalidResponse, "paymail host returned outputs not equal to requested satoshis: expected %d, got %d", satoshis, sum)
	}
	return nil
}

func (s *service) Notify(ctx context.Context, address string, p2pMetadata *paymail.P2PMetaData, reference string, tx *trx.Transaction) error {
	sanitizedPm, err := s.GetSanitizedPaymail(address)
	if err != nil {
		return spverrors.Wrapf(err, "failed to sanitize paymail")
	}
	success, _, p2pSubmitTxURL, format := s.GetP2P(ctx, sanitizedPm.Domain)

	if !success {
		return spverrors.Newf("paymail host does not support P2P transactions")
	}

	p2pTransaction := &paymail.P2PTransaction{
		MetaData:  p2pMetadata,
		Reference: reference,
	}

	switch format {
	case BeefPaymailPayloadFormat:
		p2pTransaction.Beef, err = tx.BEEFHex()
		if err != nil {
			return spverrors.Wrapf(err, "failed to convert transaction to BEEF")
		}
	case BasicPaymailPayloadFormat:
		p2pTransaction.Hex = tx.Hex()
	default:
		return spverrors.Newf("%s is unknown format", format)
	}

	res, err := s.paymailClient.SendP2PTransaction(p2pSubmitTxURL, sanitizedPm.Alias, sanitizedPm.Domain, p2pTransaction)
	if err != nil {
		return spverrors.Wrapf(err, "failed to send transaction via paymail")
	}

	s.log.Info().Str("TxID", res.TxID).Msgf("Successfully notified paymail recipient")

	return nil
}
