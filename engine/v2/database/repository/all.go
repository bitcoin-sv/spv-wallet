package repository

import "gorm.io/gorm"

// All holds all repositories.
type All struct {
	Addresses  *Addresses
	Paymails   *Paymails
	Operations *Operations
	Users      *Users
	Outputs    *Outputs
	Data       *Data
}

// NewRepositories creates a new holder for all repositories.
func NewRepositories(db *gorm.DB) *All {
	return &All{
		Addresses:  NewAddressesRepo(db),
		Paymails:   NewPaymailsRepo(db),
		Operations: NewOperationsRepo(db),
		Users:      NewUsersRepo(db),
		Outputs:    NewOutputsRepo(db),
		Data:       NewDataRepo(db),
	}
}
