package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func SelectIndex(reader *bufio.Reader, out io.Writer, total int) (int, error) {
	if total <= 0 {
		return 0, fmt.Errorf("no hay opciones disponibles")
	}

	for {
		fmt.Fprintf(out, "Elige una opcion [1-%d]: ", total)
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf("no se pudo leer la opcion: %w", err)
		}

		answer := strings.TrimSpace(line)
		option, convErr := strconv.Atoi(answer)
		if convErr == nil && option >= 1 && option <= total {
			return option - 1, nil
		}

		fmt.Fprintln(out, "[agentpack] Opcion invalida. Intenta de nuevo.")

		if err == io.EOF {
			return 0, fmt.Errorf("entrada finalizada sin una opcion valida")
		}
	}
}

func YesNo(reader *bufio.Reader, out io.Writer, text string) (bool, error) {
	for {
		fmt.Fprintf(out, "%s [y/N]: ", text)
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return false, fmt.Errorf("no se pudo leer la respuesta: %w", err)
		}

		answer := strings.ToLower(strings.TrimSpace(line))
		switch answer {
		case "", "n", "no":
			return false, nil
		case "y", "yes", "s", "si":
			return true, nil
		default:
			fmt.Fprintln(out, "[agentpack] Respuesta invalida. Usa y/n.")
		}

		if err == io.EOF {
			return false, nil
		}
	}
}
