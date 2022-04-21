package config

// graphContextKey for context key
type graphContextKey string

var (
	// GraphConfigKey is the ctx key for the
	GraphConfigKey graphContextKey = "graphql_config"

	// GraphRequestInfo is the ctx key for the request info passed down to graphql for logging
	GraphRequestInfo graphContextKey = "request_info"
)
