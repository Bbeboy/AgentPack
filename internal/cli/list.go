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
		Short: t("list.short"),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			packagesRoot, err := storage.PackagesRoot()
			if err != nil {
				return err
			}

			entries, err := os.ReadDir(packagesRoot)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintln(cmd.OutOrStdout(), out("list.empty"))
					return nil
				}
				return fmt.Errorf(t("list.read", err))
			}

			packages := make([]string, 0)
			for _, entry := range entries {
				if entry.IsDir() {
					packages = append(packages, entry.Name())
				}
			}

			if len(packages) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), out("list.empty"))
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), out("list.title", len(packages)))
			for _, name := range packages {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", name)
			}

			return nil
		},
	}
}
