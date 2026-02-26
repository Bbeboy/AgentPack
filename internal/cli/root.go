package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:          "agentpack",
	Short:        t("root.short"),
	SilenceUsage: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
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
	rootCmd.AddCommand(newCompletionCmd(rootCmd))
}
