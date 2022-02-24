package pmail

import (
	"errors"
)

// ErrMissingPaymail missing paymail
var ErrMissingPaymail = errors.New("missing paymail")

// ErrMissingPaymailID missing id in paymail
var ErrMissingPaymailID = errors.New("missing id in paymail")

// ErrMissingPaymailAddress missing alias in paymail
var ErrMissingPaymailAddress = errors.New("missing alias in paymail")

// ErrMissingPaymailDomain missing domain in paymail
var ErrMissingPaymailDomain = errors.New("missing domain in paymail")

// ErrMissingPaymailExternalXPub missing external xPub in paymail
var ErrMissingPaymailExternalXPub = errors.New("missing external xPub in paymail")

// ErrMissingPaymailXPubID missing xpub_id in paymail
var ErrMissingPaymailXPubID = errors.New("missing xpub_id in paymail")
