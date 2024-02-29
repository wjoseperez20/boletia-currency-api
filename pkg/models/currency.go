package models

type Currency struct {
	ID    int     `json:"id" gorm:"primary_key"`
	Name  string  `json:"name" gorm:"unique"`
	Code  string  `json:"code" gorm:"unique"`
	Value float64 `json:"value" gorm:"not null"`
}
