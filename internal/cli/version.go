package cli

import (
	"fmt"

	"github.com/Bbeboy/AgentPack/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: t("version.short"),
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), out("version.output", version.Value()))
		},
	}
}
