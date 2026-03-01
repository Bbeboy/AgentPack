package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newLegacyAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:                "add",
		Hidden:             true,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf(t("add.legacy.deprecated"))
		},
	}
}
