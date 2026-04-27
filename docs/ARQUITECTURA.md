# Arquitectura de AgentPack — Especificación para implementación

**Audiencia primaria:** un agente IA que va a ejecutar la migración.
**Audiencia secundaria:** revisores humanos del proyecto.

Este documento es prescriptivo y normativo. No discute alternativas; las alternativas se descartaron en la conversación previa de diseño. Si una situación durante la implementación no está cubierta aquí, **detente y pregunta** (ver §11).

---

## 0. Cómo leer este documento

### 0.1 Convenciones de marcado

Cada regla, decisión o instrucción de este documento lleva uno de los siguientes marcadores. **Trata cada marcador como una directiva ejecutable**, no como un adorno tipográfico.

| Marcador | Significado | Acción del agente |
|---|---|---|
| 🔒 **NO-NEGOCIABLE** | Invariante arquitectónica. Violarla rompe el diseño. | Implementa exactamente como se describe. Si la implementación parece imposible, detente y pregunta. **No improvises.** |
| ❓ **PREGUNTAR** | Decisión donde el usuario tiene preferencia pero existen opciones legítimas. | **Detente y pregunta al usuario** antes de proceder. No asumas. La pregunta exacta a hacer está incluida. |
| 💡 **RECOMENDADO** | Forma sugerida pero hay flexibilidad. | Sigue la recomendación salvo que tengas una razón concreta para desviarte. Si te desvías, documenta por qué en el commit. |
| ⛔ **ANTI-PATRÓN** | Práctica explícitamente prohibida. | Nunca lo hagas. Si te encuentras tentado a hacerlo, hay un problema de diseño en otra parte. |
| ✅ **EJEMPLO CORRECTO** | Forma canónica de hacer algo. | Cópialo, adáptalo. |
| ❌ **EJEMPLO INCORRECTO** | Forma que parece razonable pero es errónea. | Reconócelo y evítalo. |

### 0.2 Cuando una situación no está cubierta

Si durante la implementación encuentras un caso que:

- No está descrito en este documento, **y**
- No es claramente derivable de los principios de §3, **y**
- No es trivialmente convencional en Go

**Detente. Pregunta al usuario.** Formula la pregunta así:

> "Encontré la situación X durante la implementación de Y. Las opciones que veo son A, B, C. Mi recomendación tentativa es B porque [razón]. ¿Cómo procedo?"

No asumas. No elijas la opción que parezca más rápida. No mezcles opciones. Pregunta.

### 0.3 Cuando un test o linter falla

Cuando encuentres un fallo de test, lint o build durante la implementación:

1. **No deshabilites tests ni reglas de lint** para hacer pasar el build. Eso es 🔒 **NO-NEGOCIABLE**.
2. Lee el error, identifica la causa raíz.
3. Si la causa raíz es una regla de este documento que el código actual viola, ese código se debe corregir, no la regla.
4. Si la causa raíz parece ser una regla de este documento que está mal formulada, **detente y pregunta**.

---

## 1. Resumen ejecutivo

**Qué es AgentPack:** CLI en Go para crear, almacenar e instalar paquetes reutilizables de skills/agents/commands/rules para asistentes de codificación con IA. Repo: `github.com/Bbeboy/AgentPack`. Binario distribuido como ejecutable estático para Linux/macOS/Windows × amd64/arm64.

**Patrón arquitectónico adoptado:** Hexagonal (Ports & Adapters) con cuatro capas: `domain`, `app`, `adapter/<entrada>`, `adapter/<salida>`. Composition root en `cmd/agentpack/main.go`.

**Por qué hexagonal y no otro patrón:** decidido en la conversación de diseño. Los argumentos están en el documento de discusión previo. Para esta implementación, **el patrón está cerrado**. No reabrir la decisión.

**Estado del proyecto:** portafolio personal + aprendizaje. Esto implica que la **calidad de lectura del código** y la **claridad arquitectónica** importan tanto como la funcionalidad. Optimizar para legibilidad y demostración de oficio.

---

## 2. Estructura de directorios objetivo

🔒 **NO-NEGOCIABLE:** la siguiente estructura es exacta. Nombres de carpetas, ubicaciones y separación son obligatorios.

