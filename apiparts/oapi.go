package epicrm_apiparts

import (
	"log"

	"github.com/getkin/kin-openapi/openapi3"
)

func LogSpecErrorAndQuit(err error) {
	log.Fatalf("Error loading OpenAPI spec: %s\n", err)
}

func HandleSpec(spec *openapi3.T) {
	// To skip validation that checks if the server names match.
	// See https://github.com/deepmap/oapi-codegen/blob/e238df58329025ae89c0cd11171280171c3974db/examples/petstore-expanded/chi/petstore.go#L42
	spec.Servers = nil
}
