package internal

import "fmt"

func PaymailAddress(alias, domain string) string {
	return fmt.Sprintf("%s@%s", alias, domain)
}
