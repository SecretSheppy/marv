package cmds

import (
	"fmt"
	"slices"
	"strings"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/fws"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"frameworks"},
	Short:   "lists all installed frameworks",
	Long:    "lists all installed frameworks by name",
	Run: func(cmd *cobra.Command, args []string) {
		frameworks := fws.Frameworks()
		slices.SortFunc(frameworks, func(a, b fwlib.Framework) int {
			return strings.Compare(a.Meta().Name, b.Meta().Name)
		})
		for _, fw := range frameworks {
			fmt.Println(fw.Meta().Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
