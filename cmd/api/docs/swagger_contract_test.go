package docs

import (
	"encoding/json"
	"os"
	"testing"
)

type swaggerDocument struct {
	Definitions map[string]struct {
		Properties map[string]json.RawMessage `json:"properties"`
	} `json:"definitions"`
	Paths map[string]map[string]struct {
		Responses map[string]struct {
			Schema struct {
				Reference string `json:"$ref"`
			} `json:"schema"`
		} `json:"responses"`
	} `json:"paths"`
}

func TestPublicErrorResponsesUseMiddlewareErrorResponseContract(t *testing.T) {
	contents, err := os.ReadFile("swagger.json")
	if err != nil {
		t.Fatal(err)
	}

	var document swaggerDocument
	if err := json.Unmarshal(contents, &document); err != nil {
		t.Fatal(err)
	}

	const errorReference = "#/definitions/go-clean-arch_internal_shared_adapter_delivery_http_middleware.ErrorResponse"
	contract, ok := document.Definitions["go-clean-arch_internal_shared_adapter_delivery_http_middleware.ErrorResponse"]
	if !ok {
		t.Fatalf("missing %s definition", errorReference)
	}
	for _, property := range []string{"code", "message"} {
		if _, ok := contract.Properties[property]; !ok {
			t.Errorf("error contract missing %q property", property)
		}
	}

	for path, methods := range document.Paths {
		for method, operation := range methods {
			for status, response := range operation.Responses {
				if status[0] == '4' || status[0] == '5' {
					if response.Schema.Reference != errorReference {
						t.Errorf(
							"%s %s status %s references %q, want %q",
							method,
							path,
							status,
							response.Schema.Reference,
							errorReference,
						)
					}
				}
			}
		}
	}
}
