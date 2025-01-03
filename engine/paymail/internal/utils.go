package internal

import "fmt"

// PaymailAddress returns a paymail address from an alias and domain.
func PaymailAddress(alias, domain string) string {
	return fmt.Sprintf("%s@%s", alias, domain)
}
