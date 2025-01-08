package contact_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/contact/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func Test_ClientService_AdminCreateContact(t *testing.T) {
	tests := []struct {
		name           string
		contactPaymail string
		creatorPaymail string
		fullName       string
		metadata       *engine.Metadata
		storedContacts []engine.Contact
		setupMocks     []struct {
			Paymail string
			PubKey  string
		}
		engine           error
		expectedStatus   engine.ContactStatus
		expectedFullName string
	}{
		{
			name:           "Happy path without metadata",
			contactPaymail: fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail: fixtures.Sender.DefaultPaymail(),
			fullName:       "John Doe",
			metadata:       nil,
			engine:           nil,
			expectedStatus:   engine.ContactNotConfirmed,
			expectedFullName: "John Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given:
			// this should also create a client instance
			// we should be able to manipulate PKI and contacts in DB
			// it should also be able to add paymails and xpubs in a db
			given, then := testabilities.New(t)

			service, cleanup := given.Engine()
			defer cleanup()

			// when:
			res, err := service.AdminCreateContact(context.Background(), tt.contactPaymail, tt.creatorPaymail, tt.fullName, tt.metadata)

			// then:
			if tt.engine != nil {
				then.ErrorIs(err, tt.engine).WithNilResponse(res)
			} else {
				then.NoError(err).WithResponse(res).WithStatus(tt.expectedStatus).WithFullName(tt.expectedFullName)
			}

			// given
			// pt := &paymailTestMock{}
			// if tt.setupMocks != nil {
			// 	tt.setupMocks(pt)
			// }
			// defer pt.cleanup()
			//
			// ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup(), WithPaymailClient(pt.paymailClient))
			// defer cleanup()
			//
			// _, err := client.NewXpub(ctx, csXpub, client.DefaultModelOptions()...)
			// require.NoError(t, err)
			//
			// if tt.creatorPaymail != "unknown@example.com" && tt.creatorPaymail != "" {
			// 	_, err = client.NewPaymailAddress(ctx, csXpub, tt.creatorPaymail, "Jane Doe", "", client.DefaultModelOptions()...)
			// 	require.NoError(t, err)
			// }
			//
			// // when
			// res, err := client.AdminCreateContact(ctx, tt.contactPaymail, tt.creatorPaymail, tt.fullName, tt.metadata)

			// then
			// if tt.expectedError != nil {
			// 	require.ErrorIs(t, err, tt.expectedError)
			// 	require.Nil(t, res)
			// } else {
			// 	require.NoError(t, err)
			// 	require.NotNil(t, res)
			// 	require.Equal(t, tt.expectedStatus, res.Status)
			// 	require.Equal(t, tt.expectedFullName, res.FullName)
			// }
		})
	}
}

// func Test_ClientService_AdminCreateContact_ContactAlreadyExists(t *testing.T) {
// 	pt := &paymailTestMock{}
// 	pt.setup("example.com", true)
// 	pt.mockPki("user2@example.com", "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")
// 	pt.mockPki("user1@example.com", "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")
// 	pt.mockPike("user1@example.com")
// 	defer pt.cleanup()
//
// 	ctx, client, cleanup := CreateTestSQLiteClient(t, false, false, withTaskManagerMockup(), WithPaymailClient(pt.paymailClient))
// 	defer cleanup()
//
// 	_, err := client.NewXpub(ctx, csXpub, client.DefaultModelOptions()...)
// 	require.NoError(t, err)
//
// 	_, err = client.NewPaymailAddress(ctx, csXpub, "user2@example.com", "Jane Doe", "", client.DefaultModelOptions()...)
// 	require.NoError(t, err)
//
// 	contact := &Contact{
// 		ID:          uuid.NewString(),
// 		Model:       *NewBaseModel(ModelContact, client.DefaultModelOptions()...),
// 		FullName:    "Existing Contact",
// 		Paymail:     "user1@example.com",
// 		OwnerXpubID: csXpubHash,
// 		PubKey:      csXpub,
// 		Status:      ContactConfirmed,
// 	}
// 	err = contact.Save(ctx)
// 	require.NoError(t, err)
//
// 	// when
// 	res, err := client.AdminCreateContact(ctx, "user1@example.com", "user2@example.com", "John Doe", nil)
//
// 	// then
// 	require.ErrorIs(t, err, spverrors.ErrContactAlreadyExists)
// 	require.Nil(t, res)
// }
//
//
// 		{
// 	name:           "Happy path with metadata",
// 	contactPaymail: "user1@example.com",
// 	creatorPaymail: "user2@example.com",
// 	fullName:       "John Doe",
// 	metadata: &engine.Metadata{
// 		"key1": "value1",
// 		"key2": 42,
// 	},
// 	// setupMocks: func(pt *paymailTestMock) {
// 	// 	pt.setup("example.com", true)
// 	// 	pt.mockPki("user2@example.com", "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")
// 	// 	pt.mockPki("user1@example.com", "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")
// 	// 	pt.mockPike("user1@example.com")
// 	// },
// 	expectedError:    nil,
// 	expectedStatus:   engine.ContactNotConfirmed,
// 	expectedFullName: "John Doe",
// },
// {
// 	name:           "Edge case: Creator paymail not found",
// 	contactPaymail: "user1@example.com",
// 	creatorPaymail: "unknown@example.com",
// 	fullName:       "John Doe",
// 	metadata:       nil,
// 	// setupMocks: func(pt *paymailTestMock) {
// 	// 	pt.setup("example.com", true)
// 	// 	pt.mockPki("unknown@example.com", "")
// 	// },
// 	expectedError:    spverrors.ErrCouldNotFindPaymail,
// 	expectedStatus:   engine.ContactNotConfirmed,
// 	expectedFullName: "",
// },
// {
// 	name:           "Edge case: PKI retrieval fails",
// 	contactPaymail: "user1@example.com",
// 	creatorPaymail: "user2@example.com",
// 	fullName:       "John Doe",
// 	metadata:       nil,
// 	// setupMocks: func(pt *paymailTestMock) {
// 	// 	pt.setup("example.com", true)
// 	// 	pt.mockPki("user2@example.com", "04c85162f06f5391028211a3683d669301fc72085458ce94d0a9e77ba4ff61f90a")
// 	// },
// 	expectedError:    spverrors.ErrGettingPKIFailed,
// 	expectedStatus:   engine.ContactNotConfirmed,
// 	expectedFullName: "",
// },
// {
// 	name:           "Edge case: missing creator paymail",
// 	contactPaymail: "user1@example.com",
// 	creatorPaymail: "",
// 	fullName:       "John Doe",
// 	metadata:       nil,
// 	expectedError:  spverrors.ErrMissingContactCreatorPaymail,
// },
// {
// 	name:           "Edge case: missing contact full name",
// 	contactPaymail: "user1@example.com",
// 	creatorPaymail: "user2@example.com",
// 	fullName:       "",
// 	metadata:       nil,
// 	expectedError:  spverrors.ErrMissingContactFullName,
// },
//
