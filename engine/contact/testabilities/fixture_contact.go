package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/go-paymail"
	// "github.com/bitcoin-sv/spv-wallet/engine"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
)

type contactFixture struct {
	paymailsFixture *paymailsFixture
	engineFixture   testengine.EngineFixture
	t               testing.TB
}

type paymailsFixture struct {
	serverURL     string
	paymailClient paymail.ClientInterface
}

type ContactFixture interface {
	// Paymails() PaymailsFixture
	// Engine() (engine.ClientInterface, func())
	Engine() (testengine.EngineWithConfig, func())
}

// type PaymailsFixture interface {
// 	RegisterPaymail(domain string, supportPike bool) (PaymailsFixture, func())
// 	MockPki(paymail, pubkey string) PaymailsFixture
// 	MockPike(paymail string) PaymailsFixture
// }

func given(t testing.TB) ContactFixture {
	return &contactFixture{
		t: t,
		// paymailsFixture: &paymailsFixture{},
		engineFixture: testengine.Given(t),
	}
}

// func (cf *contactFixture) Engine() (engine.ClientInterface, func()) {
func (cf *contactFixture) Engine() (testengine.EngineWithConfig, func()) {
	// engine, cleanup := cf.engineFixture.Engine()

	// cf.engineFixture.PaymailClient().
	return cf.engineFixture.Engine()
	// return engine.Engine, cleanup

}

// func (cf *contactFixture) Paymails() PaymailsFixture {
// 	return cf.paymailsFixture
// }
//
// func (p *paymailsFixture) RegisterPaymail(domain string, supportPike bool) (PaymailsFixture, func()) {
// 	httpmock.Reset()
// 	serverURL := "https://" + domain + "/api/v1/" + paymail.DefaultServiceName
//
// 	wellKnownURL := fmt.Sprintf("https://%s:443/.well-known/%s", domain, paymail.DefaultServiceName)
// 	wellKnownBody := paymail.CapabilitiesPayload{
// 		BsvAlias:     paymail.DefaultBsvAliasVersion,
// 		Capabilities: map[string]interface{}{paymail.BRFCPki: fmt.Sprintf("%s/id/{alias}@{domain.tld}", serverURL)},
// 	}
//
// 	if supportPike {
// 		wellKnownBody.Capabilities[paymail.BRFCPike] = map[string]string{
// 			paymail.BRFCPikeInvite:  fmt.Sprintf("%s/contact/invite/{alias}@{domain.tld}", serverURL),
// 			paymail.BRFCPikeOutputs: fmt.Sprintf("%s/pike/outputs/{alias}@{domain.tld}", serverURL),
// 		}
// 	}
//
// 	wellKnownResponse, _ := json.Marshal(wellKnownBody)
// 	wellKnownResponder := httpmock.NewStringResponder(http.StatusOK, string(wellKnownResponse))
// 	httpmock.RegisterResponder(http.MethodGet, wellKnownURL, wellKnownResponder)
//
// 	p.serverURL = serverURL
// 	p.paymailClient = xtester.MockClient(domain)
//
// 	cleanup := func() {
//
// 		httpmock.Reset()
// 		p.serverURL = ""
// 	}
//
// 	return p, cleanup
// }
//
// func (p *paymailsFixture) MockPki(paymail, pubkey string) PaymailsFixture {
// 	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s/id/%s", p.serverURL, paymail),
// 		httpmock.NewStringResponder(
// 			200,
// 			`{"bsvalias":"1.0","handle":"`+paymail+`","pubkey":"`+pubkey+`"}`,
// 		),
// 	)
//
// 	return p
// }
//
// func (p *paymailsFixture) MockPike(paymail string) PaymailsFixture {
// 	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/contact/invite/%s", p.serverURL, paymail),
// 		httpmock.NewStringResponder(
// 			200,
// 			"{}",
// 		),
// 	)
// 	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/pike/outputs%s", p.serverURL, paymail),
// 		httpmock.NewStringResponder(
// 			200,
// 			"{}",
// 		),
// 	)
//
// 	return p
// }
//
