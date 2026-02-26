package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:          "agentpack",
	Short:        "Gestiona paquetes de skills para agentes de IA",
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
	rootCmd.AddCommand(newRemoveCmd())
	rootCmd.AddCommand(newRemoveSkillCmd())
	rootCmd.AddCommand(newCompletionCmd(rootCmd))
}
