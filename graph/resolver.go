package graph

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/config"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the resolver
type Resolver struct{}

// GQLConfig GraphQL config
type GQLConfig struct {
	AppConfig *config.AppConfig
	Services  *config.AppServices
	Signed    bool
	XPub      string
	XPubID    string
	Auth      *bux.AuthPayload
	AuthError error
}

// GetConfigFromContext get the AppConfig, Services and rawXPubKey from the context
func GetConfigFromContext(ctx context.Context) (*GQLConfig, error) {
	ctxConfig := ctx.Value(config.GraphConfigKey).(*GQLConfig)
	if ctxConfig == nil {
		return nil, errors.New("could not find config in context")
	}
	if ctxConfig.AuthError != nil {
		return nil, ctxConfig.AuthError
	}

	return ctxConfig, nil
}

// GetConfigFromContextSigned get the AppConfig, Services and rawXPubKey from the context
// require that the request is properly signed
func GetConfigFromContextSigned(ctx context.Context) (*GQLConfig, error) {
	ctxConfig := ctx.Value(config.GraphConfigKey).(*GQLConfig)
	if ctxConfig == nil {
		return nil, errors.New("could not find config in context")
	}
	if ctxConfig.AuthError != nil {
		return nil, ctxConfig.AuthError
	}

	if !ctxConfig.Signed && !ctxConfig.AppConfig.Authentication.SigningDisabled {
		return nil, bux.ErrMissingSignature
	}

	return ctxConfig, nil
}

// GetConfigFromContextAdmin get the AppConfig, Services and rawXPubKey from the context + check for admin
func GetConfigFromContextAdmin(ctx context.Context) (*GQLConfig, error) {
	ctxConfig, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ctxConfig.AuthError != nil {
		return nil, ctxConfig.AuthError
	}

	// Check that we are the right key
	// todo we also need to check that the request is signed, otherwise the admin action
	// will not be checking the signature, since signatures are always on or always off
	if !ctxConfig.AppConfig.Authentication.IsAdmin(ctxConfig.XPub) {
		return nil, errors.New("invalid admin authentication")
	}

	return ctxConfig, nil
}

// ConditionsParseGraphQL parse the conditions passed from GraphQL
func ConditionsParseGraphQL(conditions map[string]interface{}) *map[string]interface{} {

	c, _ := json.Marshal(conditions) //nolint:errchkjson // ignore for now

	// string replace all keys "__...." -> "$..."
	m := regexp.MustCompile("\"__")
	cc := m.ReplaceAllString(string(c), "\"$")

	var parsedConditions *map[string]interface{}
	_ = json.Unmarshal([]byte(cc), &parsedConditions)

	return parsedConditions
}
