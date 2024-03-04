package engine

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

var ErrInvalidRequesterPaymail = errors.New("invalid requester paymail address")

func (c *Client) AddContact(ctx context.Context, ctcFName, ctcPaymail, requesterPKey, requesterFName, requesterPaymail string, opts ...ModelOps) (*Contact, error) {
	requesterXPubId := utils.Hash(requesterPKey)

	reqPaymail, err := getPaymailAddress(ctx, requesterPaymail)
	if err != nil {
		return nil, err
	}
	if reqPaymail == nil || reqPaymail.XpubID != requesterXPubId {
		return nil, ErrInvalidRequesterPaymail
	}

	pmSrvnt := &PaymailServant{
		cs: c.Cachestore(),
		pc: c.PaymailClient(),
	}

	contactPaymail := pmSrvnt.GetSanitizedPaymail(ctcPaymail)
	contactPki, err := pmSrvnt.GetPkiForPaymail(ctx, contactPaymail)
	if err != nil {
		return nil, err
	}

	data := newContactData{
		fullName: ctcFName,
		paymail:  contactPaymail,
		pubKey:   contactPki.PubKey,
		status:   ContactStatusNotConf,
		opts:     opts,
	}

	contact, err := c.addContact(ctx, &data, requesterXPubId)
	if err != nil {
		return nil, err
	}

	// request new contact
	requesterContactRequest := paymail.PikeContactRequestPayload{
		FullName:      requesterFName,
		PaymailAdress: requesterPaymail,
	}
	if _, err = pmSrvnt.AddContactRequest(ctx, contactPaymail, &requesterContactRequest); err != nil {
		c.Logger().Warn().
			Str("requesterPaymil", requesterPaymail).
			Str("requestedContact", ctcPaymail).
			Msgf("adding contact request failed: %s", err.Error())
	}

	return contact, nil
}

func (c *Client) AddContactRequest(ctx context.Context, fullName, paymailAdress, requesterXPubID string, opts ...ModelOps) (*Contact, error) {
	pmSrvnt := &PaymailServant{
		cs: c.Cachestore(),
		pc: c.PaymailClient(),
	}

	contactPaymail := pmSrvnt.GetSanitizedPaymail(paymailAdress)
	contactPki, err := pmSrvnt.GetPkiForPaymail(ctx, contactPaymail)
	if err != nil {
		return nil, err
	}

	// add contact request
	data := newContactData{
		fullName: fullName,
		paymail:  contactPaymail,
		pubKey:   contactPki.PubKey,
		status:   ContactStatusAwaitAccept,
		opts:     opts,
	}

	contactRequest, err := c.addContact(ctx, &data, requesterXPubID)
	if err != nil {
		return nil, err
	}

	return contactRequest, nil
}

func (c *Client) addContact(ctx context.Context, data *newContactData, requesterXPubId string) (*Contact, error) {
	// check if exists already
	contact, err := getContact(ctx, data.paymail.adress, requesterXPubId)
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
