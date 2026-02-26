package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bbeboy/AgentPack/internal/fsutil"
	"github.com/Bbeboy/AgentPack/internal/platform"
	"github.com/Bbeboy/AgentPack/internal/prompt"
	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <package-name> <skills-path>",
		Short: t("create.short"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]
			sourceArg := args[1]

			sourcePath, err := resolveCreateSource(sourceArg, cmd)
			if err != nil {
				return err
			}

			sourcePath, err = filepath.Abs(sourcePath)
			if err != nil {
				return fmt.Errorf(t("create.path.resolve", err))
			}

			sourceInfo, err := os.Stat(sourcePath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf(t("create.source.missing", sourcePath))
				}
				return fmt.Errorf(t("add.source.read", err))
			}
			if !sourceInfo.IsDir() {
				return fmt.Errorf(t("create.source.notdir", sourcePath))
			}

			packagesRoot, err := storage.PackagesRoot()
			if err != nil {
				return err
			}

			if err := os.MkdirAll(packagesRoot, 0o755); err != nil {
				return fmt.Errorf(t("create.root.create", err))
			}

			packagePath, err := storage.PackagePath(packageName)
			if err != nil {
				return err
			}
			if _, err := os.Stat(packagePath); err == nil {
				return fmt.Errorf(t("create.exists", packageName))
			} else if !os.IsNotExist(err) {
				return fmt.Errorf(t("create.dest.validate", err))
			}

			if err := os.MkdirAll(packagePath, 0o755); err != nil {
				return fmt.Errorf(t("create.dest.create", err))
			}

			created := true
			defer func() {
				if created {
					_ = os.RemoveAll(packagePath)
				}
			}()

			fmt.Fprintln(cmd.OutOrStdout(), out("create.start", packageName))

			if err := fsutil.CopyDirContents(sourcePath, packagePath); err != nil {
				return fmt.Errorf(t("create.copy", err))
			}

			created = false

			fmt.Fprintln(cmd.OutOrStdout(), out("create.done", packagePath))
			return nil
		},
	}
}

func resolveCreateSource(sourceArg string, cmd *cobra.Command) (string, error) {
	if sourceArg != "." {
		return sourceArg, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf(t("create.cwd", err))
	}

	candidates := platform.CandidateSkillsPaths()

	found := make([]string, 0, len(candidates))
	for _, rel := range candidates {
		full := filepath.Join(cwd, rel)
		info, err := os.Stat(full)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", fmt.Errorf(t("create.candidate.check", rel, err))
		}
		if info.IsDir() {
			found = append(found, rel)
		}
	}

	if len(found) == 0 {
		return "", fmt.Errorf(t("create.no.candidates", strings.Join(candidates, ", ")))
	}

	if len(found) == 1 {
		selected := filepath.Join(cwd, found[0])
		fmt.Fprintln(cmd.OutOrStdout(), out("create.detected", found[0]))
		return selected, nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), out("create.multiple"))
	for i, option := range found {
		fmt.Fprintf(cmd.OutOrStdout(), "  %d) %s\n", i+1, option)
	}

	reader := bufio.NewReader(cmd.InOrStdin())
	idx, err := prompt.SelectIndex(reader, cmd.OutOrStdout(), len(found))
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, found[idx]), nil
}
