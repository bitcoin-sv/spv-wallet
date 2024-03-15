package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/mrz1836/go-datastore"
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

	data := contactData{
		fullName: ctcFName,
		paymail:  contactPm,
		pubKey:   contactPki.PubKey,
		status:   ContactNotConfirmed,
		opts:     opts,
	}

	contact, err := c.saveContact(ctx, &data, reqXPubID)
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

	data := contactData{
		fullName: fullName,
		paymail:  contactPm,
		pubKey:   contactPki.PubKey,
		status:   ContactAwaitAccept,
		opts:     opts,
	}

	contactRequest, err := c.saveContact(ctx, &data, requesterXPubID)
	if err != nil {
		return nil, fmt.Errorf("adding %s contact failed. Reason: %w", paymailAdress, err)
	}

	return contactRequest, nil
}

func (c *Client) saveContact(ctx context.Context, data *contactData, requesterXPubId string) (*Contact, error) {
	// check if exists already
	contact, err := getContact(ctx, data.paymail.adress, requesterXPubId, c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	if contact == nil { // insert
		contact = newContact(
			data.fullName,
			data.paymail.adress,
			data.pubKey,
			requesterXPubId,
			data.status,
			c.DefaultModelOptions(append(data.opts, New())...)...,
		)

	} else { // update
		contact.FullName = data.fullName
		contact.SetOptions(data.opts...) // add metada if exists
	}

	if err = contact.Save(ctx); err != nil {
		return nil, err
	}
	return contact, nil
}

type contactData struct {
	fullName string
	paymail  *SanitizedPaymail
	pubKey   string
	status   ContactStatus
	opts     []ModelOps
}

func (c *Client) UpdateContact(ctx context.Context, fullName, pubKey, xPubID, paymailAddr string, status ContactStatus, opts ...ModelOps) (*Contact, error) {
	contact, err := getContact(ctx, paymailAddr, xPubID, opts...)

	if err != nil {
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	if contact == nil {
		return nil, fmt.Errorf("contact not found")
	}

	if fullName != "" {
		contact.FullName = fullName
	}

	if pubKey != "" {
		contact.PubKey = pubKey
	}

	if status != "" {
		contact.Status = status
	}

	if paymailAddr != "" {
		contact.Paymail = paymailAddr
	}

	if err = contact.Save(ctx); err != nil {
		return nil, err
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
