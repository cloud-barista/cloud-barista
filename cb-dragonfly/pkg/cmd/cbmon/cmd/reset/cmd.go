package reset

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "",
		Long:  "",
	}
	cmd.AddCommand(newResetConfigCmd())
	return cmd
}
