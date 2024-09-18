package paymailmock

import (
	"github.com/aws/aws-sdk-go/service/appmesh"
	"github.com/bitcoin-sv/go-paymail"
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
		endpoint: endpoint(appmesh.HttpMethodPost, P2PDestinationsForSats(1000).response()),
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
