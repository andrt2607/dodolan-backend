package models

import (
	"dodolan/utils/token"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Seller struct {
	SellerId  int       `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Products  []Product `json:"-"`
}

// func VerifyPassword(password, hashedPassword string) error {
// 	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
// }

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

// func (u *Seller) SaveSellerAsUser(db *gorm.DB) (*Seller, error) {
// 	//turn password into hash
// 	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
// 	if errPassword != nil {
// 		return &Seller{}, errPassword
// 	}
// 	u.Password = string(hashedPassword)
// 	//remove spaces in username
// 	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

// 	var err error = db.Create(&u).Error
// 	if err != nil {
// 		return &Seller{}, err
// 	}
// 	return u, nil
// }
