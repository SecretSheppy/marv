package cmds

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/pkg/mutations"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "transform",
	Short: "transforms framework output into standardised JSON",
	Long:  "transforms the output data from the configured frameworks into the marv internal format (JSON) and then exports it.",
	Run: func(cmd *cobra.Command, args []string) {
		exportCommand()
	},
}

func exportCommand() {
	_, activeFws := getConfigAndFws()

	if output != "" {
		p, err := os.Stat(output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if !p.IsDir() {
			output = path.Base(output)
		}
	}

	if mergeOutput {
		mergeAndExport(activeFws)
		return
	}
	individualExport(activeFws)
}

func individualExport(activeFws []fwlib.Framework) {
	for _, fw := range activeFws {
		ms := fw.Mutations()
		marshal, err := json.Marshal(ms)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if output == "" {
			fmt.Println(string(marshal))
			continue
		}
		out := path.Join(output, fw.Meta().Name+".json")
		os.WriteFile(out, marshal, 0644)
	}
}

func mergeAndExport(activeFws []fwlib.Framework) {
	masterMs := mutations.Mutations{}
	for _, fw := range activeFws {
		masterMs.Merge(fw.Mutations())
	}
	marshal, err := json.Marshal(masterMs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if output == "" {
		fmt.Println(string(marshal))
		return
	}
	out := path.Join(output, "merged-fws-ouput.json")
	os.WriteFile(out, marshal, 0644)
}

func init() {
	exportCmd.Flags().BoolVarP(&mergeOutput, "merge-output", "m", false, "merges all frameworks output into one large json")
	exportCmd.Flags().StringVarP(&output, "output-path", "o", "", "specifies the output path")
	rootCmd.AddCommand(exportCmd)
}
