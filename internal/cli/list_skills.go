package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newListSkillsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-skills <nombre-paquete>",
		Short: "Lista las skills de un paquete",
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

			entries, err := os.ReadDir(packagePath)
			if err != nil {
				return fmt.Errorf("no se pudo listar el contenido del paquete '%s': %w", packageName, err)
			}

			skills := make([]string, 0)
			for _, entry := range entries {
				if entry.IsDir() {
					skills = append(skills, entry.Name())
				}
			}

			if len(skills) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] El paquete '%s' no tiene skills.\n", packageName)
				fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Ruta: %s\n", packagePath)
				return nil
			}

			sort.Strings(skills)
			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Skills del paquete '%s' (%d):\n", packageName, len(skills))
			for _, skill := range skills {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", skill)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Ruta: %s\n", packagePath)

			return nil
		},
	}
}
