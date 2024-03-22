package model

import "time"

type Price struct {
	ItemCode string `json:"item_code"`
	Low      int    `json:"low"`
	Average  int    `json:"average"`
	High     int    `json:"high"`

	RegionCode *string   `json:"region_code,omitempty"`
	MarketCode *string   `json:"market_code,omitempty"`
	UpdateTime time.Time `json:"update_time"`
}

func (p Price) StringDateFmt(fmt string) string {
	return p.UpdateTime.Format(fmt)
}