```
agentpack/
├── cmd/
│   └── agentpack/
│       └── main.go              # composition root: cablea adaptadores y casos de uso
├── internal/
│   ├── domain/                  # núcleo, sin imports del proyecto
│   │   ├── pkg/                 # entidad Package + invariantes
│   │   ├── skill/               # entidad Skill
│   │   └── platform/            # PlatformTarget como concepto de dominio
│   ├── app/                     # casos de uso + puertos
│   │   ├── ports.go             # TODAS las interfaces de puertos
│   │   ├── <use_case>.go        # uno por caso de uso
│   │   └── <use_case>_test.go
│   ├── adapter/
│   │   ├── cli/                 # entrada: Cobra
│   │   │   ├── root.go
│   │   │   ├── <command>.go
│   │   │   ├── render/
│   │   │   │   ├── human.go
│   │   │   │   └── json.go
│   │   │   └── exit/
│   │   │       └── codes.go
│   │   ├── mcp/                 # entrada: servidor MCP (Fase 3)
│   │   ├── tui/                 # entrada: TUI (Fase 4)
│   │   ├── fs/                  # salida: filesystem real
│   │   ├── store/               # salida: persistencia en ~/.agentpack/
│   │   ├── platform/            # salida: detector de plataformas
│   │   ├── prompt/              # salida: prompts TTY
│   │   └── i18n/                # salida: catálogo de mensajes
│   └── pkg/
│       └── testutil/            # fakes y builders compartidos
├── docs/
│   ├── adr/                     # Architecture Decision Records
│   └── ARQUITECTURA.md          # este archivo (mover aquí desde la raíz)
├── go.mod
├── go.sum
├── .golangci.yml
├── .goreleaser.yaml
└── README.md
```

⛔ **ANTI-PATRÓN:** crear paquetes "utils", "helpers" o "common" en `internal/`. Si necesitas helpers compartidos, viven en el paquete específico al que sirven o en `pkg/testutil` si son solo para tests. **Nunca** crees `internal/utils/`.

⛔ **ANTI-PATRÓN:** poner archivos sueltos en la raíz de `internal/`. Cada archivo `.go` vive dentro de un subpaquete con propósito claro.

---

## 3. Principios rectores

Estos seis principios resuelven el 90% de las dudas de diseño que aparecerán durante la implementación. Cuando dudes, vuelve aquí.

### 3.1 🔒 Las dependencias apuntan hacia adentro

`domain` ← `app` ← `adapter/<x>` ← `cmd/agentpack`

`domain` no importa nada del proyecto. `app` importa solo `domain`. Los adaptadores importan `app` y `domain` pero **nunca otros adaptadores**. `cmd/agentpack` es el único lugar que importa todo.

### 3.2 🔒 Separación de mecanismo y política

La **política** (qué hacer) vive en `domain` y `app`. El **mecanismo** (cómo hacerlo en este sistema, con esta tecnología) vive en `adapter/`. Si una regla de negocio aparece en un adaptador, está mal ubicada. Si una llamada a `os.ReadFile` aparece en `app`, está mal ubicada.

### 3.3 🔒 Las tres etapas del CLI están separadas

Cada invocación de comando atraviesa: (1) **deserialización** de argumentos, (2) **selección** del caso de uso, (3) **ejecución**. Estas tres etapas viven en componentes distintos: la deserialización en `adapter/cli/<comando>.go`, la selección en el cableado de Cobra (estática, no en runtime), la ejecución en `app/<use_case>.go`.

### 3.4 🔒 La descripción ES la interfaz

El texto de ayuda de cada comando (descripción, flags, ejemplos) y el comportamiento del parser son la **misma fuente de verdad**, declarada una sola vez en el `cobra.Command`. Documentación externa (README, web) se genera o referencia desde aquí. No duplicar.

### 3.5 🔒 Conocimiento como datos, no como código

Cuando una regla del sistema dependa de un catálogo enumerable (plataformas soportadas, códigos de error, mensajes), expresarla como **dato** (JSON, TOML, struct literal de tabla) y no como una cadena de `if/else` o `switch`. Aplicación específica: las definiciones de plataforma viven en un archivo JSON embebido al binario, sobrescribible por el usuario.

### 3.6 🔒 Salida estructurada disponible siempre

Todo comando que produzca datos legibles (lista, info, search, stats) debe tener un equivalente `--json` con salida parseable. Esto solo es viable si `app/<use_case>.go` devuelve un `Output` estructurado y la presentación humana/JSON vive en `adapter/cli/render/`.

