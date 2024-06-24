package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/mrz1836/go-cachestore"
)

var (
	ErrCapabilitiesPkiUnsupported  = errors.New("server doesn't support PKI")
	ErrCapabilitiesPikeUnsupported = errors.New("server doesn't support PIKE")
)

type PaymailServant struct {
	cs cachestore.ClientInterface
	pc paymail.ClientInterface
}

func (s *PaymailServant) GetSanitizedPaymail(addr string) (*paymail.SanitisedPaymail, error) {
	if err := paymail.ValidatePaymail(addr); err != nil {
		return nil, err
	}

	sanitized := &paymail.SanitisedPaymail{}
	sanitised.Alias, sanitised.Domain, sanitised.Address = paymail.SanitizePaymail(addr)

	return sanitized, nil
}

func (s *PaymailServant) GetPkiForPaymail(ctx context.Context, sPaymail *paymail.SanitisedPaymail) (*paymail.PKIResponse, error) {
	capabilities, err := getCapabilities(ctx, s.cs, s.pc, sPaymail.Domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get paymail capability: %w", err)
	}

	if !capabilities.Has(paymail.BRFCPki, paymail.BRFCPkiAlternate) {
		return nil, ErrCapabilitiesPkiUnsupported
	}

	url := capabilities.GetString(paymail.BRFCPki, paymail.BRFCPkiAlternate)
	pki, err := s.pc.GetPKI(url, sPaymail.Alias, sPaymail.Domain)
	if err != nil {
		return nil, fmt.Errorf("error getting PKI: %w", err)
	}

	return pki, nil
}

func (s *PaymailServant) AddContactRequest(ctx context.Context, receiverPaymail *paymail.SanitisedPaymail, contactData *paymail.PikeContactRequestPayload) (*paymail.PikeContactRequestResponse, error) {
	capabilities, err := getCapabilities(ctx, s.cs, s.pc, receiverPaymail.Domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get paymail capability: %w", err)
	}

	if !capabilities.Has(paymail.BRFCPike, "") {
		return nil, ErrCapabilitiesPikeUnsupported
	}

	url := capabilities.ExtractPikeInviteURL()
	response, err := s.pc.AddContactRequest(url, receiverPaymail.Alias, receiverPaymail.Domain, contactData)
	if err != nil {
		return nil, fmt.Errorf("error during requesting new contact: %w", err)
	}

	return response, nil
}
