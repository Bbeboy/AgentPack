package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var packageNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,63}$`)
var skillNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,63}$`)

func PackagesRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener HOME: %w", err)
	}
	return filepath.Join(home, ".agentpack", "packages-skills"), nil
}

func PackagePath(name string) (string, error) {
	if !packageNamePattern.MatchString(name) {
		return "", fmt.Errorf("nombre de paquete invalido: usa solo letras, numeros, '.', '_' o '-' (max 64 caracteres)")
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
		return "", fmt.Errorf("nombre de skill invalido: usa solo letras, numeros, '.', '_' o '-' (max 64 caracteres)")
	}

	return filepath.Join(packagePath, skillName), nil
}