---

## 4. Capas en detalle

### 4.1 `internal/domain/`

**Responsabilidad:** representar los conceptos del problema y sus reglas, independientes de cualquier tecnología.

🔒 **NO-NEGOCIABLE — restricciones de imports:**
- ❌ `os`, `io`, `io/ioutil`, `path/filepath` para tocar el FS
- ❌ Cualquier paquete de `internal/`
- ❌ `github.com/spf13/cobra` o cualquier librería de CLI
- ❌ `fmt.Println` o cualquier I/O a stdout/stderr
- ✅ `errors`, `fmt.Errorf` (solo para construir errores), `strings`, `time`, `regexp`
- ✅ Otros subpaquetes de `domain/`

**Contenido obligatorio:**

- `pkg.Package`: entidad raíz con `Name`, `CreatedAt`, lista de `SkillRef`. Métodos `AddSkill`, `RemoveSkill`, `Rename`. Toda mutación valida invariantes.
- `pkg.Name`: value object inmutable. Constructor `NewName(string) (Name, error)` que aplica las reglas: ≤64 caracteres, debe empezar con letra o número, charset `[a-zA-Z0-9._-]`. **Estas reglas viven aquí, no en el comando ni en el storage.**
- `skill.Skill`: entidad. Path relativo dentro del paquete + frontmatter validado (opcional).
- `platform.Target`: tipo enumerable + métodos. Cada target conoce su path raíz relativo. **Esto se rediseñará en Fase 2 para cargar desde `platforms.json`** (ver §6.5).

✅ **EJEMPLO CORRECTO** (estructura de `pkg/name.go`):

```go
package pkg

import (
    "errors"
    "regexp"
)

var ErrInvalidName = errors.New("invalid package name")

type Name struct {
    value string
}

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)

func NewName(s string) (Name, error) {
    if len(s) == 0 || len(s) > 64 {
        return Name{}, ErrInvalidName
    }
    if !nameRegex.MatchString(s) {
        return Name{}, ErrInvalidName
    }
    return Name{value: s}, nil
}

func (n Name) String() string { return n.value }
```

**Tests del dominio:**

🔒 **NO-NEGOCIABLE:** tests del dominio son table-driven, sin tmpdirs, sin mocks, sin I/O. Si un test del dominio necesita un mock o un `t.TempDir()`, la abstracción está mal puesta y debe moverse fuera de `domain/`.

### 4.2 `internal/app/`

**Responsabilidad:** orquestar las operaciones que el usuario invoca. Un caso de uso = una operación de alto nivel = un comando del CLI (más adelante también del MCP y TUI).

🔒 **NO-NEGOCIABLE — patrón obligatorio:** cada caso de uso es una struct con dependencias inyectadas y un único método público `Execute(ctx, input) (output, error)`.

✅ **EJEMPLO CORRECTO:**

```go
// internal/app/install_package.go
package app

import "context"

type InstallPackage struct {
    Store    PackageStore
    FS       FileSystem
    Detector PlatformDetector
    Prompter Prompter
}

type InstallInput struct {
    PackageName string
    DryRun      bool
    OnConflict  ConflictStrategy
    Platforms   []string // empty = autodetect
    Global      bool
}

type InstallOutput struct {
    Target          platform.Target
    InstalledSkills []skill.Skill
    SkippedSkills   []skill.Skill
    DryRun          bool
}

func (uc *InstallPackage) Execute(ctx context.Context, in InstallInput) (InstallOutput, error) {
    // 1. validar input (construir pkg.Name, etc.)
    // 2. cargar paquete del store
    // 3. detectar plataforma o usar la solicitada
    // 4. para cada skill: decidir conflicto, aplicar
    // 5. devolver Output estructurado
}
```

**Puertos (`ports.go`):**

🔒 **NO-NEGOCIABLE:** todos los puertos viven en `app/ports.go`. No dispersos por archivos. Si un puerto solo lo usa un caso de uso, igual va en `ports.go`.

🔒 **NO-NEGOCIABLE:** los puertos son **interfaces pequeñas y específicas al caso de uso**. No hay un mega-`FileSystem` con 40 métodos. Hay `FileReader`, `FileWriter`, `DirCopier` separados si los casos de uso los usan separados. Aplica el principio de segregación de interfaces.

