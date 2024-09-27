package paymailaddress

import "context"

// Service is a component that provides methods for working with paymail address.
type Service interface {
	HasPaymailAddress(ctx context.Context, xPubID string, address string) (bool, error)
	GetDefaultPaymailAddress(ctx context.Context, xPubID string) (string, error)
}
