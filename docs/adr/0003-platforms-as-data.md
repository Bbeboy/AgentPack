# 0003 - Model platform support as data

## Context

Platform detection is currently encoded in Go literals. The target architecture requires platforms to become data so they can evolve without expanding conditional logic.

## Decision

Represent supported platforms in an embedded `platforms.json` file, with user overrides merged from the AgentPack config location in later phases.

## Consequences

- Platform definitions become easier to extend and review.
- The application layer can depend on structured platform data instead of direct condition chains.
- A manifest-compatible future is easier to support without redesigning detection again.
