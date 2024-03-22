package fetch

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wonwooseo/panawa/pkg/price"
	"github.com/wonwooseo/panawa/pkg/price/kamis"
)

func Command(baseLogger zerolog.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "fetch data",
		Long:  "fetch current date's data",
		Run: func(cmd *cobra.Command, args []string) {
			logger := baseLogger.With().Str("caller", "cmd/fetch").Logger()
			logger.Info().Msg("not implemented")

			apiURL := viper.GetString("api.kamis.url")
			var priceFetcher price.DataClient = kamis.NewDataClient(baseLogger, apiURL)

			datePrice, regionalMarketPrices, err := priceFetcher.GetDatePrices(context.Background(), time.Date(2024, 3, 20, 6, 0, 0, 0, time.UTC), "0000")
			if err != nil {
				logger.Error().Err(err).Msg("failed to fetch price data")
				return
			}

			logger.Info().Any("date_price", datePrice).Send()
			for region, marketPrices := range regionalMarketPrices {
				logger.Info().Str("region_code", region).Any("market_prices", marketPrices).Send()
			}
		},
	}
}
