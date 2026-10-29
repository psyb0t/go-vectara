# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.2.5, 2026-10-29

- Vectara's spec now marks the per-LLM request and token limits as deprecated. The generated doc comments say the platform still enforces them and that Vectara support sets limits. The fields carry a Go `Deprecated:` marker: `RequestsPerSecond` on `LLM`, `InputTokensPerDay`, `InputTokensPerMinute` and `InputTokensPerMonth` on `LLMQuota` and the Create and Update LLM request types, and the matching output token and request rate fields. Existing code still compiles, but linters such as staticcheck now flag uses of these fields.
- Updated the `LlmQuota` doc comments on `Agent`, `CreateAgentRequest`, `UpdateAgentRequest` and `ReplaceAgentRequest` to say an omitted field sets no limit of its own and that the platform may apply its own limits.
- Moved the `UpdateWidgetConnectorConfiguration` type to another generated file. Its fields are unchanged.

## v0.2.4, 2026-10-28

- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.2.3, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.
- Raised the minimum Go version in `go.mod` from 1.26.9 to 1.26.10, which fixes the `crypto/x509` vulnerability GO-2026-6809 found by govulncheck.

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