❓ **PREGUNTAR antes de implementar puertos:** ¿prefieres que los puertos estén separados al máximo (un puerto por capability, ej. `Reader`, `Writer`, `Existence`, `Copier`) o agrupados por dominio funcional (ej. un solo `FileSystem` con los métodos que use AgentPack)? Mi recomendación tentativa: agrupados por dominio funcional para simplicidad, separados solo si un test específico se beneficia. Pero confirma.

**Tests de aplicación:**

🔒 **NO-NEGOCIABLE:** los tests de `app/` se prueban con **fakes en memoria** que viven en `internal/pkg/testutil/`. Cada test cubre una historia de usuario completa sin tocar disco. No `t.TempDir()` en tests de app.

⛔ **ANTI-PATRÓN:** usar `gomock`, `mockery` u otra librería de mocks generados. Los fakes a mano son más legibles, más mantenibles, y demuestran mejor oficio. Esto es una regla deliberada de este proyecto.

### 4.3 `internal/adapter/cli/`

**Responsabilidad:** las tres etapas de Zamora —deserializar, enrutar, presentar— y nada más.

**Patrón por comando:**

1. Declarar el `cobra.Command` con flags, args y descripción.
2. En el `RunE`: leer flags y args, construir el `Input` del caso de uso, llamar `Execute`, formatear la salida con `render`, mapear errores a exit codes.

🔒 **NO-NEGOCIABLE — lo que NO debe hacer este paquete:**
- ❌ Tocar el filesystem directamente
- ❌ Validar reglas de dominio (nombres, formatos)
- ❌ Conocer dónde viven los paquetes en disco
- ❌ Formar paths del store
- ❌ Llamar directamente a `internal/adapter/store/` o `internal/adapter/fs/`

Toda esa lógica vive en dominio o en adaptadores de salida. El adapter `cli` solo conoce `app/` y se le inyectan los casos de uso ya construidos.

✅ **EJEMPLO CORRECTO** (estructura de un comando):

```go
// internal/adapter/cli/install.go
package cli

func NewInstallCmd(uc *app.InstallPackage, cat *i18n.Catalog) *cobra.Command {
    var dryRun bool
    var jsonOut bool
    var onConflict string
    var platforms []string

    cmd := &cobra.Command{
        Use:   "install <package>",
        Short: cat.T("install.short"),
        Long:  cat.T("install.long"),
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            in := app.InstallInput{
                PackageName: args[0],
                DryRun:      dryRun,
                OnConflict:  parseConflict(onConflict),
                Platforms:   platforms,
            }
            out, err := uc.Execute(cmd.Context(), in)
            if err != nil {
                return err // mapeo a exit code en main
            }
            if jsonOut {
                return render.JSON(cmd.OutOrStdout(), out)
            }
            return render.InstallHuman(cmd.OutOrStdout(), out, cat)
        },
    }
    cmd.Flags().BoolVar(&dryRun, "dry-run", false, cat.T("install.flag.dry_run"))
    cmd.Flags().BoolVar(&jsonOut, "json", false, cat.T("flag.json"))
    cmd.Flags().StringVar(&onConflict, "on-conflict", "prompt", cat.T("install.flag.on_conflict"))
    cmd.Flags().StringSliceVar(&platforms, "platforms", nil, cat.T("install.flag.platforms"))
    return cmd
}
```

### 4.4 `internal/adapter/<salida>/`

Cada subdirectorio implementa exactamente un puerto definido en `app/ports.go`.

🔒 **NO-NEGOCIABLE:** un adaptador de salida no importa otros adaptadores. Si dos adaptadores comparten lógica, esa lógica sube a `app/` o a `domain/`.

**Adaptadores específicos a implementar:**

- `adapter/fs/osfs.go` — `FileSystem` con `os` y `path/filepath`. Todo path proveniente de input de usuario pasa por `filepath.Clean` y se valida que no escape del directorio esperado (mitigación de path traversal).
- `adapter/store/homedir_store.go` — `PackageStore` contra `~/.agentpack/packages-skills`. Resuelve el directorio respetando `$XDG_DATA_HOME` cuando está definido. **El formato en disco es un detalle interno** y debe poder cambiarse sin tocar `app/` ni `domain/`.
- `adapter/platform/detector.go` — `PlatformDetector`. Lee `platforms.json` (embebido + override de usuario, ver §6.5).
- `adapter/prompt/survey.go` — `Prompter` con una librería como `survey` o `huh`. Detecta TTY; sin TTY, devuelve `ErrNoTTY` (no se cuelga esperando input).
- `adapter/i18n/catalog.go` — catálogo de mensajes. Claves estables (`install.success.target`), nunca el texto en sí como clave.

