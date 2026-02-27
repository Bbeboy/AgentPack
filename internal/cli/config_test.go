package cli

import (
	"strings"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/config"
)

func TestConfigSetLanguageCommandPersistsValue(t *testing.T) {
	setupCLITest(t)

	cmd := newConfigCmd()
	cmd.SetArgs([]string{"set", "language", "es"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("config set command returned error: %v", err)
	}

	settings, err := config.LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings returned error: %v", err)
	}
	if settings.Language != "es" {
		t.Fatalf("expected persisted language es, got %q", settings.Language)
	}
}

func TestConfigSetRejectsUnknownKey(t *testing.T) {
	setupCLITest(t)

	cmd := newConfigCmd()
	cmd.SetArgs([]string{"set", "unknown", "value"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected unknown key error")
	}
	if !strings.Contains(err.Error(), "unknown") && !strings.Contains(err.Error(), "desconocida") {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}
