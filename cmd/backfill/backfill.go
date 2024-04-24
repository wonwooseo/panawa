package backfill

import (
	"context"
	"errors"
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
	cmd := &cobra.Command{
		Use:   "backfill",
		Short: "backfill data",
		Long:  "backfill data of previous dates",
		PreRun: func(cmd *cobra.Command, args []string) {
			sDate, _ := cmd.Flags().GetString("backfill.startdate")
			if sDate != "" {
				viper.Set("backfill.startdate", sDate)
			}
			eDate, _ := cmd.Flags().GetString("backfill.enddate")
			if eDate != "" {
				viper.Set("backfill.enddate", eDate)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger := baseLogger.With().Str("caller", "cmd/backfill").Logger()

			var regionCodeResolver code.Resolver = code.NewRegionCodeResolver()
			var priceFetcher price.DataClient = kamis.NewDataClient(baseLogger, regionCodeResolver)
			var repository db.Repository = mongodb.NewRepository(baseLogger)

			kstLoc := time.FixedZone("KST", 9*60*60) // UTC+09:00
			sDate, err := time.ParseInLocation("2006-01-02", viper.GetString("backfill.startdate"), kstLoc)
			if err != nil {
				logger.Error().Err(err).Str("startdate", viper.GetString("backfill.startdate")).Msg("failed to parse start date")
				return
			}
			eDate, err := time.ParseInLocation("2006-01-02", viper.GetString("backfill.enddate"), kstLoc)
			if err != nil {
				logger.Error().Err(err).Str("enddate", viper.GetString("backfill.enddate")).Msg("failed to parse end date")
				return
			}

			fetchCodes := viper.GetStringSlice("backfill.codes")
			for _, itemCode := range fetchCodes {
				for cDate := sDate; cDate.Unix() <= eDate.Unix(); cDate = cDate.AddDate(0, 0, 1) {
					datePrice, regionalMarketPrices, err := priceFetcher.GetDatePrices(context.Background(), cDate, itemCode)
					if err != nil {
						if errors.Is(err, price.ErrPriceDataNotFound) {
							continue
						}
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

					logger.Info().Str("item_code", itemCode).Str("date", cDate.Format("2006-01-02")).Msg("saved prices to DB")
				}
			}
		},
	}
	cmd.PersistentFlags().String("backfill.startdate", "", "start date of backfill")
	cmd.PersistentFlags().String("backfill.enddate", "", "end date of backfill")

	return cmd
}
