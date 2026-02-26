package cli

import (
	"fmt"

	"github.com/Bbeboy/AgentPack/internal/config"
	"github.com/spf13/cobra"
)

func newLangCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lang <en|es>",
		Short: t("lang.short"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			lang := args[0]
			if !config.IsSupportedLanguage(lang) {
				return fmt.Errorf(t("language.invalid", lang))
			}

			if err := config.SaveLanguage(lang); err != nil {
				return err
			}

			currentLang = lang
			fmt.Fprintln(cmd.OutOrStdout(), out("language.set", lang))
			return nil
		},
	}
}
