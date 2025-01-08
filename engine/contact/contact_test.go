package contact_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/contact/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func Test_ClientService_AdminCreateContact_Success(t *testing.T) {
	tests := []struct {
		name             string
		contactPaymail   string
		creatorPaymail   string
		fullName         string
		metadata         *engine.Metadata
		expectedError    error
		expectedStatus   engine.ContactStatus
		expectedFullName string
	}{
		{
			name:             "Happy path without metadata",
			contactPaymail:   fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail:   fixtures.Sender.DefaultPaymail(),
			fullName:         "John Doe",
			metadata:         nil,
			expectedError:    nil,
			expectedStatus:   engine.ContactNotConfirmed,
			expectedFullName: "John Doe",
		},
		{
			name:           "Happy path with metadata",
			contactPaymail: fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail: fixtures.Sender.DefaultPaymail(),
			fullName:       "John Doe",
			metadata: &engine.Metadata{
				"key1": "value1",
				"key2": 420,
			},
			expectedError:    nil,
			expectedStatus:   engine.ContactNotConfirmed,
			expectedFullName: "John Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given:
			given, then := testabilities.New(t)

			service, cleanup := given.Engine()
			defer cleanup()

			// when:
			res, err := service.AdminCreateContact(context.Background(), tt.contactPaymail, tt.creatorPaymail, tt.fullName, tt.metadata)

			then.NoError(err).WithResponse(res).WithStatus(tt.expectedStatus).WithFullName(tt.expectedFullName)
		})
	}
}

func Test_ClientService_AdminCreateContact_ContactAlreadtExists(t *testing.T) {
	t.Run("Should fail the second time due to contact already exists", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)

		service, cleanup := given.Engine()
		defer cleanup()

		//and:
		_, err := service.AdminCreateContact(context.Background(), fixtures.RecipientExternal.DefaultPaymail(), fixtures.Sender.DefaultPaymail(), "John Doe", nil)
		then.NoError(err)

		//when:
		res, err := service.AdminCreateContact(context.Background(), fixtures.RecipientExternal.DefaultPaymail(), fixtures.Sender.DefaultPaymail(), "John Doe", nil)

		//then:
		then.ErrorIs(err, spverrors.ErrContactAlreadyExists).WithNilResponse(res)
	})
}

// func Test_ClientService_AdminCreateContact_ContactAlreadyExists(t *testing.T) {
//
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
