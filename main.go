package main

import (
    "github.com/beego/beego/v2/client/orm"
    _ "github.com/lib/pq"
    "github.com/beego/beego/v2/server/web"
    //"booking-app/models"
    _ "booking-app/routers"
)

func main() {
    orm.RegisterDriver("postgres", orm.DRPostgres)

    orm.RegisterDataBase("default", "postgres", "user=postgres password=postgres host=localhost port=5432 dbname=booking_db sslmode=disable")

    orm.RunSyncdb("default", false, true)

    web.Run()
}