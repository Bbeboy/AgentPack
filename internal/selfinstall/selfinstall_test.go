package selfinstall

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestMaybeInstallTriggersAndReplacesTarget(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "agentpack-v1.2.3-linux-amd64")
	targetPath := filepath.Join(tempDir, "agentpack")

	if err := os.WriteFile(sourcePath, []byte("new-binary"), 0o755); err != nil {
		t.Fatalf("write source: %v", err)
	}
	if err := os.WriteFile(targetPath, []byte("old-binary"), 0o755); err != nil {
		t.Fatalf("write target: %v", err)
	}

	result, err := MaybeInstall("en", nil, sourcePath, "linux", "v1.2.3", Dependencies{
		LookupPath: func(name string) (string, error) {
			if name != "agentpack" {
				t.Fatalf("expected lookup for agentpack, got %q", name)
			}
			return targetPath, nil
		},
		ReplaceFile: ReplaceFile,
	})
	if err != nil {
		t.Fatalf("expected install success, got error: %v", err)
	}
	if !result.Triggered {
		t.Fatalf("expected flow to trigger")
	}
	if result.Target != targetPath {
		t.Fatalf("expected target %q, got %q", targetPath, result.Target)
	}

	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if string(content) != "new-binary" {
		t.Fatalf("expected target to be replaced, got %q", string(content))
	}
}

func TestMaybeInstallSkipsWhenArgsProvided(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "agentpack-v1.2.3-linux-amd64")

	if err := os.WriteFile(sourcePath, []byte("new-binary"), 0o755); err != nil {
		t.Fatalf("write source: %v", err)
	}

	result, err := MaybeInstall("en", []string{"version"}, sourcePath, "linux", "v1.2.3", Dependencies{
		LookupPath: func(name string) (string, error) {
			t.Fatalf("lookup should not be called, got %q", name)
			return "", nil
		},
		ReplaceFile: func(source, target string) error {
			t.Fatalf("replace should not be called: %s -> %s", source, target)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("expected skip without error, got %v", err)
	}
	if result.Triggered {
		t.Fatalf("expected flow to skip when args are provided")
	}
}

func TestMaybeInstallSkipsWhenBinaryNameIsNotReleaseArtifact(t *testing.T) {
	result, err := MaybeInstall("en", nil, filepath.Join("/tmp", "agentpack-custom"), "linux", "v1.2.3", Dependencies{
		LookupPath: func(name string) (string, error) {
			t.Fatalf("lookup should not be called, got %q", name)
			return "", nil
		},
		ReplaceFile: func(source, target string) error {
			t.Fatalf("replace should not be called: %s -> %s", source, target)
			return nil
		},
	})
	if err != nil {
		t.Fatalf("expected skip without error, got %v", err)
	}
	if result.Triggered {
		t.Fatalf("expected flow to skip for non-release binary name")
	}
}

func TestMaybeInstallWindowsUsesExeTargetName(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "agentpack-v1.2.3-windows-amd64.exe")

	if err := os.WriteFile(sourcePath, []byte("new-binary"), 0o755); err != nil {
		t.Fatalf("write source: %v", err)
	}

	called := false
	_, err := MaybeInstall("en", nil, sourcePath, "windows", "v1.2.3", Dependencies{
		LookupPath: func(name string) (string, error) {
			called = true
			if name != "agentpack.exe" {
				t.Fatalf("expected lookup for agentpack.exe, got %q", name)
			}
			return filepath.Join(tempDir, "agentpack.exe"), nil
		},
		ReplaceFile: func(source, target string) error {
			return nil
		},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !called {
		t.Fatalf("expected lookup to be called")
	}
}

func TestMaybeInstallAcceptsArm64ArtifactName(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "agentpack-v1.2.3-linux-arm64")
	targetPath := filepath.Join(tempDir, "agentpack")

	if err := os.WriteFile(sourcePath, []byte("new-binary"), 0o755); err != nil {
		t.Fatalf("write source: %v", err)
	}
	if err := os.WriteFile(targetPath, []byte("old-binary"), 0o755); err != nil {
		t.Fatalf("write target: %v", err)
	}

	result, err := MaybeInstall("en", nil, sourcePath, "linux", "v1.2.3", Dependencies{
		LookupPath: func(name string) (string, error) {
			if name != "agentpack" {
				t.Fatalf("expected lookup for agentpack, got %q", name)
			}
			return targetPath, nil
		},
		ReplaceFile: ReplaceFile,
	})
	if err != nil {
		t.Fatalf("expected install success, got %v", err)
	}
	if !result.Triggered {
		t.Fatalf("expected flow to trigger for arm64 artifact")
	}
}

func TestMaybeInstallReturnsActionableErrorWhenTargetMissing(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "agentpack-v1.2.3-linux-amd64")

	if err := os.WriteFile(sourcePath, []byte("new-binary"), 0o755); err != nil {
		t.Fatalf("write source: %v", err)
	}

	result, err := MaybeInstall("en", nil, sourcePath, "linux", "v1.2.3", Dependencies{
		LookupPath: func(name string) (string, error) {
			return "", errors.New("not found")
		},
		ReplaceFile: ReplaceFile,
	})
	if err == nil {
		t.Fatalf("expected error when target is missing")
	}
	if !result.Triggered {
		t.Fatalf("expected flow to trigger for release artifact")
	}
}

