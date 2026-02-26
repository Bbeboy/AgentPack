package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Bbeboy/AgentPack/internal/fsutil"
	"github.com/Bbeboy/AgentPack/internal/platform"
	"github.com/Bbeboy/AgentPack/internal/prompt"
	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install <package-name>",
		Short: t("install.short"),
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
					return fmt.Errorf(t("install.pkg.missing", packageName, packagePath))
				}
				return fmt.Errorf(t("install.pkg.read", err))
			}
			if !packageInfo.IsDir() {
				return fmt.Errorf(t("install.pkg.notdir", packageName))
			}

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf(t("install.cwd", err))
			}

			destinationRoot, _, _, err := platform.ResolveSkillsDestination(cwd)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), out("install.start", packageName, destinationRoot))

			if _, err := os.Stat(destinationRoot); os.IsNotExist(err) {
				if err := os.MkdirAll(destinationRoot, 0o755); err != nil {
					return fmt.Errorf(t("install.dest.create", err))
				}
			} else if err != nil {
				return fmt.Errorf(t("install.dest.validate", err))
			}

			entries, err := os.ReadDir(packagePath)
			if err != nil {
				return fmt.Errorf(t("install.list", err))
			}

			installed := 0
			overwritten := 0
			skipped := 0

			conflicts := make([]string, 0)
			skills := make([]os.DirEntry, 0)
			for _, entry := range entries {
				if !entry.IsDir() {
					fmt.Fprintln(cmd.OutOrStdout(), out("install.skip.non.dir", entry.Name()))
					continue
				}

				skills = append(skills, entry)
				target := filepath.Join(destinationRoot, entry.Name())
				if _, err := os.Stat(target); err == nil {
					conflicts = append(conflicts, entry.Name())
				} else if !os.IsNotExist(err) {
					return fmt.Errorf(t("install.conflict.check", entry.Name(), err))
				}
			}

			if len(conflicts) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), out("install.conflict.title"))
				for _, conflict := range conflicts {
					fmt.Fprintln(cmd.OutOrStdout(), out("install.conflict.item", conflict))
				}
			}

			reader := bufio.NewReader(cmd.InOrStdin())
			for _, skill := range skills {
				sourceSkillPath := filepath.Join(packagePath, skill.Name())
				targetSkillPath := filepath.Join(destinationRoot, skill.Name())

				targetInfo, statErr := os.Stat(targetSkillPath)
				targetExists := statErr == nil
				if statErr != nil && !os.IsNotExist(statErr) {
					return fmt.Errorf(t("install.target.check", skill.Name(), statErr))
				}

				if targetExists {
					overwrite, err := prompt.YesNo(reader, cmd.OutOrStdout(), t("install.overwrite.ask", skill.Name()))
					if err != nil {
						return err
					}
					if !overwrite {
						skipped++
						fmt.Fprintln(cmd.OutOrStdout(), out("install.skip.skill", skill.Name()))
						continue
					}

					if targetInfo.IsDir() {
						if err := fsutil.MergeDir(sourceSkillPath, targetSkillPath); err != nil {
							return fmt.Errorf(t("install.overwrite.fail", skill.Name(), err))
						}
					} else {
						if err := os.Remove(targetSkillPath); err != nil {
							return fmt.Errorf(t("install.replace.fail", skill.Name(), err))
						}
						if err := fsutil.CopyDir(sourceSkillPath, targetSkillPath); err != nil {
							return fmt.Errorf(t("install.overwrite.fail", skill.Name(), err))
						}
					}
					overwritten++
					fmt.Fprintln(cmd.OutOrStdout(), out("install.overwrite", skill.Name()))
					continue
				}

				if err := fsutil.CopyDir(sourceSkillPath, targetSkillPath); err != nil {
					return fmt.Errorf(t("install.skill.fail", skill.Name(), err))
				}
				installed++
			}

			fmt.Fprintln(cmd.OutOrStdout(), out("install.summary", installed, overwritten, skipped))
			return nil
		},
	}
}
