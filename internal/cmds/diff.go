package cmds

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/SecretSheppy/marv/internal/mutations"
	mapset "github.com/deckarep/golang-set/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var before, after string

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "displays the Marv mutations diff web interface",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: this will ultimately display the mutations diff web interface + export mutations diff file
	},
}

type Change struct {
	ID        mutations.MutationID `json:"id"`
	OldStatus mutations.Status     `json:"old_status"`
	NewStatus mutations.Status     `json:"new_status"`
}

var diffExportCmd = &cobra.Command{
	Use:   "export",
	Short: "produces a Marv diff file from two Marv mutation schema files",
	Run: func(cmd *cobra.Command, args []string) {
		if !strings.HasSuffix(before, ".json") {
			log.Error().Msg("--before (-b) must be a valid JSON file")
			os.Exit(1)
		}

		if !strings.HasSuffix(after, ".json") {
			log.Error().Msg("--after (-a) must be a valid JSON file")
			os.Exit(1)
		}

		b, err := os.ReadFile(before)
		if err != nil {
			log.Fatal().Err(err).Msg("error reading file")
		}

		a, err := os.ReadFile(after)
		if err != nil {
			log.Fatal().Err(err).Msg("error reading file")
		}

		var beforeResults, afterResults mutations.Mutations
		if err := json.Unmarshal(b, &beforeResults); err != nil {
			log.Fatal().Err(err).Msg("error unmarshalling JSON")
		}
		if err := json.Unmarshal(a, &afterResults); err != nil {
			log.Fatal().Err(err).Msg("error unmarshalling JSON")
		}

		beforeIDs := mapset.NewSet[mutations.MutationID](beforeResults.GetIDs()...)
		afterIDs := mapset.NewSet[mutations.MutationID](afterResults.GetIDs()...)

		intersection := beforeIDs.Intersect(afterIDs)
		changes := make([]Change, 0)
		for _, id := range intersection.ToSlice() {
			beforeStatus := beforeResults.GetMutationByID(id).Status
			afterStatus := afterResults.GetMutationByID(id).Status
			if beforeStatus != afterStatus {
				changes = append(changes, Change{
					ID:        id,
					OldStatus: beforeStatus,
					NewStatus: afterStatus,
				})
			}
		}

		marhsal, err := json.Marshal(struct {
			Removed []mutations.MutationID `json:"removed"`
			Added   []mutations.MutationID `json:"added"`
			Changed []Change               `json:"changed"`
		}{
			Removed: beforeIDs.Difference(intersection).ToSlice(),
			Added:   afterIDs.Difference(intersection).ToSlice(),
			Changed: changes,
		})

		if err = os.WriteFile("diff.json", marhsal, 0644); err != nil {
			log.Fatal().Err(err).Msg("error writing file")
		}
	},
}

func init() {
	diffCmd.PersistentFlags().StringVarP(&before, "before", "b", "", "before schema file name")
	diffCmd.PersistentFlags().StringVarP(&after, "after", "a", "", "after schema file name")

	diffCmd.AddCommand(diffExportCmd)

	rootCmd.AddCommand(diffCmd)
}
