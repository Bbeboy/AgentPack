package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/testutil"
)

func TestRenameCommandRenamesPackage(t *testing.T) {
	home := setupCLITest(t)
	oldPath := testutil.EnsurePackageDir(t, home, "oldpack")
	testutil.WriteFile(t, filepath.Join(oldPath, "docker", "SKILL.md"), "name: docker\n")

	cmd := newRenameCmd()
	cmd.SetArgs([]string{"oldpack", "newpack"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("rename command returned error: %v", err)
	}

	newPath := filepath.Join(home, ".agentpack", "packages-skills", "newpack")
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected renamed package at %q: %v", newPath, err)
	}

	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("expected old package path removed, stat err=%v", err)
	}
}

func TestRenameCommandFailsWhenTargetExists(t *testing.T) {
	home := setupCLITest(t)
	testutil.EnsurePackageDir(t, home, "oldpack")
	testutil.EnsurePackageDir(t, home, "newpack")

	cmd := newRenameCmd()
	cmd.SetArgs([]string{"oldpack", "newpack"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when target package exists")
	}
}
