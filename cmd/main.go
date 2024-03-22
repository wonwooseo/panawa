package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wonwooseo/panawa/build"
	"github.com/wonwooseo/panawa/cmd/backfill"
	"github.com/wonwooseo/panawa/cmd/fetch"
)

func main() {
	baseLogger := log.Logger
	logger := baseLogger.With().Str("caller", "main").Logger()

	logger.Info().Str("version", build.Version).Str("build_time", build.BuildTime).Msg("batch is starting..")

	rootCmd := cobra.Command{
		Short: "panawa",
		Long:  "panawa batch",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cfgPath, err := cmd.Flags().GetString("config")
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to get config flag")
			}
			viper.SetConfigFile(cfgPath)
			if err := viper.ReadInConfig(); err != nil {
				logger.Fatal().Err(err).Msg("failed to read in config")
			}
		},
	}
	rootCmd.PersistentFlags().String("config", "", "path to config file")

	rootCmd.AddCommand(fetch.Command(baseLogger))
	rootCmd.AddCommand(backfill.Command(baseLogger))
	rootCmd.Execute()
}
