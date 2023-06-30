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
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	PostalCode string    `json:"postal_code"`
	Phone      string    `json:"phone"`
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

// func (u *Customer) SaveCustomerAsUser(db *gorm.DB) (*Customer, error) {
// 	//turn password into hash
// 	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
// 	if errPassword != nil {
// 		return &Customer{}, errPassword
// 	}
// 	u.Password = string(hashedPassword)
// 	//remove spaces in username
// 	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

// 	var err error = db.Create(&u).Error
// 	if err != nil {
// 		return &Customer{}, err
// 	}
// 	return u, nil
// }
