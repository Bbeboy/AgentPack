package cli

import (
	"fmt"
	"os"

	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lista los paquetes de skills guardados",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			packagesRoot, err := storage.PackagesRoot()
			if err != nil {
				return err
			}

			entries, err := os.ReadDir(packagesRoot)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintln(cmd.OutOrStdout(), "[agentpack] No hay paquetes guardados.")
					fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Ruta: %s\n", packagesRoot)
					return nil
				}
				return fmt.Errorf("no se pudo leer la ruta de paquetes: %w", err)
			}

			packages := make([]string, 0)
			for _, entry := range entries {
				if entry.IsDir() {
					packages = append(packages, entry.Name())
				}
			}

			if len(packages) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "[agentpack] No hay paquetes guardados.")
				fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Ruta: %s\n", packagesRoot)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Paquetes disponibles (%d):\n", len(packages))
			for _, name := range packages {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", name)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Ruta: %s\n", packagesRoot)

			return nil
		},
	}
}
