package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Bbeboy/AgentPack/internal/fsutil"
	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export <package-name>",
		Short: t("export.short"),
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
					return fmt.Errorf(t("export.package.missing", packageName, packagePath))
				}
				return fmt.Errorf(t("export.package.read", err))
			}
			if !packageInfo.IsDir() {
				return fmt.Errorf(t("export.package.notdir", packageName))
			}

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf(t("export.cwd", err))
			}

			destinationPath := filepath.Join(cwd, packageName)
			if _, err := os.Stat(destinationPath); err == nil {
				return fmt.Errorf(t("export.destination.exists", destinationPath))
			} else if !os.IsNotExist(err) {
				return fmt.Errorf(t("export.destination.read", err))
			}

			if err := os.MkdirAll(destinationPath, 0o755); err != nil {
				return fmt.Errorf(t("export.destination.create", err))
			}

			exported := true
			defer func() {
				if exported {
					_ = os.RemoveAll(destinationPath)
				}
			}()

			fmt.Fprintln(cmd.OutOrStdout(), out("export.start", packageName, destinationPath))

			if err := fsutil.CopyDirContents(packagePath, destinationPath); err != nil {
				return fmt.Errorf(t("export.copy", packageName, err))
			}

			exported = false

			fmt.Fprintln(cmd.OutOrStdout(), out("export.done", packageName, destinationPath))
			return nil
		},
	}
}
