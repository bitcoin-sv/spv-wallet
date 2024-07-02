package engine

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/mrz1836/go-cachestore"
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
	sanitized.Alias, sanitized.Domain, sanitized.Address = paymail.SanitizePaymail(addr)

	return sanitized, nil
}

func (s *PaymailServant) GetPkiForPaymail(ctx context.Context, sPaymail *paymail.SanitisedPaymail) (*paymail.PKIResponse, error) {
	capabilities, err := getCapabilities(ctx, s.cs, s.pc, sPaymail.Domain)
	if err != nil {
		return nil, spverrors.ErrGetCapabilities
	}

	if !capabilities.Has(paymail.BRFCPki, paymail.BRFCPkiAlternate) {
		return nil, spverrors.ErrCapabilitiesPkiUnsupported
	}

	url := capabilities.GetString(paymail.BRFCPki, paymail.BRFCPkiAlternate)
	pki, err := s.pc.GetPKI(url, sPaymail.Alias, sPaymail.Domain)
	if err != nil {
		return nil, err
	}

	return pki, nil
}

func (s *PaymailServant) AddContactRequest(ctx context.Context, receiverPaymail *paymail.SanitisedPaymail, contactData *paymail.PikeContactRequestPayload) (*paymail.PikeContactRequestResponse, error) {
	capabilities, err := getCapabilities(ctx, s.cs, s.pc, receiverPaymail.Domain)
	if err != nil {
		return nil, spverrors.ErrGetCapabilities
	}

	if !capabilities.Has(paymail.BRFCPike, "") {
		return nil, spverrors.ErrCapabilitiesPikeUnsupported
	}

	url := capabilities.ExtractPikeInviteURL()
	response, err := s.pc.AddContactRequest(url, receiverPaymail.Alias, receiverPaymail.Domain, contactData)
	if err != nil {
		return nil, err
	}

	return response, nil
}
