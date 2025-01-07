package main

import (
	"booking-app/models"
	_ "booking-app/routers"

	"github.com/beego/beego/v2/server/web"
	"log"
)

func main() {
	// Initialize database and ensure the locations table exists
	err := models.InitializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	web.Run()
}

