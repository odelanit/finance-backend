package Models

import "time"

type Candle struct {
	ID        uint      `json:"id"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Time      uint      `json:"time"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *Candle) TableName() string {
	return "stock_data"
}

type Price struct {
	ID        uint      `json:"id"`
	Litecoin  float64   `json:"litecoin"`
	Monero    float64   `json:"monero"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
