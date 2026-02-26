package main

import (
	"fmt"
	"os"

	"github.com/Bbeboy/AgentPack/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "[agentpack] Error: %v\n", err)
		os.Exit(1)
	}
}
