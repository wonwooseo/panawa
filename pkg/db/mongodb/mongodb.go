package mongodb

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/wonwooseo/panawa/pkg/db/model"
)

type Repository struct {
	logger zerolog.Logger
	cli    *mongo.Client

	database string
}

func NewRepository(baseLogger zerolog.Logger) *Repository {
	logger := baseLogger.With().Str("caller", "db/mongodb").Logger()

	url := viper.GetString("mongodb.url")
	database := viper.GetString("mongodb.database")

	cli, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create MongoDB client")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Fatal().Err(err).Msg("MongoDB server is not responding")
	}

	return &Repository{
		logger:   logger,
		cli:      cli,
		database: database,
	}
}

func (r *Repository) SaveDatePrice(ctx context.Context, p *model.Price) error {
	// TODO
	return nil
}

func (r *Repository) SaveRegionalMarketPrices(ctx context.Context, ps []*model.Price) error {
	// TODO
	return nil
}

func (r *Repository) Close() error {
	return r.cli.Disconnect(context.Background())
}