❓ **PREGUNTAR antes de elegir librería de prompts:** ¿prefieres `survey/v2` (más establecida, API clásica) o `huh` (de Charm, integra mejor con el ecosistema Bubbletea que vas a usar para la TUI)? Mi recomendación tentativa: `huh`, por consistencia con la TUI futura. Pero confirma.

### 4.5 `cmd/agentpack/main.go`

🔒 **NO-NEGOCIABLE:** este es el ÚNICO archivo del proyecto que conoce simultáneamente todas las capas. Aquí se construyen instancias concretas y se cablea.

✅ **EJEMPLO CORRECTO** (esqueleto):

```go
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    // ...
)

func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    // Adaptadores concretos
    fs := osfs.New()
    store := homedirstore.New(resolveDataDir())
    detector := platformdetect.New(loadPlatformsConfig())
    prompter := huhprompt.New(os.Stdin, os.Stdout)
    catalog := i18n.Load(currentLang())

    // Casos de uso
    install := &app.InstallPackage{Store: store, FS: fs, Detector: detector, Prompter: prompter}
    create  := &app.CreatePackage{Store: store, FS: fs, Detector: detector}
    list    := &app.ListPackages{Store: store}
    // ...

    // Cableado CLI
    root := cli.NewRootCmd(catalog)
    root.AddCommand(cli.NewInstallCmd(install, catalog))
    root.AddCommand(cli.NewCreateCmd(create, catalog))
    root.AddCommand(cli.NewListCmd(list, catalog))
    // ...

    err := root.ExecuteContext(ctx)
    os.Exit(exit.From(err))
}
```

---

## 5. Reglas de dependencia (enforcement mecánico)

🔒 **NO-NEGOCIABLE:** la siguiente tabla se hace cumplir vía `depguard` en `.golangci.yml`. Una violación bloquea el merge.

| De \ A | `domain` | `app` | `adapter/cli` | `adapter/<otros>` | `cmd/agentpack` |
|---|---|---|---|---|---|
| `domain` | ✅ | ❌ | ❌ | ❌ | ❌ |
| `app` | ✅ | ✅ | ❌ | ❌ | ❌ |
| `adapter/cli` | ✅ | ✅ | ✅ | ❌ | ❌ |
| `adapter/<otros>` | ✅ | ✅ | ❌ | ✅ (mismo paquete) | ❌ |
| `cmd/agentpack` | ✅ | ✅ | ✅ | ✅ | ✅ |

Configuración mínima en `.golangci.yml`:

```yaml
linters:
  enable: [depguard, gofmt, govet, staticcheck, errcheck, ineffassign]

linters-settings:
  depguard:
    rules:
      domain:
        list-mode: lax
        files: ["**/internal/domain/**"]
        deny:
          - pkg: "github.com/Bbeboy/AgentPack/internal/app"
          - pkg: "github.com/Bbeboy/AgentPack/internal/adapter"
          - pkg: "github.com/Bbeboy/AgentPack/cmd"
          - pkg: "os"
          - pkg: "path/filepath"
      app:
        list-mode: lax
        files: ["**/internal/app/**"]
        deny:
          - pkg: "github.com/Bbeboy/AgentPack/internal/adapter"
          - pkg: "github.com/Bbeboy/AgentPack/cmd"
          - pkg: "github.com/spf13/cobra"
```

---

## 6. Convenciones transversales

### 6.1 🔒 Naming

- **Comandos del CLI:** siempre en inglés (`install`, `create`, `add-skill`). No traducibles.
- **Identificadores en Go:** inglés. Sin excepciones.
- **Mensajes al usuario (help, errores, output):** traducibles vía catalog i18n.
- **Claves del catalog:** snake_case con jerarquía por punto: `install.success.target`, `error.package_not_found`.

### 6.2 ❓ PREGUNTAR — Idioma de los comentarios del código

