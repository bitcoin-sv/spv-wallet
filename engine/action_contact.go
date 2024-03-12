package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

var (
	ErrInvalidRequesterXpub = errors.New("invalid requester xpub")
	ErrAddingContactRequest = errors.New("adding contact request failed")
)

func (c *Client) UpsertContact(ctx context.Context, ctcFName, ctcPaymail, requesterXpub string, opts ...ModelOps) (*Contact, error) {
	reqXPubID := utils.Hash(requesterXpub)

	reqPms, err := c.GetPaymailAddressesByXPubID(ctx, reqXPubID, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(reqPms) == 0 {
		return nil, ErrInvalidRequesterXpub
	}

	reqPm := reqPms[0]

	pmSrvnt := &PaymailServant{
		cs: c.Cachestore(),
		pc: c.PaymailClient(),
	}

	contactPm := pmSrvnt.GetSanitizedPaymail(ctcPaymail)
	contactPki, err := pmSrvnt.GetPkiForPaymail(ctx, contactPm)
	if err != nil {
		return nil, fmt.Errorf("geting PKI for %s failed. Reason: %w", ctcPaymail, err)
	}

	data := newContactData{
		fullName: ctcFName,
		paymail:  contactPm,
		pubKey:   contactPki.PubKey,
		status:   ContactNotConfirmed,
		opts:     opts,
	}

	contact, err := c.addContact(ctx, &data, reqXPubID)
	if err != nil {
		return nil, fmt.Errorf("adding %s contact failed. Reason: %w", ctcPaymail, err)
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

	// add contact request
	data := newContactData{
		fullName: fullName,
		paymail:  contactPm,
		pubKey:   contactPki.PubKey,
		status:   ContactAwaitAccept,
		opts:     opts,
	}

	contactRequest, err := c.addContact(ctx, &data, requesterXPubID)
	if err != nil {
		return nil, fmt.Errorf("adding %s contact failed. Reason: %w", paymailAdress, err)
	}

	return contactRequest, nil
}

func (c *Client) addContact(ctx context.Context, data *newContactData, requesterXPubId string) (*Contact, error) {
	// check if exists already
	contact, err := getContact(ctx, data.paymail.adress, requesterXPubId, c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}
	if contact != nil {
		return contact, nil
	}

	contact = newContact(
		data.fullName,
		data.paymail.adress,
		data.pubKey,
		requesterXPubId,
		data.status,
		c.DefaultModelOptions(append(data.opts, New())...)...,
	)

	if err = contact.Save(ctx); err != nil {
		return nil, err
	}
	return contact, nil
}

type newContactData struct {
	fullName string
	paymail  *SanitizedPaymail
	pubKey   string
	status   ContactStatus
	opts     []ModelOps
}
