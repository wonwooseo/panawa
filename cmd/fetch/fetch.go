package fetch

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func Command(baseLogger zerolog.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "fetch data",
		Long:  "fetch current date's data",
		Run: func(cmd *cobra.Command, args []string) {
			logger := baseLogger.With().Str("caller", "cmd/fetch").Logger()
			logger.Info().Msg("not implemented")
			// TODO: get data
		},
	}
}
