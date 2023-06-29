package models

import "time"

type Product struct {
	ProductId   int         `gorm:"primary_key" json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       float64     `json:"price"`
	Stock       int64       `json:"stock"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	SellerId    int         `json:"seller_id"`
	Seller      Seller      `json:"-"`
	OrderItems  []OrderItem `json:"-"`
}
