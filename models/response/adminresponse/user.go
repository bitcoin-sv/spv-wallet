package adminresponse

import "time"

// User represents the normal user entity from the admin perspective
type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	PublicKey string    `json:"publicKey"`
	Paymails  []Paymail `json:"paymails"`
}
