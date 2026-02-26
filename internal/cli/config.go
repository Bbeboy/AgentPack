package cli

import (
	"fmt"

	"github.com/Bbeboy/AgentPack/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: t("config.short"),
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(newConfigSetCmd())

	return cmd
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: t("config.set.short"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			switch key {
			case "language":
				if !config.IsSupportedLanguage(value) {
					return fmt.Errorf(t("language.invalid", value))
				}

				if err := config.SaveLanguage(value); err != nil {
					return err
				}

				currentLang = value
				fmt.Fprintln(cmd.OutOrStdout(), out("language.set", value))
				return nil
			default:
				return fmt.Errorf(t("config.key.unknown", key))
			}
		},
	}
}
