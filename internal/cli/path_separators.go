package cli

import (
	"os"
	"strings"
)

func normalizePathSeparators(value string) string {
	if os.PathSeparator == '/' {
		return strings.ReplaceAll(value, "\\", "/")
	}

	return strings.ReplaceAll(value, "/", "\\")
}

func containsPathSeparator(value string) bool {
	if strings.ContainsRune(value, os.PathSeparator) {
		return true
	}

	if os.PathSeparator == '/' {
		return strings.Contains(value, "\\")
	}

	return strings.Contains(value, "/")
}
