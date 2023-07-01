package models

import (
	"dodolan/utils/token"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Customer struct {
	CustomerId int       `gorm:"primary_key" json:"id"`
	FirstName  string    `json:"first_name" validate:"required, min=4,max=15"`
	LastName   string    `json:"last_name" validate:"required, min=4,max=15"`
	Email      string    `json:"email" validate:"required, email"`
	Username   string    `json:"username" validate:"required, min=4,max=16"`
	Password   string    `json:"password" validate:"required max=8"`
	Address    string    `json:"address" validate:"required"`
	City       string    `json:"city" validate:"required"`
	Country    string    `json:"country" validate:"required"`
	PostalCode string    `json:"postal_code" validate:"required, len=5"`
	Phone      string    `json:"phone" validate:"required, numeric, len=12"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Orders     []Order   `json:"-"`
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheckCustomer(username string, password string, db *gorm.DB) (string, error) {

	var err error

	u := Customer{}

	err = db.Model(Customer{}).Where("username = ?", username).Take(&u).Error

	fmt.Println("ini pw dari db : ", u.Password)
	fmt.Println("ini pw input : ", password)

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(uint(u.CustomerId), "CUSTOMER")

	if err != nil {
		return "", err
	}

	return token, nil

}
