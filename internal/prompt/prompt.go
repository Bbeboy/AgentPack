package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Bbeboy/AgentPack/internal/i18n"
)

func SelectIndex(reader *bufio.Reader, out io.Writer, total int) (int, error) {
	lang := i18n.ResolveLanguage()

	if total <= 0 {
		return 0, fmt.Errorf(i18n.Message(lang, "prompt.select.eof"))
	}

	for {
		fmt.Fprintf(out, i18n.Message(lang, "prompt.select", total))
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf(i18n.Message(lang, "prompt.select.readerror", err))
		}

		answer := strings.TrimSpace(line)
		option, convErr := strconv.Atoi(answer)
		if convErr == nil && option >= 1 && option <= total {
			return option - 1, nil
		}

		fmt.Fprintln(out, i18n.Message(lang, "prompt.select.invalid"))

		if err == io.EOF {
			return 0, fmt.Errorf(i18n.Message(lang, "prompt.select.eof"))
		}
	}
}

func YesNo(reader *bufio.Reader, out io.Writer, text string) (bool, error) {
	lang := i18n.ResolveLanguage()

	for {
		fmt.Fprintf(out, "%s [y/N]: ", text)
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return false, fmt.Errorf(i18n.Message(lang, "prompt.yesno.readerror", err))
		}

		answer := strings.ToLower(strings.TrimSpace(line))
		switch answer {
		case "", "n", "no":
			return false, nil
		case "y", "yes", "s", "si":
			return true, nil
		default:
			fmt.Fprintln(out, i18n.Message(lang, "prompt.yesno.invalid"))
		}

		if err == io.EOF {
			return false, nil
		}
	}
}
