package models

import "time"

type Currency struct {
	ID        int       `json:"id" gorm:"type:integer;autoIncrement:true"`
	Name      string    `json:"name" gorm:"index; not null"`
	Code      string    `json:"code" gorm:"not null"`
	Value     float64   `json:"value" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type CurrencyAPIResponse struct {
	Meta struct {
		LastUpdatedAt time.Time `json:"last_updated_at"`
	} `json:"meta"`
	Data map[string]struct {
		Code  string  `json:"code"`
		Value float64 `json:"value"`
	} `json:"data"`
}

func (Currency) TableName() string {
	return "currency"
}
