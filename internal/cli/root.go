package cli

import (
	"fmt"

	"github.com/Bbeboy/AgentPack/internal/version"
	"github.com/spf13/cobra"
)

var rootShowVersion bool

var rootCmd = &cobra.Command{
	Use:          "agentpack",
	Short:        t("root.short"),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if rootShowVersion {
			fmt.Fprintln(cmd.OutOrStdout(), out("version.output", version.Value()))
			return nil
		}
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVarP(&rootShowVersion, "version", "v", false, t("version.short"))

	rootCmd.AddCommand(newCreateCmd())
	rootCmd.AddCommand(newInstallCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newListSkillsCmd())
	rootCmd.AddCommand(newRenameCmd())
	rootCmd.AddCommand(newRemoveCmd())
	rootCmd.AddCommand(newRemoveSkillCmd())
	rootCmd.AddCommand(newAddCmd())
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newLangCmd())
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newCompletionCmd(rootCmd))
}
