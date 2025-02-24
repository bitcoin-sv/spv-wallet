package manualtests

import "github.com/bitcoin-sv/spv-wallet/api/manualtests/client"

type ClientFactory = func(state *State) (*client.ClientWithResponses, error)

func AdminClientFactory(state *State) (*client.ClientWithResponses, error) {
	return state.AdminClient()
}

func UserClientFactory(state *State) (*client.ClientWithResponses, error) {
	return state.UserClient()
}

func AnonymousClientFactory(state *State) (*client.ClientWithResponses, error) {
	return state.AnonymousClient()
}

func UnknownUserClientFactory(state *State) (*client.ClientWithResponses, error) {
	return state.UnknownUserClient()
}
