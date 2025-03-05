package contacts

import (
	"context"

	goPaymail "github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/rs/zerolog"
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

// UpsertContact creates or updates a contact
func (s *Service) UpsertContact(ctx context.Context, newContact contactsmodels.NewContact) (*contactsmodels.Contact, error) {
	rAlias, rDomain, rAddress := goPaymail.SanitizePaymail(newContact.RequesterPaymail)
	rPaymail, err := s.paymailAddressService.Find(ctx, rAlias, rDomain)
	if err != nil || rPaymail == nil {
		return nil, spverrors.ErrCouldNotFindPaymail
	}

	if rPaymail.UserID != newContact.UserID {
		return nil, spverrors.ErrUserDoNotOwnPaymail
	}

	contactPaymail, err := s.paymailService.GetSanitizedPaymail(newContact.NewContactPaymail)
	if err != nil {
		return nil, spverrors.ErrContactInvalidPaymail
	}

	contactPKI, err := s.paymailService.GetPkiForPaymail(ctx, contactPaymail)
	if err != nil {
		return nil, spverrors.ErrGettingPKIFailed
	}

	contact, err := s.Find(ctx, newContact.UserID, newContact.NewContactPaymail)
	if err != nil {
		return nil, spverrors.ErrGetContact.Wrap(err)
	}

	newContact.NewContactPubKey = contactPKI.PubKey

	if contact != nil {
		newContact.Status = contact.Status
		c, err := s.contactsRepo.Update(ctx, newContact)
		if err != nil {
			return nil, spverrors.ErrUpdateContact.WithTrace(err)
		}

		return c, nil
	}

	newContact.Status = contactsmodels.ContactNotConfirmed

	c, err := s.contactsRepo.Create(ctx, newContact)
	if err != nil {
		return nil, spverrors.ErrSaveContact.WithTrace(err)
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

		return c, spverrors.ErrAddingContactRequest
	}

	return c, nil
}

// AddContactRequest adds a new contact based on request data
func (s *Service) AddContactRequest(ctx context.Context, fullName, paymail, userID string) (*contactsmodels.Contact, error) {
	contactPaymail, err := s.paymailService.GetSanitizedPaymail(paymail)
	if err != nil {
		return nil, spverrors.ErrRequestedContactInvalid
	}

	contactPki, err := s.paymailService.GetPkiForPaymail(ctx, contactPaymail)
	if err != nil {
		return nil, spverrors.ErrGettingPKIFailed
	}

	contact, err := s.Find(ctx, userID, contactPaymail.Address)
	if err != nil {
		return nil, err
	}

	if contact != nil {
		contact, err = s.updateContactPubKey(ctx, contact, contactPki.PubKey)
		if err != nil {
			return nil, err
		}
		return contact, nil
	}

	contact, err = s.contactsRepo.Create(ctx, contactsmodels.NewContact{
		UserID:            userID,
		FullName:          fullName,
		NewContactPaymail: paymail,
		NewContactPubKey:  contactPki.PubKey,
		Status:            contactsmodels.ContactNotConfirmed,
	})
	if err != nil {
		return nil, spverrors.ErrSaveContact.WithTrace(err)
	}

	return contact, nil
}

// Find returns a paymail by alias and domain
func (s *Service) Find(ctx context.Context, userID, paymail string) (*contactsmodels.Contact, error) {
	contact, err := s.contactsRepo.Find(ctx, userID, paymail)
	if err != nil {
		return nil, spverrors.ErrGetContact.WithTrace(err)
	}

	return contact, nil
}

// PaginatedForUser returns contacts for a user based on userID and the provided paging options and db conditions.
func (s *Service) PaginatedForUser(ctx context.Context, userID string, page filter.Page, conditions map[string]interface{}) (*models.PagedResult[contactsmodels.Contact], error) {
	entities, err := s.contactsRepo.PaginatedForUser(ctx, userID, page, conditions)
	if err != nil {
		return nil, spverrors.ErrContactFailedToGetPaginatedResults.WithTrace(err)
	}

	return entities, nil
}

// PaginatedForAdmin returns all contacts based on the provided paging options and db conditions.
func (s *Service) PaginatedForAdmin(ctx context.Context, page filter.Page, conditions map[string]interface{}) (*models.PagedResult[contactsmodels.Contact], error) {
	entities, err := s.contactsRepo.PaginatedForAdmin(ctx, page, conditions)
	if err != nil {
		return nil, spverrors.ErrContactFailedToGetPaginatedResults.WithTrace(err)
	}

	return entities, nil
}

// UpdateFullNameByID updates the full name of a contact by ID
func (s *Service) UpdateFullNameByID(ctx context.Context, contactID uint, fullName string) (*contactsmodels.Contact, error) {
	c, err := s.contactsRepo.UpdateByID(ctx, contactID, fullName)
	if err != nil {
		return nil, spverrors.ErrUpdateContactStatus.WithTrace(err)
	}

	return c, nil
}

// RemoveContact deletes a contact
func (s *Service) RemoveContact(ctx context.Context, userID, paymail string) error {
	err := s.contactsRepo.Delete(ctx, userID, paymail)
	if err != nil {
		return spverrors.ErrDeleteContact.WithTrace(err)
	}

	return nil
}

// RemoveContactByID deletes a contact by ID
func (s *Service) RemoveContactByID(ctx context.Context, contactID uint) error {
	err := s.contactsRepo.DeleteByID(ctx, contactID)
	if err != nil {
		return spverrors.ErrDeleteContact.WithTrace(err)
	}

	return nil
}

// ConfirmContact confirms a contact
func (s *Service) ConfirmContact(ctx context.Context, userID, paymail string) error {
	contact, err := s.Find(ctx, userID, paymail)
	if err != nil {
		return err
	} else if contact == nil {
		return spverrors.ErrContactNotFound.WithTrace(err)
	}

	if contact.Status != contactsmodels.ContactNotConfirmed {
		s.log.Error().Msgf("cannot confirm contact. Reason: status: %s, expected: %s", contact.Status, contactsmodels.ContactNotConfirmed)
		return spverrors.ErrContactInWrongStatus.WithTrace(err)
	}

	err = s.contactsRepo.UpdateStatus(ctx, userID, paymail, contactsmodels.ContactConfirmed)
	if err != nil {
		return spverrors.ErrUpdateContactStatus.WithTrace(err)
	}

	return nil
}

// UnconfirmContact unconfirms a contact
func (s *Service) UnconfirmContact(ctx context.Context, userID, paymail string) error {
	contact, err := s.Find(ctx, userID, paymail)
	if err != nil {
		return err
	} else if contact == nil {
		return spverrors.ErrContactNotFound.WithTrace(err)
	}

	if contact.Status != contactsmodels.ContactConfirmed {
		s.log.Error().Msgf("cannot unconfirm contact. Reason: status: %s, expected: %s", contact.Status, contactsmodels.ContactConfirmed)
		return spverrors.ErrContactInWrongStatus.WithTrace(err)
	}

	err = s.contactsRepo.UpdateStatus(ctx, userID, paymail, contactsmodels.ContactNotConfirmed)
	if err != nil {
		return spverrors.ErrUpdateContactStatus.WithTrace(err)
	}

	return nil
}

// AcceptContact accept a contact
func (s *Service) AcceptContact(ctx context.Context, userID, paymail string) error {
	contact, err := s.Find(ctx, userID, paymail)
	if err != nil {
		return err
	} else if contact == nil {
		return spverrors.ErrContactNotFound.WithTrace(err)
	}

	if contact.Status != contactsmodels.ContactAwaitAccept {
		s.log.Error().Msgf("cannot accept contact. Reason: status: %s, expected: %s", contact.Status, contactsmodels.ContactAwaitAccept)
		return spverrors.ErrContactInWrongStatus.WithTrace(err)
	}

	err = s.contactsRepo.UpdateStatus(ctx, userID, paymail, contactsmodels.ContactNotConfirmed)
	if err != nil {
		return spverrors.ErrUpdateContactStatus.WithTrace(err)
	}

	return nil
}

// AcceptContactByID accept a contact by ID
func (s *Service) AcceptContactByID(ctx context.Context, id uint) (*contactsmodels.Contact, error) {
	contact, err := s.contactsRepo.UpdateStatusByID(ctx, id, contactsmodels.ContactNotConfirmed)
	if err != nil {
		return nil, spverrors.ErrUpdateContactStatus.WithTrace(err)
	}

	return contact, nil
}

// RejectContact reject a contact
func (s *Service) RejectContact(ctx context.Context, userID, paymail string) error {
	contact, err := s.Find(ctx, userID, paymail)
	if err != nil {
		return err
	} else if contact == nil {
		return spverrors.ErrContactNotFound.WithTrace(err)
	}

	if contact.Status != contactsmodels.ContactAwaitAccept {
		s.log.Error().Msgf("cannot reject contact. Reason: status: %s, expected: %s", contact.Status, contactsmodels.ContactAwaitAccept)
		return spverrors.ErrContactInWrongStatus.WithTrace(err)
	}

	err = s.contactsRepo.Delete(ctx, userID, paymail)
	if err != nil {
		return spverrors.ErrUpdateContactStatus.WithTrace(err)
	}

	return nil
}

// RejectContactByID reject a contact by ID
func (s *Service) RejectContactByID(ctx context.Context, id uint) (*contactsmodels.Contact, error) {
	contact, err := s.contactsRepo.UpdateStatusByID(ctx, id, contactsmodels.ContactRejected)
	if err != nil {
		return nil, spverrors.ErrUpdateContactStatus.WithTrace(err)
	}

	return contact, nil
}

// AdminCreateContact creates a new contact for the provided paymail
func (s *Service) AdminCreateContact(ctx context.Context, newContact contactsmodels.NewContact) (*contactsmodels.Contact, error) {
	err := validateNewContact(newContact)
	if err != nil {
		return nil, err
	}

	rAlias, rDomain, _ := goPaymail.SanitizePaymail(newContact.RequesterPaymail)
	rPaymailAddr, err := s.paymailAddressService.Find(ctx, rAlias, rDomain)
	if err != nil {
		return nil, spverrors.ErrCouldNotFindPaymail.WithTrace(err)
	} else if rPaymailAddr == nil {
		return nil, spverrors.ErrCouldNotFindPaymail
	}

	contact, err := s.Find(ctx, rPaymailAddr.UserID, newContact.NewContactPaymail)
	if err != nil {
		return nil, spverrors.ErrGetContact.WithTrace(err)
	} else if contact != nil {
		return nil, spverrors.ErrContactAlreadyExists
	}

	contactPaymail, err := s.paymailService.GetSanitizedPaymail(newContact.NewContactPaymail)
	if err != nil {
		return nil, spverrors.ErrRequestedContactInvalid
	}

	contactPKI, err := s.paymailService.GetPkiForPaymail(ctx, contactPaymail)
	if err != nil {
		return nil, spverrors.ErrGettingPKIFailed
	}

	newContact.NewContactPubKey = contactPKI.PubKey

	contact, err = s.contactsRepo.Create(ctx, newContact)
	if err != nil {
		return nil, spverrors.ErrSaveContact.Wrap(err)
	}

	return contact, nil
}

// AdminConfirmContacts confirms provided contacts.
func (s *Service) AdminConfirmContacts(ctx context.Context, paymailA, paymailB string) error {
	contactA, contactB, err := s.retrieveContacts(ctx, paymailA, paymailB)
	if err != nil {
		return spverrors.ErrGetContact.WithTrace(err)
	}

	err = s.contactsRepo.UpdateStatus(ctx, contactA.UserID, contactA.Paymail, contactsmodels.ContactConfirmed)
	if err != nil {
		return spverrors.ErrUpdateContactStatus.WithTrace(err)
	}

	err = s.contactsRepo.UpdateStatus(ctx, contactB.UserID, contactB.Paymail, contactsmodels.ContactConfirmed)
	if err != nil {
		return spverrors.ErrUpdateContactStatus.WithTrace(err)
	}

	return nil
}

func (s *Service) updateContactPubKey(ctx context.Context, contact *contactsmodels.Contact, pubKey string) (*contactsmodels.Contact, error) {
	contactToUpdate := contactsmodels.NewContact{}
	if contact.PubKey != pubKey {
		contactToUpdate.NewContactPubKey = pubKey

		if contact.Status == contactsmodels.ContactConfirmed {
			contactToUpdate.Status = contactsmodels.ContactNotConfirmed
		}

		c, err := s.contactsRepo.Update(ctx, contactToUpdate)
		if err != nil {
			return nil, spverrors.ErrUpdateContact.WithTrace(err)
		}
		return c, nil
	}

	return contact, nil
}

func (s *Service) retrieveContacts(ctx context.Context, paymailA, paymailB string) (*contactsmodels.Contact, *contactsmodels.Contact, error) {
	aAlias, aDomain, _ := goPaymail.SanitizePaymail(paymailA)
	aPaymail, err := s.paymailAddressService.Find(ctx, aAlias, aDomain)
	if err != nil || aPaymail == nil {
		return nil, nil, spverrors.ErrCouldNotFindPaymail
	}

	bAlias, bDomain, _ := goPaymail.SanitizePaymail(paymailB)
	bPaymail, err := s.paymailAddressService.Find(ctx, bAlias, bDomain)
	if err != nil || bPaymail == nil {
		return nil, nil, spverrors.ErrCouldNotFindPaymail
	}

	contactA, err := s.Find(ctx, aPaymail.UserID, paymailB)
	if err != nil {
		return nil, nil, err
	} else if contactA == nil {
		return nil, nil, spverrors.ErrContactNotFound
	}

	contactB, err := s.Find(ctx, bPaymail.UserID, paymailA)
	if err != nil {
		return nil, nil, err
	} else if contactB == nil {
		return nil, nil, spverrors.ErrContactNotFound
	}

	return contactA, contactB, nil
}

func validateNewContact(newContact contactsmodels.NewContact) error {
	if newContact.FullName == "" {
		return spverrors.ErrMissingContactFullName
	}

	if newContact.NewContactPaymail == "" {
		return spverrors.ErrMissingContactPaymailParam
	}

	if newContact.RequesterPaymail == "" {
		return spverrors.ErrMissingContactCreatorPaymail
	}

	return nil
}
