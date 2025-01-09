package paymailmock

import (
	"net/http"
	"strings"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/jarcoal/httpmock"
)

func capabilitySenderValidation() *CapabilityMock {
	return &CapabilityMock{
		name: paymail.BRFCSenderValidation,
		value: func(_ paymailDomainName) any {
			return false
		},
	}
}

func capabilityPki() *CapabilityMock {
	return &CapabilityMock{
		name: paymail.BRFCPki,
		value: func(dn paymailDomainName) any {
			return dn.PKI()
		},
		endpoint: endpoint(http.MethodGet, func(request *http.Request) (*http.Response, error) {
			paymailAddress := request.URL.Path[strings.LastIndex(request.URL.Path, "/")+1:]
			alias := paymailAddress[:strings.Index(paymailAddress, "@")]
			pki := ""

			if alias == "recipient" {
				pki = fixtures.RecipientExternalPKI
			}

			if alias == "sender" {
				pki = fixtures.SenderPKI
			}

			if pki != "" {
				resp := obj{
					"bsvalias": "1.0",
					"handle":   paymailAddress,
					"pubkey":   pki,
				}
				return httpmock.NewJsonResponse(http.StatusOK, resp)
			} else {
				return httpmock.NewJsonResponse(http.StatusNotFound, obj{
					"message": "Not Found",
				})
			}
		}),
	}
}

func capabilityPaymentDestination() *CapabilityMock {
	return &CapabilityMock{
		name: paymail.BRFCPaymentDestination,
		value: func(dn paymailDomainName) any {
			return dn.PaymentDestination()
		},
	}
}

func capabilityP2PTransaction() *CapabilityMock {
	return &CapabilityMock{
		name: paymail.BRFCP2PTransactions,
		value: func(dn paymailDomainName) any {
			return dn.P2PTransaction()
		},
	}
}

func capabilityP2PPaymentDestination() *CapabilityMock {
	return &CapabilityMock{
		name: paymail.BRFCP2PPaymentDestination,
		value: func(dn paymailDomainName) any {
			return dn.P2PPaymentDestination()
		},
		endpoint: endpointWithStaticResponse(http.MethodPost, P2PDestinationsForSats(1000).response()),
	}
}

func capabilityBEEFTransaction() *CapabilityMock {
	return &CapabilityMock{
		name: paymail.BRFCBeefTransaction,
		value: func(dn paymailDomainName) any {
			return dn.BEEFTransaction()
		},
	}
}
