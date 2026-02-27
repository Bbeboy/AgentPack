# AgentPack

Go CLI to create, store, and install reusable agent skill packages.

Spanish documentation is available at `README.es.md`.

## Highlights

- Create packages from an existing skills directory.
- Auto-detect source folder when using `.`.
- Install into detected platform skills directory (for example `.opencode/skills`, `.agents/skills`, `.cla/skills`).
- Handle skill conflicts interactively during install.
- Add paths into an existing package with `add`.
- Remove a package path with `remove <path> --from <package>`.
- Rename packages with `rename`.
- CLI output and help in English or Spanish (`config set language`, `lang`).
- CI runs `go test ./...` on push and pull requests.

## Requirements

- Go `1.23+`
- Linux (current support)

## Install

Recommended:

```bash
go install github.com/Bbeboy/AgentPack/cmd/agentpack@latest
agentpack --help
```

If `agentpack` is resolved from `~/.local/bin` in your `PATH`, install directly there to avoid stale binaries:

```bash
GOBIN="$HOME/.local/bin" go install github.com/Bbeboy/AgentPack/cmd/agentpack@latest
agentpack -v
```

To verify which binary is being executed:

```bash
which -a agentpack
agentpack -v
```

From source:

```bash
git clone https://github.com/Bbeboy/AgentPack.git
cd AgentPack
go mod tidy
go build -o agentpack ./cmd/agentpack
install -m 755 ./agentpack ~/.local/bin/agentpack
agentpack --help
```

## Quick Start

```bash
# 1) create package from current project skills folder
agentpack create backend-base .

# 2) list saved packages
agentpack list

# 3) install package into detected platform skills destination
agentpack install backend-base

# 4) add a new path into package
agentpack add ./skills/docker --to backend-base

# 5) remove path from package
agentpack remove docker/SKILL.md --from backend-base

# 6) rename package
agentpack rename backend-base backend-v2
```

## Key Concepts

### Package storage location

Packages are stored locally at:

```text
~/.agentpack/packages-skills/<package-name>
```

### Install destination

`install` resolves a platform root in the current project, then installs into `<platform-root>/skills`.

Common examples:

- `.opencode/skills`
- `.agents/skills`
- `.cla/skills`
- `.cursor/skills`

Fallback when no platform is detected:

```text
.agents/skills
```

GitHub Copilot detection uses `.github/skills` (not just `.github`) to avoid false positives in repositories that only use GitHub workflows.

### Package name rules

- Max 64 chars.
- Must start with a letter or number.
- Allowed chars: letters, numbers, `.`, `_`, `-`.

## Commands

| Command | Description | Example |
| --- | --- | --- |
| `agentpack create <package-name> <skills-path>` | Create package from a skills path. | `agentpack create backend-base .` |
| `agentpack install <package-name>` | Install package into detected platform skills destination. | `agentpack install backend-base` |
| `agentpack add <file-or-folder> --to <package-name>` | Add a file/folder to an existing package. | `agentpack add ./skills/docker --to backend-base` |
| `agentpack list` | List saved packages. | `agentpack list` |
| `agentpack list-skills <package-name>` | List skills inside a package. | `agentpack list-skills backend-base` |
| `agentpack rename <current-name> <new-name>` | Rename an existing package. | `agentpack rename backend-base backend-v2` |
| `agentpack remove <package-name>` | Remove an entire package (with confirmation). | `agentpack remove backend-base` |
| `agentpack remove <path> --from <package-name>` | Remove a specific package path. | `agentpack remove docker/SKILL.md --from backend-base` |
| `agentpack remove-skill <package-name> <skill-name>` | Remove one skill folder from package. | `agentpack remove-skill backend-base docker` |
| `agentpack config set language <en\|es>` | Set global CLI language. | `agentpack config set language es` |
| `agentpack lang <en\|es>` | Language shortcut command. | `agentpack lang en` |
| `agentpack -v` | Show installed CLI version (short flag). | `agentpack -v` |
| `agentpack version` | Show installed CLI version. | `agentpack version` |
| `agentpack completion [bash\|zsh\|fish\|powershell]` | Generate shell completion script. | `agentpack completion fish` |

## Language

Commands stay in English. Help text and runtime feedback are localized.

```bash
agentpack config set language es
agentpack --help
agentpack lang en
```

Default language is `en`.

## Notes on `create .`

When you run:

```bash
agentpack create <package-name> .
```

AgentPack searches known platform skills paths in priority order, plus `skills` fallback. If one path is found, it is selected automatically. If multiple paths are found, you can choose interactively.

## Development

Run in development mode:

```bash
go run ./cmd/agentpack --help
```

Format, test, build:

```bash
go fmt ./...
go test ./...
go build -o agentpack ./cmd/agentpack
```

## Testing and CI

- Test files are colocated with package code (`*_test.go`).
- Shared test helpers live in `internal/testutil`.
- CI workflow: `.github/workflows/test.yml`.
- CI gate runs:

```bash
go test ./...
```

## Project Structure

```text
cmd/
  agentpack/
    main.go
internal/
  cli/
    *.go
    *_test.go
  config/
    settings.go
    settings_test.go
  fsutil/
    copy.go
  i18n/
    messages.go
  platform/
    skills.go
    skills_test.go
  prompt/
    prompt.go
  storage/
    storage.go
    storage_test.go
  testutil/
    fs.go
```

## Troubleshooting

### `agentpack: command not found`

- Confirm binary exists in `~/.local/bin/agentpack` or your `GOBIN`.
- Confirm that directory is in `PATH`.

### Unexpected old output after reinstall

If output still shows old messages, you likely have multiple binaries installed:

```bash
which -a agentpack
GOBIN="$HOME/.local/bin" go install ./cmd/agentpack
agentpack -v
```

## Binary Versioning Plan (No Go Required)

To make installation direct (without `go install`), we can implement this phased release pipeline:

1. **Tag-based versioning**
   - Use semantic tags like `v0.2.0`.
   - Build binaries with injected metadata (`version`, `commit`, `date`) via `-ldflags`.

2. **GitHub Releases artifacts**
   - Build matrix: `linux`, `darwin`, `windows` x `amd64`, `arm64`.
   - Publish `.tar.gz` / `.zip` artifacts + `checksums.txt`.

3. **Install script**
   - Provide `install.sh` that detects OS/ARCH and downloads the correct release binary.
   - Keep a one-liner install flow for users without Go.

4. **Package managers (optional next step)**
   - Homebrew tap (macOS/Linux), Scoop/winget (Windows).
   - Reuse same GitHub Release artifacts as source of truth.

### Package not found

If `install`, `add`, `remove`, `list-skills`, or `remove-skill` reports package not found, verify:

- exact package name (`agentpack list`)
- storage path `~/.agentpack/packages-skills`

### Dry-run

```bash
agentpack remove <package-name> --dry-run
agentpack remove-skill <package-name> <skill-name> --dry-run
agentpack remove <path> --from <package-name> --dry-run
```

## Roadmap

- Extend platform support to `rules`, `commands`, `agents`, and `MCP`.
- Add `config get` and `config list` for runtime settings visibility.
- Expand CI with race checks and optional integration test stage.
- Validate `SKILL.md` frontmatter and conventions (optional mode).
- Add command to rename skills inside a package.
- Add cross-OS support improvements (Windows/macOS).
- Publish binary GitHub releases.

## Contributing

1. Fork repository.
2. Create a feature branch.
3. Run `go fmt ./...` and `go test ./...`.
4. Open a pull request with clear scope and rationale.

## License

MIT. See `LICENSE`.
