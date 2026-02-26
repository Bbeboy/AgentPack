package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Bbeboy/AgentPack/internal/prompt"
	"github.com/Bbeboy/AgentPack/internal/storage"
	"github.com/spf13/cobra"
)

func newRemoveSkillCmd() *cobra.Command {
	var force bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "remove-skill <nombre-paquete> <nombre-skill>",
		Short: "Elimina una skill de un paquete",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]
			skillName := args[1]

			skillPath, err := storage.SkillPath(packageName, skillName)
			if err != nil {
				return err
			}

			skillInfo, err := os.Stat(skillPath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("no se encontro la skill '%s' en el paquete '%s'\n[agentpack] Buscado en: %s", skillName, packageName, skillPath)
				}
				return fmt.Errorf("no se pudo leer la skill '%s': %w", skillName, err)
			}

			if !skillInfo.IsDir() {
				return fmt.Errorf("la skill '%s' en el paquete '%s' no es un directorio valido", skillName, packageName)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Skill objetivo: %s\n", skillPath)

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Dry-run: se eliminaria la skill '%s' del paquete '%s'.\n", skillName, packageName)
				return nil
			}

			if !force {
				reader := bufio.NewReader(cmd.InOrStdin())
				question := fmt.Sprintf("Eliminar la skill '%s' del paquete '%s'?", skillName, packageName)
				confirm, err := prompt.YesNo(reader, cmd.OutOrStdout(), question)
				if err != nil {
					return err
				}
				if !confirm {
					fmt.Fprintln(cmd.OutOrStdout(), "[agentpack] Operacion cancelada.")
					return nil
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Eliminando skill '%s' del paquete '%s'...\n", skillName, packageName)
			if err := os.RemoveAll(skillPath); err != nil {
				return fmt.Errorf("no se pudo eliminar la skill '%s' del paquete '%s': %w", skillName, packageName, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "[agentpack] Skill eliminada: %s\n", skillName)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Elimina la skill sin pedir confirmacion")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Muestra que se eliminaria sin borrar nada")

	return cmd
}
