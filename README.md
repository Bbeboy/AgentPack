# AgentPack

CLI en Go para crear, guardar e instalar paquetes reutilizables de skills para agentes de IA.

Con AgentPack puedes tomar las skills de un proyecto, empaquetarlas localmente y volver a instalarlas en cualquier otro proyecto con un solo comando.

## Caracteristicas principales

- Crea paquetes desde una ruta de skills existente.
- Soporta deteccion de ruta cuando usas `.` como origen.
- Instala en la carpeta `skills` de la plataforma detectada (por ejemplo `.opencode/skills`, `.agents/skills`, `.cla/skills`, etc.).
- Detecta conflictos por nombre de skill y pregunta si sobrescribir o ignorar.
- Mantiene archivos existentes que no estan en conflicto.
- Permite borrar paquetes con confirmacion o de forma forzada.
- Permite eliminar rutas internas de un paquete con `remove <ruta> --from <paquete>`.
- Permite agregar archivos o carpetas a un paquete con `add <ruta> --to <paquete>`.
- Permite listar skills por paquete y eliminar una skill especifica.
- Soporta idioma EN/ES para ayuda y feedback (`config set language` o `lang`).
- Incluye modo seguro `--dry-run` para revisar borrados sin ejecutar.
- Soporta autocompletado para `bash`, `zsh`, `fish` y `powershell`.

## Tabla de contenido

