package main

import (
	"dodolan/config"
	"dodolan/routes"
)

func main() {
	db := config.ConnectDataBase()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	r := routes.SetupRouter(db)
	r.Run()
}
