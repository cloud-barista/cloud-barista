package set

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "",
		Long:  "",
	}
	cmd.AddCommand(newSetConfigCmd())
	return cmd
}
