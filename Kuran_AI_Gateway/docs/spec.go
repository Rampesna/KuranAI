package docs

import (
	_ "embed"
)

//go:embed openapi.json
var openAPIJSON []byte

// OpenAPI3 returns the API OpenAPI 3.0 spec as JSON bytes.
func OpenAPI3() ([]byte, error) {
	return openAPIJSON, nil
}
