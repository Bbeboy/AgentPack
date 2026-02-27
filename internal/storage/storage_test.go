package storage

import (
	"path/filepath"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/testutil"
)

func TestPackagesRootUsesHomeDirectory(t *testing.T) {
	home := testutil.SetupHome(t)

	root, err := PackagesRoot()
	if err != nil {
		t.Fatalf("PackagesRoot returned error: %v", err)
	}

	expected := filepath.Join(home, ".agentpack", "packages-skills")
	if root != expected {
		t.Fatalf("expected %q, got %q", expected, root)
	}
}

func TestPackagePathValidatesName(t *testing.T) {
	testutil.SetupHome(t)

	if _, err := PackagePath("backend-base"); err != nil {
		t.Fatalf("expected valid package name, got error: %v", err)
	}

	if _, err := PackagePath("../bad"); err == nil {
		t.Fatal("expected error for invalid package name")
	}
}

func TestSkillPathValidatesNames(t *testing.T) {
	testutil.SetupHome(t)

	path, err := SkillPath("backend", "docker")
	if err != nil {
		t.Fatalf("expected valid skill path, got error: %v", err)
	}

	expectedSuffix := filepath.Join("backend", "docker")
	if filepath.Base(filepath.Dir(path)) != "backend" || filepath.Base(path) != "docker" {
		t.Fatalf("expected path ending in %q, got %q", expectedSuffix, path)
	}

	if _, err := SkillPath("backend", "../bad"); err == nil {
		t.Fatal("expected error for invalid skill name")
	}
}