func TestMaybeInstallTriggersForRegularBinaryNameWhenPathDiffers(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "downloads", "agentpack")
	targetPath := filepath.Join(tempDir, "bin", "agentpack")

	if err := os.MkdirAll(filepath.Dir(sourcePath), 0o755); err != nil {
		t.Fatalf("create source dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		t.Fatalf("create target dir: %v", err)
	}
	if err := os.WriteFile(sourcePath, []byte("new-binary"), 0o755); err != nil {
		t.Fatalf("write source: %v", err)
	}
	if err := os.WriteFile(targetPath, []byte("old-binary"), 0o755); err != nil {
		t.Fatalf("write target: %v", err)
	}

	result, err := MaybeInstall("en", nil, sourcePath, "linux", "v1.2.3", Dependencies{
		LookupPath: func(name string) (string, error) {
			if name != "agentpack" {
				t.Fatalf("expected lookup for agentpack, got %q", name)
			}
			return targetPath, nil
		},
		ReplaceFile: ReplaceFile,
	})
	if err != nil {
		t.Fatalf("expected install success, got error: %v", err)
	}
	if !result.Triggered {
		t.Fatalf("expected flow to trigger")
	}

	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if string(content) != "new-binary" {
		t.Fatalf("expected target to be replaced, got %q", string(content))
	}
}

func TestMaybeInstallSkipsWhenExecutableAlreadyMatchesTargetPath(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "agentpack")

	if err := os.WriteFile(targetPath, []byte("same-binary"), 0o755); err != nil {
		t.Fatalf("write target: %v", err)
	}

	result, err := MaybeInstall("en", nil, targetPath, "linux", "v1.2.3", Dependencies{
		LookupPath: func(name string) (string, error) {
			return targetPath, nil
		},
		ReplaceFile: func(source, target string) error {
			t.Fatalf("replace should not run when source and target are the same path")
			return nil
		},
	})
	if err != nil {
		t.Fatalf("expected skip without error, got %v", err)
	}
	if result.Triggered {
		t.Fatalf("expected flow to skip when executable already matches installed binary")
	}
}

func TestMaybeInstallSkipsForDevVersion(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "agentpack-v1.2.3-linux-amd64")

	if err := os.WriteFile(sourcePath, []byte("new-binary"), 0o755); err != nil {
		t.Fatalf("write source: %v", err)
	}

	result, err := MaybeInstall("en", nil, sourcePath, "linux", "dev", Dependencies{
		LookupPath: func(name string) (string, error) {
			t.Fatalf("lookup should not be called for dev version")
			return "", nil
		},
		ReplaceFile: func(source, target string) error {
			t.Fatalf("replace should not be called for dev version")
			return nil
		},
	})
	if err != nil {
		t.Fatalf("expected skip without error, got %v", err)
	}
	if result.Triggered {
		t.Fatalf("expected flow to skip for dev version")
	}
}