¿Prefieres comentarios en inglés (consistente con el ecosistema Go global, mejor para colaboradores externos) o en español (consistente con tu audiencia primaria, demuestra el carácter bilingüe del proyecto)? Mi recomendación tentativa: inglés en código exportado/público, español opcional en comentarios internos extensos. Pero confirma.

### 6.3 🔒 Errores

- Errores tipados en dominio y aplicación. Sentinels (`var ErrPackageNotFound = errors.New(...)`) o tipos cuando carguen contexto (`type ConflictError struct { Path string }`).
- En el adaptador CLI, mapeo de error → exit code en `adapter/cli/exit/codes.go`.
- Mensajes al usuario con formato: `Error (E0003): <mensaje>. → <sugerencia accionable>.`
- Códigos de error (E0003) son **estables** entre versiones. Una vez asignados, no se reusan.

✅ **EJEMPLO CORRECTO** de mensaje de error:

```
Error (E0003): package "backend-base" not found.
  → Run `agentpack list` to see available packages,
    or `agentpack create backend-base <skills-path>` to create it.
```

### 6.4 🔒 Streams (stdout vs stderr)

- **stdout:** datos de salida (lista, JSON, info). Lo que un script `>` redirigiría.
- **stderr:** logs, errores, prompts, mensajes informativos.
- `agentpack list --json > pkgs.json` debe producir un archivo limpio sin contaminación de mensajes.

### 6.5 🔒 Plataformas como datos

Las plataformas soportadas viven en `platforms.json` en la raíz del proyecto, embebido al binario vía `embed.FS`. Estructura mínima:

```json
{
  "schema_version": "1",
  "platforms": [
    {
      "id": "opencode",
      "name": "OpenCode",
      "directory": ".opencode/",
      "subdirs": {
        "skills": "skills/",
        "commands": "commands/",
        "agents": "agents/"
      },
      "mcp_file": "opencode.json",
      "priority": 10
    },
    {
      "id": "claude_code",
      "name": "Claude Code",
      "directory": ".claude/",
      "root_file": "CLAUDE.md",
      "subdirs": {
        "skills": "skills/",
        "commands": "commands/",
        "agents": "agents/",
        "rules": "rules/"
      },
      "mcp_file": ".mcp.json",
      "priority": 20
    }
  ]
}
```

🔒 **NO-NEGOCIABLE:** el archivo del usuario en `~/.agentpack/platforms.json` (o `$XDG_CONFIG_HOME/agentpack/platforms.json`) hace **deep-merge** sobre el embebido. Permite añadir plataformas sin recompilar.

❓ **PREGUNTAR antes de fijar el schema:** los campos arriba son una propuesta. ¿Quieres que tome inspiración exacta del `platforms.jsonc` de OpenPackage para facilitar interoperabilidad futura, o diseñas un schema propio? Mi recomendación: schema propio, simple, pero con un campo `compat: { openpackage_id: "..." }` para mapeo futuro. Confirma.

### 6.6 🔒 Configuración con precedencia

Orden de precedencia, mayor a menor prioridad:

1. Flags de comando (`--language es`)
2. Variables de entorno (`AGENTPACK_LANG`, `NO_COLOR`, `DEBUG`, `XDG_DATA_HOME`)
3. Config de proyecto (`.agentpack/config.toml` en cwd o ancestros, opcional)
4. Config de usuario (`$XDG_CONFIG_HOME/agentpack/config.toml`)
5. Config del sistema (`/etc/agentpack/config.toml`, opcional)
6. Defaults compilados

Implementación: un `config.Resolver` en aplicación que recibe las distintas fuentes ya cargadas por adaptadores y aplica la cadena.

### 6.7 🔒 Logging y debug

- 💡 **RECOMENDADO:** `log/slog` (stdlib desde Go 1.21). No traer dependencias externas para logging.
- Nivel de log controlado por `--debug` (flag global) y `AGENTPACK_DEBUG=1`.
- Logs en debug van a stderr. **Nunca** a stdout.

### 6.8 🔒 Telemetría

⛔ **ANTI-PATRÓN:** AgentPack **no envía telemetría**, no hace "phone home", no rastrea uso. No hay un flag para habilitarla, ni siquiera opt-in. Esta es una decisión de producto fija.

---

## 7. Testing

### 7.1 🔒 Pirámide

