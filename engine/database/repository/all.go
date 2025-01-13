package repository

import "gorm.io/gorm"

type All struct {
	Addresses  *Addresses
	Paymails   *Paymails
	Operations *Operations
	Users      *Users
	Outputs    *Outputs
}

func NewRepositories(db *gorm.DB) *All {
	return &All{
		Addresses:  NewAddressesRepo(db),
		Paymails:   NewPaymailsRepo(db),
		Operations: NewOperationsRepo(db),
		Users:      NewUsersRepo(db),
		Outputs:    NewOutputsRepo(db),
	}
}
