package cli

import "github.com/Bbeboy/AgentPack/internal/i18n"

var currentLang = i18n.ResolveLanguage()

func t(key string, args ...any) string {
	return i18n.Message(currentLang, key, args...)
}

func out(key string, args ...any) string {
	return i18n.Message(currentLang, "prefix", t(key, args...))
}