| Capa | Tipo | Cobertura mínima |
|---|---|---|
| `domain/` | Unitario puro, table-driven | 90% |
| `app/` | Unitario con fakes en memoria | 85% |
| `adapter/<x>/` | Integración con tmpdir | 70% |
| `cmd/agentpack/` | Smoke E2E (binario compilado) | ≥1 test/comando |

### 7.2 🔒 Tabla de tests

Tests obligatoriamente table-driven cuando hay >2 casos. Estructura:

```go
func TestNewName(t *testing.T) {
    cases := []struct{
        name    string
        input   string
        wantErr bool
    }{
        {"valid simple",     "mypkg",        false},
        {"valid with dots",  "my.pkg.v2",    false},
        {"empty",            "",             true},
        {"too long",         strings.Repeat("a", 65), true},
        {"starts with dash", "-pkg",         true},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            _, err := NewName(tc.input)
            if (err != nil) != tc.wantErr {
                t.Errorf("got err=%v, wantErr=%v", err, tc.wantErr)
            }
        })
    }
}
```

### 7.3 🔒 Golden files

La salida de cada comando (humana y JSON) tiene un golden file en `testdata/golden/<command>_<scenario>.{txt,json}`. Tests E2E comparan output contra el golden. Actualizar con `go test ./... -update` (flag custom).

### 7.4 🔒 Locales en tests

Tests E2E corren con `LANG=C` o `LANG=en_US.UTF-8` explícitamente fijado. Nunca confiar en el locale del entorno.

---

## 8. CI / cross-platform

### 8.1 🔒 Matriz mínima de tests

```yaml
strategy:
  fail-fast: false
  matrix:
    os: [ubuntu-latest, ubuntu-24.04-arm, macos-13, macos-latest, windows-latest]
```

Donde:
- `ubuntu-latest` = Linux amd64
- `ubuntu-24.04-arm` = Linux arm64
- `macos-13` = macOS amd64 (Intel)
- `macos-latest` = macOS arm64 (Apple Silicon)
- `windows-latest` = Windows amd64

### 8.2 🔒 Smoke test del binario compilado

Por cada combinación del matrix, además de `go test ./...`, ejecutar:

```yaml
- run: go build -o agentpack ./cmd/agentpack
- run: ./agentpack version
- run: ./agentpack list
- run: ./agentpack --help
```

(En Windows: `agentpack.exe`).

### 8.3 ❓ PREGUNTAR — Windows ARM64

Windows ARM64 no tiene runners nativos baratos en GitHub Actions. ¿Aceptas que se cross-compile pero no se teste nativamente en CI? Mi recomendación: sí, con nota explícita en el README. Confirma.

### 8.4 🔒 Cross-compilation check

El build `release` (goreleaser) debe producir binarios para **6 targets**:

- `linux_amd64`
- `linux_arm64`
- `darwin_amd64`
- `darwin_arm64`
- `windows_amd64`
- `windows_arm64`

CI verifica que los 6 compilan en cada PR (no solo en release).

---

## 9. Versionado

### 9.1 🔒 SemVer estricto

El **contrato público** es: nombres de comandos, flags, formato de salida JSON, formato de almacenamiento en `~/.agentpack/`. Romper cualquiera = bump mayor.

### 9.2 ❓ PREGUNTAR — Manejo del breaking change `remove`/`uninstall`

Decisión pendiente de la conversación previa. Tres opciones:

A. Romper directamente: añadir `uninstall <pkg>` (workspace) y `delete <pkg>` (store), eliminar `remove <pkg>` y `remove-skill`. Bump 0.2 → 0.3.

B. Romper con grace period: añadir nuevos, mantener `remove` y `remove-skill` como deprecados con warning durante 1-2 versiones, eliminar después.

C. No romper: añadir solo `uninstall <pkg>` (workspace), mantener `remove <pkg>` para el store.

Mi recomendación: **B**. Demuestra disciplina de versionado, da grace period, deja la decisión documentada en CHANGELOG. Pero confirma cuál prefieres.

### 9.3 🔒 Migración del formato de storage

Cuando cambie el layout en disco, `homedir_store.go` detecta la versión y aplica migraciones idempotentes al arranque. **Nunca pedir al usuario migrar manualmente.**

---

## 10. Documentación obligatoria del repo

🔒 **NO-NEGOCIABLE:** después de la migración, el repo debe contener:

