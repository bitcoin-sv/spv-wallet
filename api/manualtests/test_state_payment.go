package manualtests

import "github.com/joomcode/errorx"

var PaymentConfigError = StateError.NewSubtype("payment")

type Payment struct {
	Amount          uint64 `mapstructure:"amount" yaml:"amount"`
	RecipientID     string `mapstructure:"recipient_id" yaml:"recipient_id"`
	ExternalPaymail string `mapstructure:"external_paymail" yaml:"external_paymail"`
	state           *State
}

func (p *Payment) ShouldGetInternalRecipient() (*User, error) {
	err := p.validateInternalRecipient()
	if err != nil {
		return nil, err
	}

	user, _ := p.state.GetUserById(p.RecipientID)

	err = user.init()
	if err != nil {
		return nil, errorx.Decorate(err, "couldn't prepare internal recipient")
	}

	return user, nil
}

func (p *Payment) ShouldGetInternalRecipientPaymail() (string, error) {
	user, err := p.ShouldGetInternalRecipient()
	if err != nil {
		return "", err
	}
	return user.PaymailAddress(), nil
}

func (p *Payment) ShouldGetExternalRecipientPaymail() (string, error) {
	if p.ExternalPaymail == "" || p.ExternalPaymail == notConfiguredExternalPaymail {
		return "", PaymentConfigError.New("Configure payment.external_paymail in file://%s", p.state.configFilePath)
	}

	return p.ExternalPaymail, nil
}

func (p *Payment) validateInternalRecipient() error {
	if p.RecipientID == "" || p.RecipientID == notConfiguredRecipientID {
		var suggestion string
		if len(p.state.OldUsers) == 0 {
			suggestion = "No old users to suggest, create new user with admin API"
		} else if p.state.OldUsers[len(p.state.OldUsers)-1].ID != p.state.User.ID {
			suggestion = "Use ID " + p.state.OldUsers[len(p.state.OldUsers)-1].ID
		} else {
			suggestion = "Use ID of one of old users (users from zzz_old_users)"
		}

		return PaymentConfigError.New("Configure payment.recipient_id in file://%s \n Suggestion: %s \n", p.state.configFilePath, suggestion)
	}

	_, err := p.state.GetOldUserById(p.RecipientID)
	if err != nil {
		return PaymentConfigError.Wrap(err, "Configured payment.recipient_id is not id of one of zzz_old_users in file://%s", p.state.configFilePath)
	}

	return nil
}
