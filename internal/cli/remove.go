package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Bbeboy/AgentPack/internal/prompt"
	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newRemoveCmd() *cobra.Command {
	var force bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "remove <nombre-paquete>",
		Short: "Elimina un paquete guardado",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]

			packagePath, err := storage.PackagePath(packageName)
			if err != nil {
				return err
			}

			info, err := os.Stat(packagePath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("no se encontro el paquete '%s'\n[agentpack] Buscado en: %s", packageName, packagePath)
				}
				return fmt.Errorf("no se pudo leer el paquete: %w", err)
			}

			if !info.IsDir() {
				return fmt.Errorf("el paquete '%s' no es un directorio valido", packageName)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Paquete objetivo: %s\n", packagePath)

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Dry-run: se eliminaria el paquete '%s'.\n", packageName)
				return nil
			}

			if !force {
				reader := bufio.NewReader(cmd.InOrStdin())
				confirm, err := prompt.YesNo(reader, cmd.OutOrStdout(), fmt.Sprintf("Eliminar el paquete '%s' y todo su contenido?", packageName))
				if err != nil {
					return err
				}
				if !confirm {
					fmt.Fprintln(cmd.OutOrStdout(), "[agentpack] Operacion cancelada.")
					return nil
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Eliminando paquete '%s'...\n", packageName)
			if err := os.RemoveAll(packagePath); err != nil {
				return fmt.Errorf("no se pudo eliminar el paquete '%s': %w", packageName, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Paquete eliminado: %s\n", packageName)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Elimina el paquete sin pedir confirmacion")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Muestra que se eliminaria sin borrar nada")

	return cmd
}
