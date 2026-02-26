package cmds

import (
	"fmt"
	"os"

	"github.com/SecretSheppy/marv/internal/marvinfo"
	"github.com/spf13/cobra"
)

var (
	port                     int
	src, review              string
	toolRunner, fileWatchers bool

	rootCmd = &cobra.Command{
		Use:   "marv",
		Short: "marv is a tool that allows for efficient analysis and review of mutations through visualisations",
		Long: `Mutations Analysis, Review and Visualisation (Marv) is a tool that allows for efficient analysis and 
review of mutations through visualisations - it can be used 'as is' or can be integrated into a
third party application to streamline review processes`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: start main application server here
		},
	}
)

func Execute() {
	rootCmd.Version = marvinfo.Get().Version
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "port to listen on")
	rootCmd.Flags().StringVarP(&src, "src", "s", "", "source files directory")
	rootCmd.Flags().StringVarP(&review, "review", "r", "", "review output directory")
	rootCmd.Flags().BoolVarP(&toolRunner, "enable-tool-runner", "t", false, "enable tool runner")
	rootCmd.Flags().BoolVarP(&fileWatchers, "enable-watchers", "w", false, "enable file watchers")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
