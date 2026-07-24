package cmds

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/SecretSheppy/marv/internal/config"
	"github.com/SecretSheppy/marv/internal/themes"
	"github.com/SecretSheppy/marv/web"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var themesCmd = &cobra.Command{
	Use:   "themes",
	Short: "manage themes for marv",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listThemesCmd = &cobra.Command{
	Use:   "list",
	Short: "list all available themes for marv",
	Run: func(cmd *cobra.Command, args []string) {
		persistent := config.GetPersistentData(config.PersistentTheme)
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)

		for _, theme := range themes.List(web.ThemesFS) {
			var isSelected, isDefault string
			if theme == persistent {
				isSelected = "   ->   "
			}
			if theme == config.DefaultTheme {
				isDefault = "(default)"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", isSelected, theme, isDefault)
		}
		w.Flush()
	},
}

var setThemeCmd = &cobra.Command{
	Use:   "set",
	Short: "sets a persistent theme for marv",
	Long:  "sets a persistent theme for marv that is used every time you run `marv`. any persistent theme can still be overridden by the --theme flag",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error().Msg("must specify a theme")
			return
		}
		for _, theme := range themes.List(web.ThemesFS) {
			if theme == args[0] {
				config.SetPersistentData(config.PersistentTheme, args[0])
				return
			}
		}
		log.Error().Msgf("invalid theme '%s'", args[0])
	},
}

func init() {
	themesCmd.AddCommand(listThemesCmd)
	themesCmd.AddCommand(setThemeCmd)

	rootCmd.AddCommand(themesCmd)
}
