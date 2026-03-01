package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/testutil"
)

func TestExportCommandCopiesPackageIntoCWD(t *testing.T) {
	home := setupCLITest(t)
	packagePath := testutil.EnsurePackageDir(t, home, "demo")
	testutil.WriteFile(t, filepath.Join(packagePath, "docker", "SKILL.md"), "name: docker\n")
	testutil.WriteFile(t, filepath.Join(packagePath, "README.md"), "package docs\n")

	projectDir := t.TempDir()
	withWorkingDir(t, projectDir)

	output, err := runCommand(newExportCmd(), []string{"demo"}, "")
	if err != nil {
		t.Fatalf("export command returned error: %v", err)
	}

	exportedSkill := filepath.Join(projectDir, "demo", "docker", "SKILL.md")
	data, readErr := os.ReadFile(exportedSkill)
	if readErr != nil {
		t.Fatalf("expected exported skill file: %v", readErr)
	}
	if string(data) != "name: docker\n" {
		t.Fatalf("unexpected exported skill content: %q", string(data))
	}

	if !strings.Contains(output, "exported") {
		t.Fatalf("expected export success output, got %q", output)
	}
}

func TestExportCommandFailsWhenDestinationAlreadyExists(t *testing.T) {
	home := setupCLITest(t)
	testutil.EnsurePackageDir(t, home, "demo")

	projectDir := t.TempDir()
	withWorkingDir(t, projectDir)
	if err := os.MkdirAll(filepath.Join(projectDir, "demo"), 0o755); err != nil {
		t.Fatalf("could not create existing destination: %v", err)
	}

	_, err := runCommand(newExportCmd(), []string{"demo"}, "")
	if err == nil {
		t.Fatal("expected error when destination already exists")
	}
	if !strings.Contains(err.Error(), "destination already exists") {
		t.Fatalf("unexpected destination exists error: %q", err.Error())
	}
}

func TestExportCommandFailsWhenPackageMissing(t *testing.T) {
	setupCLITest(t)

	_, err := runCommand(newExportCmd(), []string{"missing"}, "")
	if err == nil {
		t.Fatal("expected error when package does not exist")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("unexpected package missing error: %q", err.Error())
	}
}
