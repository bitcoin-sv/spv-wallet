package manualtests

import (
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
)

type ClientFactory = func(state *State) (*client.ClientWithResponses, error)

func AdminClientFactory(state *State) (*client.ClientWithResponses, error) {
	return state.AdminClient()
}

func CurrentUserClientFactory(state *State) (*client.ClientWithResponses, error) {
	return state.CurrentUserClient()
}

func RecipientClientFactory(state *State) (*client.ClientWithResponses, error) {
	return state.RecipientClient()
}

func UserClientFactoryWithID(userID string) ClientFactory {
	return func(state *State) (*client.ClientWithResponses, error) {
		return state.UserClient(userID)
	}
}

func AnonymousClientFactory(state *State) (*client.ClientWithResponses, error) {
	return state.AnonymousClient()
}

func UnknownUserClientFactory(state *State) (*client.ClientWithResponses, error) {
	return state.UnknownUserClient()
}
