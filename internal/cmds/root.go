package cmds

import (
	"encoding/json"
	"os"
	"path"
	"strconv"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/fws"
	"github.com/SecretSheppy/marv/internal/config"
	"github.com/SecretSheppy/marv/internal/marvinfo"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/internal/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	port, configFile, outputPath string
	mergeOutput                  bool
	frameworks                   []string

	rootCmd = &cobra.Command{
		Use:   "marv",
		Short: "marv is a tool that allows for efficient analysis and review of mutations through visualisations",
		Long: `Mutations Analysis, Review and Visualisation (Marv) is a tool that allows for efficient analysis and 
review of mutations through visualisations - it can be used 'as is' or can be integrated into a
third party application to streamline review processes`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCommand()
		},
	}
)

func getMarvYml() []byte {
	yml, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to find or read file")
		os.Exit(1)
	}
	return yml
}

func getConfig(yml []byte) *config.Config {
	cfg := &config.Config{}
	if err := yaml.Unmarshal(yml, cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
		os.Exit(1)
	}

	if err := mergeFlagsWithConfig(cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to merge config and flags")
		os.Exit(1)
	}

	return cfg
}

func mergeFlagsWithConfig(cfg *config.Config) error {
	if port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		cfg.Marv.Port = p
	}
	if cfg.Marv.Port == 0 {
		cfg.Marv.Port = config.DefaultPort
	}
	if outputPath != "" {
		cfg.Marv.Output.Path = outputPath
	}
	if cfg.Marv.Output.Path == "" {
		cfg.Marv.Output.Path = config.DefaultPath
	}
	if mergeOutput && !cfg.Marv.Output.Merge {
		cfg.Marv.Output.Merge = mergeOutput
	}
	return nil
}

func transformMutations(activeFws []fwlib.Framework) error {
	for _, fw := range activeFws {
		if decompiling, ok := fw.(fwlib.Decompiling); ok {
			decompiling.SetDecompiler()
		}

		if err := fw.TransformResults(); err != nil {
			return err
		}

		fw.Mutations().GenerateIDs()
	}
	return nil
}

func export(conf *config.Config, activeFws []fwlib.Framework) error {
	for _, fw := range activeFws {
		marshal, err := json.Marshal(fw.Mutations())
		if err != nil {
			return err
		}

		out := path.Join(conf.Marv.Output.Path, fw.Meta().Name+".json")
		if err := os.WriteFile(out, marshal, 0644); err != nil {
			return err
		}
	}
	return nil
}

func mergeAndExport(conf *config.Config, activeFws []fwlib.Framework) error {
	merged := mutations.Mutations{}
	for _, fw := range activeFws {
		merged.Merge(fw.Mutations())
	}

	marshal, err := json.Marshal(merged)
	if err != nil {
		return err
	}

	out := path.Join(conf.Marv.Output.Path, "mutations.json")
	return os.WriteFile(out, marshal, 0644)
}

func exportTransformedMutations(conf *config.Config, activeFws []fwlib.Framework) error {
	if err := os.MkdirAll(conf.Marv.Output.Path, 0755); err != nil {
		return err
	}
	if conf.Marv.Output.Merge {
		return mergeAndExport(conf, activeFws)
	}
	return export(conf, activeFws)
}

func exportCommand() (*config.Config, []fwlib.Framework) {
	yml := getMarvYml()
	conf := getConfig(yml)

	activeFws, err := fws.ActiveFrameworks(yml)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get active frameworks")
		os.Exit(1)
	}

	if err := transformMutations(activeFws); err != nil {
		log.Fatal().Err(err).Msg("Failed to transform results")
		os.Exit(1)
	}

	if err := exportTransformedMutations(conf, activeFws); err != nil {
		log.Fatal().Err(err).Msg("Failed to export mutations")
		os.Exit(1)
	}
	return conf, activeFws
}

func rootCommand() {
	conf, activeFws := exportCommand()

	log.Info().Msgf("Starting server at http://localhost:%d/", conf.Marv.Port)
	if err := server.NewServer(conf.Marv.Port, activeFws).Serve(); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
		os.Exit(1)
	}
}

func Execute() {
	rootCmd.Version = marvinfo.Get().Version
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", ".marv.yml file path")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", "", "specifies the output path")
	rootCmd.PersistentFlags().BoolVarP(&mergeOutput, "merge", "m", false, "merges all frameworks output into one large json")

	rootCmd.Flags().StringVarP(&port, "port", "p", "", "port to listen on")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute marv command")
		os.Exit(1)
	}
}
