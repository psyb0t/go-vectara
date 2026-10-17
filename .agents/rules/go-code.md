# Go code

Refines `~/.claude/rule-details/go/wrapping-errors.md` and `~/.claude/rule-details/coding-style.md`.

- **GO1** `client_*.gen.go` files are generated. Never edit them by hand. Change `config.yaml` or the spec patch, then run `make generate`.
- **GO2** `scripts/patch-spec.sh` invokes the YAML tool in `internal/cmd/patch-spec` before generation. Each patch explains the spec problem it fixes.
- **GO3** The downloaded spec in `api/` is gitignored and never committed.
- **GO4** Unformatted number fields stay `float64` through `config.yaml`. Explicit `format: float` fields retain the spec's precision.
- **GO5** Hand-written Go (`gen.go` and the tests) wraps errors with `ctxerrors.Wrap` against `ctxerrors/commerr` sentinels. No `fmt.Errorf`, no `errors.New` outside a sentinel declaration.
- **GO6** The exported API is whatever oapi-codegen emits. A config change that renames or removes an exported name is a breaking release (see `releases.md`).
- **GO7** `make lint` passes with zero issues. Hand-written lines are at most 80 columns, the `.golangci.yml` limit.
