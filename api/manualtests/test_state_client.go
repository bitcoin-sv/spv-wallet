package manualtests

import (
	"context"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/joomcode/errorx"
)

func (s *State) AdminClient() (*client.ClientWithResponses, error) {
	c, err := client.NewClientWithResponses(s.ServerURL, client.WithRequestEditorFn(s.authAsAdmin()))
	if err != nil {
		return nil, StateError.Wrap(err, "could not create admin client")
	}
	return c, nil
}

func (s *State) CurrentUserClient() (*client.ClientWithResponses, error) {
	c, err := client.NewClientWithResponses(s.ServerURL, client.WithRequestEditorFn(s.authAsCurrentUser()))
	if err != nil {
		return nil, StateError.Wrap(err, "could not create admin client")
	}
	return c, nil
}

func (s *State) RecipientClient() (*client.ClientWithResponses, error) {
	err := s.Payment.validateInternalRecipient()
	if err != nil {
		return nil, StateError.Wrap(err, "couldn't create client for recipient user")
	}

	userClient, err := s.UserClient(s.Payment.RecipientID)
	if err != nil {
		return nil, errorx.Decorate(err, "couldn't create client for recipient user")
	}

	return userClient, nil
}

func (s *State) UserClient(userID string) (*client.ClientWithResponses, error) {
	user, err := s.GetUserById(userID)
	if err != nil {
		return nil, StateError.Wrap(err, "couldn't create client for user %s", userID)
	}

	err = user.init()
	if err != nil {
		return nil, StateError.Wrap(err, "couldn't create client for user %s", userID)
	}

	c, err := client.NewClientWithResponses(s.ServerURL, client.WithRequestEditorFn(authWithKey(user.Xpub)))
	if err != nil {
		return nil, StateError.Wrap(err, "could not create admin client")
	}
	return c, nil
}

func (s *State) AnonymousClient() (*client.ClientWithResponses, error) {
	c, err := client.NewClientWithResponses(s.ServerURL)
	if err != nil {
		return nil, StateError.Wrap(err, "could not create admin client")
	}
	return c, nil
}

func (s *State) UnknownUserClient() (*client.ClientWithResponses, error) {
	c, err := client.NewClientWithResponses(s.ServerURL, client.WithRequestEditorFn(s.authAsUnknown()))
	if err != nil {
		return nil, StateError.Wrap(err, "could not create admin client")
	}
	return c, nil
}

func (s *State) authAsAdmin() client.RequestEditorFn {
	return authWithKey(s.AdminXpub)
}

func (s *State) authAsCurrentUser() client.RequestEditorFn {
	return authWithKey(s.User.Xpub)
}

func (s *State) authAsUnknown() client.RequestEditorFn {
	// unknownXpriv := "xprv9s21ZrQH143K3jw372AgTDppRcaFZZiikzCbZodX4CqzZM8SZdqkBgcmmp5DKtmaKFNcJq9AcpFd35X2oyKXJogmnziU5h7tb72qSFghNnA"
	unknownXpub := "xpub661MyMwAqRbcGE1WD3hgpMmYyeQjy2Sa8D8CNC38cYNyS9Tb7B9zjUwFd5ZSzFYjKMWsNzh1sxb5AspdaSp31Zqz9CDM5GL74Eib22h25Pv"
	return authWithKey(unknownXpub)
}

func authWithKey(xpub string) client.RequestEditorFn {
	return func(_ context.Context, req *http.Request) error {
		req.Header.Add("x-auth-xpub", xpub)
		return nil
	}
}
