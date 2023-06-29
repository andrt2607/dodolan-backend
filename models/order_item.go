package models

import "time"

type OrderItem struct {
	OrderItemId  int `gorm:"primary_key" json:"id"`
	UidOrderItem string
	Quantity     int64     `json:"quantity"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	OrderId      int       `json:"order_id"`
	Order        Order     `json:"-"`
	ProductId    int       `json:"product_id"`
	Product      Product   `json:"-"`
}
