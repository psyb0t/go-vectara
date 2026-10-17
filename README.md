# go-vectara

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-vectara.svg)](https://pkg.go.dev/github.com/psyb0t/go-vectara)
[![CI](https://github.com/psyb0t/go-vectara/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-vectara/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-vectara/badges/version.svg)](https://github.com/psyb0t/go-vectara/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-vectara/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-vectara/badges/importers.svg)](https://github.com/psyb0t/go-vectara/blob/badges/importers.md)

Unofficial Go client for the [Vectara API v2](https://docs.vectara.com), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document Vectara publishes. I am not affiliated with Vectara.

The package is `vectara`, in the repository root. Everything in `client_*.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-vectara
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields. The full API, with the spec's descriptions as doc comments, is in the [Go reference](https://pkg.go.dev/github.com/psyb0t/go-vectara) and in `go doc github.com/psyb0t/go-vectara`.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-vectara"
)

func main() {
	auth := vectara.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("x-api-key", os.Getenv("VECTARA_API_KEY"))

		return nil
	})

	client, err := vectara.NewClientWithResponses("https://api.vectara.io", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.ListCorporaWithResponse(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles.

Method names come from the spec's operation IDs. The client covers corpora, documents, retrieval and RAG queries, chats, agents and sessions, tools, pipelines, account administration and observability. Streaming responses need raw client methods and caller-managed decoding and body closure.

## Auth

The base URL is `https://api.vectara.io`. Generated paths already include `/v2`.

Pass an API key in `x-api-key` with `WithRequestEditorFn`, as shown above. Keys need the roles required by each operation. Read credentials from the environment and never log authentication headers.

For OAuth 2.0, obtain an access token using the client credentials flow at `https://auth.vectara.io/oauth2/token`, then set `req.Header.Set("Authorization", "Bearer "+os.Getenv("VECTARA_ACCESS_TOKEN"))` in the request editor. The client does not obtain or refresh tokens automatically.

The spec also declares `VisitorToken` (`X-Visitor-Id`) and `FederatedSignIn` (`Authorization: Bearer <token>`) for the end-user operations that explicitly accept them.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `https://docs.vectara.com/vectara-oas-v2.yaml` to `api/openapi.yaml`.
- `scripts/patch-spec.sh` runs `go tool patch-spec` to write `api/openapi.fixed.yaml`, removing invalid `required: false` from `HallucinationCorrectionRequest.query`. OpenAPI 3.0 requires a list of required properties on the containing object instead.
- The patch names `JsonSchemaObject.additionalProperties` as `AllowsAdditionalProperties` in Go, avoiding a collision with the generator's unknown-fields map. The JSON key and unknown keyword handling stay intact.
- The patch replaces the spec's dummy Slack bot token example with `xoxb-EXAMPLE-BOT-TOKEN`, because GitHub push protection rejects commits that contain the original value.
- oapi-codegen runs with `config.yaml` over the full patched spec. The config maps unformatted `number` fields to `float64` and resolves generated type-name collisions.
- `scripts/split-gen.sh` splits the large output into `client_*.gen.go` files.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Agent integrations

The [skill](.agents/skills/go-vectara) teaches a coding agent how to use go-vectara when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-vectara@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-vectara@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-vectara
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Vectara's published API specification. The specification itself is not redistributed here.
