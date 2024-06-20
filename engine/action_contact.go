package engine

import (
	"context"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/spverrors"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
)

func (c *Client) UpsertContact(ctx context.Context, ctcFName, ctcPaymail, requesterXPubID, requesterPaymail string, opts ...ModelOps) (*Contact, error) {
	reqPm, err := c.getPaymail(ctx, requesterXPubID, requesterPaymail)
	if err != nil {
		return nil, err
	}

	pmSrvnt := &PaymailServant{
		cs: c.Cachestore(),
		pc: c.PaymailClient(),
	}
	contactPm, err := pmSrvnt.GetSanitizedPaymail(ctcPaymail)
	if err != nil {
		return nil, fmt.Errorf("requested contact paymail is invalid. Reason: %w", err)
	}

	contact, err := c.upsertContact(ctx, pmSrvnt, requesterXPubID, ctcFName, contactPm, opts...)
	if err != nil {
		return nil, err
	}

	// request new contact
	requesterContactRequest := paymail.PikeContactRequestPayload{
		FullName: reqPm.PublicName,
		Paymail:  reqPm.String(),
	}
	if _, err = pmSrvnt.AddContactRequest(ctx, contactPm, &requesterContactRequest); err != nil {
		c.Logger().Warn().
			Str("requesterPaymail", reqPm.String()).
			Str("requestedContact", ctcPaymail).
			Msgf("adding contact request failed: %s", err.Error())

		return contact, spverrors.ErrAddingContactRequest
	}

	return contact, nil
}

