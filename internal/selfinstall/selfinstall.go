package selfinstall

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bbeboy/AgentPack/internal/i18n"
)

type Result struct {
	Triggered bool
	Target    string
}

type Dependencies struct {
	LookupPath  func(string) (string, error)
	ReplaceFile func(string, string) error
}

func MaybeInstall(lang string, args []string, executablePath, runtimeOS, currentVersion string, deps Dependencies) (Result, error) {
	if len(args) > 0 {
		return Result{}, nil
	}

	if strings.TrimSpace(currentVersion) == "" || currentVersion == "dev" {
		return Result{}, nil
	}

	baseName := filepath.Base(executablePath)
	if !isReleaseArtifactName(baseName) && !matchesInstalledBinaryName(baseName, runtimeOS) {
		return Result{}, nil
	}

	binaryName := installedBinaryName(runtimeOS)
	targetPath, err := deps.LookupPath(binaryName)
	if err != nil {
		return Result{Triggered: true}, fmt.Errorf(i18n.Message(lang, "selfinstall.target.notfound", binaryName))
	}

	if samePath(runtimeOS, executablePath, targetPath) {
		return Result{}, nil
	}

	if err := deps.ReplaceFile(executablePath, targetPath); err != nil {
		return Result{Triggered: true}, fmt.Errorf(i18n.Message(lang, "selfinstall.replace.fail", targetPath, err))
	}

	return Result{Triggered: true, Target: targetPath}, nil
}

func ReplaceFile(sourcePath, targetPath string) error {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("inspect source binary: %w", err)
	}

	in, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source binary: %w", err)
	}
	defer in.Close()

	tempFile, err := os.CreateTemp(filepath.Dir(targetPath), ".agentpack-self-install-*")
	if err != nil {
		return fmt.Errorf("create temp file in target directory: %w", err)
	}
	tempPath := tempFile.Name()

	defer func() {
		_ = os.Remove(tempPath)
	}()

	if _, err := io.Copy(tempFile, in); err != nil {
		tempFile.Close()
		return fmt.Errorf("copy binary content: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Chmod(tempPath, sourceInfo.Mode().Perm()); err != nil {
		return fmt.Errorf("set executable permissions: %w", err)
	}

	if err := replaceWithRename(tempPath, targetPath); err != nil {
		return fmt.Errorf("replace installed binary: %w", err)
	}

	return nil
}

func replaceWithRename(tempPath, targetPath string) error {
	if err := os.Rename(tempPath, targetPath); err == nil {
		return nil
	}

	if err := os.Remove(targetPath); err != nil {
		return err
	}

	if err := os.Rename(tempPath, targetPath); err != nil {
		return err
	}

	return nil
}

func installedBinaryName(runtimeOS string) string {
	if runtimeOS == "windows" {
		return "agentpack.exe"
	}
	return "agentpack"
}

func matchesInstalledBinaryName(baseName, runtimeOS string) bool {
	return strings.EqualFold(baseName, installedBinaryName(runtimeOS))
}

func samePath(runtimeOS, leftPath, rightPath string) bool {
	leftAbs, leftErr := filepath.Abs(leftPath)
	rightAbs, rightErr := filepath.Abs(rightPath)
	if leftErr != nil || rightErr != nil {
		return false
	}

	leftClean := filepath.Clean(leftAbs)
	rightClean := filepath.Clean(rightAbs)

	if runtimeOS == "windows" {
		return strings.EqualFold(leftClean, rightClean)
	}

	return leftClean == rightClean
}

func isReleaseArtifactName(name string) bool {
	lower := strings.ToLower(name)
	hasExe := strings.HasSuffix(lower, ".exe")
	if hasExe {
		lower = strings.TrimSuffix(lower, ".exe")
	}

	if !strings.HasPrefix(lower, "agentpack-v") {
		return false
	}

	if !strings.HasSuffix(lower, "-amd64") && !strings.HasSuffix(lower, "-arm64") {
		return false
	}

	return strings.Contains(lower, "-linux-") || strings.Contains(lower, "-darwin-") || strings.Contains(lower, "-windows-")
}
