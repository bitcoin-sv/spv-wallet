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
		contactPaymail string
		creatorPaymail string
		fullName       string
		metadata       *engine.Metadata
	}{
		"Create contact without metadata": {
			contactPaymail: fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail: fixtures.Sender.DefaultPaymail(),
			fullName:       "John Doe",
			metadata:       nil,
		},
		"Create contact with metadata": {
			contactPaymail: fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail: fixtures.Sender.DefaultPaymail(),
			fullName:       "John Doe",
			metadata: &engine.Metadata{
				"key1": "value1",
				"key2": 420,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//given:
			given, then := testabilities.New(t)

			service, cleanup := given.Engine()
			defer cleanup()

			//when:
			contact, err := service.AdminCreateContact(context.Background(),
				tt.contactPaymail,
				tt.creatorPaymail,
				tt.fullName,
				tt.metadata,
			)

			//then:
			then.
				Contact(contact).
				WithNoError(err).
				AsNotConfirmed().
				WithFullName(tt.fullName)
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
		given.
			ExternalPaymailServer().
			WillRespondOnCapability(paymail.BRFCPki).
			WithInternalServerError()

		//when:
		contact, err := service.AdminCreateContact(context.Background(),
			fixtures.RecipientExternal.DefaultPaymail(),
			fixtures.Sender.DefaultPaymail(),
			"John Doe",
			nil,
		)

		//then:
		then.Contact(contact).WithError(err).ThatIs(spverrors.ErrGettingPKIFailed)
	})
}

func Test_ClientService_AdminCreateContact_Fail(t *testing.T) {
	tests := map[string]struct {
		contactPaymail string
		creatorPaymail string
		fullName       string
		expectedError  error
	}{
		"Should fail when creator paymail not found": {
			contactPaymail: fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail: "not_exist@example.com",
			fullName:       "John Doe",
			expectedError:  spverrors.ErrCouldNotFindPaymail,
		},
		"Should fail when missing creator paymail": {
			contactPaymail: fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail: "",
			fullName:       "John Doe",
			expectedError:  spverrors.ErrMissingContactCreatorPaymail,
		},
		"Should fail when missing contact full name": {
			contactPaymail: fixtures.RecipientExternal.DefaultPaymail(),
			creatorPaymail: fixtures.Sender.DefaultPaymail(),
			fullName:       "",
			expectedError:  spverrors.ErrMissingContactFullName,
		},
		"Should fail when missing contact paymail": {
			contactPaymail: "",
			creatorPaymail: fixtures.Sender.DefaultPaymail(),
			fullName:       "John Doe",
			expectedError:  spverrors.ErrMissingContactPaymailParam,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//given:
			given, then := testabilities.New(t)

			service, cleanup := given.Engine()
			defer cleanup()

			//when:
			contact, err := service.AdminCreateContact(context.Background(),
				tt.contactPaymail,
				tt.creatorPaymail,
				tt.fullName,
				nil,
			)

			//then:
			then.Contact(contact).WithError(err).ThatIs(tt.expectedError)
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
		contact, err := service.AdminCreateContact(context.Background(),
			fixtures.RecipientExternal.DefaultPaymail(),
			fixtures.Sender.DefaultPaymail(),
			"John Doe",
			nil,
		)
		then.Contact(contact).WithNoError(err)

		//when:
		contact, err = service.AdminCreateContact(context.Background(),
			fixtures.RecipientExternal.DefaultPaymail(),
			fixtures.Sender.DefaultPaymail(),
			"John Doe",
			nil,
		)

		//then:
		then.Contact(contact).WithError(err).ThatIs(spverrors.ErrContactAlreadyExists)
	})
}
