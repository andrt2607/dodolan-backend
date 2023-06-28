package models

type OrderItem struct {
	OrderItemId int     `gorm:"primary_key" json:"id"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
	OrderId     int     `json:"order_id"`
	Order       Order   `json:"-"`
	ProductId   int     `json:"product_id"`
	Product     Product `json:"-"`
}
