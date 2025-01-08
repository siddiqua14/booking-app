package database

import (
	"fmt"

	"github.com/astaxie/beego"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=postgres dbname=booking_db port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Migrate tables
	DB.AutoMigrate(&Location{}, &RentalProperty{}, &PropertyDetails{})
	fmt.Println("Database connected and tables migrated!")
}

type Location struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string
	Hotels   string // JSON string of hotels
}

type RentalProperty struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string
	Type     string
	Bedrooms int
	Bathrooms int
	Amenities string // JSON string of amenities
}

type PropertyDetails struct {
	ID          uint   `gorm:"primaryKey"`
	PropertyID  uint
	Description string
	Images      string // JSON string of image URLs
}