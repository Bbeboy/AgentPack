package config

import (
	"path/filepath"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/testutil"
)

func TestLoadSettingsReturnsDefaultWhenMissing(t *testing.T) {
	testutil.SetupHome(t)

	settings, err := LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings returned error: %v", err)
	}
	if settings.Language != "en" {
		t.Fatalf("expected default language en, got %q", settings.Language)
	}
}

func TestSaveAndLoadLanguage(t *testing.T) {
	testutil.SetupHome(t)

	if err := SaveLanguage("es"); err != nil {
		t.Fatalf("SaveLanguage returned error: %v", err)
	}

	settings, err := LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings returned error: %v", err)
	}
	if settings.Language != "es" {
		t.Fatalf("expected language es, got %q", settings.Language)
	}
}

func TestLoadSettingsFallsBackWhenLanguageIsInvalid(t *testing.T) {
	home := testutil.SetupHome(t)
	testutil.WriteFile(t, filepath.Join(home, ".agentpack", "config.json"), `{"language":"fr"}`)

	settings, err := LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings returned error: %v", err)
	}
	if settings.Language != "en" {
		t.Fatalf("expected fallback language en, got %q", settings.Language)
	}
}

func TestSaveLanguageRejectsUnsupportedLanguage(t *testing.T) {
	testutil.SetupHome(t)

	if err := SaveLanguage("pt"); err == nil {
		t.Fatal("expected error for unsupported language")
	}
}
