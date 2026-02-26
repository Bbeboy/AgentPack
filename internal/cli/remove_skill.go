package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Bbeboy/AgentPack/internal/prompt"
	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newRemoveSkillCmd() *cobra.Command {
	var force bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "remove-skill <package-name> <skill-name>",
		Short: t("removeskill.short"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]
			skillName := args[1]

			skillPath, err := storage.SkillPath(packageName, skillName)
			if err != nil {
				return err
			}

			skillInfo, err := os.Stat(skillPath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf(t("remove.skill.missing", skillName, packageName, skillPath))
				}
				return fmt.Errorf(t("remove.skill.read", skillName, err))
			}

			if !skillInfo.IsDir() {
				return fmt.Errorf(t("remove.skill.notdir", skillName, packageName))
			}

			if dryRun {
				fmt.Fprintln(cmd.OutOrStdout(), out("remove.skill.dryrun", skillName, packageName))
				return nil
			}

			if !force {
				reader := bufio.NewReader(cmd.InOrStdin())
				question := t("remove.skill.ask", skillName, packageName)
				confirm, err := prompt.YesNo(reader, cmd.OutOrStdout(), question)
				if err != nil {
					return err
				}
				if !confirm {
					fmt.Fprintln(cmd.OutOrStdout(), out("remove.cancelled"))
					return nil
				}
			}

			fmt.Fprintln(cmd.OutOrStdout(), out("remove.skill.start", skillName, packageName))
			if err := os.RemoveAll(skillPath); err != nil {
				return fmt.Errorf(t("remove.skill.fail", skillName, packageName, err))
			}

			fmt.Fprintln(cmd.OutOrStdout(), out("remove.skill.done", skillName))
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, t("flag.force"))
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, t("flag.dryrun"))

	return cmd
}
