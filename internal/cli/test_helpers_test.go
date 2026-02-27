package cli

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/Bbeboy/AgentPack/internal/testutil"
	"github.com/spf13/cobra"
)

func setupCLITest(t *testing.T) string {
	t.Helper()
	home := testutil.SetupHome(t)
	prev := currentLang
	currentLang = "en"
	t.Cleanup(func() {
		currentLang = prev
	})
	return home
}

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("could not change working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(prev)
	})
}

func runCommand(cmd *cobra.Command, args []string, input string) (string, error) {
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	if input != "" {
		cmd.SetIn(bytes.NewBufferString(input))
	} else {
		cmd.SetIn(bytes.NewBuffer(nil))
	}
	cmd.SetArgs(args)
	err := cmd.Execute()
	if err != nil {
		if output.Len() > 0 {
			return output.String(), err
		}
		return "", err
	}

	data, readErr := io.ReadAll(&output)
	if readErr != nil {
		return output.String(), nil
	}
	return string(data), nil
}
