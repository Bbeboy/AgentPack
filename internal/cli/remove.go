package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bbeboy/AgentPack/internal/prompt"
	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newRemoveCmd() *cobra.Command {
	var force bool
	var dryRun bool
	var fromPackage string

	cmd := &cobra.Command{
		Use:   "remove <package-or-path>",
		Short: t("remove.short"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			if fromPackage != "" {
				return removePathFromPackage(cmd, fromPackage, target, force, dryRun)
			}

			if containsPathSeparator(target) {
				return fmt.Errorf(t("remove.from.required"))
			}

			return removePackage(cmd, target, force, dryRun)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, t("flag.force"))
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, t("flag.dryrun"))
	cmd.Flags().StringVar(&fromPackage, "from", "", t("flag.from"))

	return cmd
}

func removePackage(cmd *cobra.Command, packageName string, force bool, dryRun bool) error {
	packagePath, err := storage.PackagePath(packageName)
	if err != nil {
		return err
	}

	info, err := os.Stat(packagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(t("remove.pkg.missing", packageName, packagePath))
		}
		return fmt.Errorf(t("remove.pkg.read", err))
	}

	if !info.IsDir() {
		return fmt.Errorf(t("remove.pkg.notdir", packageName))
	}

	if dryRun {
		fmt.Fprintln(cmd.OutOrStdout(), out("remove.pkg.dryrun", packageName))
		return nil
	}

	if !force {
		reader := bufio.NewReader(cmd.InOrStdin())
		confirm, err := prompt.YesNo(reader, cmd.OutOrStdout(), t("remove.pkg.ask", packageName))
		if err != nil {
			return err
		}
		if !confirm {
			fmt.Fprintln(cmd.OutOrStdout(), out("remove.cancelled"))
			return nil
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), out("remove.pkg.start", packageName))
	if err := os.RemoveAll(packagePath); err != nil {
		return fmt.Errorf(t("remove.pkg.fail", packageName, err))
	}

	fmt.Fprintln(cmd.OutOrStdout(), out("remove.pkg.done", packageName))
	return nil
}

func removePathFromPackage(cmd *cobra.Command, packageName string, relativePath string, force bool, dryRun bool) error {
	packagePath, err := storage.PackagePath(packageName)
	if err != nil {
		return err
	}

	if _, err := os.Stat(packagePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(t("remove.pkg.missing", packageName, packagePath))
		}
		return fmt.Errorf(t("remove.pkg.read", err))
	}

	clean, err := validateRelativePath(relativePath)
	if err != nil {
		return err
	}

	targetPath := filepath.Join(packagePath, clean)
	if _, err := os.Stat(targetPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(t("remove.path.notfound", clean, packageName))
		}
		return fmt.Errorf(t("remove.path.fail", clean, packageName, err))
	}

	if dryRun {
		fmt.Fprintln(cmd.OutOrStdout(), out("remove.path.dryrun", clean, packageName))
		return nil
	}

	if !force {
		reader := bufio.NewReader(cmd.InOrStdin())
		confirm, err := prompt.YesNo(reader, cmd.OutOrStdout(), t("remove.path.ask", clean, packageName))
		if err != nil {
			return err
		}
		if !confirm {
			fmt.Fprintln(cmd.OutOrStdout(), out("remove.cancelled"))
			return nil
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), out("remove.path.start", clean, packageName))
	if err := os.RemoveAll(targetPath); err != nil {
		return fmt.Errorf(t("remove.path.fail", clean, packageName, err))
	}

	fmt.Fprintln(cmd.OutOrStdout(), out("remove.path.done", clean, packageName))
	return nil
}

func validateRelativePath(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf(t("remove.path.empty"))
	}
	normalized := normalizePathSeparators(value)
	if filepath.IsAbs(normalized) {
		return "", fmt.Errorf(t("remove.path.absolute"))
	}

	clean := filepath.Clean(normalized)
	if clean == "." {
		return "", fmt.Errorf(t("remove.path.empty"))
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf(t("remove.path.escape", value))
	}

	return clean, nil
}
