package cmds

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all installed framework extensions",
	Long:  "lists all installed framework extensions by frameworks name",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: complete when extension loader is written
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
