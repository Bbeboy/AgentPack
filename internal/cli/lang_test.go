package cli

import (
	"strings"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/config"
)

func TestLangCommandUpdatesCurrentLanguageAndConfig(t *testing.T) {
	setupCLITest(t)

	cmd := newLangCmd()
	cmd.SetArgs([]string{"es"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("lang command returned error: %v", err)
	}

	if currentLang != "es" {
		t.Fatalf("expected currentLang es, got %q", currentLang)
	}

	settings, err := config.LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings returned error: %v", err)
	}
	if settings.Language != "es" {
		t.Fatalf("expected persisted language es, got %q", settings.Language)
	}
}

func TestLangCommandRejectsInvalidLanguage(t *testing.T) {
	setupCLITest(t)

	cmd := newLangCmd()
	cmd.SetArgs([]string{"pt"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected invalid language error")
	}
	if !strings.Contains(err.Error(), "invalid") && !strings.Contains(err.Error(), "invalido") {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}
