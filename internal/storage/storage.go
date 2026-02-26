package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Bbeboy/AgentPack/internal/i18n"
)

var packageNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,63}$`)
var skillNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,63}$`)

func PackagesRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		lang := i18n.ResolveLanguage()
		return "", fmt.Errorf(i18n.Message(lang, "storage.home", err))
	}
	return filepath.Join(home, ".agentpack", "packages-skills"), nil
}

func PackagePath(name string) (string, error) {
	if !packageNamePattern.MatchString(name) {
		lang := i18n.ResolveLanguage()
		return "", fmt.Errorf(i18n.Message(lang, "storage.package.invalid"))
	}

	root, err := PackagesRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(root, name), nil
}

func SkillPath(packageName, skillName string) (string, error) {
	packagePath, err := PackagePath(packageName)
	if err != nil {
		return "", err
	}

	if !skillNamePattern.MatchString(skillName) {
		lang := i18n.ResolveLanguage()
		return "", fmt.Errorf(i18n.Message(lang, "storage.skill.invalid"))
	}

	return filepath.Join(packagePath, skillName), nil
}
