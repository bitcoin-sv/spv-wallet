package api

import _ "embed"

// Yaml is the content of OpenAPI YAML file.
//
//go:embed gen.api.yaml
var Yaml string
