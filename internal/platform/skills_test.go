package platform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCandidateSkillsPathsIncludesFallback(t *testing.T) {
	paths := CandidateSkillsPaths()
	if len(paths) == 0 {
		t.Fatal("expected candidate paths")
	}

	if paths[0] != ".opencode/skills" {
		t.Fatalf("expected first candidate to be .opencode/skills, got %q", paths[0])
	}

	last := paths[len(paths)-1]
	if last != "skills" {
		t.Fatalf("expected last candidate to be skills fallback, got %q", last)
	}
}

func TestResolveSkillsDestinationUsesPriorityOrder(t *testing.T) {
	projectRoot := t.TempDir()

	if err := os.MkdirAll(filepath.Join(projectRoot, ".agents"), 0o755); err != nil {
		t.Fatalf("could not create .agents: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(projectRoot, ".opencode"), 0o755); err != nil {
		t.Fatalf("could not create .opencode: %v", err)
	}

	path, platformName, detected, err := ResolveSkillsDestination(projectRoot)
	if err != nil {
		t.Fatalf("ResolveSkillsDestination returned error: %v", err)
	}

	expected := filepath.Join(projectRoot, ".opencode", "skills")
	if path != expected {
		t.Fatalf("expected %q, got %q", expected, path)
	}
	if platformName != "OpenCode" {
		t.Fatalf("expected platform OpenCode, got %q", platformName)
	}
	if !detected {
		t.Fatal("expected detected=true")
	}
}

func TestResolveSkillsDestinationFallsBackToAgents(t *testing.T) {
	projectRoot := t.TempDir()

	path, platformName, detected, err := ResolveSkillsDestination(projectRoot)
	if err != nil {
		t.Fatalf("ResolveSkillsDestination returned error: %v", err)
	}

	expected := filepath.Join(projectRoot, ".agents", "skills")
	if path != expected {
		t.Fatalf("expected %q, got %q", expected, path)
	}
	if platformName != "Amp" {
		t.Fatalf("expected platform Amp, got %q", platformName)
	}
	if detected {
		t.Fatal("expected detected=false for fallback")
	}
}
