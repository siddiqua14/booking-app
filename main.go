package main

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
    _ "github.com/lib/pq"
	_ "booking-app/routers"  // Import routers to initialize routes
	//"booking-app/models"      // Import models for database
)

func main() {
    orm.RegisterDriver("postgres", orm.DRPostgres)

    // Database connection setup
    orm.RegisterDataBase("default", "postgres", "user=postgres password=postgres host=localhost port=5432 dbname=booking_db sslmode=disable")

    orm.RunSyncdb("default", false, true) // Automatically create the tables in the database

    // Start the Beego server
    web.Run()
}
