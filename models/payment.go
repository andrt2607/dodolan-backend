package models

import "time"

type Payment struct {
	PaymentId   int       `gorm:"primary_key" json:"id"`
	PaymentDate time.Time `json:"payment_date"`
	Amount      float64   `json:"amount"`
	OrderId     int       `json:"order_id"`
	Order       Order     `json:"-"`
}
