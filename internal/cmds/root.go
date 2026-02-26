package cmds

import (
	"fmt"
	"os"

	"github.com/SecretSheppy/marv/internal/config"
	"github.com/SecretSheppy/marv/internal/marvinfo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const defaultPort = 8080

var (
	port                     int
	src, review, configFile  string
	toolRunner, fileWatchers bool

	rootCmd = &cobra.Command{
		Use:   "marv",
		Short: "marv is a tool that allows for efficient analysis and review of mutations through visualisations",
		Long: `Mutations Analysis, Review and Visualisation (Marv) is a tool that allows for efficient analysis and 
review of mutations through visualisations - it can be used 'as is' or can be integrated into a
third party application to streamline review processes`,
		Run: func(cmd *cobra.Command, args []string) {
			yml, err := os.ReadFile(configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			_, err = mergeYmlFlagConfigs(yml)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// TODO: start main application server here
		},
	}
)

func mergeYmlFlagConfigs(yml []byte) (*config.Config, error) {
	cfg := &config.Config{}
	if err := yaml.Unmarshal(yml, cfg); err != nil {
		return nil, err
	}
	if port != defaultPort {
		cfg.Web.Port = port
	}
	if src != "" {
		cfg.Paths.Sources = src
	}
	if review != "" {
		cfg.Paths.Reviews = review
	}
	if toolRunner {
		cfg.Features.ToolRunner = true
	}
	if fileWatchers {
		cfg.Features.FileWatchers = true
	}
	return cfg, nil
}

func Execute() {
	rootCmd.Version = marvinfo.Get().Version
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "port to listen on")
	rootCmd.Flags().StringVarP(&src, "src", "s", "", "source files directory")
	rootCmd.Flags().StringVarP(&review, "review", "r", "", "review output directory")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", ".marv.yml", ".marv.yml file path")
	rootCmd.Flags().BoolVarP(&toolRunner, "enable-tool-runner", "t", false, "enable tool runner")
	rootCmd.Flags().BoolVarP(&fileWatchers, "enable-watchers", "w", false, "enable file watchers")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
