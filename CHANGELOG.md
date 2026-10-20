# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.2.2, 2026-10-20

- Added the `GenerateSessionDescriptions` type (a `bool` alias) and the optional `GenerateSessionDescriptions` field on `Agent`, `CreateAgentRequest`, `UpdateAgentRequest` and `ReplaceAgentRequest`. When true, the platform makes one extra model request after each completed turn to write a description for sessions created without one. That request uses tokens and counts against `llm_quota`.
- Added the `JSON507` field and `GetJSON507()` method, returning `*Error`, on `CreateCorpusDocumentResponse`, `UpdateCorpusDocumentResponse`, `ReplaceCorpusDocumentMetadataResponse`, `ReplaceFilterAttributesResponse2` and `UploadFileResponse`, for HTTP 507 `application/json` responses.
- Updated generated doc comments for session `Description` fields to match the spec's new wording about when the platform generates descriptions.

## v0.2.1, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.2.0, 2026-10-17

- **Breaking:** Vectara's spec renamed the operation for `GET /v2/agent_aliases/{aliasKey}/connectors/{connectorId}` to `GetLegacyPublicConnector`, with `GetLegacyPublicConnectorWithResponse`, `NewGetLegacyPublicConnectorRequest`, `ParseGetLegacyPublicConnectorResponse`, `GetLegacyPublicConnectorParams` and `GetLegacyPublicConnectorResponse`. `GetPublicConnector` and its helpers keep their signatures but now call the new `GET /v2/agent_aliases/{aliasKey}/connectors/{connectorId}/public`. Code that calls `GetPublicConnector` still compiles; switch it to `GetLegacyPublicConnector` to keep calling the old path.
- Generated doc comments keep the API's own description text. The spec patch no longer rewrites dashes, filesystem terms or wording in descriptions.
- `make sec` uses Semgrep 1.181.0. Its pinned requirements now hash every published file, so the development image also builds on arm64.

## v0.1.0, 2026-10-17

- Added `vectara`, a REST client generated with oapi-codegen from Vectara API v2 (spec version 2.0, OpenAPI 3.0.0). It covers all 206 operations across 125 paths, including corpora, documents, search, generation, chats, agents, sessions, tools, pipelines and account administration.
- Number fields without a format are `float64` instead of oapi-codegen's default `float32`.
- `make generate` downloads the current spec from https://docs.vectara.com/vectara-oas-v2.yaml, removes an invalid property-level `required: false` flag from `HallucinationCorrectionRequest.query`, and generates the client in split files. The spec is not committed.
- Added HTTP contract tests, a compiling usage example, vendored dependencies, development tooling and agent integrations.
- Normalized generated documentation punctuation and wording while preserving wire schemas and examples.
- Replaced the spec's dummy Slack bot token example with a placeholder in generated comments.
- Assigned `JsonSchemaObject.additionalProperties` the Go field name `AllowsAdditionalProperties` to avoid a collision with the generated map of unknown JSON properties. The JSON key is unchanged.
