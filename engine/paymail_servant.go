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

type SanitizedPaymail struct {
	alias, domain, adress string
}

func (s *PaymailServant) GetSanitizedPaymail(paymailAdress string) *SanitizedPaymail {
	sanitized := &SanitizedPaymail{}
	sanitized.alias, sanitized.domain, sanitized.adress = paymail.SanitizePaymail(paymailAdress)

	return sanitized
}

func (s *PaymailServant) GetPkiForPaymail(ctx context.Context, sPaymail *SanitizedPaymail) (*paymail.PKIResponse, error) {
	capabilities, err := getCapabilities(ctx, s.cs, s.pc, sPaymail.domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get paymail capability: %w", err)
	}

	if !capabilities.Has(paymail.BRFCPki, paymail.BRFCPkiAlternate) {
		return nil, ErrCapabilitiesPkiUnsupported
	}

	url := capabilities.GetString(paymail.BRFCPki, paymail.BRFCPkiAlternate)
	pki, err := s.pc.GetPKI(url, sPaymail.alias, sPaymail.domain)
	if err != nil {
		return nil, fmt.Errorf("error getting PKI: %w", err)
	}

	return pki, nil
}

func (s *PaymailServant) AddContactRequest(ctx context.Context, receiverPaymail *SanitizedPaymail, contactData *paymail.PikeContactRequestPayload) (*paymail.PikeContactRequestResponse, error) {
	capabilities, err := getCapabilities(ctx, s.cs, s.pc, receiverPaymail.domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get paymail capability: %w", err)
	}

	if !capabilities.Has(paymail.BRFCPike, "") {
		return nil, ErrCapabilitiesPikeUnsupported
	}

	url := capabilities.GetString(paymail.BRFCPike, "")
	response, err := s.pc.AddContactRequest(url, receiverPaymail.alias, receiverPaymail.domain, contactData)
	if err != nil {
		return nil, fmt.Errorf("error during requesting new contact: %w", err)
	}

	return response, nil
}
