package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const cs_xpub = "xpub661MyMwAqRbcGFL3kTp9Y2fNccswbtC6gceUtkAfo2gn6k49BQbXqxmL1zqKe1MGLrx24S2a5FmK3G8hXtyk8wQS2VRyMNBG14NuxBHhevX"

func Test_ClientService_UpsertContact(t *testing.T) {
	t.Run("insert contact", func(t *testing.T) {
		// given
		paymailAddr := "bran_the_broken@winterfell.com"

		pt := &paymailTestMock{}
		pt.setup(t, "winterfell.com", true)
		defer pt.cleanup()

		pt.mockPki(paymailAddr, "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")
		pt.mockPike(paymailAddr)

		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup(), WithPaymailClient(pt.paymailClient))
		defer cleanup()

		_, err := client.NewXpub(ctx, cs_xpub, client.DefaultModelOptions()...)
		require.NoError(t, err)

		_, err = client.NewPaymailAddress(ctx, cs_xpub, "lady_stoneheart@winterfell.com", "Catelyn Stark", "", client.DefaultModelOptions()...)
		require.NoError(t, err)

		// when
		res, err := client.UpsertContact(ctx, "Bran Stark", paymailAddr, cs_xpub, "", client.DefaultModelOptions()...)

		// then
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, ContactNotConfirmed, res.Status)

	})

	t.Run("insert contact - no xpub", func(t *testing.T) {
		// given
		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup())
		defer cleanup()

		// when
		res, err := client.UpsertContact(ctx, "Bran Stark", "bran_the_broken@winterfell.com", cs_xpub, "", client.DefaultModelOptions()...)

		// then
		require.ErrorIs(t, err, ErrInvalidRequesterXpub)
		require.Nil(t, res)

	})

	t.Run("insert contact - contact's server doesn't support PIKE", func(t *testing.T) {
		// given
		paymailAddr := "bran_the_broken@winterfell.com"

		pt := &paymailTestMock{}
		pt.setup(t, "winterfell.com", false)
		defer pt.cleanup()

		pt.mockPki(paymailAddr, "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")

		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup(), WithPaymailClient(pt.paymailClient))
		defer cleanup()

		_, err := client.NewXpub(ctx, cs_xpub, client.DefaultModelOptions()...)
		require.NoError(t, err)

		_, err = client.NewPaymailAddress(ctx, cs_xpub, "lady_stoneheart@winterfell.com", "Catelyn Stark", "", client.DefaultModelOptions()...)
		require.NoError(t, err)

		// when
		res, err := client.UpsertContact(ctx, "Bran Stark", paymailAddr, cs_xpub, "lady_stoneheart@winterfell.com", client.DefaultModelOptions()...)

		// then
		require.ErrorIs(t, err, ErrAddingContactRequest)
		require.NotNil(t, res)
		require.Equal(t, ContactNotConfirmed, res.Status)
	})

	t.Run("update contact - PKI hasn't changed", func(t *testing.T) {
		// given
		paymailAddr := "bran_the_broken@winterfell.com"
		updatedFullname := "Brandon Stark"

		pt := &paymailTestMock{}
		pt.setup(t, "winterfell.com", true)
		defer pt.cleanup()

		pt.mockPki(paymailAddr, "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")
		pt.mockPike(paymailAddr)

		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup(), WithPaymailClient(pt.paymailClient))
		defer cleanup()

		_, err := client.NewXpub(ctx, cs_xpub, client.DefaultModelOptions()...)
		require.NoError(t, err)

		_, err = client.NewPaymailAddress(ctx, cs_xpub, "lady_stoneheart@winterfell.com", "Catelyn Stark", "", client.DefaultModelOptions()...)
		require.NoError(t, err)

		contact, err := client.UpsertContact(ctx, "Bran Stark", paymailAddr, cs_xpub, "", client.DefaultModelOptions()...)
		require.NoError(t, err)
		require.NotNil(t, contact)

		// confirm contact
		contact.Status = ContactConfirmed
		err = contact.Save(ctx)
		require.NoError(t, err)

		// when
		updatedContact, err := client.UpsertContact(ctx, updatedFullname, paymailAddr, cs_xpub, "", client.DefaultModelOptions()...)

		// then
		require.NoError(t, err)
		require.NotNil(t, updatedContact)

		require.Equal(t, updatedFullname, updatedContact.FullName)
		require.Equal(t, contact.PubKey, updatedContact.PubKey)

		// status shouldn't change
		require.Equal(t, ContactConfirmed, updatedContact.Status)

	})

	t.Run("update contact - PKI has changed", func(t *testing.T) {
		// given
		paymailAddr := "bran_the_broken@winterfell.com"
		updatedPki := "03c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90b"
		updatedFullname := "Brandon Stark"

		pt := &paymailTestMock{}
		pt.setup(t, "winterfell.com", true)
		defer pt.cleanup()

		pt.mockPki(paymailAddr, "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")
		pt.mockPike(paymailAddr)

		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup(), WithPaymailClient(pt.paymailClient))
		defer cleanup()

		_, err := client.NewXpub(ctx, cs_xpub, client.DefaultModelOptions()...)
		require.NoError(t, err)

		_, err = client.NewPaymailAddress(ctx, cs_xpub, "lady_stoneheart@winterfell.com", "Catelyn Stark", "", client.DefaultModelOptions()...)
		require.NoError(t, err)

		contact, err := client.UpsertContact(ctx, "Bran Stark", paymailAddr, cs_xpub, "", client.DefaultModelOptions()...)
		require.NoError(t, err)
		require.NotNil(t, contact)

		// confirm contact
		contact.Status = ContactConfirmed
		err = contact.Save(ctx)
		require.NoError(t, err)

		// when
		// change PKI
		pt.mockPki(paymailAddr, updatedPki)

		updatedContact, err := client.UpsertContact(ctx, updatedFullname, paymailAddr, cs_xpub, "lady_stoneheart@winterfell.com", client.DefaultModelOptions()...)

		// then
		require.NoError(t, err)
		require.NotNil(t, updatedContact)

		require.Equal(t, updatedFullname, updatedContact.FullName)
		require.Equal(t, updatedPki, updatedContact.PubKey)

		// status should back to unconfirmed
		require.Equal(t, ContactNotConfirmed, updatedContact.Status)
	})
}

