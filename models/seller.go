package models

import (
	"dodolan/utils/token"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Seller struct {
	SellerId  int       `gorm:"primary_key" json:"id"`
	Name      string    `json:"name" validate:"required, min=4, max=15"`
	Username  string    `json:"username" validate:"required, min=4,max=16"`
	Password  string    `json:"password" validate:"required max=8"`
	Address   string    `json:"address" validate:"required"`
	Phone     string    `json:"phone" validate:"required, numeric, len=12"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Products  []Product `json:"-"`
}

func LoginCheckSeller(username string, password string, db *gorm.DB) (string, error) {

	var err error

	u := Seller{}

	err = db.Model(Customer{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(uint(u.SellerId), "SELLER")

	if err != nil {
		return "", err
	}

	return token, nil

}
