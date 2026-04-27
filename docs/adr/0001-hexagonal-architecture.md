# 0001 - Adopt hexagonal architecture

## Context

AgentPack started as a small CLI with direct package-to-filesystem coupling. The migration plan requires clearer boundaries before new surfaces like MCP and TUI are introduced.

## Decision

Adopt a hexagonal architecture with four layers:

- `internal/domain`
- `internal/app`
- `internal/adapter/*`
- `cmd/agentpack`

Dependencies must point inward and the composition root remains in `cmd/agentpack/main.go`.

## Consequences

- New behavior will be implemented behind ports and adapters instead of adding more direct filesystem logic to commands.
- Dependency rules will be enforced incrementally with `depguard`.
- The migration will happen in phases to avoid mixing structural change with feature deletion.
