package models

type Seller struct {
	SellerId int       `gorm:"primary_key" json:"id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Phone    string    `json:"phone"`
	Products []Product `json:"-"`
}
