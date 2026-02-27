package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/testutil"
)

func TestCreateThenInstallSkipsConflictInteractively(t *testing.T) {
	setupCLITest(t)

	sourceProject := t.TempDir()
	testutil.WriteFile(t, filepath.Join(sourceProject, ".opencode", "skills", "docker", "SKILL.md"), "from-package\n")
	testutil.WriteFile(t, filepath.Join(sourceProject, ".opencode", "skills", "git", "SKILL.md"), "git-skill\n")

	withWorkingDir(t, sourceProject)
	if _, err := runCommand(newCreateCmd(), []string{"flow-pack", "."}, ""); err != nil {
		t.Fatalf("create command returned error: %v", err)
	}

	targetProject := t.TempDir()
	testutil.WriteFile(t, filepath.Join(targetProject, ".opencode", "skills", "docker", "SKILL.md"), "already-there\n")
	withWorkingDir(t, targetProject)

	output, err := runCommand(newInstallCmd(), []string{"flow-pack"}, "n\n")
	if err != nil {
		t.Fatalf("install command returned error: %v", err)
	}

	dockerPath := filepath.Join(targetProject, ".opencode", "skills", "docker", "SKILL.md")
	dockerData, readErr := os.ReadFile(dockerPath)
	if readErr != nil {
		t.Fatalf("expected docker skill to exist: %v", readErr)
	}
	if string(dockerData) != "already-there\n" {
		t.Fatalf("expected conflict to remain unchanged when skipped, got %q", string(dockerData))
	}

	gitPath := filepath.Join(targetProject, ".opencode", "skills", "git", "SKILL.md")
	if _, statErr := os.Stat(gitPath); statErr != nil {
		t.Fatalf("expected non-conflicting skill installed: %v", statErr)
	}

	if !strings.Contains(output, "conflicts detected") {
		t.Fatalf("expected conflict output, got %q", output)
	}
	if !strings.Contains(output, "installed=1 overwritten=0 skipped=1") {
		t.Fatalf("expected summary counts in output, got %q", output)
	}

}

func TestInstallOverwritesConflictInteractively(t *testing.T) {
	home := setupCLITest(t)

	packageRoot := testutil.EnsurePackageDir(t, home, "overwrite-pack")
	testutil.WriteFile(t, filepath.Join(packageRoot, "docker", "SKILL.md"), "from-package\n")

	targetProject := t.TempDir()
	testutil.WriteFile(t, filepath.Join(targetProject, ".opencode", "skills", "docker", "SKILL.md"), "existing\n")
	withWorkingDir(t, targetProject)

	output, err := runCommand(newInstallCmd(), []string{"overwrite-pack"}, "y\n")
	if err != nil {
		t.Fatalf("install command returned error: %v", err)
	}

	dockerPath := filepath.Join(targetProject, ".opencode", "skills", "docker", "SKILL.md")
	dockerData, readErr := os.ReadFile(dockerPath)
	if readErr != nil {
		t.Fatalf("expected docker skill after overwrite: %v", readErr)
	}
	if string(dockerData) != "from-package\n" {
		t.Fatalf("expected overwrite content, got %q", string(dockerData))
	}

	if !strings.Contains(output, "overwrote 'docker'") {
		t.Fatalf("expected overwrite output, got %q", output)
	}
	if !strings.Contains(output, "installed=0 overwritten=1 skipped=0") {
		t.Fatalf("expected summary counts in output, got %q", output)
	}
}

func TestInstallUsesDetectedPlatformDestination(t *testing.T) {
	home := setupCLITest(t)
	packageRoot := testutil.EnsurePackageDir(t, home, "dest-pack")
	testutil.WriteFile(t, filepath.Join(packageRoot, "docker", "SKILL.md"), "from-package\n")

	targetProject := t.TempDir()
	if err := os.MkdirAll(filepath.Join(targetProject, ".cursor"), 0o755); err != nil {
		t.Fatalf("could not create .cursor marker: %v", err)
	}
	withWorkingDir(t, targetProject)

	output, err := runCommand(newInstallCmd(), []string{"dest-pack"}, "")
	if err != nil {
		t.Fatalf("install command returned error: %v", err)
	}

	installedPath := filepath.Join(targetProject, ".cursor", "skills", "docker", "SKILL.md")
	if _, statErr := os.Stat(installedPath); statErr != nil {
		t.Fatalf("expected install in detected platform path: %v", statErr)
	}

	if !strings.Contains(output, filepath.Join(targetProject, ".cursor", "skills")) {
		t.Fatalf("expected destination path in output, got %q", output)
	}
}
