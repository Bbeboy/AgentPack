# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project aims to follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Native multi-OS CI test matrix for Linux amd64/arm64, macOS amd64/arm64, and Windows amd64.
- Smoke test job that builds the CLI and verifies `version`, `list`, and `--help` on every native CI runner.
- `golangci-lint` baseline configuration, `CONTRIBUTING.md`, and initial ADR documents under `docs/adr/`.
- Architecture document copy under `docs/ARQUITECTURA.md`.

### Changed
- CI gate now aggregates native tests, smoke tests, lint, self-update validation, and cross-build checks.
- README and README.es now document the CI matrix, smoke coverage, branch protection, and architecture docs path.
