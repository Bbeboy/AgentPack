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

func TestAddCommandPreservesRelativePathInsidePackage(t *testing.T) {
	home := setupCLITest(t)
	packagePath := testutil.EnsurePackageDir(t, home, "demo")

	projectDir := t.TempDir()
	withWorkingDir(t, projectDir)

	relSource := filepath.Join("skills", "docker", "SKILL.md")
	testutil.WriteFile(t, filepath.Join(projectDir, relSource), "name: docker\n")

	output, err := runCommand(newAddCmd(), []string{relSource, "--to", "demo"}, "")
	if err != nil {
		t.Fatalf("add command returned error: %v", err)
	}

	targetPath := filepath.Join(packagePath, relSource)
	data, readErr := os.ReadFile(targetPath)
	if readErr != nil {
		t.Fatalf("expected copied file at preserved relative path: %v", readErr)
	}
	if string(data) != "name: docker\n" {
		t.Fatalf("unexpected copied content: %q", string(data))
	}

	if !strings.Contains(output, filepath.Join("skills", "docker", "SKILL.md")) {
		t.Fatalf("expected relative destination in output, got %q", output)
	}
}

func TestAddCommandRejectsRelativeTraversal(t *testing.T) {
	setupCLITest(t)

	parent := t.TempDir()
	child := filepath.Join(parent, "child")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatalf("could not create child dir: %v", err)
	}
	withWorkingDir(t, child)

	testutil.WriteFile(t, filepath.Join(parent, "outside.txt"), "outside\n")
	testutil.EnsurePackageDir(t, os.Getenv("HOME"), "demo")

	errOutput, err := runCommand(newAddCmd(), []string{"../outside.txt", "--to", "demo"}, "")
	if err == nil {
		t.Fatal("expected error for traversal path")
	}
	if !strings.Contains(err.Error(), "cannot escape") && !strings.Contains(err.Error(), "no puede salir") {
		t.Fatalf("unexpected traversal error: %q", err.Error())
	}
	_ = errOutput
}

func TestAddCommandRejectsWindowsStyleTraversal(t *testing.T) {
	setupCLITest(t)

	_, err := cleanAddRelativePath("..\\outside.txt")
	if err == nil {
		t.Fatal("expected error for windows-style traversal path")
	}
	if !strings.Contains(err.Error(), "cannot escape") && !strings.Contains(err.Error(), "no puede salir") {
		t.Fatalf("unexpected traversal error: %q", err.Error())
	}
}
