package paymailmock

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/service/appmesh"
	"github.com/bitcoin-sv/go-paymail"
	"github.com/jarcoal/httpmock"
)

func capabilitySenderValidation() *CapabilityMock {
	return &CapabilityMock{
		name: paymail.BRFCSenderValidation,
		value: func(_ paymailDomainName) any {
			return false
		},
		endpoint: endpoint(http.MethodGet, func(request *http.Request) (*http.Response, error) {
			paymailAddress := request.URL.Path[strings.LastIndex(request.URL.Path, "/")+1:]
			alias := paymailAddress[:strings.Index(paymailAddress, "@")]
			pki := ""

			if alias == "recipient" {
				pki = "03bf409b6b2842150142c6b92cb11ba6a06310bdacd0ff2118a9b9da60ed994c2b"
			}

			if alias == "sender" {
				pki = "02ed100a85ac774757c967e2a7a8a1c7fdef901795805b494df69d7d02f663d259"
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

func capabilityPki() *CapabilityMock {
	return &CapabilityMock{
		name: paymail.BRFCPki,
		value: func(dn paymailDomainName) any {
			return dn.PKI()
		},
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
		endpoint: endpointWithStaticResponse(appmesh.HttpMethodPost, P2PDestinationsForSats(1000).response()),
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
