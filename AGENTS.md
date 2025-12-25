# Repository Guidelines

## Project Structure & Module Organization

- Root package code lives in `ft.go` and `global.go`.
- Tests are in `ft_test.go` and follow standard Go `_test.go` naming.
- Examples live in `examples/`.
- Tooling is isolated in `tools/` with its own `go.mod` for lint dependencies.
- Local binaries install into `bin/` via the Makefile.

## Build, Test, and Development Commands

- `go test ./...` runs the full unit test suite.
- `go test ./... -run TestName` runs a focused test by name.
- `make bin` prepares `bin/` and tidies tooling deps in `tools/`.
- `make lint` installs and runs GolangCI-Lint with repo settings.

## Coding Style & Naming Conventions

- Follow standard Go formatting (`gofmt`); tabs are used for indentation.
- Keep exported identifiers `PascalCase`; unexported are `camelCase`.
- Test functions must be `TestXxx` and benchmarks `BenchmarkXxx`.
- Lint rules are defined in `.golangci.yml` (staticcheck, revive, goimports, etc.).

## Testing Guidelines

- Testing uses Go’s `testing` package with `testify` assertions.
- Keep tests focused on observable behavior of the `ft` package.
- Run `go test ./...` before submitting changes; no explicit coverage threshold.

## Commit & Pull Request Guidelines

- Commit messages follow Conventional Commits (`feat:`, `fix:`, `docs:`, `test:`, `chore:`), sometimes with a scope (e.g., `fix(test): ...`).
- PRs should include a short summary, motivation, and any relevant issue links.
- Update documentation or examples when behavior changes.

## Release Notes

- Changelog updates use `git-cliff` per `README.md` (see release steps there).
