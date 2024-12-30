package adminrequest

// CreateUser represents an admin request to create a new user as an admin.
type CreateUser struct {
	PublicKey string `json:"publicKey"`

	Paymail *AddPaymail `json:"paymail"` // creating paymail during user creation is optional
}

func (r *CreateUser) PaymailDefined() bool {
	return r.Paymail != nil
}
