package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bbeboy/AgentPack/internal/fsutil"
	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var toPackage string

	cmd := &cobra.Command{
		Use:   "add-skill <file-or-folder>",
		Short: t("addskill.short"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sourceArg := args[0]
			sourcePath, err := filepath.Abs(sourceArg)
			if err != nil {
				return fmt.Errorf(t("add.source.read", err))
			}

			sourceInfo, err := os.Stat(sourcePath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf(t("add.source.missing", sourceArg))
				}
				return fmt.Errorf(t("add.source.read", err))
			}

			packagePath, err := storage.PackagePath(toPackage)
			if err != nil {
				return err
			}

			packageInfo, err := os.Stat(packagePath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf(t("add.package.missing", toPackage, packagePath))
				}
				return fmt.Errorf(t("add.package.read", err))
			}
			if !packageInfo.IsDir() {
				return fmt.Errorf(t("add.package.notdir", toPackage))
			}

			destinationRel, err := resolveAddDestinationRelative(sourceArg, sourcePath)
			if err != nil {
				return err
			}

			targetPath := filepath.Join(packagePath, destinationRel)
			fmt.Fprintln(cmd.OutOrStdout(), out("add.start", sourceArg, toPackage))

			if sourceInfo.IsDir() {
				if err := addDirectory(sourcePath, targetPath); err != nil {
					return fmt.Errorf(t("add.path.copy", sourceArg, toPackage, err))
				}
			} else {
				if err := addFile(sourcePath, targetPath); err != nil {
					return fmt.Errorf(t("add.path.copy", sourceArg, toPackage, err))
				}
			}

			fmt.Fprintln(cmd.OutOrStdout(), out("add.done", destinationRel, toPackage))
			return nil
		},
	}

	cmd.Flags().StringVar(&toPackage, "to", "", t("flag.to"))
	_ = cmd.MarkFlagRequired("to")

	return cmd
}

func addDirectory(sourcePath string, targetPath string) error {
	targetInfo, err := os.Stat(targetPath)
	if err == nil {
		if targetInfo.IsDir() {
			return fsutil.MergeDir(sourcePath, targetPath)
		}
		if removeErr := os.Remove(targetPath); removeErr != nil {
			return removeErr
		}
		return fsutil.CopyDir(sourcePath, targetPath)
	}

	if !os.IsNotExist(err) {
		return err
	}

	return fsutil.CopyDir(sourcePath, targetPath)
}

func addFile(sourcePath string, targetPath string) error {
	targetInfo, err := os.Stat(targetPath)
	if err == nil && targetInfo.IsDir() {
		return fmt.Errorf(t("add.path.conflict", sourcePath, targetPath))
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}

	in, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer in.Close()

	sourceInfo, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, sourceInfo.Mode().Perm())
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}

	return out.Close()
}

func resolveAddDestinationRelative(sourceArg string, sourcePath string) (string, error) {
	if filepath.IsAbs(sourceArg) {
		cwd, err := os.Getwd()
		if err == nil {
			rel, relErr := filepath.Rel(cwd, sourcePath)
			if relErr == nil {
				clean, cleanErr := cleanAddRelativePath(rel)
				if cleanErr == nil {
					return clean, nil
				}
			}
		}
		return filepath.Base(sourcePath), nil
	}

	return cleanAddRelativePath(sourceArg)
}

func cleanAddRelativePath(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf(t("add.path.empty"))
	}

	normalized := normalizePathSeparators(value)
	clean := filepath.Clean(normalized)
	if clean == "." {
		return "", fmt.Errorf(t("add.path.empty"))
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf(t("add.path.escape", value))
	}

	return clean, nil
}
