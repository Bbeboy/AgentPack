package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/testutil"
)

func TestRootCommandRegistrationForAddSkill(t *testing.T) {
	setupCLITest(t)

	addSkillCmd, _, err := rootCmd.Find([]string{"add-skill"})
	if err != nil {
		t.Fatalf("expected add-skill command registered: %v", err)
	}
	if addSkillCmd == nil || addSkillCmd.Name() != "add-skill" {
		t.Fatalf("expected add-skill command, got %v", addSkillCmd)
	}

	legacyAddCmd, _, legacyErr := rootCmd.Find([]string{"add"})
	if legacyErr != nil {
		t.Fatalf("expected legacy add command to resolve: %v", legacyErr)
	}
	if legacyAddCmd == nil || legacyAddCmd.Name() != "add" {
		t.Fatalf("expected legacy add command, got %v", legacyAddCmd)
	}
	if !legacyAddCmd.Hidden {
		t.Fatal("expected legacy add command hidden from help output")
	}

}

func TestLegacyAddCommandReturnsMigrationError(t *testing.T) {
	setupCLITest(t)

	_, err := runCommand(rootCmd, []string{"add", "any", "--to", "demo"}, "")
	if err == nil {
		t.Fatal("expected legacy add command to fail")
	}
	if !strings.Contains(err.Error(), "add-skill") {
		t.Fatalf("expected migration guidance in error, got %q", err.Error())
	}
}

func TestRootAddSkillCommandIsRegistered(t *testing.T) {
	home := setupCLITest(t)
	packagePath := testutil.EnsurePackageDir(t, home, "demo")

	sourceDir := t.TempDir()
	sourceFile := filepath.Join(sourceDir, "SKILL.md")
	testutil.WriteFile(t, sourceFile, "name: demo\n")

	if _, err := runCommand(rootCmd, []string{"add-skill", sourceFile, "--to", "demo"}, ""); err != nil {
		t.Fatalf("add-skill command returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(packagePath, "SKILL.md")); err != nil {
		t.Fatalf("expected copied skill in package, got error: %v", err)
	}
}
