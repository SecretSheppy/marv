package cmds

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:     "stats",
	Aliases: []string{"statistics"},
	Short:   "generates a statistics json for all specified marv mutations schema jsons",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal().Msg("must specify at least one file name")
		}

		for _, arg := range args {
			log.Info().Str("file", arg).Msg("processing file")

			raw, err := os.ReadFile(arg)
			if err != nil {
				log.Fatal().Err(err).Msgf("failed to read file '%s'", arg)
			}

			var ms mutations.Mutations
			if err := json.Unmarshal(raw, &ms); err != nil {
				log.Fatal().Err(err).Msgf("failed to parse file '%s'", arg)
			}

			stats := ms.StatisticsFrom("")
			json, err := stats.MarshalJSON()
			if err != nil {
				log.Fatal().Err(err).Msgf("failed to marshal json for file '%s'", arg)
			}

			filename := strings.TrimRight(arg, ".json")
			if err := os.WriteFile(fmt.Sprintf("%s-stats.json", filename), json, 0644); err != nil {
				log.Error().Err(err).Msgf("failed to write json for file '%s'", arg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
