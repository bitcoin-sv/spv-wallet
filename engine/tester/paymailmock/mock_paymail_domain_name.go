package paymailmock

import (
	"github.com/bitcoin-sv/go-paymail"
)

type paymailDomainName string

// ServerURL returns the server URL for the paymail domain
func (domain paymailDomainName) ServerURL() string {
	return "https://" + string(domain) + "/" + paymail.DefaultServiceName
}

// CapabilitiesURL returns the capabilities URL for the paymail domain
func (domain paymailDomainName) CapabilitiesURL() string {
	return "https://" + string(domain) + ":443/.well-known/" + paymail.DefaultServiceName
}

// PKI returns the PKI URL for the paymail domain
func (domain paymailDomainName) PKI() string {
	return domain.template("pki")
}

// PaymentDestination returns the payment destination URL for the paymail domain
func (domain paymailDomainName) PaymentDestination() string {
	return domain.template("address")
}

// P2PTransaction returns the P2P transaction URL for the paymail domain
func (domain paymailDomainName) P2PTransaction() string {
	return domain.template("receive-transaction")
}

// BEEFTransaction returns the BEEF transaction URL for the paymail domain
func (domain paymailDomainName) BEEFTransaction() string {
	return domain.template("receive-beef")
}

// P2PPaymentDestination returns the P2P payment destination URL for the paymail domain
func (domain paymailDomainName) P2PPaymentDestination() string {
	return domain.template("p2p-payment-destination")
}

func (domain paymailDomainName) template(path string) string {
	return domain.ServerURL() + "/" + path + "/{alias}@{domain.tld}"
}
