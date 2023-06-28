package models

type Product struct {
	ProductId   int     `gorm:"primary_key" json:"id"`
	Name        int     `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int64   `json:"stock"`
	SellerId    int     `json:"seller_id"`
	Seller      Seller  `json:"-"`
	// OrderItems  []OrderItem `json:"-"`
}
