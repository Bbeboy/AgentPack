package main

import (
	"fmt"
	"os"

	"github.com/Bbeboy/AgentPack/internal/cli"
	"github.com/Bbeboy/AgentPack/internal/i18n"
)

func main() {
	if err := cli.Execute(); err != nil {
		lang := i18n.ResolveLanguage()
		fmt.Fprintf(os.Stderr, i18n.Message(lang, "main.error", err)+"\n")
		os.Exit(1)
	}
}
