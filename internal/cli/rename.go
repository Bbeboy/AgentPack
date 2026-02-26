package cli

import (
	"fmt"
	"os"

	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newRenameCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rename <current-name> <new-name>",
		Short: t("rename.short"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			currentName := args[0]
			newName := args[1]

			currentPath, err := storage.PackagePath(currentName)
			if err != nil {
				return err
			}

			newPath, err := storage.PackagePath(newName)
			if err != nil {
				return err
			}

			if currentName == newName {
				return fmt.Errorf(t("rename.same", currentName))
			}

			info, err := os.Stat(currentPath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf(t("rename.current.missing", currentName, currentPath))
				}
				return fmt.Errorf(t("rename.current.read", err))
			}
			if !info.IsDir() {
				return fmt.Errorf(t("rename.current.notdir", currentName))
			}

			if _, err := os.Stat(newPath); err == nil {
				return fmt.Errorf(t("rename.target.exists", newName))
			} else if !os.IsNotExist(err) {
				return fmt.Errorf(t("rename.target.read", err))
			}

			fmt.Fprintln(cmd.OutOrStdout(), out("rename.start", currentName, newName))
			if err := os.Rename(currentPath, newPath); err != nil {
				return fmt.Errorf(t("rename.fail", currentName, newName, err))
			}

			fmt.Fprintln(cmd.OutOrStdout(), out("rename.done", currentName, newName))
			return nil
		},
	}
}
