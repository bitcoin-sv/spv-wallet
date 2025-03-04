// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get contacts
	// (GET /api/v2/admin/contacts)
	AdminGetContacts(c *gin.Context, params AdminGetContactsParams)
	// Confirm contact
	// (POST /api/v2/admin/contacts/confirmations)
	AdminConfirmContact(c *gin.Context)
	// Delete contact
	// (DELETE /api/v2/admin/contacts/{id})
	AdminDeleteContact(c *gin.Context, id int)
	// Update contact
	// (PUT /api/v2/admin/contacts/{id})
	AdminUpdateContact(c *gin.Context, id int)
	// Create contact
	// (POST /api/v2/admin/contacts/{paymail})
	AdminCreateContact(c *gin.Context, paymail string)
	// Reject invitation
	// (DELETE /api/v2/admin/invitations/{id})
	AdminRejectInvitation(c *gin.Context, id int)
	// Accept invitation
	// (POST /api/v2/admin/invitations/{id})
	AdminAcceptInvitation(c *gin.Context, id int)
	// Get admin status
	// (GET /api/v2/admin/status)
	AdminStatus(c *gin.Context)
	// Create user
	// (POST /api/v2/admin/users)
	CreateUser(c *gin.Context)
	// Get user by id
	// (GET /api/v2/admin/users/{id})
	UserById(c *gin.Context, id string)
	// Add paymails to user
	// (POST /api/v2/admin/users/{id}/paymails)
	AddPaymailToUser(c *gin.Context, id string)
	// Get shared config
	// (GET /api/v2/configs/shared)
	SharedConfig(c *gin.Context)
	// Get contacts
	// (GET /api/v2/contacts)
	GetContacts(c *gin.Context, params GetContactsParams)
	// Remove contact
	// (DELETE /api/v2/contacts/{paymail})
	RemoveContact(c *gin.Context, paymail string)
	// Get contact
	// (GET /api/v2/contacts/{paymail})
	GetContact(c *gin.Context, paymail string)
	// Add contact
	// (PUT /api/v2/contacts/{paymail})
	UpsertContact(c *gin.Context, paymail string)
	// Unconfirm contact
	// (DELETE /api/v2/contacts/{paymail}/confirmation)
	UnconfirmContact(c *gin.Context, paymail string)
	// Confirm contact
	// (POST /api/v2/contacts/{paymail}/confirmation)
	ConfirmContact(c *gin.Context, paymail string)
	// Get data for user
	// (GET /api/v2/data/{id})
	DataById(c *gin.Context, id string)
	// Reject invitation
	// (DELETE /api/v2/invitations/{paymail})
	RejectInvitation(c *gin.Context, paymail string)
	// Accept invitation
	// (POST /api/v2/invitations/{paymail}/contacts)
	AcceptInvitation(c *gin.Context, paymail string)
	// Get operations for user
	// (GET /api/v2/operations/search)
	SearchOperations(c *gin.Context, params SearchOperationsParams)
	// Record transaction outline
	// (POST /api/v2/transactions)
	RecordTransactionOutline(c *gin.Context)
	// Create transaction outline
	// (POST /api/v2/transactions/outlines)
	CreateTransactionOutline(c *gin.Context, params CreateTransactionOutlineParams)
	// Get current user
	// (GET /api/v2/users/current)
	CurrentUser(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// AdminGetContacts operation middleware
func (siw *ServerInterfaceWrapper) AdminGetContacts(c *gin.Context) {

	var err error

	c.Set(XPubAuthScopes, []string{"admin"})

	// Parameter object where we will unmarshal all parameters from the context
	var params AdminGetContactsParams

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", c.Request.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter page: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "size" -------------

	err = runtime.BindQueryParameter("form", true, false, "size", c.Request.URL.Query(), &params.Size)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter size: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "sort" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort", c.Request.URL.Query(), &params.Sort)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sort: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "sortBy" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortBy", c.Request.URL.Query(), &params.SortBy)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sortBy: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "fullName" -------------

	err = runtime.BindQueryParameter("form", true, false, "fullName", c.Request.URL.Query(), &params.FullName)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter fullName: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "paymail" -------------

	err = runtime.BindQueryParameter("form", true, false, "paymail", c.Request.URL.Query(), &params.Paymail)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter paymail: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "id" -------------

	err = runtime.BindQueryParameter("form", true, false, "id", c.Request.URL.Query(), &params.Id)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "pubKey" -------------

	err = runtime.BindQueryParameter("form", true, false, "pubKey", c.Request.URL.Query(), &params.PubKey)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter pubKey: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "status" -------------

	err = runtime.BindQueryParameter("form", true, false, "status", c.Request.URL.Query(), &params.Status)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter status: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AdminGetContacts(c, params)
}

