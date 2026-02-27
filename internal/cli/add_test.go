package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/testutil"
)

func TestAddCommandCopiesFileIntoPackage(t *testing.T) {
	home := setupCLITest(t)
	packagePath := testutil.EnsurePackageDir(t, home, "demo")

	sourceDir := t.TempDir()
	sourceFile := filepath.Join(sourceDir, "SKILL.md")
	testutil.WriteFile(t, sourceFile, "name: demo\n")

	cmd := newAddCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{sourceFile, "--to", "demo"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("add command returned error: %v", err)
	}

	targetPath := filepath.Join(packagePath, "SKILL.md")
	data, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("could not read copied file: %v", err)
	}
	if string(data) != "name: demo\n" {
		t.Fatalf("unexpected copied content: %q", string(data))
	}

	if !strings.Contains(out.String(), "agentpack:") {
		t.Fatalf("expected prefixed output, got %q", out.String())
	}
}

func TestAddCommandFailsWhenPackageIsMissing(t *testing.T) {
	setupCLITest(t)

	sourceDir := t.TempDir()
	sourceFile := filepath.Join(sourceDir, "SKILL.md")
	testutil.WriteFile(t, sourceFile, "name: demo\n")

	cmd := newAddCmd()
	cmd.SetArgs([]string{sourceFile, "--to", "missing"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when package does not exist")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %q", err.Error())
	}
}