- [Instalacion express (2 minutos)](#instalacion-express-2-minutos)
- [Requisitos](#requisitos)
- [Instalacion](#instalacion)
  - [Opcion recomendada: go install](#opcion-recomendada-go-install)
  - [Si agentpack no aparece en PATH](#si-agentpack-no-aparece-en-path)
  - [Opcion 2: build desde fuente](#opcion-2-build-desde-fuente)
  - [Autocompletado](#autocompletado)
- [Uso rapido](#uso-rapido)
- [Conceptos clave](#conceptos-clave)
- [Referencia de comandos](#referencia-de-comandos)
- [Flujo completo con ejemplos](#flujo-completo-con-ejemplos)
- [Estructura del proyecto](#estructura-del-proyecto)
- [Desarrollo local](#desarrollo-local)
- [Troubleshooting](#troubleshooting)
- [Roadmap](#roadmap)
- [Contribuciones](#contribuciones)
- [Licencia](#licencia)

## Instalacion express (2 minutos)

Si ya tienes Go instalado, con esto puedes empezar:

```bash
go install github.com/Bbeboy/AgentPack/cmd/agentpack@latest
agentpack --help
```

Si `agentpack` no se encuentra, revisa la seccion [Si agentpack no aparece en PATH](#si-agentpack-no-aparece-en-path).

## Requisitos

- Go `1.23+`
- Linux (soporte actual)

Nota: el codigo ya usa `filepath` y `os.UserHomeDir`, por lo que esta preparado para evolucionar a multiplataforma.

## Instalacion

### Opcion recomendada: go install

```bash
go install github.com/Bbeboy/AgentPack/cmd/agentpack@latest
```

Verifica que quedo instalado:

```bash
agentpack --help
```

### Si agentpack no aparece en PATH

`go install` deja el binario en `GOBIN` o, si no esta definido, en `$(go env GOPATH)/bin`.

En la mayoria de instalaciones de Go, eso equivale a `~/go/bin`.

Bash:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

Zsh:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

Fish:

```fish
fish_add_path ~/go/bin
```

### Opcion 2: build desde fuente

```bash
git clone https://github.com/Bbeboy/AgentPack.git
cd AgentPack
go mod tidy
go build -o agentpack ./cmd/agentpack
install -m 755 ./agentpack ~/.local/bin/agentpack
agentpack --help
```

En muchos sistemas Linux, `~/.local/bin` ya esta en `PATH`. Si no lo esta:

- En Bash/Zsh, agrega `~/.local/bin` al `PATH` en tu archivo de shell.
- En Fish, usa:

```fish
fish_add_path ~/.local/bin
```

### Autocompletado

Comando rapido para la sesion actual:

```bash
source <(agentpack completion bash)
```

Para Fish (persistente):

```fish
agentpack completion fish > ~/.config/fish/completions/agentpack.fish
```

Luego reinicia la sesion de Fish o ejecuta:

```fish
exec fish
```

## Uso rapido

1. Crear un paquete desde el proyecto actual:

```bash
agentpack create backend-base .
```

2. Ver paquetes guardados:

```bash
agentpack list
```

3. Instalar el paquete en otro proyecto:

```bash
agentpack install backend-base
```

4. Eliminar un paquete guardado:

```bash
agentpack remove backend-base
```

5. Listar skills de un paquete:

```bash
agentpack list-skills backend-base
```

6. Eliminar una skill especifica:

```bash
agentpack remove-skill backend-base docker
```

7. Agregar archivo o carpeta a un paquete:

```bash
agentpack add ./mi-nueva-skill --to backend-base
```

8. Eliminar ruta interna de un paquete:

```bash
agentpack remove docker/SKILL.md --from backend-base
```

9. Cambiar idioma de salida:

```bash
agentpack config set language es
agentpack lang en
```

10. Renombrar un paquete:

```bash
agentpack rename backend-base backend-v2
```

## Conceptos clave

### Donde se guardan los paquetes

Los paquetes se guardan localmente en:

```text
~/.agentpack/packages-skills/<nombre-paquete>
```

### Donde se instalan las skills

`install` detecta la plataforma por carpetas del proyecto y copia las skills a `<plataforma>/skills`.

Ejemplos de destino:

- `.opencode/skills`
- `.agents/skills`
- `.cla/skills`
- `.cursor/skills`

Si no detecta ninguna plataforma, usa fallback a:

```text
.agents/skills
```

en el proyecto donde ejecutas el comando.

### Validacion de nombre de paquete

El nombre del paquete debe cumplir:

- Maximo 64 caracteres.
- Debe iniciar con letra o numero.
- Caracteres permitidos: letras, numeros, `.`, `_`, `-`.

Ejemplos validos:

- `backend-base`
- `pack_v1`
- `skills.2026`

## Referencia de comandos

| Comando | Descripcion | Ejemplo |
| --- | --- | --- |
| `agentpack create <nombre-paquete> <ruta-skills>` | Crea un paquete desde una ruta de skills. | `agentpack create backend-base /mi/proyecto/.agents/skills` |
| `agentpack install <nombre-paquete>` | Instala un paquete en la carpeta `skills` de la plataforma detectada (fallback `.agents/skills`). | `agentpack install backend-base` |
| `agentpack add <archivo-o-carpeta> --to <nombre-paquete>` | Agrega un archivo o carpeta a un paquete existente. | `agentpack add ./nueva-skill --to backend-base` |
| `agentpack list` | Lista paquetes guardados localmente. | `agentpack list` |
| `agentpack list-skills <nombre-paquete>` | Lista las skills dentro de un paquete. | `agentpack list-skills backend-base` |
| `agentpack rename <nombre-actual> <nombre-nuevo>` | Renombra un paquete existente. | `agentpack rename backend-base backend-v2` |
| `agentpack remove <nombre-paquete>` | Elimina un paquete guardado (con confirmacion). | `agentpack remove backend-base` |
| `agentpack remove <archivo-o-carpeta> --from <nombre-paquete>` | Elimina una ruta especifica dentro de un paquete. | `agentpack remove docker/SKILL.md --from backend-base` |
| `agentpack remove <nombre-paquete> --force` | Elimina un paquete sin preguntar confirmacion. | `agentpack remove backend-base -f` |
| `agentpack remove <nombre-paquete> --dry-run` | Simula la eliminacion de un paquete sin borrar. | `agentpack remove backend-base --dry-run` |
| `agentpack remove-skill <nombre-paquete> <nombre-skill>` | Elimina una skill especifica de un paquete (con confirmacion). | `agentpack remove-skill backend-base docker` |
| `agentpack remove-skill <nombre-paquete> <nombre-skill> -f` | Elimina una skill sin preguntar confirmacion. | `agentpack remove-skill backend-base docker -f` |
| `agentpack remove-skill <nombre-paquete> <nombre-skill> --dry-run` | Simula la eliminacion de una skill sin borrar. | `agentpack remove-skill backend-base docker --dry-run` |
| `agentpack config set language <en\|es>` | Guarda el idioma de salida global. | `agentpack config set language es` |
| `agentpack lang <en\|es>` | Atajo para cambiar idioma. | `agentpack lang en` |
| `agentpack completion [bash\|zsh\|fish\|powershell]` | Genera script de autocompletado. | `agentpack completion fish` |

## Flujo completo con ejemplos

### 1) Crear paquete desde ruta explicita

```bash
agentpack create backend-base /mi/proyecto/.agents/skills
```

Salida esperada:

```text
agentpack: creating package 'backend-base'
agentpack: package created at /home/usuario/.agentpack/packages-skills/backend-base
```

### 2) Crear paquete usando `.`

```bash
agentpack create backend-base .
```

Cuando usas `.`, AgentPack busca estas rutas en orden:

1. `.agents/skills`
2. `.opencode/skills`
3. `.agent/skills`
4. `skills`

Si encuentra varias, muestra un menu interactivo para que elijas.

### 3) Instalar paquete sin conflictos

```bash
agentpack install backend-base
```

Salida esperada (ejemplo):

```text
agentpack: installing 'backend-base' into /ruta/proyecto/.opencode/skills
agentpack: installed=4 overwritten=0 skipped=0
```

### 4) Instalar paquete con conflictos

Si ya existe una skill con el mismo nombre en el destino detectado (`<plataforma>/skills`), se considera conflicto.

Ejemplo:

```text
agentpack: conflicts detected
agentpack: 'docker' already exists
overwrite skill 'docker'? [y/N]: y
agentpack: overwrote 'docker'
```

Si respondes `n` (o Enter), la skill se ignora y no se toca.

### 5) Eliminar paquete guardado

Con confirmacion:

```bash
agentpack remove backend-base
```

Forzado (sin pregunta):

```bash
agentpack remove backend-base --force
```

Simulacion sin borrar:

```bash
agentpack remove backend-base --dry-run
```

### 6) Listar skills de un paquete

```bash
agentpack list-skills backend-base
```

Salida esperada (ejemplo):

```text
agentpack: skills in 'backend-base' (3)
- docker
- golang-api
- testing
```

### 7) Eliminar una skill especifica

Con confirmacion:

```bash
agentpack remove-skill backend-base docker
```

Forzado (sin pregunta):

```bash
agentpack remove-skill backend-base docker -f
```

Simulacion sin borrar:

```bash
agentpack remove-skill backend-base docker --dry-run
```

### 8) Agregar archivo/carpeta a un paquete

```bash
agentpack add ./skills/docker --to backend-base
```

Salida esperada (ejemplo):

```text
agentpack: adding './skills/docker' to package 'backend-base'
agentpack: added 'docker' to package 'backend-base'
```

### 9) Eliminar ruta interna de un paquete

```bash
agentpack remove docker/SKILL.md --from backend-base
```

Salida esperada (ejemplo):

```text
remove 'docker/SKILL.md' from package 'backend-base'? [y/N]: y
agentpack: removing 'docker/SKILL.md' from package 'backend-base'
agentpack: removed 'docker/SKILL.md' from package 'backend-base'
```

### 10) Cambiar idioma de salida

```bash
agentpack config set language es
agentpack --help
agentpack lang en
```

Salida esperada (ejemplo):

```text
agentpack: idioma actualizado a es
...
agentpack: language set to en
```

### 11) Renombrar un paquete

```bash
agentpack rename backend-base backend-v2
```

Salida esperada (ejemplo):

```text
agentpack: renaming package 'backend-base' to 'backend-v2'
agentpack: package renamed from 'backend-base' to 'backend-v2'
```

## Estructura del proyecto

```text
cmd/
  agentpack/
    main.go                # entrypoint
internal/
  cli/                     # comandos Cobra
    root.go
    i18n.go
    create.go
    install.go
    add.go
    list.go
    list_skills.go
    rename.go
    remove.go
    remove_skill.go
    config.go
    lang.go
    completion.go
  i18n/
    messages.go            # catalogo de mensajes EN/ES
  platform/
    skills.go              # resolucion de destino por plataforma
  config/
    settings.go            # config global (~/.agentpack/config.json)
  storage/
    storage.go             # rutas de almacenamiento local
  fsutil/
    copy.go                # copia y merge de archivos/carpetas
  prompt/
    prompt.go              # prompts interactivos
```

## Desarrollo local

### Ejecutar en modo desarrollo

```bash
go run ./cmd/agentpack --help
```

### Formatear, testear y compilar

```bash
go fmt ./...
go test ./...
go build -o agentpack ./cmd/agentpack
```

## Troubleshooting

### `agentpack: command not found`

- Verifica que el binario exista en `~/.local/bin/agentpack`.
- Verifica que `~/.local/bin` este en tu `PATH`.
- En Fish, ejecuta `fish_add_path ~/.local/bin` y reinicia shell.

### Error: paquete no encontrado

Si ejecutas `install`, `add`, `remove`, `list-skills` o `remove-skill` y no existe el paquete, verifica:

- nombre exacto del paquete (`agentpack list`),
- ruta de almacenamiento `~/.agentpack/packages-skills`.

### Eliminar sin pregunta de confirmacion

Si quieres eliminar directo un paquete, usa:

```bash
agentpack remove <nombre-paquete> --force
```

Para skill especifica:

```bash
agentpack remove-skill <nombre-paquete> <nombre-skill> -f
```

### Simular eliminacion sin borrar

Paquete:

```bash
agentpack remove <nombre-paquete> --dry-run
```

Skill especifica:

```bash
agentpack remove-skill <nombre-paquete> <nombre-skill> --dry-run
```

Ruta interna del paquete:

```bash
agentpack remove <ruta> --from <nombre-paquete> --dry-run
```

### Error: skill no encontrada

Si `remove-skill` falla por skill inexistente:

- valida el nombre de la skill con `agentpack list-skills <nombre-paquete>`,
- revisa la ruta reportada en el mensaje de error.

### Error al crear con `.`

Si no detecta ruta de skills, crea una de estas carpetas en el proyecto:

- `.agents/skills`
- `.opencode/skills`
- `.agent/skills`
- `skills`

### Cambiar idioma de ayuda y mensajes

Usa cualquiera de estas opciones:

```bash
agentpack config set language es
agentpack lang en
```

Idioma por defecto: `en`.

### Conflictos al instalar

El conflicto se resuelve por skill (carpeta). Puedes:

- sobrescribir (`y`) para merge con overwrite en conflicto,
- ignorar (`n`) para mantener la skill actual intacta.

## Roadmap

- Consolidar la matriz de plataformas para `skills` con pruebas automatizadas por entorno.
- Validacion opcional de `SKILL.md` (frontmatter y convenciones).
- Comando para renombrar skills dentro de un paquete.
- Extender soporte multi-plataforma para `rules`, `commands`, `agents` y `MCP` (ademas de `skills`).
- Mejoras de soporte multiplataforma de sistema operativo (Windows/macOS).
- Publicacion con releases binarias en GitHub.

## Contribuciones

Las contribuciones son bienvenidas. Para cambios grandes, abre primero un issue para discutir la propuesta.

Flujo recomendado:

1. Fork del repositorio.
2. Crea una rama de feature.
3. Ejecuta `go fmt ./...` y `go test ./...`.
4. Abre un Pull Request con descripcion clara.

## Licencia

Este proyecto esta bajo licencia MIT. Revisa el archivo `LICENSE`.
