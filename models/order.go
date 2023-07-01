package models

import "time"

type Order struct {
	OrderId     int `gorm:"primary_key" json:"id"`
	UidOrder    string
	OrderDate   time.Time   `json:"order_date"`
	TotalAmount float64     `json:"total_amount" validate:"required"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	CustomerId  int         `json:"customer_id" validate:"required"`
	Customer    Customer    `json:"-"`
	OrderItems  []OrderItem `json:"-"`
	Payments    []Payment   `json:"-"`
}
