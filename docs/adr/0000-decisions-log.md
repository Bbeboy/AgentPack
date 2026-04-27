# 0000 - Decisions Log

## Status

Accepted during Fase 0 preparation.

## Closed questions

| ID | Decision |
| --- | --- |
| Q1 | Ports in `app/ports.go` will be grouped by functional domain. |
| Q2 | Prompt adapter will use `huh`. |
| Q3 | Code comments will be written in Spanish. |
| Q4 | `platforms.json` will use an AgentPack-owned schema. |
| Q5 | Windows ARM64 is accepted as cross-build only for CI in this phase. |
| Q6 | Future command transition will introduce `uninstall` and `delete` with immediate break from `remove` and `remove-skill`. |
| Q7 | Minimum supported Go version remains `1.23`. |
| Q8 | The module path remains `github.com/Bbeboy/AgentPack`. |
| Q9 | Package manifest support is included in this migration, not deferred. |
| Q10 | GitHub URL installation stays out of scope for this migration and will be revisited later. |

## Notes

- The canonical architecture document now lives at `docs/ARQUITECTURA.md`.
- This log records product and architecture decisions; implementation lands in later phases.
