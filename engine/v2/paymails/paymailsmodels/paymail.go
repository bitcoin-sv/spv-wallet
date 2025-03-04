package paymailsmodels

import (
	"net/url"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailerrors"
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
		return paymailerrors.ErrInvalidAvatarURL.Wrap(err)
	}

	if URL.Scheme != "http" && URL.Scheme != "https" {
		return paymailerrors.ErrInvalidAvatarURL.Wrap(spverrors.Newf("avatarURL should have http(s) scheme"))
	}

	return nil
}
