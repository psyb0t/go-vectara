// Package vectara is a Go client for the Vectara API v2,
// generated with oapi-codegen from the OpenAPI document Vectara publishes.
package vectara

//go:generate bash scripts/fetch-spec.sh
//go:generate bash scripts/patch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.fixed.yaml
//go:generate bash scripts/split-gen.sh
