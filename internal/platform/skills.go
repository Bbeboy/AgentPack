package platform

import (
	"fmt"
	"os"
	"path/filepath"
)

type Candidate struct {
	Name string
	Root string
}

var skillCandidates = []Candidate{
	{Name: "OpenCode", Root: ".opencode"},
	{Name: "Amp", Root: ".agents"},
	{Name: "Antigravity", Root: ".agent"},
	{Name: "AdaL", Root: ".adal"},
	{Name: "Augment Code", Root: ".augment"},
	{Name: "Claude Code", Root: ".cla"},
	{Name: "Claude Code Plugin", Root: ".claude-plugin"},
	{Name: "Cline", Root: ".cline"},
	{Name: "CodeBuddy", Root: ".codebuddy"},
	{Name: "Codex CLI", Root: ".codex"},
	{Name: "Command Code", Root: ".commandcode"},
	{Name: "Continue", Root: ".continue"},
	{Name: "Crush", Root: ".config/crush"},
	{Name: "Cursor", Root: ".cursor"},
	{Name: "Factory AI", Root: ".factory"},
	{Name: "GitHub Copilot", Root: ".github"},
	{Name: "Goose", Root: ".goose"},
	{Name: "iFlow CLI", Root: ".iflow"},
	{Name: "Junie", Root: ".junie"},
	{Name: "Kilo Code", Root: ".kilocode"},
	{Name: "Kiro", Root: ".kiro"},
	{Name: "Kode", Root: ".kode"},
	{Name: "MCPJam", Root: ".mcpjam"},
	{Name: "Mistral Vibe", Root: ".vibe"},
	{Name: "Mux", Root: ".mux"},
	{Name: "Neovate", Root: ".neovate"},
	{Name: "OpenClaw", Root: ".openclaw"},
	{Name: "OpenHands", Root: ".openhands"},
	{Name: "Pi-Mono", Root: ".pi"},
	{Name: "Pochi", Root: ".pochi"},
	{Name: "Qoder", Root: ".qoder"},
	{Name: "Qwen Code", Root: ".qwen"},
	{Name: "Roo Code", Root: ".roo"},
	{Name: "Trae", Root: ".trae"},
	{Name: "Trae CN", Root: ".trae-cn"},
	{Name: "Windsurf", Root: ".windsurf"},
	{Name: "Zencoder", Root: ".zencoder"},
}

func Candidates() []Candidate {
	result := make([]Candidate, len(skillCandidates))
	copy(result, skillCandidates)
	return result
}

func CandidateSkillsPaths() []string {
	paths := make([]string, 0, len(skillCandidates)+1)
	for _, c := range skillCandidates {
		paths = append(paths, filepath.ToSlash(filepath.Join(c.Root, "skills")))
	}
	paths = append(paths, "skills")
	return paths
}

func ResolveSkillsDestination(cwd string) (path string, platform string, detected bool, err error) {
	for _, c := range skillCandidates {
		root := filepath.Join(cwd, filepath.FromSlash(c.Root))
		info, statErr := os.Stat(root)
		if statErr != nil {
			if os.IsNotExist(statErr) {
				continue
			}
			return "", "", false, fmt.Errorf("could not inspect platform root '%s': %w", c.Root, statErr)
		}
		if info.IsDir() {
			return filepath.Join(root, "skills"), c.Name, true, nil
		}
	}

	fallback := filepath.Join(cwd, ".agents", "skills")
	return fallback, "Amp", false, nil
}
