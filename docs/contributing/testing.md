# Testing

## Checks

The check command every change must pass before a PR:

```sh
make validate
```

It runs golangci-lint, `go build ./...`, the Go tests, a JSON validity check over `config/`, and the doc-rot guard. CI runs it on every pull request and every push to `main`. It takes about 10 seconds once golangci-lint is downloaded.

Smaller targets for use while working:

```sh
make test                                 # Go tests and JSON check; run before committing
go test ./internal/source -run TestFetch  # one test
make check-docs                           # doc-rot guard only
```

The doc-rot guard fails when `AGENTS.md`, `README.md` or `CONTRIBUTING.md` contains a hardcoded component count or a freshness stamp, because those go stale. Describe components without numbers.

## Writing tests

Tests sit next to the code as `*_test.go` files under `cmd/` and `internal/`. Write them table-driven, using the standard `testing` package; the repo has no assertion library. Tests that touch the filesystem work in `t.TempDir()`. Tests that need a git source create a real repository rather than faking one.
