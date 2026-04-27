# 0002 - Keep Cobra as the CLI adapter

## Context

The project already uses Cobra and the migration plan explicitly keeps command parsing, help text, and routing in the CLI adapter.

## Decision

Continue using `github.com/spf13/cobra` as the CLI framework for AgentPack.

## Consequences

- Existing commands can be migrated incrementally instead of replaced wholesale.
- Command descriptions remain the source of truth for CLI behavior and help output.
- Future adapters such as MCP and TUI can reuse application use cases without replacing the CLI surface.
