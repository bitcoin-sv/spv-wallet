// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"encoding/json"
	"time"

	"github.com/oapi-codegen/runtime"
)

const (
	XPubAuthScopes = "XPubAuth.Scopes"
)

// ApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint defines model for api_components_errors_ErrAdminAuthOnNonAdminEndpoint.
type ApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrAuthorization defines model for api_components_errors_ErrAuthorization.
type ApiComponentsErrorsErrAuthorization struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrCannotBindRequest defines model for api_components_errors_ErrCannotBindRequest.
type ApiComponentsErrorsErrCannotBindRequest struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrCreatingUser defines model for api_components_errors_ErrCreatingUser.
type ApiComponentsErrorsErrCreatingUser struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrGettingUser defines model for api_components_errors_ErrGettingUser.
type ApiComponentsErrorsErrGettingUser struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrInvalidDomain defines model for api_components_errors_ErrInvalidDomain.
type ApiComponentsErrorsErrInvalidDomain struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrInvalidPaymail defines model for api_components_errors_ErrInvalidPaymail.
type ApiComponentsErrorsErrInvalidPaymail struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrInvalidPubKey defines model for api_components_errors_ErrInvalidPubKey.
type ApiComponentsErrorsErrInvalidPubKey struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrPaymailInconsistent defines model for api_components_errors_ErrPaymailInconsistent.
type ApiComponentsErrorsErrPaymailInconsistent struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrUnauthorized defines model for api_components_errors_ErrUnauthorized.
type ApiComponentsErrorsErrUnauthorized struct {
	union json.RawMessage
}

// ApiComponentsErrorsErrUserAuthOnNonUserEndpoint defines model for api_components_errors_ErrUserAuthOnNonUserEndpoint.
type ApiComponentsErrorsErrUserAuthOnNonUserEndpoint struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrWrongAuthScopeFormat defines model for api_components_errors_ErrWrongAuthScopeFormat.
type ApiComponentsErrorsErrWrongAuthScopeFormat struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrorSchema defines model for api_components_errors_ErrorSchema.
type ApiComponentsErrorsErrorSchema struct {
	// Code Error code
	Code string `json:"code"`

	// Message Error message
	Message string `json:"message"`
}

// ApiComponentsModelsPaymail defines model for api_components_models_Paymail.
type ApiComponentsModelsPaymail struct {
	Alias      string  `json:"alias"`
	Avatar     string  `json:"avatar"`
	Domain     string  `json:"domain"`
	Id         float32 `json:"id"`
	Paymail    string  `json:"paymail"`
	PublicName string  `json:"publicName"`
}

// ApiComponentsModelsSharedConfig Shared config
type ApiComponentsModelsSharedConfig struct {
	ExperimentalFeatures *map[string]bool `json:"experimentalFeatures,omitempty"`
	PaymailDomains       *[]string        `json:"paymailDomains,omitempty"`
}

// ApiComponentsModelsUser defines model for api_components_models_User.
type ApiComponentsModelsUser struct {
	CreatedAt time.Time                    `json:"createdAt"`
	Id        string                       `json:"id"`
	Paymails  []ApiComponentsModelsPaymail `json:"paymails"`
	PublicKey string                       `json:"publicKey"`
	UpdatedAt time.Time                    `json:"updatedAt"`
}

// ApiComponentsRequestsAddPaymail defines model for api_components_requests_AddPaymail.
type ApiComponentsRequestsAddPaymail struct {
	Address    string `json:"address"`
	Alias      string `json:"alias"`
	AvatarURL  string `json:"avatarURL"`
	Domain     string `json:"domain"`
	PublicName string `json:"publicName"`
}

// ApiComponentsRequestsCreateUser defines model for api_components_requests_CreateUser.
type ApiComponentsRequestsCreateUser struct {
	Paymail   *ApiComponentsRequestsAddPaymail `json:"paymail,omitempty"`
	PublicKey string                           `json:"publicKey"`
}

// ApiComponentsResponsesAdminAddPaymailSuccess defines model for api_components_responses_AdminAddPaymailSuccess.
type ApiComponentsResponsesAdminAddPaymailSuccess = ApiComponentsModelsPaymail

