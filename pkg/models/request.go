package models

import "time"

type RequestHistory struct {
	ID           uint `gorm:"primaryKey"`
	Timestamp    time.Time
	Endpoint     string
	ResponseTime time.Duration
	StatusCode   int
	RequestBody  string `gorm:"type:json"`
	ResponseBody string `gorm:"type:json"`
}

func (RequestHistory) TableName() string {
	return "request_history"
}
