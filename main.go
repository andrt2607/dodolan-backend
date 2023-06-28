package main

import "dodolan/config"

func main() {
	db := config.ConnectDataBase()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}
