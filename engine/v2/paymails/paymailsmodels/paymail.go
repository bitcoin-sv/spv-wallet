package paymailsmodels

import (
	"regexp"
	"time"
)

var explicitHTTPURLRegex = regexp.MustCompile(`^https?:\/\/(?:www\.)?[a-zA-Z0-9-]+(?:\.[a-zA-Z]{2,})+(?:\/[^\s]*)?`)

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

// CheckAvatarURL checks if avatar is either empty string or a proper url link
func (np *NewPaymail) CheckAvatarURL() bool {
	return np.Avatar == "" || explicitHTTPURLRegex.MatchString(np.Avatar)
}
