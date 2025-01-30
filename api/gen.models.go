// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"time"
)

// ApiComponentsErrorsErrUserNotFound defines model for api_components_errors_ErrUserNotFound.
type ApiComponentsErrorsErrUserNotFound struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
}

// ApiComponentsErrorsErrorSchema defines model for api_components_errors_ErrorSchema.
type ApiComponentsErrorsErrorSchema struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error message
	Message string `json:"message"`
}

// ApiComponentsModelsUser defines model for api_components_models_User.
type ApiComponentsModelsUser struct {
	// Id User ID
	Id uint64 `json:"id"`

	// Name User name
	Name string `json:"name"`
}

// ApiComponentsRequestsAdminRequest defines model for api_components_requests_AdminRequest.
type ApiComponentsRequestsAdminRequest struct {
	// Id Example of admin request body
	Id uint64 `json:"id"`
}

// ApiComponentsRequestsGetPaymails defines model for api_components_requests_GetPaymails.
type ApiComponentsRequestsGetPaymails struct {
	// Domain Paymail domain
	Domain *string `json:"domain,omitempty"`
}

// ApiComponentsResponsesCommonResponse Common response object
type ApiComponentsResponsesCommonResponse struct {
	Timestamp time.Time `json:"timestamp"`
}

// ApiComponentsResponsesUserExampleResponse defines model for api_components_responses_UserExampleResponse.
type ApiComponentsResponsesUserExampleResponse struct {
	// AdditionalPropertyExample The user model additional property example
	AdditionalPropertyExample *string `json:"additionalPropertyExample,omitempty"`

	// Id User ID
	Id uint64 `json:"id"`

	// Name User name
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

// ApiComponentsErrorsErrorUserNotFoundResponse defines model for api_components_errors_ErrorUserNotFoundResponse.
type ApiComponentsErrorsErrorUserNotFoundResponse = ApiComponentsErrorsErrUserNotFound

// GETAdminJSONRequestBody defines body for GETAdmin for application/json ContentType.
type GETAdminJSONRequestBody = ApiComponentsRequestsAdminRequest

// GetApiV1AdminPaymailsJSONRequestBody defines body for GetApiV1AdminPaymails for application/json ContentType.
type GetApiV1AdminPaymailsJSONRequestBody = ApiComponentsRequestsGetPaymails
