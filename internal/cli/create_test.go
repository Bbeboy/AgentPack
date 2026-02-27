package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/testutil"
)

func TestCreateCommandWithDotSingleCandidate(t *testing.T) {
	home := setupCLITest(t)

	projectDir := t.TempDir()
	withWorkingDir(t, projectDir)
	testutil.WriteFile(t, filepath.Join(projectDir, ".opencode", "skills", "docker", "SKILL.md"), "name: docker\n")

	output, err := runCommand(newCreateCmd(), []string{"pack-single", "."}, "")
	if err != nil {
		t.Fatalf("create command returned error: %v", err)
	}

	packageSkill := filepath.Join(home, ".agentpack", "packages-skills", "pack-single", "docker", "SKILL.md")
	data, readErr := os.ReadFile(packageSkill)
	if readErr != nil {
		t.Fatalf("expected copied skill file: %v", readErr)
	}
	if string(data) != "name: docker\n" {
		t.Fatalf("unexpected copied content: %q", string(data))
	}
	if !strings.Contains(output, "detected skills folder") {
		t.Fatalf("expected auto-detected folder output, got: %q", output)
	}
}

func TestCreateCommandWithDotMultipleCandidatesUsesSelection(t *testing.T) {
	home := setupCLITest(t)

	projectDir := t.TempDir()
	withWorkingDir(t, projectDir)
	testutil.WriteFile(t, filepath.Join(projectDir, ".opencode", "skills", "from-opencode", "SKILL.md"), "source: opencode\n")
	testutil.WriteFile(t, filepath.Join(projectDir, ".agents", "skills", "from-agents", "SKILL.md"), "source: agents\n")

	output, err := runCommand(newCreateCmd(), []string{"pack-multi", "."}, "2\n")
	if err != nil {
		t.Fatalf("create command returned error: %v", err)
	}

	agentsSkill := filepath.Join(home, ".agentpack", "packages-skills", "pack-multi", "from-agents", "SKILL.md")
	if _, statErr := os.Stat(agentsSkill); statErr != nil {
		t.Fatalf("expected selected skills source copied: %v", statErr)
	}

	opencodeSkill := filepath.Join(home, ".agentpack", "packages-skills", "pack-multi", "from-opencode", "SKILL.md")
	if _, statErr := os.Stat(opencodeSkill); !os.IsNotExist(statErr) {
		t.Fatalf("expected non-selected source not copied, stat err=%v", statErr)
	}

	if !strings.Contains(output, "multiple skills folders found") {
		t.Fatalf("expected multiple choice prompt output, got: %q", output)
	}
}

func TestCreateCommandWithDotFailsWhenNoCandidates(t *testing.T) {
	setupCLITest(t)

	projectDir := t.TempDir()
	withWorkingDir(t, projectDir)

	_, err := runCommand(newCreateCmd(), []string{"pack-empty", "."}, "")
	if err == nil {
		t.Fatal("expected error when no skills folders are present")
	}
	if !strings.Contains(err.Error(), "no skills folder") {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}
