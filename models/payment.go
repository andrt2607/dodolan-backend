package models

import "time"

type Payment struct {
	PaymentId   int `gorm:"primary_key" json:"id"`
	UidPayment  string
	PaymentDate time.Time `json:"payment_date"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	OrderId     int       `json:"order_id"`
	Order       Order     `json:"-"`
}
