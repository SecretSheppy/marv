package cmds

import (
	"os"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/fws"
	"github.com/SecretSheppy/marv/internal/config"
	"github.com/SecretSheppy/marv/internal/marvinfo"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const defaultPort = 8080

var (
	port                                  int
	review, configFile, output            string
	frameworks                            []string
	mergeOutput, toolRunner, fileWatchers bool

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

func rootCommand() {
	_, activeFws := getConfigAndFws()

	for _, fw := range activeFws {
		if decompiling, ok := fw.(fwlib.Decompiling); ok {
			decompiling.SetDecompiler()
		}

		if err := fw.TransformResults(); err != nil {
			log.Fatal().Err(err)
			os.Exit(0)
		}
	}

	// TODO: start main application server here
}

func getConfigAndFws() (*config.Config, []fwlib.Framework) {
	yml, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
	}

	cfg, err := mergeYmlFlagConfigs(yml)
	if err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
	}

	activeFws := make([]fwlib.Framework, 0)
	for _, fw := range fws.Frameworks() {
		loaded, err := fw.Yaml().Load(yml)
		if err != nil {
			log.Fatal().Err(err)
			os.Exit(1)
		}
		if !loaded {
			continue
		}
		if err := fw.LoadResults(); err != nil {
			log.Fatal().Err(err)
			os.Exit(1)
		}
		activeFws = append(activeFws, fw)
	}
	return cfg, activeFws
}

func mergeYmlFlagConfigs(yml []byte) (*config.Config, error) {
	cfg := &config.Config{}
	if err := yaml.Unmarshal(yml, cfg); err != nil {
		return nil, err
	}
	if port != defaultPort {
		cfg.Marv.Port = port
	}
	if review != "" {
		cfg.Marv.ReviewDir = review
	}
	if toolRunner {
		cfg.Marv.Features.ToolRunner = true
	}
	if fileWatchers {
		cfg.Marv.Features.FileWatchers = true
	}
	return cfg, nil
}

func Execute() {
	rootCmd.Version = marvinfo.Get().Version
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", ".marv.yml", ".marv.yml file path")
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "port to listen on")
	rootCmd.Flags().StringVarP(&review, "review", "r", "", "review output directory")
	rootCmd.Flags().BoolVarP(&toolRunner, "enable-tool-runner", "t", false, "enable tool runner")
	rootCmd.Flags().BoolVarP(&fileWatchers, "enable-watchers", "w", false, "enable file watchers")

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err)
		os.Exit(1)
	}
}
