package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Bbeboy/AgentPack/internal/fsutil"
	"github.com/Bbeboy/AgentPack/internal/prompt"
	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <nombre-paquete> <ruta-skills>",
		Short: "Crea un paquete de skills desde una ruta",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]
			sourceArg := args[1]

			sourcePath, err := resolveCreateSource(sourceArg, cmd)
			if err != nil {
				return err
			}

			sourcePath, err = filepath.Abs(sourcePath)
			if err != nil {
				return fmt.Errorf("no se pudo resolver la ruta origen: %w", err)
			}

			sourceInfo, err := os.Stat(sourcePath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("la ruta origen no existe: %s", sourcePath)
				}
				return fmt.Errorf("no se pudo leer la ruta origen: %w", err)
			}
			if !sourceInfo.IsDir() {
				return fmt.Errorf("la ruta origen no es un directorio: %s", sourcePath)
			}

			packagesRoot, err := storage.PackagesRoot()
			if err != nil {
				return err
			}

			if err := os.MkdirAll(packagesRoot, 0o755); err != nil {
				return fmt.Errorf("no se pudo crear el directorio base de paquetes: %w", err)
			}

			packagePath, err := storage.PackagePath(packageName)
			if err != nil {
				return err
			}
			if _, err := os.Stat(packagePath); err == nil {
				return fmt.Errorf("ya existe un paquete con el nombre '%s'", packageName)
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("no se pudo validar el paquete destino: %w", err)
			}

			if err := os.MkdirAll(packagePath, 0o755); err != nil {
				return fmt.Errorf("no se pudo crear el directorio del paquete: %w", err)
			}

			created := true
			defer func() {
				if created {
					_ = os.RemoveAll(packagePath)
				}
			}()

			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Creando paquete '%s'...\n", packageName)
			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Origen: %s\n", sourcePath)
			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Destino: %s\n", packagePath)
			fmt.Fprintln(cmd.OutOrStdout(), "[agentpack] Copiando skills...")

			if err := fsutil.CopyDirContents(sourcePath, packagePath); err != nil {
				return fmt.Errorf("no se pudo copiar el contenido de skills: %w", err)
			}

			created = false

			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Listo. Paquete creado: %s\n", packageName)
			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Ruta: %s\n", packagePath)
			return nil
		},
	}
}

func resolveCreateSource(sourceArg string, cmd *cobra.Command) (string, error) {
	if sourceArg != "." {
		return sourceArg, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio actual: %w", err)
	}

	candidates := []string{
		".agents/skills",
		".opencode/skills",
		".agent/skills",
		"skills",
	}

	found := make([]string, 0, len(candidates))
	for _, rel := range candidates {
		full := filepath.Join(cwd, rel)
		info, err := os.Stat(full)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", fmt.Errorf("no se pudo revisar la ruta candidata '%s': %w", rel, err)
		}
		if info.IsDir() {
			found = append(found, rel)
		}
	}

	if len(found) == 0 {
		return "", fmt.Errorf("con '.' no se encontro ninguna carpeta de skills en este proyecto (buscado: .agents/skills, .opencode/skills, .agent/skills, skills)")
	}

	if len(found) == 1 {
		selected := filepath.Join(cwd, found[0])
		fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Carpeta detectada automaticamente: %s\n", found[0])
		return selected, nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), "[agentpack] Se encontraron multiples carpetas de skills:")
	for i, option := range found {
		fmt.Fprintf(cmd.OutOrStdout(), "  %d) %s\n", i+1, option)
	}

	reader := bufio.NewReader(cmd.InOrStdin())
	idx, err := prompt.SelectIndex(reader, cmd.OutOrStdout(), len(found))
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, found[idx]), nil
}
