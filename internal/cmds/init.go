package cmds

import (
	"fmt"
	"os"

	"github.com/SecretSheppy/marv/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialises a new default marv.yml file",
	Long:  "initialises a new default marv.yml file in the current working directory",
	Run: func(cmd *cobra.Command, args []string) {
		marshal, err := yaml.Marshal(config.Init())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(".marv.yml", marshal, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
