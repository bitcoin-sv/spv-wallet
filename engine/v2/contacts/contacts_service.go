package contacts

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails"
	"github.com/rs/zerolog"

	goPaymail "github.com/bitcoin-sv/go-paymail"
)

// Service for contacts
type Service struct {
	contactsRepo          ContactRepo
	paymailService        paymail.ServiceClient
	paymailAddressService *paymails.Service
	log                   zerolog.Logger
}

// NewService creates a new contacts service
func NewService(repo ContactRepo, paymailAddressService *paymails.Service, paymailService paymail.ServiceClient, log zerolog.Logger) *Service {
	return &Service{
		contactsRepo:          repo,
		paymailService:        paymailService,
		paymailAddressService: paymailAddressService,
		log:                   log,
	}
}

func (s *Service) UpsertContact(ctx context.Context, newContact *contactsmodels.NewContact) error {
	rAlias, rDomain, rAddress := goPaymail.SanitizePaymail(newContact.RequesterPaymail)
	rPaymail, err := s.paymailAddressService.Find(ctx, rAlias, rDomain)
	if err != nil {
		return err
	}

	if rPaymail.UserID != newContact.UserID {
		return spverrors.ErrUserDoNotOwnPaymail
	}

	contactPaymail, err := s.paymailService.GetSanitizedPaymail(newContact.NewContactPaymail)
	if err != nil {
		return spverrors.ErrContactInvalidPaymail
	}

	contact, err := s.Find(ctx, newContact.UserID, newContact.NewContactPaymail)
	if err != nil {
		return err
	}

	contactPKI, err := s.paymailService.GetPkiForPaymail(ctx, contactPaymail)
	if err != nil {
		return spverrors.ErrGettingPKIFailed
	}

	newContact.NewContactPubKey = contactPKI.PubKey

	if contact != nil {
		return s.contactsRepo.Update(ctx, newContact)
	}

	err = s.contactsRepo.Create(ctx, newContact)
	if err != nil {
		return err
	}

	requesterContactRequest := goPaymail.PikeContactRequestPayload{
		FullName: rPaymail.PublicName,
		Paymail:  rAddress,
	}
	if _, err = s.paymailService.AddContactRequest(ctx, contactPaymail, &requesterContactRequest); err != nil {
		s.log.Warn().
			Str("requesterPaymail", rAddress).
			Str("requestedContact", newContact.NewContactPaymail).
			Msgf("adding contact request failed: %s", err.Error())

		return spverrors.ErrAddingContactRequest
	}

	return nil
}

// Find returns a paymail by alias and domain
func (s *Service) Find(ctx context.Context, userID, paymail string) (*contactsmodels.Contact, error) {
	contact, err := s.contactsRepo.Find(ctx, userID, paymail)
	if err != nil {
		return nil, err
	}

	return contact, nil
}