// AdminConfirmContact operation middleware
func (siw *ServerInterfaceWrapper) AdminConfirmContact(c *gin.Context) {

	c.Set(XPubAuthScopes, []string{"admin"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AdminConfirmContact(c)
}

// AdminDeleteContact operation middleware
func (siw *ServerInterfaceWrapper) AdminDeleteContact(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"admin"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AdminDeleteContact(c, id)
}

// AdminUpdateContact operation middleware
func (siw *ServerInterfaceWrapper) AdminUpdateContact(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"admin"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AdminUpdateContact(c, id)
}

// AdminCreateContact operation middleware
func (siw *ServerInterfaceWrapper) AdminCreateContact(c *gin.Context) {

	var err error

	// ------------- Path parameter "paymail" -------------
	var paymail string

	err = runtime.BindStyledParameterWithOptions("simple", "paymail", c.Param("paymail"), &paymail, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter paymail: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"admin"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AdminCreateContact(c, paymail)
}

// AdminRejectInvitation operation middleware
func (siw *ServerInterfaceWrapper) AdminRejectInvitation(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"admin"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AdminRejectInvitation(c, id)
}

// AdminAcceptInvitation operation middleware
func (siw *ServerInterfaceWrapper) AdminAcceptInvitation(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"admin"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AdminAcceptInvitation(c, id)
}

// AdminStatus operation middleware
func (siw *ServerInterfaceWrapper) AdminStatus(c *gin.Context) {

	c.Set(XPubAuthScopes, []string{"admin"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AdminStatus(c)
}

// CreateUser operation middleware
func (siw *ServerInterfaceWrapper) CreateUser(c *gin.Context) {

	c.Set(XPubAuthScopes, []string{"admin"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateUser(c)
}

// UserById operation middleware
func (siw *ServerInterfaceWrapper) UserById(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"admin"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.UserById(c, id)
}

// AddPaymailToUser operation middleware
func (siw *ServerInterfaceWrapper) AddPaymailToUser(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"admin"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AddPaymailToUser(c, id)
}

// SharedConfig operation middleware
func (siw *ServerInterfaceWrapper) SharedConfig(c *gin.Context) {

	c.Set(XPubAuthScopes, []string{"admin", "user"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.SharedConfig(c)
}

// GetContacts operation middleware
func (siw *ServerInterfaceWrapper) GetContacts(c *gin.Context) {

	var err error

	c.Set(XPubAuthScopes, []string{"user"})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetContactsParams

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", c.Request.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter page: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "size" -------------

	err = runtime.BindQueryParameter("form", true, false, "size", c.Request.URL.Query(), &params.Size)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter size: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "sort" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort", c.Request.URL.Query(), &params.Sort)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sort: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "sortBy" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortBy", c.Request.URL.Query(), &params.SortBy)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sortBy: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "fullName" -------------

	err = runtime.BindQueryParameter("form", true, false, "fullName", c.Request.URL.Query(), &params.FullName)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter fullName: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "paymail" -------------

	err = runtime.BindQueryParameter("form", true, false, "paymail", c.Request.URL.Query(), &params.Paymail)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter paymail: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "id" -------------

	err = runtime.BindQueryParameter("form", true, false, "id", c.Request.URL.Query(), &params.Id)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "pubKey" -------------

	err = runtime.BindQueryParameter("form", true, false, "pubKey", c.Request.URL.Query(), &params.PubKey)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter pubKey: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "status" -------------

	err = runtime.BindQueryParameter("form", true, false, "status", c.Request.URL.Query(), &params.Status)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter status: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetContacts(c, params)
}

// RemoveContact operation middleware
func (siw *ServerInterfaceWrapper) RemoveContact(c *gin.Context) {

	var err error

	// ------------- Path parameter "paymail" -------------
	var paymail string

	err = runtime.BindStyledParameterWithOptions("simple", "paymail", c.Param("paymail"), &paymail, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter paymail: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"user"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.RemoveContact(c, paymail)
}

// GetContact operation middleware
func (siw *ServerInterfaceWrapper) GetContact(c *gin.Context) {

	var err error

	// ------------- Path parameter "paymail" -------------
	var paymail string

	err = runtime.BindStyledParameterWithOptions("simple", "paymail", c.Param("paymail"), &paymail, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter paymail: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"user"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetContact(c, paymail)
}

// UpsertContact operation middleware
func (siw *ServerInterfaceWrapper) UpsertContact(c *gin.Context) {

	var err error

	// ------------- Path parameter "paymail" -------------
	var paymail string

	err = runtime.BindStyledParameterWithOptions("simple", "paymail", c.Param("paymail"), &paymail, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter paymail: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"user"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.UpsertContact(c, paymail)
}

// UnconfirmContact operation middleware
func (siw *ServerInterfaceWrapper) UnconfirmContact(c *gin.Context) {

	var err error

	// ------------- Path parameter "paymail" -------------
	var paymail string

	err = runtime.BindStyledParameterWithOptions("simple", "paymail", c.Param("paymail"), &paymail, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter paymail: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"user"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.UnconfirmContact(c, paymail)
}

// ConfirmContact operation middleware
func (siw *ServerInterfaceWrapper) ConfirmContact(c *gin.Context) {

	var err error

	// ------------- Path parameter "paymail" -------------
	var paymail string

	err = runtime.BindStyledParameterWithOptions("simple", "paymail", c.Param("paymail"), &paymail, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter paymail: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"user"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ConfirmContact(c, paymail)
}

// DataById operation middleware
func (siw *ServerInterfaceWrapper) DataById(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"user"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DataById(c, id)
}

// RejectInvitation operation middleware
func (siw *ServerInterfaceWrapper) RejectInvitation(c *gin.Context) {

	var err error

	// ------------- Path parameter "paymail" -------------
	var paymail string

	err = runtime.BindStyledParameterWithOptions("simple", "paymail", c.Param("paymail"), &paymail, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter paymail: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"user"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.RejectInvitation(c, paymail)
}

// AcceptInvitation operation middleware
func (siw *ServerInterfaceWrapper) AcceptInvitation(c *gin.Context) {

	var err error

	// ------------- Path parameter "paymail" -------------
	var paymail string

	err = runtime.BindStyledParameterWithOptions("simple", "paymail", c.Param("paymail"), &paymail, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter paymail: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(XPubAuthScopes, []string{"user"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AcceptInvitation(c, paymail)
}

// SearchOperations operation middleware
func (siw *ServerInterfaceWrapper) SearchOperations(c *gin.Context) {

	var err error

	c.Set(XPubAuthScopes, []string{"user"})

	// Parameter object where we will unmarshal all parameters from the context
	var params SearchOperationsParams

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", c.Request.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter page: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "size" -------------

	err = runtime.BindQueryParameter("form", true, false, "size", c.Request.URL.Query(), &params.Size)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter size: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "sort" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort", c.Request.URL.Query(), &params.Sort)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sort: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "sortBy" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortBy", c.Request.URL.Query(), &params.SortBy)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter sortBy: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.SearchOperations(c, params)
}

// RecordTransactionOutline operation middleware
func (siw *ServerInterfaceWrapper) RecordTransactionOutline(c *gin.Context) {

	c.Set(XPubAuthScopes, []string{"user"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.RecordTransactionOutline(c)
}

// CreateTransactionOutline operation middleware
func (siw *ServerInterfaceWrapper) CreateTransactionOutline(c *gin.Context) {

	var err error

	c.Set(XPubAuthScopes, []string{"user"})

	// Parameter object where we will unmarshal all parameters from the context
	var params CreateTransactionOutlineParams

	// ------------- Optional query parameter "format" -------------

	err = runtime.BindQueryParameter("form", true, false, "format", c.Request.URL.Query(), &params.Format)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter format: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateTransactionOutline(c, params)
}

// CurrentUser operation middleware
func (siw *ServerInterfaceWrapper) CurrentUser(c *gin.Context) {

	c.Set(XPubAuthScopes, []string{"user"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CurrentUser(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/api/v2/admin/contacts", wrapper.AdminGetContacts)
	router.POST(options.BaseURL+"/api/v2/admin/contacts/confirmations", wrapper.AdminConfirmContact)
	router.DELETE(options.BaseURL+"/api/v2/admin/contacts/:id", wrapper.AdminDeleteContact)
	router.PUT(options.BaseURL+"/api/v2/admin/contacts/:id", wrapper.AdminUpdateContact)
	router.POST(options.BaseURL+"/api/v2/admin/contacts/:paymail", wrapper.AdminCreateContact)
	router.DELETE(options.BaseURL+"/api/v2/admin/invitations/:id", wrapper.AdminRejectInvitation)
	router.POST(options.BaseURL+"/api/v2/admin/invitations/:id", wrapper.AdminAcceptInvitation)
	router.GET(options.BaseURL+"/api/v2/admin/status", wrapper.AdminStatus)
	router.POST(options.BaseURL+"/api/v2/admin/users", wrapper.CreateUser)
	router.GET(options.BaseURL+"/api/v2/admin/users/:id", wrapper.UserById)
	router.POST(options.BaseURL+"/api/v2/admin/users/:id/paymails", wrapper.AddPaymailToUser)
	router.GET(options.BaseURL+"/api/v2/configs/shared", wrapper.SharedConfig)
	router.GET(options.BaseURL+"/api/v2/contacts", wrapper.GetContacts)
	router.DELETE(options.BaseURL+"/api/v2/contacts/:paymail", wrapper.RemoveContact)
	router.GET(options.BaseURL+"/api/v2/contacts/:paymail", wrapper.GetContact)
	router.PUT(options.BaseURL+"/api/v2/contacts/:paymail", wrapper.UpsertContact)
	router.DELETE(options.BaseURL+"/api/v2/contacts/:paymail/confirmation", wrapper.UnconfirmContact)
	router.POST(options.BaseURL+"/api/v2/contacts/:paymail/confirmation", wrapper.ConfirmContact)
	router.GET(options.BaseURL+"/api/v2/data/:id", wrapper.DataById)
	router.DELETE(options.BaseURL+"/api/v2/invitations/:paymail", wrapper.RejectInvitation)
	router.POST(options.BaseURL+"/api/v2/invitations/:paymail/contacts", wrapper.AcceptInvitation)
	router.GET(options.BaseURL+"/api/v2/operations/search", wrapper.SearchOperations)
	router.POST(options.BaseURL+"/api/v2/transactions", wrapper.RecordTransactionOutline)
	router.POST(options.BaseURL+"/api/v2/transactions/outlines", wrapper.CreateTransactionOutline)
	router.GET(options.BaseURL+"/api/v2/users/current", wrapper.CurrentUser)
}
