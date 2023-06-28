package models

import "time"

type Order struct {
	OrderId     int         `gorm:"primary_key" json:"id"`
	OrderDate   time.Time   `json:"order_date"`
	TotalAmount float64     `json:"total_amount"`
	CustomerId  int         `json:"customer_id"`
	Customer    Customer    `json:"-"`
	OrderItems  []OrderItem `json:"-"`
	Payments    []Payment   `json:"-"`
}