func (c *Client) AddContactRequest(ctx context.Context, fullName, paymailAdress, requesterXPubID string, opts ...ModelOps) (*Contact, error) {
	pmSrvnt := &PaymailServant{
		cs: c.Cachestore(),
		pc: c.PaymailClient(),
	}

	contactPm, err := pmSrvnt.GetSanitizedPaymail(paymailAdress)
	if err != nil {
		return nil, fmt.Errorf("requested contact paymail is invalid. Reason: %w", err)
	}

	contactPki, err := pmSrvnt.GetPkiForPaymail(ctx, contactPm)
	if err != nil {
		return nil, fmt.Errorf("geting PKI for %s failed. Reason: %w", paymailAdress, err)
	}

	// check if exists already
	contact, err := getContact(ctx, contactPm.Address, requesterXPubID, c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	save := false
	if contact != nil {
		save = contact.UpdatePubKey(contactPki.PubKey)
	} else {
		contact = newContact(
			fullName,
			contactPm.Address,
			contactPki.PubKey,
			requesterXPubID,
			ContactAwaitAccept,
			c.DefaultModelOptions(append(opts, New())...)...,
		)

		save = true
	}

	if save {
		if err = contact.Save(ctx); err != nil {
			return nil, fmt.Errorf("adding %s contact failed. Reason: %w", paymailAdress, err)
		}
	}

	return contact, nil
}

func (c *Client) GetContacts(ctx context.Context, metadata *Metadata, conditions map[string]interface{}, queryParams *datastore.QueryParams) ([]*Contact, error) {
	contacts, err := getContacts(ctx, metadata, conditions, queryParams, c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (c *Client) GetContactsByXpubID(ctx context.Context, xPubID string, metadata *Metadata, conditions map[string]interface{}, queryParams *datastore.QueryParams) ([]*Contact, error) {
	contacts, err := getContactsByXpubID(ctx, xPubID, metadata, conditions, queryParams, c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (c *Client) GetContactsByXPubIDCount(ctx context.Context, xPubID string, metadata *Metadata, conditions map[string]interface{}, opts ...ModelOps) (int64, error) {
	count, err := getContactsByXPubIDCount(
		ctx,
		xPubID,
		metadata,
		conditions,
		c.DefaultModelOptions(opts...)...,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Client) GetContactsCount(ctx context.Context, metadata *Metadata, conditions map[string]interface{}, opts ...ModelOps) (int64, error) {
	count, err := getModelCountByConditions(ctx, ModelContact, Contact{}, metadata, conditions, c.DefaultModelOptions(opts...)...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Client) UpdateContact(ctx context.Context, id, fullName string, metadata *Metadata) (*Contact, error) {
	contact, err := getContactByID(ctx, id, c.DefaultModelOptions()...)
	if err != nil {
		c.logContactError("", "", fmt.Sprintf("error while geting contact: %s", err.Error()))
		return nil, err
	}

	if contact == nil {
		return nil, spverrors.ErrContactNotFound
	}

	contact.FullName = fullName
	contact.UpdateMetadata(*metadata)

	if err = contact.Save(ctx); err != nil {
		c.logContactError(contact.OwnerXpubID, contact.Paymail, fmt.Sprintf("unexpected error while saving contact: %s", err.Error()))
		return nil, err
	}

	return contact, nil
}

func (c *Client) AdminChangeContactStatus(ctx context.Context, id string, status ContactStatus) (*Contact, error) {
	contact, err := getContactByID(ctx, id, c.DefaultModelOptions()...)
	if err != nil {
		c.logContactError("", "", fmt.Sprintf("error while geting contact: %s", err.Error()))
		return nil, err
	}

	if contact == nil {
		return nil, spverrors.ErrContactNotFound
	}

	switch status {
	case ContactNotConfirmed:
		err = contact.Accept()
	case ContactRejected:
		err = contact.Reject()
	case ContactConfirmed:
		err = contact.Confirm()
	}

	if err != nil {
		c.logContactError(contact.OwnerXpubID, contact.Paymail, fmt.Sprintf("error while changing contact status: %s", err.Error()))
		return nil, err
	}

	if err = contact.Save(ctx); err != nil {
		c.logContactError(contact.OwnerXpubID, contact.Paymail, fmt.Sprintf("unexpected error while saving contact: %s", err.Error()))
		return nil, err
	}
	return contact, nil
}

func (c *Client) DeleteContact(ctx context.Context, id string) error {
	contact, err := getContactByID(ctx, id, c.DefaultModelOptions()...)
	if err != nil {
		c.logContactError("", "", fmt.Sprintf("error while geting contact: %s", err.Error()))
		return err
	}

	if contact == nil {
		return spverrors.ErrContactNotFound
	}

	contact.Delete()

	if err = contact.Save(ctx); err != nil {
		c.logContactError(contact.OwnerXpubID, contact.Paymail, fmt.Sprintf("unexpected error while saving contact: %s", err.Error()))
		return err
	}

	return nil
}

func (c *Client) AcceptContact(ctx context.Context, xPubID, paymail string) error {
	contact, err := getContact(ctx, paymail, xPubID, c.DefaultModelOptions()...)
	if err != nil {
		c.logContactError(xPubID, paymail, fmt.Sprintf("unexpected error while geting contact: %s", err.Error()))
		return err
	}
	if contact == nil {
		return spverrors.ErrContactNotFound
	}

	if err = contact.Accept(); err != nil {
		c.logContactWarining(xPubID, paymail, err.Error())
		return spverrors.ErrContactIncorrectStatus
	}

	if err = contact.Save(ctx); err != nil {
		c.logContactError(xPubID, paymail, fmt.Sprintf("unexpected error while saving contact: %s", err.Error()))
		return err
	}

	return nil
}

func (c *Client) RejectContact(ctx context.Context, xPubID, paymail string) error {
	contact, err := getContact(ctx, paymail, xPubID, c.DefaultModelOptions()...)
	if err != nil {
		c.logContactError(xPubID, paymail, fmt.Sprintf("unexpected error while geting contact: %s", err.Error()))
		return err
	}
	if contact == nil {
		return spverrors.ErrContactNotFound
	}

	if err = contact.Reject(); err != nil {
		c.logContactWarining(xPubID, paymail, err.Error())
		return spverrors.ErrContactIncorrectStatus
	}

	if err = contact.Save(ctx); err != nil {
		c.logContactError(xPubID, paymail, fmt.Sprintf("unexpected error while saving contact: %s", err.Error()))
		return err
	}

	return nil
}

func (c *Client) ConfirmContact(ctx context.Context, xPubID, paymail string) error {
	contact, err := getContact(ctx, paymail, xPubID, c.DefaultModelOptions()...)
	if err != nil {
		c.logContactError(xPubID, paymail, fmt.Sprintf("unexpected error while geting contact: %s", err.Error()))
		return err
	}
	if contact == nil {
		return spverrors.ErrContactNotFound
	}

	if err = contact.Confirm(); err != nil {
		c.logContactWarining(xPubID, paymail, err.Error())
		return spverrors.ErrContactIncorrectStatus
	}

	if err = contact.Save(ctx); err != nil {
		c.logContactError(xPubID, paymail, fmt.Sprintf("unexpected error while saving contact: %s", err.Error()))
		return err
	}

	return nil
}

func (c *Client) UnconfirmContact(ctx context.Context, xPubID, paymail string) error {
	contact, err := getContact(ctx, paymail, xPubID, c.DefaultModelOptions()...)
	if err != nil {
		c.logContactError(xPubID, paymail, fmt.Sprintf("unexpected error while geting contact: %s", err.Error()))
		return err
	}
	if contact == nil {
		return spverrors.ErrContactNotFound
	}

	if err = contact.Unconfirm(); err != nil {
		c.logContactWarining(xPubID, paymail, err.Error())
		return spverrors.ErrContactIncorrectStatus
	}

	if err = contact.Save(ctx); err != nil {
		c.logContactError(xPubID, paymail, fmt.Sprintf("unexpected error while saving contact: %s", err.Error()))
		return err
	}

	return nil
}

func (c *Client) getPaymail(ctx context.Context, xpubID, paymailAddr string) (*PaymailAddress, error) {
	if paymailAddr != "" {
		res, err := c.GetPaymailAddress(ctx, paymailAddr, c.DefaultModelOptions()...)
		if err != nil {
			return nil, err
		}

		if res == nil || res.XpubID != xpubID {
			return nil, spverrors.ErrInvalidRequesterXpub
		}

		return res, nil
	}

	emptyConditions := make(map[string]interface{})

	paymails, err := c.GetPaymailAddressesByXPubID(ctx, xpubID, nil, emptyConditions, nil)
	if err != nil {
		return nil, err
	}
	if len(paymails) == 0 {
		return nil, spverrors.ErrInvalidRequesterXpub
	} else if len(paymails) > 1 {
		return nil, spverrors.ErrMoreThanOnePaymailRegistered
	}

	return paymails[0], nil
}

func (c *Client) upsertContact(ctx context.Context, pmSrvnt *PaymailServant, reqXPubID, ctcFName string, ctcPaymail *paymail.SanitisedPaymail, opts ...ModelOps) (*Contact, error) {
	contactPki, err := pmSrvnt.GetPkiForPaymail(ctx, ctcPaymail)
	if err != nil {
		return nil, fmt.Errorf("geting PKI for %s failed. Reason: %w", ctcPaymail.Address, err)
	}

	// check if exists already
	contact, err := getContact(ctx, ctcPaymail.Address, reqXPubID, c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	if contact == nil { // insert
		contact = newContact(
			ctcFName,
			ctcPaymail.Address,
			contactPki.PubKey,
			reqXPubID,
			ContactNotConfirmed,
			c.DefaultModelOptions(append(opts, New())...)...,
		)
	} else { // update
		contact.FullName = ctcFName
		contact.SetOptions(opts...)

		contact.UpdatePubKey(contactPki.PubKey)
	}

	if err = contact.Save(ctx); err != nil {
		return nil, fmt.Errorf("adding %s contact failed. Reason: %w", ctcPaymail, err)
	}

	return contact, nil
}

func (c *Client) logContactWarining(xPubID, cPaymail, warning string) {
	c.Logger().Warn().
		Str("xPubID", xPubID).
		Str("contact", cPaymail).
		Msg(warning)
}

func (c *Client) logContactError(xPubID, cPaymail, errorMsg string) {
	c.Logger().Error().
		Str("xPubID", xPubID).
		Str("contact", cPaymail).
		Msg(errorMsg)
}
