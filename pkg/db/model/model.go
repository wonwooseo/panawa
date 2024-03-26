package model

type Price struct {
	ItemCode string `bson:"item_code"`
	DateUnix int64  `bson:"date_unix"`
	Low      int    `bson:"low"`
	Average  int    `bson:"average"`
	High     int    `bson:"high"`

	RegionCode     *string `bson:"region_code,omitempty"`
	MarketCode     *string `bson:"market_code,omitempty"`
	UpdateTimeUnix int64   `bson:"update_time_unix"`
}
