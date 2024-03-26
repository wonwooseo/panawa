package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
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
	coll := r.cli.Database(r.database).Collection("date_prices")
	id := fmt.Sprintf("price#%s#%d", p.ItemCode, p.DateUnix)
	_, err := coll.UpdateByID(ctx, id, bson.D{{"$set", p}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) SaveRegionalMarketPrices(ctx context.Context, ps []*model.Price) error {
	coll := r.cli.Database(r.database).Collection("regional_market_prices")
	input := make([]mongo.WriteModel, len(ps))
	for i, p := range ps {
		id := fmt.Sprintf("price#%s#%d#%s#%s", p.ItemCode, p.DateUnix, *p.RegionCode, *p.MarketCode)
		input[i] = mongo.NewUpdateOneModel().SetFilter(bson.D{{"_id", id}}).SetUpdate(bson.D{{"$set", p}}).SetUpsert(true)
	}
	_, err := coll.BulkWrite(ctx, input)
	if err != nil {
		return err
	}
	return nil
}
