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
		Use:   "list-skills <package-name>",
		Short: t("listskills.short"),
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
					return fmt.Errorf(t("remove.pkg.missing", packageName, packagePath))
				}
				return fmt.Errorf(t("remove.pkg.read", err))
			}

			if !packageInfo.IsDir() {
				return fmt.Errorf(t("remove.pkg.notdir", packageName))
			}

			entries, err := os.ReadDir(packagePath)
			if err != nil {
				return fmt.Errorf(t("install.list", err))
			}

			skills := make([]string, 0)
			for _, entry := range entries {
				if entry.IsDir() {
					skills = append(skills, entry.Name())
				}
			}

			if len(skills) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), out("listskills.empty", packageName))
				return nil
			}

			sort.Strings(skills)
			fmt.Fprintln(cmd.OutOrStdout(), out("listskills.title", packageName, len(skills)))
			for _, skill := range skills {
				fmt.Fprintf(cmd.OutOrStdout(), "- %s\n", skill)
			}

			return nil
		},
	}
}
