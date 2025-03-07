package manualtests

import (
	"encoding/json"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/stretchr/testify/require"
)

type (
	GenericCall[R Result]          = func(c *client.ClientWithResponses) (R, error)
	Call                           = GenericCall[Result]
	CallWithT                      = func(t testing.TB, c *client.ClientWithResponses) (Result, error)
	GenericCallWithState[R Result] = func(state StateForCall, c *client.ClientWithResponses) (R, error)
	CallWithState                  = GenericCallWithState[Result]
)

// ToCall because go generics are stupid, we need to have this wrapping.
func ToCall[R Result](f GenericCall[R]) Call {
	return func(c *client.ClientWithResponses) (Result, error) {
		result, err := f(c)
		return result, err
	}
}

type APICall struct {
	t      testing.TB
	state  *State
	client *client.ClientWithResponses
}

type StateForCall struct {
	*State
	T testing.TB
}

func APICallForAdmin(t testing.TB) *APICall {
	return APICallFor(t, AdminClientFactory)
}

func APICallForCurrentUser(t testing.TB) *APICall {
	return APICallFor(t, CurrentUserClientFactory)
}

func APICallForRecipient(t testing.TB) *APICall {
	return APICallFor(t, RecipientClientFactory)
}

func APICallForUserWithID(t testing.TB, userID string) *APICall {
	return APICallFor(t, UserClientFactoryWithID(userID))
}

func APICallForAnonymous(t testing.TB) *APICall {
	return APICallFor(t, AnonymousClientFactory)
}

func APICallForUnknownUser(t testing.TB) *APICall {
	return APICallFor(t, UnknownUserClientFactory)
}

func APICallFor(t testing.TB, clientFactory ClientFactory) *APICall {
	state := NewState()
	err := state.Load()
	require.NoError(t, err)

	apiClient, err := clientFactory(state)
	require.NoError(t, err)

	return &APICall{
		t:      t,
		state:  state,
		client: apiClient,
	}
}

//goland:noinspection GoMixedReceiverTypes // This is intentional
func (a APICall) WithT(t testing.TB) *APICall {
	a.t = t
	return &a
}

// CallWithStateForSuccess alias for CallWithState(callback).RequireSuccess()
func (a *APICall) CallWithStateForSuccess(callback CallWithState) *APICallResponse {
	call := a.CallWithState(callback)
	call.RequireSuccess()
	return call
}

// CallForSuccess alias for Call(callback).RequireSuccess()
func (a *APICall) CallForSuccess(callback Call) *APICallResponse {
	call := a.Call(callback)
	call.RequireSuccess()
	return call
}

func (a *APICall) Call(callback Call) *APICallResponse {
	return a.CallWithState(func(_ StateForCall, client *client.ClientWithResponses) (Result, error) {
		return callback(client)
	})
}

func (a *APICall) CallWithT(callback CallWithT) *APICallResponse {
	return a.CallWithState(func(state StateForCall, client *client.ClientWithResponses) (Result, error) {
		return callback(state.T, client)
	})
}

func (a *APICall) CallWithState(callback CallWithState) *APICallResponse {
	res, err := callback(StateForCall{
		State: a.state,
		T:     a.t,
	},
		a.client,
	)
	require.NoError(a.t, err)

	Print(res)

	return &APICallResponse{
		result:  res,
		apiCall: a,
		t:       a.t,
	}
}

// CallWithUpdateState calls an API and requires it to success
// ALSO if success this method is saving the state.
func (a *APICall) CallWithUpdateState(callback CallWithState) *APICallResponse {
	res, err := callback(StateForCall{
		State: a.state,
		T:     a.t,
	},
		a.client,
	)
	require.NoError(a.t, err)

	Print(res)

	err = a.state.SaveOnSuccess(res)
	require.NoError(a.t, err)

	apiCallResponse := &APICallResponse{
		result:  res,
		apiCall: a,
		t:       a.t,
	}

	apiCallResponse.RequireSuccess()

	return apiCallResponse
}

func (a *APICall) State() *State {
	return a.state
}

type APICallResponse struct {
	t       testing.TB
	result  Result
	apiCall *APICall
}

func (r *APICallResponse) RequireSuccess() {
	r.t.Helper()
	RequireSuccess(r.t, r.result)
}

func (r *APICallResponse) RequireBadRequestWithCode(errorCode string) {
	r.t.Helper()
	r.RequireBadRequest()
	r.RequireErrorCode(errorCode)
}

func (r *APICallResponse) RequireBadRequest() {
	r.t.Helper()
	r.RequireStatus(400)
}

func (r *APICallResponse) RequireUnauthorizedForAnonymous() {
	r.t.Helper()
	r.RequireUnauthorizedWithCode("error-unauthorized-auth-header-missing")
}

func (r *APICallResponse) RequireUnauthorizedForUnknownUser() {
	r.t.Helper()
	r.RequireUnauthorizedWithCode("error-unauthorized")
}

func (r *APICallResponse) RequireUnauthorizedForUserOnAdminAPI() {
	r.t.Helper()
	r.RequireUnauthorizedWithCode("error-unauthorized-xpub-not-an-admin-key")
}

func (r *APICallResponse) RequireUnauthorizedForAdminOnUserAPI() {
	r.t.Helper()
	r.RequireUnauthorizedWithCode("error-admin-auth-on-user-endpoint")
}

func (r *APICallResponse) RequireUnauthorizedWithCode(errorCode string) {
	r.t.Helper()
	r.RequireUnauthorized()
	r.RequireErrorCode(errorCode)
}

func (r *APICallResponse) RequireUnauthorized() {
	r.t.Helper()
	r.RequireStatus(401)
}

func (r *APICallResponse) RequireStatus(status int) {
	r.t.Helper()
	require.Equal(r.t, status, r.result.StatusCode(), "Unexpected status code")
}

func (r *APICallResponse) RequireErrorCode(code string) {
	r.t.Helper()
	body := r.ToMap()
	require.Equal(r.t, code, body["code"], "Unexpected error code")
}

func (r *APICallResponse) ToMap() map[string]any {
	r.t.Helper()
	body, _ := ExtractBody(r.result)

	response := make(map[string]any)
	err := json.Unmarshal([]byte(body), &response)
	require.NoError(r.t, err)

	return response
}
