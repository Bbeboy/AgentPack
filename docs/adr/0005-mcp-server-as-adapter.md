# 0005 - Treat the MCP server as an input adapter

## Context

The roadmap adds an MCP server after the CLI migration. That server must expose existing use cases instead of reimplementing business logic.

## Decision

Implement the MCP server as an input adapter that calls the same application layer use cases as the CLI.

## Consequences

- MCP support will not bypass validation and orchestration rules already defined in the app layer.
- CLI, MCP, and future TUI behavior can stay consistent by sharing use case outputs.
- Adapter-specific serialization stays outside the domain and application layers.
