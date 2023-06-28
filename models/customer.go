package models

type Customer struct {
	CustomerId int     `gorm:"primary_key" json:"id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Address    string  `json:"address"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	PostalCode string  `json:"postal_code"`
	Phone      string  `json:"phone"`
	Orders     []Order `json:"-"`
}
