package cli

import (
	"strings"
	"testing"
)

func TestVersionCommandPrintsVersion(t *testing.T) {
	setupCLITest(t)

	output, err := runCommand(newVersionCmd(), nil, "")
	if err != nil {
		t.Fatalf("version command returned error: %v", err)
	}

	if !strings.Contains(output, "agentpack: version") {
		t.Fatalf("expected version output, got %q", output)
	}
}

func TestRootVersionFlagPrintsVersion(t *testing.T) {
	setupCLITest(t)

	output, err := runCommand(rootCmd, []string{"-v"}, "")
	if err != nil {
		t.Fatalf("root -v returned error: %v", err)
	}

	if !strings.Contains(output, "agentpack: version") {
		t.Fatalf("expected root version output, got %q", output)
	}
}
