package cmds

import (
	"os"

	"github.com/SecretSheppy/marv/fws"
	"github.com/SecretSheppy/marv/internal/config"
	"github.com/rs/zerolog/log"
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
		log.Warn().Msg(".marv.yml already exists in this directory")
		os.Exit(0)
	}

	marshal, err := yaml.Marshal(config.Init())
	if err != nil {
		log.Error().Err(err)
		os.Exit(1)
	}

	fwsMap := fws.FrameworksMap()
	for _, fw := range frameworks {
		if fwsMap[fw] == nil {
			log.Warn().Msgf("skipping unknown framework %s\n", fw)
			continue
		}

		fwMarshal, err := yaml.Marshal(fwsMap[fw].Yaml().Init())
		if err != nil {
			log.Error().Err(err)
			os.Exit(1)
		}

		marshal = append(marshal, fwMarshal...)
	}

	if err := os.WriteFile(marvYml, marshal, 0644); err != nil {
		log.Error().Err(err)
		os.Exit(1)
	}
}

func init() {
	initCmd.Flags().StringArrayVarP(&frameworks, "frameworks", "f", []string{}, "a list of frameworks to include in the initialisation")
	rootCmd.AddCommand(initCmd)
}