// ApiComponentsResponsesAdminCreateUserInternalServerError defines model for api_components_responses_AdminCreateUserInternalServerError.
type ApiComponentsResponsesAdminCreateUserInternalServerError = ApiComponentsErrorsErrCreatingUser

// ApiComponentsResponsesAdminCreateUserSuccess defines model for api_components_responses_AdminCreateUserSuccess.
type ApiComponentsResponsesAdminCreateUserSuccess = ApiComponentsModelsUser

// ApiComponentsResponsesAdminGetUser defines model for api_components_responses_AdminGetUser.
type ApiComponentsResponsesAdminGetUser = ApiComponentsModelsUser

// ApiComponentsResponsesAdminGetUserInternalServerError defines model for api_components_responses_AdminGetUserInternalServerError.
type ApiComponentsResponsesAdminGetUserInternalServerError = ApiComponentsErrorsErrGettingUser

// ApiComponentsResponsesAdminUserBadRequest defines model for api_components_responses_AdminUserBadRequest.
type ApiComponentsResponsesAdminUserBadRequest struct {
	union json.RawMessage
}

// ApiComponentsResponsesNotAuthorized defines model for api_components_responses_NotAuthorized.
type ApiComponentsResponsesNotAuthorized = ApiComponentsErrorsErrUnauthorized

// ApiComponentsResponsesSharedConfig Shared config
type ApiComponentsResponsesSharedConfig = ApiComponentsModelsSharedConfig

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = ApiComponentsRequestsCreateUser

// AddPaymailToUserJSONRequestBody defines body for AddPaymailToUser for application/json ContentType.
type AddPaymailToUserJSONRequestBody = ApiComponentsRequestsAddPaymail

