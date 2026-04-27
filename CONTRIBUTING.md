# Contributing

Thanks for helping improve AgentPack.

## Requirements

- Go `1.23.x`
- `golangci-lint` available locally for lint checks

## Local checks

Run the standard development checks before opening a pull request:

```bash
go test ./...
golangci-lint run
```

If you are changing command behavior or build wiring, also verify the CLI still builds:

```bash
go build -o agentpack ./cmd/agentpack
./agentpack --help
```

On Windows, replace the final command with:

```powershell
go build -o agentpack.exe ./cmd/agentpack
.\agentpack.exe --help
```

## Pull request rules

- Keep one feature per pull request.
- Keep CI green before merge.
- Add an entry to `CHANGELOG.md` under `## [Unreleased]`.
- Do not delete old behavior in the same PR that introduces its replacement.

## CI overview

The required `ci-gate` check aggregates:

- native `go test ./...` runs across Linux, macOS, and Windows
- smoke binary checks on the same native matrix
- `golangci-lint`
- self-update validation
- cross-build verification for Linux/macOS/Windows on amd64 and arm64
