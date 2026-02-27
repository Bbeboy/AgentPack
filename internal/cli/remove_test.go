package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/testutil"
)

func TestValidateRelativePath(t *testing.T) {
	absPath := filepath.Join(string(filepath.Separator), "tmp", "file")

	tests := []struct {
		name      string
		value     string
		wantErr   bool
		wantClean string
	}{
		{name: "empty", value: "", wantErr: true},
		{name: "dot", value: ".", wantErr: true},
		{name: "absolute", value: absPath, wantErr: true},
		{name: "traversal", value: "../secret", wantErr: true},
		{name: "valid", value: "docker/SKILL.md", wantErr: false, wantClean: filepath.Join("docker", "SKILL.md")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := validateRelativePath(tc.value)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for value %q", tc.value)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for value %q: %v", tc.value, err)
			}
			if got != tc.wantClean {
				t.Fatalf("expected clean path %q, got %q", tc.wantClean, got)
			}
		})
	}
}

func TestRemoveCommandPathModeRemovesTarget(t *testing.T) {
	home := setupCLITest(t)
	packagePath := testutil.EnsurePackageDir(t, home, "demo")
	targetFile := filepath.Join(packagePath, "docker", "SKILL.md")
	testutil.WriteFile(t, targetFile, "name: docker\n")

	cmd := newRemoveCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"docker/SKILL.md", "--from", "demo", "--force"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("remove command returned error: %v", err)
	}

	if _, err := os.Stat(targetFile); !os.IsNotExist(err) {
		t.Fatalf("expected target file removed, stat err=%v", err)
	}
}

func TestRemoveCommandDryRunKeepsPackage(t *testing.T) {
	home := setupCLITest(t)
	packagePath := testutil.EnsurePackageDir(t, home, "demo")

	cmd := newRemoveCmd()
	cmd.SetArgs([]string{"demo", "--dry-run"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("remove dry-run returned error: %v", err)
	}

	if _, err := os.Stat(packagePath); err != nil {
		t.Fatalf("expected package to remain after dry-run, err=%v", err)
	}
}

func TestRemoveCommandRequiresFromForPathTarget(t *testing.T) {
	setupCLITest(t)

	cmd := newRemoveCmd()
	cmd.SetArgs([]string{"docker/SKILL.md"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for path target without --from")
	}
}
