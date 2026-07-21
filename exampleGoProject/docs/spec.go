package docs

import (
	_ "embed"
)

//go:embed openapi.json
var openAPIJSON []byte

// OpenAPI3 returns the gateway API OpenAPI 3.0 spec as JSON bytes (embed; openapi.json değişince yeniden build gerekir).
func OpenAPI3() ([]byte, error) {
	return openAPIJSON, nil
}
