package price

import (
	"context"
	"time"

	"github.com/wonwooseo/panawa/pkg/db/model"
)

type DataClient interface {
	GetDatePrices(ctx context.Context, date time.Time, itemCode string) (datePrice *model.Price, regionalMarketPrices map[string][]*model.Price, err error)
}
