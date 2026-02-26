package cmds

import (
	"fmt"
	"os"

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
		if f, _ := os.Stat(marvYml); f != nil {
			fmt.Println(".marv.yml already exists in this directory")
			os.Exit(0)
		}
		marshal, err := yaml.Marshal(config.Init())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(marvYml, marshal, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