func TestClientService_AddContactRequest(t *testing.T) {
	t.Run("add contact - new", func(t *testing.T) {
		// given
		paymailAddr := "sansa_stark@winterfell.com"

		pt := &paymailTestMock{}
		pt.setup(t, "winterfell.com", true)
		defer pt.cleanup()

		pt.mockPki(paymailAddr, "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")
		pt.mockPike(paymailAddr)

		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup(), WithPaymailClient(pt.paymailClient))
		defer cleanup()

		_, err := client.NewXpub(ctx, cs_xpub, client.DefaultModelOptions()...)
		require.NoError(t, err)

		_, err = client.NewPaymailAddress(ctx, cs_xpub, "lady_stoneheart@winterfell.com", "Catelyn Stark", "", client.DefaultModelOptions()...)
		require.NoError(t, err)

		// when
		res, err := client.AddContactRequest(ctx, "Sansa Stark", paymailAddr, cs_xpub, client.DefaultModelOptions()...)

		// then
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, ContactAwaitAccept, res.Status)
	})

	t.Run("add contact - already exist, PKI hasn't changed", func(t *testing.T) {
		// given
		paymailAddr := "sansa_stark@winterfell.com"

		pt := &paymailTestMock{}
		pt.setup(t, "winterfell.com", true)
		defer pt.cleanup()

		pt.mockPki(paymailAddr, "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")
		pt.mockPike(paymailAddr)

		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup(), WithPaymailClient(pt.paymailClient))
		defer cleanup()

		_, err := client.NewXpub(ctx, cs_xpub, client.DefaultModelOptions()...)
		require.NoError(t, err)

		_, err = client.NewPaymailAddress(ctx, cs_xpub, "lady_stoneheart@winterfell.com", "Catelyn Stark", "", client.DefaultModelOptions()...)
		require.NoError(t, err)

		contact, err := client.AddContactRequest(ctx, "Sansa Stark", paymailAddr, cs_xpub, client.DefaultModelOptions()...)
		require.NoError(t, err)
		require.NotNil(t, contact)

		// mark request as accepted
		contact.Status = ContactNotConfirmed
		err = contact.Save(ctx)
		require.NoError(t, err)

		// when
		updatedContact, err := client.AddContactRequest(ctx, "Alayne Stone", paymailAddr, cs_xpub, client.DefaultModelOptions()...)

		// then
		require.NoError(t, err)
		require.NotNil(t, updatedContact)

		require.Equal(t, contact.FullName, updatedContact.FullName)
		// status shouldn't change
		require.Equal(t, ContactNotConfirmed, updatedContact.Status)
	})

	t.Run("add contact - already exist, PKI has changed", func(t *testing.T) {
		// given
		paymailAddr := "sansa_stark@winterfell.com"
		updatedPki := "03c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90b"

		pt := &paymailTestMock{}
		pt.setup(t, "winterfell.com", true)
		defer pt.cleanup()

		pt.mockPki(paymailAddr, "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")
		pt.mockPike(paymailAddr)

		ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup(), WithPaymailClient(pt.paymailClient))
		defer cleanup()

		_, err := client.NewXpub(ctx, cs_xpub, client.DefaultModelOptions()...)
		require.NoError(t, err)

		_, err = client.NewPaymailAddress(ctx, cs_xpub, "lady_stoneheart@winterfell.com", "Catelyn Stark", "", client.DefaultModelOptions()...)
		require.NoError(t, err)

		contact, err := client.AddContactRequest(ctx, "Sansa Stark", paymailAddr, cs_xpub, client.DefaultModelOptions()...)
		require.NoError(t, err)
		require.NotNil(t, contact)

		// mark request as accepted
		contact.Status = ContactNotConfirmed
		err = contact.Save(ctx)
		require.NoError(t, err)

		// when
		// change PKI
		pt.mockPki(paymailAddr, updatedPki)

		updatedContact, err := client.AddContactRequest(ctx, "Alayne Stone", paymailAddr, cs_xpub, client.DefaultModelOptions()...)

		// then
		require.NoError(t, err)
		require.NotNil(t, updatedContact)

		require.Equal(t, contact.FullName, updatedContact.FullName)
		require.Equal(t, updatedPki, updatedContact.PubKey)

		// status should back to awaiting
		require.Equal(t, ContactAwaitAccept, updatedContact.Status)
	})
}

