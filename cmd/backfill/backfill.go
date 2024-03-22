package backfill

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func Command(baseLogger zerolog.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "backfill",
		Short: "backfill data",
		Long:  "backfill data of previous dates",
		Run: func(cmd *cobra.Command, args []string) {
			logger := baseLogger.With().Str("caller", "cmd/backfill").Logger()
			logger.Info().Msg("not implemented")
			// TODO: get previous data
		},
	}
}
