package models

type RequestHistory struct {
	ID           uint `gorm:"primaryKey"`
	Endpoint     string
	ResponseTime float64
	StatusCode   int
}

func (RequestHistory) TableName() string {
	return "request_history"
}
