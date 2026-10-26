// Package api embeds the OpenAPI contract for serving at runtime.
package api

import _ "embed"

// OpenAPI is the embedded OpenAPI 3 specification.
//
//go:embed openapi.yaml
var OpenAPI []byte
