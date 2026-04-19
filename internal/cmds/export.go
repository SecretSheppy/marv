package cmds

import (
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "exports framework output into standardised JSON",
	Long:  "exports the output data from the configured frameworks into the marv internal format (JSON) and then exports it.",
	Run: func(cmd *cobra.Command, args []string) {
		exportCommand()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
