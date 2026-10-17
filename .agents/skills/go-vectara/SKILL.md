---
name: go-vectara
description: Unofficial Go client for the Vectara API v2. Package vectara (github.com/psyb0t/go-vectara) has typed methods for corpora, documents, search, RAG queries, chats, agents, tools, pipelines and account administration. Use when Go code needs Vectara or already imports go-vectara.
homepage: https://github.com/psyb0t/go-vectara
user-invocable: true
metadata:
  openclaw:
    emoji: "💹"
    requires:
      bins:
        - go
      env:
        - VECTARA_API_KEY
permissions:
  filesystem: "Guides edits to the Go source of the current project. The library reads and writes no files."
  shell: "The Go toolchain (`go get`, `go build`, `go test`) to add and compile the dependency."
  network: "`go get` fetches the module from your Go module proxy. At runtime the code you write talks HTTPS to api.vectara.io only."
---

# go-vectara

Unofficial, typed Go client for the Vectara API v2, in package `vectara` at the module root.

For installation, operations, pagination and response fields, see [references/setup.md](references/setup.md).

## Security & safety

- Read an API key from `VECTARA_API_KEY` and send it in `x-api-key`. Never hardcode, commit or log credentials or authentication headers.
- Use `https://api.vectara.io` as the base URL. Generated paths include `/v2`.
- OAuth uses a bearer access token. Visitor tokens and federated bearer credentials work only on operations declaring those schemes.
- Creating, updating, deleting or resetting resources changes account data. Run those operations only when the user requests them. Queries and agent execution may consume account resources or invoke configured tools.

## When to use

- Go code needs Vectara corpora, documents, retrieval or RAG queries.
- Go code needs Vectara chats, agents, tools or pipelines.
- The project imports `github.com/psyb0t/go-vectara`.

## When not to use

- The code is not Go. This module ships only a Go package.
- You need a client maintained by Vectara. This client is unofficial.
- You want to change the generated client. See `AGENTS.md` in the repo.

## Usage

Use the compiling example in [example_test.go](../../../example_test.go) and the matching README example: set `x-api-key` with a request editor, call `vectara.NewClientWithResponses("https://api.vectara.io", auth)`, then call `ListCorporaWithResponse(context.Background(), nil)`.

Check `err` for request construction, transport and decoding failures. Check `StatusCode()` and the expected success field such as `JSON200` for API failures. A non-2xx response can have a nil Go error. Optional fields may be pointers; check before dereferencing.

The client does not paginate, retry or refresh OAuth tokens automatically. For server-sent events, call the raw method, decode the stream yourself and close the response body. Parsed response methods read the full body.
