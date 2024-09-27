package fixtures

// User is a fixture that is representing a user of the system.
type User struct {
	Paymails []string
	XPubID   string
}

const (
	// PaymailDomain is the "our" paymail domain in the tests.
	PaymailDomain = "example.com"
	// PaymailDomainExternal is the "their"/external paymail domain in the tests.
	PaymailDomainExternal = "external.example.com"
)

var (
	// UserWithMorePaymails is a user with more than one paymail.
	UserWithMorePaymails = User{
		Paymails: []string{
			"tester@" + PaymailDomain,
			"secondPm@" + PaymailDomain,
		},
		XPubID: "user_multi_paymail_xpub_id",
	}

	// UserWithoutPaymail is a user without any paymail.
	UserWithoutPaymail = User{
		XPubID: "user_no_paymail_xpub_id",
	}

	// Sender is a user that is a sender in the tests.
	Sender = User{
		Paymails: []string{
			"sender@" + PaymailDomain,
		},
		XPubID: "sender_xpub_id",
	}

	// RecipientInternal is a user that is a recipient from "our" server in the tests.
	RecipientInternal = User{
		Paymails: []string{
			"recipient@" + PaymailDomain,
		},
		XPubID: "recipient_xpub_id",
	}

	// RecipientExternal is a user that is a recipient from external server in the tests.
	RecipientExternal = User{
		Paymails: []string{
			"recipient@" + PaymailDomainExternal,
		},
		XPubID: "",
	}
)

// DefaultPaymail returns the default paymail of this user.
func (f *User) DefaultPaymail() string {
	return f.Paymails[0]
}

// All returns all fixtures.
func All() []User {
	return []User{
		UserWithoutPaymail,
		UserWithMorePaymails,
		Sender,
		RecipientInternal,
		RecipientExternal,
	}
}
