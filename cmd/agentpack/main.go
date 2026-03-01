package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/Bbeboy/AgentPack/internal/cli"
	"github.com/Bbeboy/AgentPack/internal/i18n"
	"github.com/Bbeboy/AgentPack/internal/selfinstall"
	"github.com/Bbeboy/AgentPack/internal/version"
)

func main() {
	lang := i18n.ResolveLanguage()

	if selfPath, err := os.Executable(); err == nil {
		result, installErr := selfinstall.MaybeInstall(lang, os.Args[1:], selfPath, runtime.GOOS, version.Value(), selfinstall.Dependencies{
			LookupPath:  exec.LookPath,
			ReplaceFile: selfinstall.ReplaceFile,
		})
		if installErr != nil {
			fmt.Fprintf(os.Stderr, i18n.Message(lang, "main.error", installErr)+"\n")
			os.Exit(1)
		}
		if result.Triggered {
			fmt.Fprintln(os.Stdout, i18n.Message(lang, "selfinstall.success", result.Target))
			return
		}
	}

	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, i18n.Message(lang, "main.error", err)+"\n")
		os.Exit(1)
	}
}
