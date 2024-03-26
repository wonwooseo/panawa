package model

import "time"

type Price struct {
	ItemCode string `bson:"item_code"`
	Low      int    `bson:"low"`
	Average  int    `bson:"average"`
	High     int    `bson:"high"`

	RegionCode *string   `bson:"region_code,omitempty"`
	MarketCode *string   `bson:"market_code,omitempty"`
	UpdateTime time.Time `bson:"update_time"`
}

func (p Price) StringDateFmt(fmt string) string {
	return p.UpdateTime.Format(fmt)
}
