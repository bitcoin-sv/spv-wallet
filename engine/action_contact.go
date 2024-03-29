package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

var (
	ErrInvalidRequesterXpub         = errors.New("invalid requester xpub")
	ErrAddingContactRequest         = errors.New("adding contact request failed")
	ErrMoreThanOnePaymailRegistered = errors.New("there are more than one paymail assigned to the xpub")
	ErrContactNotFound              = errors.New("contact not found")
	ErrContactIncorrectStatus       = errors.New("contact is in incorrect status to proceed")
)

func (c *Client) UpsertContact(ctx context.Context, ctcFName, ctcPaymail, requesterXpub, requesterPaymail string, opts ...ModelOps) (*Contact, error) {
	reqXPubID := utils.Hash(requesterXpub)
	reqPm, err := c.getPaymail(ctx, reqXPubID, requesterPaymail)
	if err != nil {
		return nil, err
	}

	pmSrvnt := &PaymailServant{
		cs: c.Cachestore(),
		pc: c.PaymailClient(),
	}
	contactPm := pmSrvnt.GetSanitizedPaymail(ctcPaymail)

	contact, err := c.upsertContact(ctx, pmSrvnt, reqXPubID, ctcFName, contactPm, opts...)
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

		return contact, ErrAddingContactRequest
	}

	return contact, nil
}

func (c *Client) AddContactRequest(ctx context.Context, fullName, paymailAdress, requesterXPubID string, opts ...ModelOps) (*Contact, error) {
	pmSrvnt := &PaymailServant{
		cs: c.Cachestore(),
		pc: c.PaymailClient(),
	}

	contactPm := pmSrvnt.GetSanitizedPaymail(paymailAdress)
	contactPki, err := pmSrvnt.GetPkiForPaymail(ctx, contactPm)
	if err != nil {
		return nil, fmt.Errorf("geting PKI for %s failed. Reason: %w", paymailAdress, err)
	}

	// check if exists already
	contact, err := getContact(ctx, contactPm.adress, requesterXPubID, c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	save := false
	if contact != nil {
		save = contact.UpdatePubKey(contactPki.PubKey)
	} else {
		contact = newContact(
			fullName,
			contactPm.adress,
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

func (c *Client) getPaymail(ctx context.Context, xpubID, paymailAddr string) (*PaymailAddress, error) {
	if paymailAddr != "" {
		res, err := c.GetPaymailAddress(ctx, paymailAddr, c.DefaultModelOptions()...)
		if err != nil {
			return nil, err
		}

		if res == nil || res.XpubID != xpubID {
			return nil, ErrInvalidRequesterXpub
		}

		return res, nil
	}

	emptyConditions := make(map[string]interface{})

	paymails, err := c.GetPaymailAddressesByXPubID(ctx, xpubID, nil, &emptyConditions, nil)
	if err != nil {
		return nil, err
	}
	if len(paymails) == 0 {
		return nil, ErrInvalidRequesterXpub
	} else if len(paymails) > 1 {
		return nil, ErrMoreThanOnePaymailRegistered
	}

	return paymails[0], nil
}

func (c *Client) upsertContact(ctx context.Context, pmSrvnt *PaymailServant, reqXPubID, ctcFName string, ctcPaymail *SanitizedPaymail, opts ...ModelOps) (*Contact, error) {

	contactPki, err := pmSrvnt.GetPkiForPaymail(ctx, ctcPaymail)
	if err != nil {
		return nil, fmt.Errorf("geting PKI for %s failed. Reason: %w", ctcPaymail.adress, err)
	}

	// check if exists already
	contact, err := getContact(ctx, ctcPaymail.adress, reqXPubID, c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	if contact == nil { // insert
		contact = newContact(
			ctcFName,
			ctcPaymail.adress,
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

func (c *Client) GetContacts(ctx context.Context, metadata *Metadata, conditions *map[string]interface{}, queryParams *datastore.QueryParams, opts ...ModelOps) ([]*Contact, error) {
	ctx = c.GetOrStartTxn(ctx, "get_contacts")

	contacts, err := getContacts(ctx, metadata, conditions, queryParams, c.DefaultModelOptions(opts...)...)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (c *Client) AcceptContact(ctx context.Context, xPubID, paymail string) error {

	contact, err := getContact(ctx, paymail, xPubID, c.DefaultModelOptions()...)
	if err != nil {
		c.logContactError(xPubID, paymail, fmt.Sprintf("unexpected error while geting contact: %s", err.Error()))
		return err
	}
	if contact == nil {
		return ErrContactNotFound
	}

	if err = contact.Accept(); err != nil {
		c.logContactWarining(xPubID, paymail, err.Error())
		return ErrContactIncorrectStatus
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
		return ErrContactNotFound
	}

	if err = contact.Reject(); err != nil {
		c.logContactWarining(xPubID, paymail, err.Error())
		return ErrContactIncorrectStatus
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
		return ErrContactNotFound
	}

	if err = contact.Confirm(); err != nil {
		c.logContactWarining(xPubID, paymail, err.Error())
		return ErrContactIncorrectStatus
	}

	if err = contact.Save(ctx); err != nil {
		c.logContactError(xPubID, paymail, fmt.Sprintf("unexpected error while saving contact: %s", err.Error()))
		return err
	}

	return nil
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
