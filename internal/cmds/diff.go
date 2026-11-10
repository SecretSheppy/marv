package cmds

import "github.com/spf13/cobra"

var before, after string

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "displays the Marv mutations diff web interface",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: this will ultimately display the mutations diff web interface + export mutations diff file
	},
}

var diffExportCmd = &cobra.Command{
	Use:   "export",
	Short: "produces a Marv diff file from two Marv mutation schema files",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: export the Marv mutations diff file

		// TODO: read in before file
		// TODO: read in after file

		// TODO: produce list of IDs in before file
		// TODO: produce list of IDs in after file

		// TODO: produce list of IDs from before not in after
		// TODO: produce list of IDs from after not in before

		// TODO: go through each ID in both files and determine whether their status is still the same.
		//  - if status is not the same, put the ID in the changes list with the old and new status.
	},
}

func init() {
	diffCmd.PersistentFlags().StringVarP(&before, "before", "b", "", "before schema file name")
	diffCmd.PersistentFlags().StringVarP(&after, "after", "a", "", "after schema file name")

	diffCmd.AddCommand(diffExportCmd)

	rootCmd.AddCommand(diffCmd)
}