// AsApiComponentsErrorsErrAuthorization returns the union data inside the ApiComponentsErrorsErrUnauthorized as a ApiComponentsErrorsErrAuthorization
func (t ApiComponentsErrorsErrUnauthorized) AsApiComponentsErrorsErrAuthorization() (ApiComponentsErrorsErrAuthorization, error) {
	var body ApiComponentsErrorsErrAuthorization
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromApiComponentsErrorsErrAuthorization overwrites any union data inside the ApiComponentsErrorsErrUnauthorized as the provided ApiComponentsErrorsErrAuthorization
func (t *ApiComponentsErrorsErrUnauthorized) FromApiComponentsErrorsErrAuthorization(v ApiComponentsErrorsErrAuthorization) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeApiComponentsErrorsErrAuthorization performs a merge with any union data inside the ApiComponentsErrorsErrUnauthorized, using the provided ApiComponentsErrorsErrAuthorization
func (t *ApiComponentsErrorsErrUnauthorized) MergeApiComponentsErrorsErrAuthorization(v ApiComponentsErrorsErrAuthorization) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsApiComponentsErrorsErrWrongAuthScopeFormat returns the union data inside the ApiComponentsErrorsErrUnauthorized as a ApiComponentsErrorsErrWrongAuthScopeFormat
func (t ApiComponentsErrorsErrUnauthorized) AsApiComponentsErrorsErrWrongAuthScopeFormat() (ApiComponentsErrorsErrWrongAuthScopeFormat, error) {
	var body ApiComponentsErrorsErrWrongAuthScopeFormat
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromApiComponentsErrorsErrWrongAuthScopeFormat overwrites any union data inside the ApiComponentsErrorsErrUnauthorized as the provided ApiComponentsErrorsErrWrongAuthScopeFormat
func (t *ApiComponentsErrorsErrUnauthorized) FromApiComponentsErrorsErrWrongAuthScopeFormat(v ApiComponentsErrorsErrWrongAuthScopeFormat) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeApiComponentsErrorsErrWrongAuthScopeFormat performs a merge with any union data inside the ApiComponentsErrorsErrUnauthorized, using the provided ApiComponentsErrorsErrWrongAuthScopeFormat
func (t *ApiComponentsErrorsErrUnauthorized) MergeApiComponentsErrorsErrWrongAuthScopeFormat(v ApiComponentsErrorsErrWrongAuthScopeFormat) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint returns the union data inside the ApiComponentsErrorsErrUnauthorized as a ApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint
func (t ApiComponentsErrorsErrUnauthorized) AsApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint() (ApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint, error) {
	var body ApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint overwrites any union data inside the ApiComponentsErrorsErrUnauthorized as the provided ApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint
func (t *ApiComponentsErrorsErrUnauthorized) FromApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint(v ApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint performs a merge with any union data inside the ApiComponentsErrorsErrUnauthorized, using the provided ApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint
func (t *ApiComponentsErrorsErrUnauthorized) MergeApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint(v ApiComponentsErrorsErrAdminAuthOnNonAdminEndpoint) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsApiComponentsErrorsErrUserAuthOnNonUserEndpoint returns the union data inside the ApiComponentsErrorsErrUnauthorized as a ApiComponentsErrorsErrUserAuthOnNonUserEndpoint
func (t ApiComponentsErrorsErrUnauthorized) AsApiComponentsErrorsErrUserAuthOnNonUserEndpoint() (ApiComponentsErrorsErrUserAuthOnNonUserEndpoint, error) {
	var body ApiComponentsErrorsErrUserAuthOnNonUserEndpoint
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromApiComponentsErrorsErrUserAuthOnNonUserEndpoint overwrites any union data inside the ApiComponentsErrorsErrUnauthorized as the provided ApiComponentsErrorsErrUserAuthOnNonUserEndpoint
func (t *ApiComponentsErrorsErrUnauthorized) FromApiComponentsErrorsErrUserAuthOnNonUserEndpoint(v ApiComponentsErrorsErrUserAuthOnNonUserEndpoint) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeApiComponentsErrorsErrUserAuthOnNonUserEndpoint performs a merge with any union data inside the ApiComponentsErrorsErrUnauthorized, using the provided ApiComponentsErrorsErrUserAuthOnNonUserEndpoint
func (t *ApiComponentsErrorsErrUnauthorized) MergeApiComponentsErrorsErrUserAuthOnNonUserEndpoint(v ApiComponentsErrorsErrUserAuthOnNonUserEndpoint) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t ApiComponentsErrorsErrUnauthorized) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *ApiComponentsErrorsErrUnauthorized) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsApiComponentsErrorsErrCannotBindRequest returns the union data inside the ApiComponentsResponsesAdminUserBadRequest as a ApiComponentsErrorsErrCannotBindRequest
func (t ApiComponentsResponsesAdminUserBadRequest) AsApiComponentsErrorsErrCannotBindRequest() (ApiComponentsErrorsErrCannotBindRequest, error) {
	var body ApiComponentsErrorsErrCannotBindRequest
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromApiComponentsErrorsErrCannotBindRequest overwrites any union data inside the ApiComponentsResponsesAdminUserBadRequest as the provided ApiComponentsErrorsErrCannotBindRequest
func (t *ApiComponentsResponsesAdminUserBadRequest) FromApiComponentsErrorsErrCannotBindRequest(v ApiComponentsErrorsErrCannotBindRequest) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeApiComponentsErrorsErrCannotBindRequest performs a merge with any union data inside the ApiComponentsResponsesAdminUserBadRequest, using the provided ApiComponentsErrorsErrCannotBindRequest
func (t *ApiComponentsResponsesAdminUserBadRequest) MergeApiComponentsErrorsErrCannotBindRequest(v ApiComponentsErrorsErrCannotBindRequest) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsApiComponentsErrorsErrInvalidPubKey returns the union data inside the ApiComponentsResponsesAdminUserBadRequest as a ApiComponentsErrorsErrInvalidPubKey
func (t ApiComponentsResponsesAdminUserBadRequest) AsApiComponentsErrorsErrInvalidPubKey() (ApiComponentsErrorsErrInvalidPubKey, error) {
	var body ApiComponentsErrorsErrInvalidPubKey
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromApiComponentsErrorsErrInvalidPubKey overwrites any union data inside the ApiComponentsResponsesAdminUserBadRequest as the provided ApiComponentsErrorsErrInvalidPubKey
func (t *ApiComponentsResponsesAdminUserBadRequest) FromApiComponentsErrorsErrInvalidPubKey(v ApiComponentsErrorsErrInvalidPubKey) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeApiComponentsErrorsErrInvalidPubKey performs a merge with any union data inside the ApiComponentsResponsesAdminUserBadRequest, using the provided ApiComponentsErrorsErrInvalidPubKey
func (t *ApiComponentsResponsesAdminUserBadRequest) MergeApiComponentsErrorsErrInvalidPubKey(v ApiComponentsErrorsErrInvalidPubKey) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsApiComponentsErrorsErrInvalidPaymail returns the union data inside the ApiComponentsResponsesAdminUserBadRequest as a ApiComponentsErrorsErrInvalidPaymail
func (t ApiComponentsResponsesAdminUserBadRequest) AsApiComponentsErrorsErrInvalidPaymail() (ApiComponentsErrorsErrInvalidPaymail, error) {
	var body ApiComponentsErrorsErrInvalidPaymail
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromApiComponentsErrorsErrInvalidPaymail overwrites any union data inside the ApiComponentsResponsesAdminUserBadRequest as the provided ApiComponentsErrorsErrInvalidPaymail
func (t *ApiComponentsResponsesAdminUserBadRequest) FromApiComponentsErrorsErrInvalidPaymail(v ApiComponentsErrorsErrInvalidPaymail) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeApiComponentsErrorsErrInvalidPaymail performs a merge with any union data inside the ApiComponentsResponsesAdminUserBadRequest, using the provided ApiComponentsErrorsErrInvalidPaymail
func (t *ApiComponentsResponsesAdminUserBadRequest) MergeApiComponentsErrorsErrInvalidPaymail(v ApiComponentsErrorsErrInvalidPaymail) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsApiComponentsErrorsErrPaymailInconsistent returns the union data inside the ApiComponentsResponsesAdminUserBadRequest as a ApiComponentsErrorsErrPaymailInconsistent
func (t ApiComponentsResponsesAdminUserBadRequest) AsApiComponentsErrorsErrPaymailInconsistent() (ApiComponentsErrorsErrPaymailInconsistent, error) {
	var body ApiComponentsErrorsErrPaymailInconsistent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromApiComponentsErrorsErrPaymailInconsistent overwrites any union data inside the ApiComponentsResponsesAdminUserBadRequest as the provided ApiComponentsErrorsErrPaymailInconsistent
func (t *ApiComponentsResponsesAdminUserBadRequest) FromApiComponentsErrorsErrPaymailInconsistent(v ApiComponentsErrorsErrPaymailInconsistent) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeApiComponentsErrorsErrPaymailInconsistent performs a merge with any union data inside the ApiComponentsResponsesAdminUserBadRequest, using the provided ApiComponentsErrorsErrPaymailInconsistent
func (t *ApiComponentsResponsesAdminUserBadRequest) MergeApiComponentsErrorsErrPaymailInconsistent(v ApiComponentsErrorsErrPaymailInconsistent) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsApiComponentsErrorsErrInvalidDomain returns the union data inside the ApiComponentsResponsesAdminUserBadRequest as a ApiComponentsErrorsErrInvalidDomain
func (t ApiComponentsResponsesAdminUserBadRequest) AsApiComponentsErrorsErrInvalidDomain() (ApiComponentsErrorsErrInvalidDomain, error) {
	var body ApiComponentsErrorsErrInvalidDomain
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromApiComponentsErrorsErrInvalidDomain overwrites any union data inside the ApiComponentsResponsesAdminUserBadRequest as the provided ApiComponentsErrorsErrInvalidDomain
func (t *ApiComponentsResponsesAdminUserBadRequest) FromApiComponentsErrorsErrInvalidDomain(v ApiComponentsErrorsErrInvalidDomain) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeApiComponentsErrorsErrInvalidDomain performs a merge with any union data inside the ApiComponentsResponsesAdminUserBadRequest, using the provided ApiComponentsErrorsErrInvalidDomain
func (t *ApiComponentsResponsesAdminUserBadRequest) MergeApiComponentsErrorsErrInvalidDomain(v ApiComponentsErrorsErrInvalidDomain) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t ApiComponentsResponsesAdminUserBadRequest) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *ApiComponentsResponsesAdminUserBadRequest) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
