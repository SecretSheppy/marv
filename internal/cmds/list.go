package cmds

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

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
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		for _, fw := range frameworks {
			fmt.Fprintf(w, "%s\t%s\n", fw.Meta().Name, fw.Meta().URL)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
