# AgentPack

CLI en Go para crear, guardar e instalar paquetes reutilizables de skills para agentes de IA.

Con AgentPack puedes tomar las skills de un proyecto, empaquetarlas localmente y volver a instalarlas en cualquier otro proyecto con un solo comando.

## Caracteristicas principales

- Crea paquetes desde una ruta de skills existente.
- Soporta deteccion de ruta cuando usas `.` como origen.
- Instala en `.agents/skills` del proyecto actual.
- Detecta conflictos por nombre de skill y pregunta si sobrescribir o ignorar.
- Mantiene archivos existentes que no estan en conflicto.
- Permite borrar paquetes con confirmacion o de forma forzada.
- Permite listar skills por paquete y eliminar una skill especifica.
- Incluye modo seguro `--dry-run` para revisar borrados sin ejecutar.
- Soporta autocompletado para `bash`, `zsh`, `fish` y `powershell`.

## Tabla de contenido

- [Requisitos](#requisitos)
- [Instalacion](#instalacion)
  - [Opcion 1: go install](#opcion-1-go-install)
  - [Opcion 2: build desde fuente](#opcion-2-build-desde-fuente)
  - [Usar agentpack globalmente](#usar-agentpack-globalmente)
  - [Autocompletado (Fish)](#autocompletado-fish)
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

## Requisitos

- Go `1.23+`
- Linux (soporte actual)

Nota: el codigo ya usa `filepath` y `os.UserHomeDir`, por lo que esta preparado para evolucionar a multiplataforma.

## Instalacion

### Opcion 1: go install

Cuando publiques en GitHub, instala asi:

```bash
go install github.com/tu-usuario/agentpack/cmd/agentpack@latest
```

### Opcion 2: build desde fuente

```bash
git clone https://github.com/tu-usuario/agentpack.git
cd agentpack
go mod tidy
go build -o agentpack ./cmd/agentpack
```

### Usar agentpack globalmente

Para poder ejecutar `agentpack` desde cualquier carpeta:

```bash
mkdir -p ~/.local/bin
cp ./agentpack ~/.local/bin/agentpack
chmod +x ~/.local/bin/agentpack
```

En muchos sistemas Linux, `~/.local/bin` ya esta en `PATH`. Si no lo esta:

- En Bash/Zsh, agrega `~/.local/bin` al `PATH` en tu archivo de shell.
- En Fish, usa:

```fish
fish_add_path ~/.local/bin
```

### Autocompletado (Fish)

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

## Conceptos clave

### Donde se guardan los paquetes

Los paquetes se guardan localmente en:

```text
~/.agentpack/packages-skills/<nombre-paquete>
```

### Donde se instalan las skills

`install` copia las skills a:

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
| `agentpack install <nombre-paquete>` | Instala un paquete en `.agents/skills` del proyecto actual. | `agentpack install backend-base` |
| `agentpack list` | Lista paquetes guardados localmente. | `agentpack list` |
| `agentpack list-skills <nombre-paquete>` | Lista las skills dentro de un paquete. | `agentpack list-skills backend-base` |
| `agentpack remove <nombre-paquete>` | Elimina un paquete guardado (con confirmacion). | `agentpack remove backend-base` |
| `agentpack remove <nombre-paquete> --force` | Elimina un paquete sin preguntar confirmacion. | `agentpack remove backend-base -f` |
| `agentpack remove <nombre-paquete> --dry-run` | Simula la eliminacion de un paquete sin borrar. | `agentpack remove backend-base --dry-run` |
| `agentpack remove-skill <nombre-paquete> <nombre-skill>` | Elimina una skill especifica de un paquete (con confirmacion). | `agentpack remove-skill backend-base docker` |
| `agentpack remove-skill <nombre-paquete> <nombre-skill> -f` | Elimina una skill sin preguntar confirmacion. | `agentpack remove-skill backend-base docker -f` |
| `agentpack remove-skill <nombre-paquete> <nombre-skill> --dry-run` | Simula la eliminacion de una skill sin borrar. | `agentpack remove-skill backend-base docker --dry-run` |
| `agentpack completion [bash\|zsh\|fish\|powershell]` | Genera script de autocompletado. | `agentpack completion fish` |

## Flujo completo con ejemplos

### 1) Crear paquete desde ruta explicita

```bash
agentpack create backend-base /mi/proyecto/.agents/skills
```

Salida esperada:

```text
[agentpack] Creando paquete 'backend-base'...
[agentpack] Origen: /mi/proyecto/.agents/skills
[agentpack] Destino: /home/tu-usuario/.agentpack/packages-skills/backend-base
[agentpack] Copiando skills...
[agentpack] Listo. Paquete creado: backend-base
[agentpack] Ruta: /home/tu-usuario/.agentpack/packages-skills/backend-base
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
[agentpack] Instalando paquete 'backend-base'...
[agentpack] Paquete: /home/tu-usuario/.agentpack/packages-skills/backend-base
[agentpack] Destino: /ruta/proyecto/.agents/skills
[agentpack] La carpeta destino no existe. Creando...
[agentpack] Instaladas: 4
[agentpack] Sobrescritas: 0
[agentpack] Ignoradas: 0
[agentpack] Listo.
```

### 4) Instalar paquete con conflictos

Si ya existe una skill con el mismo nombre en `.agents/skills`, se considera conflicto.

Ejemplo:

```text
[agentpack] Detectando conflictos...
[agentpack] Conflicto: la skill 'docker' ya existe en .agents/skills/docker
Sobrescribir la skill 'docker'? [y/N]: y
[agentpack] 'docker' sobrescrita (solo archivos en conflicto).
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
[agentpack] Skills del paquete 'backend-base' (3):
- docker
- golang-api
- testing
[agentpack] Ruta: /home/tu-usuario/.agentpack/packages-skills/backend-base
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

## Estructura del proyecto

```text
cmd/
  agentpack/
    main.go                # entrypoint
internal/
  cli/                     # comandos Cobra
    root.go
    create.go
    install.go
    list.go
    list_skills.go
    remove.go
    remove_skill.go
    completion.go
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

Si ejecutas `install`, `remove`, `list-skills` o `remove-skill` y no existe el paquete, verifica:

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

### Conflictos al instalar

El conflicto se resuelve por skill (carpeta). Puedes:

- sobrescribir (`y`) para merge con overwrite en conflicto,
- ignorar (`n`) para mantener la skill actual intacta.

## Roadmap

- Validacion opcional de `SKILL.md` (frontmatter y convenciones).
- Comando para renombrar skills dentro de un paquete.
- Mejoras de soporte multiplataforma (Windows/macOS).
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
