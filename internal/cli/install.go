package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"agentpack/internal/fsutil"
	"agentpack/internal/prompt"
	"agentpack/internal/storage"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install <nombre-paquete>",
		Short: "Instala un paquete de skills en el proyecto actual",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]

			packagePath, err := storage.PackagePath(packageName)
			if err != nil {
				return err
			}

			packageInfo, err := os.Stat(packagePath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("no se encontro el paquete '%s'\n[agentpack] Buscado en: %s", packageName, packagePath)
				}
				return fmt.Errorf("no se pudo leer el paquete: %w", err)
			}
			if !packageInfo.IsDir() {
				return fmt.Errorf("el paquete '%s' no es un directorio valido", packageName)
			}

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("no se pudo obtener el directorio actual: %w", err)
			}

			destinationRoot := filepath.Join(cwd, ".agents", "skills")

			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Instalando paquete '%s'...\n", packageName)
			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Paquete: %s\n", packagePath)
			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Destino: %s\n", destinationRoot)

			if _, err := os.Stat(destinationRoot); os.IsNotExist(err) {
				fmt.Fprintln(cmd.OutOrStdout(), "[agentpack] La carpeta destino no existe. Creando...")
				if err := os.MkdirAll(destinationRoot, 0o755); err != nil {
					return fmt.Errorf("no se pudo crear la carpeta destino: %w", err)
				}
			} else if err != nil {
				return fmt.Errorf("no se pudo validar la carpeta destino: %w", err)
			}

			entries, err := os.ReadDir(packagePath)
			if err != nil {
				return fmt.Errorf("no se pudo listar el contenido del paquete: %w", err)
			}

			installed := 0
			overwritten := 0
			skipped := 0

			conflicts := make([]string, 0)
			skills := make([]os.DirEntry, 0)
			for _, entry := range entries {
				if !entry.IsDir() {
					fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Aviso: se ignora '%s' porque no es una carpeta de skill.\n", entry.Name())
					continue
				}

				skills = append(skills, entry)
				target := filepath.Join(destinationRoot, entry.Name())
				if _, err := os.Stat(target); err == nil {
					conflicts = append(conflicts, entry.Name())
				} else if !os.IsNotExist(err) {
					return fmt.Errorf("no se pudo revisar conflicto para '%s': %w", entry.Name(), err)
				}
			}

			if len(conflicts) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "[agentpack] Detectando conflictos...")
				for _, conflict := range conflicts {
					fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Conflicto: la skill '%s' ya existe en .agents/skills/%s\n", conflict, conflict)
				}
			}

			reader := bufio.NewReader(cmd.InOrStdin())
			for _, skill := range skills {
				sourceSkillPath := filepath.Join(packagePath, skill.Name())
				targetSkillPath := filepath.Join(destinationRoot, skill.Name())

				targetInfo, statErr := os.Stat(targetSkillPath)
				targetExists := statErr == nil
				if statErr != nil && !os.IsNotExist(statErr) {
					return fmt.Errorf("no se pudo revisar la skill destino '%s': %w", skill.Name(), statErr)
				}

				if targetExists {
					overwrite, err := prompt.YesNo(reader, cmd.OutOrStdout(), fmt.Sprintf("Sobrescribir la skill '%s'?", skill.Name()))
					if err != nil {
						return err
					}
					if !overwrite {
						skipped++
						fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] '%s' ignorada.\n", skill.Name())
						continue
					}

					if targetInfo.IsDir() {
						if err := fsutil.MergeDir(sourceSkillPath, targetSkillPath); err != nil {
							return fmt.Errorf("no se pudo sobrescribir la skill '%s': %w", skill.Name(), err)
						}
					} else {
						if err := os.Remove(targetSkillPath); err != nil {
							return fmt.Errorf("no se pudo reemplazar '%s': %w", skill.Name(), err)
						}
						if err := fsutil.CopyDir(sourceSkillPath, targetSkillPath); err != nil {
							return fmt.Errorf("no se pudo sobrescribir la skill '%s': %w", skill.Name(), err)
						}
					}
					overwritten++
					fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] '%s' sobrescrita (solo archivos en conflicto).\n", skill.Name())
					continue
				}

				if err := fsutil.CopyDir(sourceSkillPath, targetSkillPath); err != nil {
					return fmt.Errorf("no se pudo instalar la skill '%s': %w", skill.Name(), err)
				}
				installed++
			}

			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Instaladas: %d\n", installed)
			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Sobrescritas: %d\n", overwritten)
			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Ignoradas: %d\n", skipped)
			fmt.Fprintln(cmd.OutOrStdout(), "[agentpack] Listo.")
			return nil
		},
	}
}