type paymailTestMock struct {
	serverUrl     string
	paymailClient paymail.ClientInterface
}

func (p *paymailTestMock) setup(t *testing.T, domain string, supportPike bool) {
	httpmock.Reset()
	serverURL := "https://" + domain + "/api/v1/" + paymail.DefaultServiceName

	wellKnownUrl := fmt.Sprintf("https://%s:443/.well-known/%s", domain, paymail.DefaultServiceName)
	wellKnownBody := paymail.CapabilitiesPayload{
		BsvAlias:     paymail.DefaultBsvAliasVersion,
		Capabilities: map[string]interface{}{paymail.BRFCPki: fmt.Sprintf("%s/id/{alias}@{domain.tld}", serverURL)},
	}

	if supportPike {
		wellKnownBody.Capabilities[paymail.BRFCPike] = fmt.Sprintf("%s/pike/{alias}@{domain.tld}", serverURL)
	}

	wellKnownReponse, _ := json.Marshal(wellKnownBody)
	wellKnonwResponder := httpmock.NewStringResponder(http.StatusOK, string(wellKnownReponse))
	httpmock.RegisterResponder(http.MethodGet, wellKnownUrl, wellKnonwResponder)

	p.serverUrl = serverURL
	p.paymailClient = newTestPaymailClient(t, []string{domain})
}

func (p *paymailTestMock) cleanup() {
	httpmock.Reset()
	p.serverUrl = ""
}

func (p *paymailTestMock) mockPki(paymail, pubkey string) {
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s/id/%s", p.serverUrl, paymail),
		httpmock.NewStringResponder(
			200,
			`{"bsvalias":"1.0","handle":"`+paymail+`","pubkey":"`+pubkey+`"}`,
		),
	)
}

func (p *paymailTestMock) mockPike(paymail string) {
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("%s/pike/%s", p.serverUrl, paymail),
		httpmock.NewStringResponder(
			200,
			"{}",
		),
	)
}
