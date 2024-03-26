package fetch

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wonwooseo/panawa/pkg/code"
	"github.com/wonwooseo/panawa/pkg/db"
	"github.com/wonwooseo/panawa/pkg/db/mongodb"
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

			var regionCodeResolver code.Resolver = code.NewRegionCodeResolver()
			var priceFetcher price.DataClient = kamis.NewDataClient(baseLogger, regionCodeResolver)
			var repository db.Repository = mongodb.NewRepository(baseLogger)

			kstLoc := time.FixedZone("KST", 9*60*60) // UTC+09:00
			kstYesterday := time.Now().UTC().In(kstLoc).AddDate(0, 0, -1)

			fetchCodes := viper.GetStringSlice("fetch.codes")
			for _, itemCode := range fetchCodes {
				datePrice, regionalMarketPrices, err := priceFetcher.GetDatePrices(context.Background(), kstYesterday, itemCode)
				if err != nil {
					logger.Error().Str("item_code", itemCode).Err(err).Msg("failed to fetch price data")
					return
				}

				if err := repository.SaveDatePrice(context.Background(), datePrice); err != nil {
					logger.Error().Str("item_code", itemCode).Err(err).Msg("failed to save date price to DB")
					return
				}
				for region, marketPrices := range regionalMarketPrices {
					if err := repository.SaveRegionalMarketPrices(context.Background(), marketPrices); err != nil {
						logger.Error().Str("item_code", itemCode).Str("region_code", region).Err(err).Msg("failed to save regional market prices to DB")
						return
					}
				}

				logger.Info().Str("item_code", itemCode).Msg("saved prices to DB")
			}
		},
	}
}
