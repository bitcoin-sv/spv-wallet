package manualtests

import (
	"context"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
)

func (s *State) AdminClient() (*client.ClientWithResponses, error) {
	c, err := client.NewClientWithResponses(s.ServerURL, client.WithRequestEditorFn(s.authAsAdmin))
	if err != nil {
		return nil, StateError.Wrap(err, "could not create admin client")
	}
	return c, nil
}

func (s *State) UserClient() (*client.ClientWithResponses, error) {
	c, err := client.NewClientWithResponses(s.ServerURL, client.WithRequestEditorFn(s.authAsUser))
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
	c, err := client.NewClientWithResponses(s.ServerURL, client.WithRequestEditorFn(s.authAsUnknown))
	if err != nil {
		return nil, StateError.Wrap(err, "could not create admin client")
	}
	return c, nil
}

func (s *State) authAsAdmin(_ context.Context, req *http.Request) error {
	req.Header.Add("x-auth-xpub", s.AdminXpub)
	return nil
}

func (s *State) authAsUser(_ context.Context, req *http.Request) error {
	req.Header.Add("x-auth-xpub", s.User.Xpub)
	return nil
}

func (s *State) authAsUnknown(_ context.Context, req *http.Request) error {
	// unknownXpriv := "xprv9s21ZrQH143K3jw372AgTDppRcaFZZiikzCbZodX4CqzZM8SZdqkBgcmmp5DKtmaKFNcJq9AcpFd35X2oyKXJogmnziU5h7tb72qSFghNnA"
	unknownXpub := "xpub661MyMwAqRbcGE1WD3hgpMmYyeQjy2Sa8D8CNC38cYNyS9Tb7B9zjUwFd5ZSzFYjKMWsNzh1sxb5AspdaSp31Zqz9CDM5GL74Eib22h25Pv"

	req.Header.Add("x-auth-xpub", unknownXpub)
	return nil
}
