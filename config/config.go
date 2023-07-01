package config

import (
	"dodolan/models"
	"dodolan/utils"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDataBase() *gorm.DB {
	username := utils.Getenv("DATABASE_USERNAME", "root")
	password := utils.Getenv("DATABASE_PASSWORD", "")
	host := utils.Getenv("DATABASE_HOST", "127.0.0.1")
	port := utils.Getenv("DATABASE_PORT", "3306")
	database := utils.Getenv("DATABASE_NAME", "db_dodolan")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.Customer{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Seller{},
		&models.Product{})

	return db
}
