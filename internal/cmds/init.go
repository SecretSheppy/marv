package cmds

import (
	"fmt"
	"os"

	"github.com/SecretSheppy/marv/fws"
	"github.com/SecretSheppy/marv/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const marvYml = ".marv.yml"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialises a new default marv.yml file",
	Long:  "initialises a new default marv.yml file in the current working directory",
	Run: func(cmd *cobra.Command, args []string) {
		initCommand()
	},
}

func initCommand() {
	if f, _ := os.Stat(marvYml); f != nil {
		fmt.Println(".marv.yml already exists in this directory")
		os.Exit(0)
	}

	marshal, err := yaml.Marshal(config.Init())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fwsMap := fws.FrameworksMap()
	for _, fw := range frameworks {
		if fwsMap[fw] == nil {
			fmt.Fprintf(os.Stderr, "Err: framework %s not recognised\n", fw)
		}

		fwMarshal, err := yaml.Marshal(fwsMap[fw].YamlInit())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		marshal = append(marshal, fwMarshal...)
	}

	if err := os.WriteFile(marvYml, marshal, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	initCmd.Flags().StringArrayVarP(&frameworks, "frameworks", "f", []string{}, "a list of frameworks to include in the initialisation")
	rootCmd.AddCommand(initCmd)
}
