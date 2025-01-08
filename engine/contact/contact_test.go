package contact_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/contact/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func Test_ClientService_AdminCreateContact_Success(t *testing.T) {
	tests := map[string]struct {
		contactPaymail   string
		creatorPaymail   string
		fullName         string
		metadata         *engine.Metadata
		expectedError    error
		expectedStatus   engine.ContactStatus
		expectedFullName string
	}{
		"Happy path without metadata": {
			contactPaymail:   fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail:   fixtures.Sender.DefaultPaymail(),
			fullName:         "John Doe",
			metadata:         nil,
			expectedError:    nil,
			expectedStatus:   engine.ContactNotConfirmed,
			expectedFullName: "John Doe",
		},
		"Happy path with metadata": {
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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//given:
			given, then := testabilities.New(t)

			service, cleanup := given.Engine()
			defer cleanup()

			//when:
			res, err := service.AdminCreateContact(context.Background(), tt.contactPaymail, tt.creatorPaymail, tt.fullName, tt.metadata)

			//then:
			then.NoError(err).WithResponse(res).WithStatus(tt.expectedStatus).WithFullName(tt.expectedFullName)
		})
	}
}

func Test_ClientService_AdminCreateContact_PKIRetrievalFail(t *testing.T) {
	t.Run("Should fail with PKI retrieval", func(t *testing.T) {
		//given:
		given, then := testabilities.New(t)

		service, cleanup := given.Engine()
		defer cleanup()

		//and:
		given.PaymailClient().WillRespondOnCapability(paymail.BRFCPki).WithInternalServerError()

		//when:
		res, err := service.AdminCreateContact(context.Background(), fixtures.RecipientExternal.DefaultPaymail(), fixtures.Sender.DefaultPaymail(), "John Doe", nil)

		//then:
		then.ErrorIs(err, spverrors.ErrGettingPKIFailed).WithNilResponse(res)
	})
}

func Test_ClientService_AdminCreateContact_Fail(t *testing.T) {
	tests := map[string]struct {
		contactPaymail   string
		creatorPaymail   string
		fullName         string
		expectedError    error
		expectedStatus   engine.ContactStatus
		expectedFullName string
	}{
		"Edge case: Creator paymail not found": {
			contactPaymail:   fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail:   "not_exist@example.com",
			fullName:         "John Doe",
			expectedError:    spverrors.ErrCouldNotFindPaymail,
			expectedStatus:   engine.ContactNotConfirmed,
			expectedFullName: "",
		},
		"Edge case: missing creator paymail": {
			contactPaymail: fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail: "",
			fullName:       "John Doe",
			expectedError:  spverrors.ErrMissingContactCreatorPaymail,
		},
		"Edge case: missing contact full name": {
			contactPaymail: fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail: fixtures.Sender.DefaultPaymail(),
			fullName:       "",
			expectedError:  spverrors.ErrMissingContactFullName,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//given:
			given, then := testabilities.New(t)

			service, cleanup := given.Engine()
			defer cleanup()

			//when:
			res, err := service.AdminCreateContact(context.Background(), tt.contactPaymail, tt.creatorPaymail, tt.fullName, nil)

			//then:
			then.ErrorIs(err, tt.expectedError).WithNilResponse(res)
		})
	}

}

func Test_ClientService_AdminCreateContact_ContactAlreadyExists(t *testing.T) {
	t.Run("Should fail the second time due to contact already exists", func(t *testing.T) {
		//given:
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
