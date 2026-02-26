package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const defaultLanguage = "en"

type Settings struct {
	Language string `json:"language"`
}

func IsSupportedLanguage(lang string) bool {
	return lang == "en" || lang == "es"
}

func LoadSettings() (Settings, error) {
	path, err := settingsPath()
	if err != nil {
		return Settings{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Settings{Language: defaultLanguage}, nil
		}
		return Settings{}, fmt.Errorf("could not read settings: %w", err)
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return Settings{}, fmt.Errorf("could not parse settings: %w", err)
	}

	if !IsSupportedLanguage(settings.Language) {
		settings.Language = defaultLanguage
	}

	return settings, nil
}

func SaveLanguage(lang string) error {
	if !IsSupportedLanguage(lang) {
		return fmt.Errorf("unsupported language '%s' (allowed: en, es)", lang)
	}

	path, err := settingsPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	settings := Settings{Language: lang}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("could not encode settings: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("could not write settings: %w", err)
	}

	return nil
}

func settingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not resolve home directory: %w", err)
	}

	return filepath.Join(home, ".agentpack", "config.json"), nil
}
