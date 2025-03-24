package paymailsmodels

import (
	"net/url"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailerrors"
	"github.com/bitcoin-sv/spv-wallet/errdef"
)

// Paymail represents a domain model from paymails service
type Paymail struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	// TODO: Handle DeletedAt

	Alias  string
	Domain string

	PublicName string
	Avatar     string

	UserID string
}

// NewPaymail represents data for creating a new paymail
type NewPaymail struct {
	Alias      string
	Domain     string
	PublicName string
	Avatar     string
	UserID     string
}

// ValidateAvatar checks if avatar is either empty string or a proper url link
func (np *NewPaymail) ValidateAvatar() error {
	if np.Avatar == "" {
		return nil
	}

	URL, err := url.Parse(np.Avatar)
	if err != nil {
		return paymailerrors.InvalidAvatarURL.
			Wrap(err, "avatarURL is not a valid URL: %s", np.Avatar).
			WithProperty(errdef.PropPublicHint, "Avatar should be a valid URL")
	}

	if URL.Scheme != "http" && URL.Scheme != "https" {
		return paymailerrors.InvalidAvatarURL.
			New("avatar has not valid scheme (http or https): %s", np.Avatar).
			WithProperty(errdef.PropPublicHint, "Avatar should have a valid scheme (http or https)")
	}

	return nil
}
