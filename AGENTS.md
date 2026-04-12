# AGENTS.md

## Purpose

This repository contains a Go CLI for running internet speed tests against two providers:

- `speedtest.net`, implemented under `speedtestdotnet/`
- `fast.com`, implemented under `fastdotcom/`

The shipped binary entrypoint is [`cmd/cli/main.go`](/Users/fopina/Documents/speedtest-cli/cmd/cli/main.go).

This guide is for coding agents making changes in this repo. Keep changes small, package-local, and aligned with the current separation between protocol logic, CLI wiring, and shared probing utilities.

## Repository Map

- [`cmd/cli/main.go`](/Users/fopina/Documents/speedtest-cli/cmd/cli/main.go): Cobra root command and subcommand registration.
- [`cmd/cli/internal/speedtestdotnet/`](/Users/fopina/Documents/speedtest-cli/cmd/cli/internal/speedtestdotnet): CLI flags, server selection, progress output, and JSON/plain-text orchestration for the `st` command.
- [`cmd/cli/internal/fastdotcom/`](/Users/fopina/Documents/speedtest-cli/cmd/cli/internal/fastdotcom): CLI flags, progress output, and JSON/plain-text orchestration for the `f` command.
- [`cmd/cli/internal/types.go`](/Users/fopina/Documents/speedtest-cli/cmd/cli/internal/types.go): shared JSON result struct and speed formatting helpers used by both providers.
- [`speedtestdotnet/`](/Users/fopina/Documents/speedtest-cli/speedtestdotnet): speedtest.net protocol/client/config/server/probing logic.
- [`fastdotcom/`](/Users/fopina/Documents/speedtest-cli/fastdotcom): fast.com client, manifest loading, and probing logic.
- [`fastdotcom/internal/`](/Users/fopina/Documents/speedtest-cli/fastdotcom/internal): fast.com token + manifest internals; keep provider-specific wire details here.
- [`prober/`](/Users/fopina/Documents/speedtest-cli/prober): shared concurrent transfer orchestration.
- [`prober/proberutil/`](/Users/fopina/Documents/speedtest-cli/prober/proberutil): aggregation helpers that convert transfer totals into speeds.
- [`oututil/`](/Users/fopina/Documents/speedtest-cli/oututil): TTY-aware live progress printer.
- [`units/`](/Users/fopina/Documents/speedtest-cli/units): speed units and human-readable formatting.
- [`geo/`](/Users/fopina/Documents/speedtest-cli/geo): coordinate and distance helpers used by server selection.
- [`version/`](/Users/fopina/Documents/speedtest-cli/version/version.go): build-time version variable.

## Architecture Rules

1. Put CLI concerns in `cmd/cli/internal/...`.
   That includes Cobra flags, help text, provider command behavior, JSON vs plain-text output, and user-facing logging.

2. Put network/protocol behavior in the provider packages.
   `speedtestdotnet/` and `fastdotcom/` should own request construction, response parsing, and probe mechanics.

3. Reuse `prober.Group` and `proberutil.SpeedCollect` for concurrent transfer measurements.
   Both providers already follow this pattern for download/upload probes; keep new throughput probes consistent with it.

4. Keep shared formatting generic.
   If logic applies to both providers, prefer `cmd/cli/internal/types.go`, `units/`, or `oututil/` over duplicating behavior.

5. Preserve the current package boundary style.
   The CLI packages coordinate tests; provider packages expose methods like `ProbeDownloadSpeed`, `ProbeUploadSpeed`, `Config`, `LoadAllServers`, and `GetManifest`.

## Existing Conventions

- Context timeouts are controlled by CLI flags and passed into provider methods with `context.WithTimeout`.
- Most fatal user-facing failures currently use `log.Fatalf` in CLI-layer orchestration.
- Progress printing is optional and TTY-sensitive through [`oututil/printer.go`](/Users/fopina/Documents/speedtest-cli/oututil/printer.go).
- Speed formatting defaults to bits for display, with `--bytes` switching to bytes.
- JSON output for both providers flows through `internal.Result`, even though provider-specific metadata differs.
- Tests are package-local `_test.go` files placed next to the implementation they cover.

## Change Guidance

### Adding or changing CLI flags

- Update the relevant `InitFlags` function under `cmd/cli/internal/.../flags.go`.
- Wire the new behavior in that provider's `main.go` or `probe.go`.
- If help output changes, update [`cmd/cli/main_test.go`](/Users/fopina/Documents/speedtest-cli/cmd/cli/main_test.go) when the command help text is asserted there.

### Adding provider behavior

- `speedtest.net` changes usually belong in `speedtestdotnet/` plus its CLI wrapper under `cmd/cli/internal/speedtestdotnet/`.
- `fast.com` changes usually belong in `fastdotcom/` plus its CLI wrapper under `cmd/cli/internal/fastdotcom/`.
- Avoid leaking provider-specific parsing or request details into the top-level Cobra command.

### Adjusting output

- For human-readable speed strings, prefer `units.BytesPerSecond` / `BitsPerSecond`.
- For live updating terminal output, use `oututil.StartPrinting()` rather than custom terminal control logic.
- For JSON output shared by both commands, extend `cmd/cli/internal/types.go` only when the shape should be common.

### Working on concurrency/probing

- Keep probe functions cancellation-aware.
- Preserve the `stream chan<- units.BytesPerSecond` pattern for incremental output.
- Be careful with closures inside loops that schedule probe functions; bind per-iteration values explicitly when needed.

## Testing

Primary local verification:

- `go test ./...`
- `go test -race ./...`

Repo-provided commands:

- `make build`
- `make test`

Notes:

- `make test` deletes `coverage.out` and `dist/`.
- The current worktree already has a user change in `Makefile`; do not overwrite or revert it unless explicitly asked.
- CI currently runs `go mod tidy` and `go test -race -coverprofile="coverage.out" -covermode=atomic ./...` in [`.github/workflows/test.yml`](/Users/fopina/Documents/speedtest-cli/.github/workflows/test.yml).

## Safe Defaults For Agents

- Read package-local tests before changing behavior in that package.
- Prefer minimal edits over broad refactors; the repo is small and intentionally direct.
- Do not introduce new dependencies unless the gain is clear.
- Preserve Cobra command names and aliases unless the task explicitly asks for a CLI breaking change.
- Keep README examples and command help in sync with real CLI behavior when changing flags or output.
- Avoid touching release workflow or Docker packaging unless the task is specifically about build/distribution.

## Watchouts

- The Makefile currently injects `-X main.version=...`, but the version variable actually lives in `version.Version`. If you are asked to fix release/build version stamping, inspect that path carefully.
- `speedtestdotnet` and `fastdotcom` have similar CLI probe/output code, but their protocol layers differ; do not force them into a shared abstraction unless there is a clear payoff.
- Network-facing tests should remain deterministic. Prefer `httptest` or existing test patterns over live external calls.

## When Unsure

- Follow the nearest existing pattern in the same package.
- Keep provider logic in provider packages and CLI logic in `cmd/cli/internal/...`.
- Add or update tests in the same area as the behavior change.