- `README.md` — quickstart, install, ejemplos, badges de CI por plataforma, link a docs.
- `README.es.md` — mismo en español (ya existe, mantener sincronizado).
- `docs/ARQUITECTURA.md` — este documento (puede vivir aquí).
- `docs/adr/` — Architecture Decision Records, uno por decisión mayor:
  - `0001-hexagonal-architecture.md`
  - `0002-cobra-for-cli.md`
  - `0003-platforms-as-data.md`
  - `0004-bubbletea-for-tui.md`
  - `0005-mcp-server-as-adapter.md`
  - …
- `CHANGELOG.md` — formato Keep a Changelog.
- `CONTRIBUTING.md` — cómo correr tests, cómo abrir PR, convenciones.

Cada ADR es corto: contexto, decisión, consecuencias, alternativas consideradas. 1-2 páginas máximo.

---

## 11. Lista consolidada de cosas a preguntar antes de empezar

Antes de escribir la primera línea de código de la migración, el agente IA debe haber recibido respuesta de Daniel a las siguientes preguntas. **No avanzar sin estas respuestas.**

| # | Pregunta | Recomendación tentativa |
|---|---|---|
| Q1 | ¿Puertos agrupados por dominio funcional o segregados al máximo? (§4.2) | Agrupados, segregar solo si hace falta |
| Q2 | ¿Librería de prompts: `survey/v2` o `huh`? (§4.4) | `huh` (consistencia con TUI futura) |
| Q3 | ¿Idioma de comentarios en código? (§6.2) | Inglés en exportado, español opcional internos |
| Q4 | ¿Schema de `platforms.json` propio o inspirado en OpenPackage? (§6.5) | Propio, con campo `compat` para mapeo |
| Q5 | ¿Aceptas que Windows ARM64 sea cross-compile sin test nativo? (§8.3) | Sí, con nota en README |
| Q6 | ¿Cómo manejar el rename `remove` → `uninstall`/`delete`? (§9.2) | Opción B: grace period con deprecation warning |
| Q7 | ¿Versión mínima de Go a soportar? | Go 1.23 (matchea el actual) |
| Q8 | ¿Nombre del módulo en `go.mod`? Confirmar `github.com/Bbeboy/AgentPack` | Mantener |
| Q9 | ¿Manifiesto por paquete (`package.yml` o similar) en esta migración o en fase posterior? | Posterior, no introducir aún |
| Q10 | ¿Soportar instalación desde GitHub URLs en algún momento? | No en esta migración. Discutir en fase posterior |

Cuando haya respuesta a las 10 preguntas, registrarlas en `docs/adr/0000-decisions-log.md` y comenzar.

---

## 12. Lista consolidada de invariantes (resumen)

Las siguientes reglas son 🔒 **NO-NEGOCIABLES**. Una violación bloquea el merge:

1. Las dependencias apuntan hacia adentro (§3.1, §5).
2. El dominio no toca filesystem ni stdlib I/O (§4.1).
3. Tests del dominio son puros, sin tmpdirs ni mocks (§4.1).
4. Cada caso de uso = struct con dependencias inyectadas + `Execute(ctx, in) (out, err)` (§4.2).
5. Todos los puertos en `app/ports.go` (§4.2).
6. No usar librerías de mocks generados; fakes a mano en `pkg/testutil/` (§4.2).
7. CLI no toca filesystem ni implementa reglas de dominio (§4.3).
8. Un adaptador no importa otros adaptadores (§4.4).
9. `cmd/agentpack/main.go` es el único composition root (§4.5).
10. Comandos en inglés, mensajes traducibles, identificadores Go en inglés (§6.1).
11. Errores tipados con códigos estables (§6.3).
12. stdout para datos, stderr para todo lo demás (§6.4).
13. Plataformas como datos en `platforms.json` embebido + override (§6.5).
14. No telemetría (§6.8).
15. Cobertura mínima por capa (§7.1).
16. Tests table-driven cuando >2 casos (§7.2).
17. Golden files para output de comandos (§7.3).
18. Matrix de CI cubre 5+ targets (§8.1).
19. Smoke test del binario compilado en CI (§8.2).
20. SemVer estricto, contrato público definido (§9.1).

---

*Este documento es la fuente de verdad. Cualquier desviación durante la implementación debe registrarse como ADR (`docs/adr/`), no quedarse en la conversación de un PR.*
