package db

import (
	"context"

	"github.com/wonwooseo/panawa/pkg/db/model"
)

type Repository interface {
	SaveDatePrice(ctx context.Context, price *model.Price) error
	SaveRegionalMarketPrices(ctx context.Context, prices []*model.Price) error
}
