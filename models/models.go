// models/models.go
package models

import (
    "github.com/beego/beego/v2/client/orm"
)

type Location struct {
    Id            int64  `orm:"auto;pk"`
    DestId        string `orm:"size(100);column(dest_id)"`
    DestType      string `orm:"size(100);column(dest_type)"`
    Value         string `orm:"size(255);column(value)"`
}

type Hotel struct {
    Id          int64     `orm:"auto;pk"`
    LocationId  *Location `orm:"rel(fk)"`
    Name        string    `orm:"size(255);null"`
    HotelId     string    `orm:"size(100);null"`
    Rating      float64   
    ReviewScore float64   
}

type RentalProperty struct {
    ID            int     `orm:"auto"`
    IDHotel       int     `orm:"index"`
    HotelID       string  `orm:"size(100)"`
    HotelName     string  `orm:"size(255)"`
    DestID        string  `orm:"size(100)"`
    Location      string  `orm:"size(255)"`
    Rating        float64 `orm:"digits(2);decimals(1)"`
    ReviewCount   int
    Price         string  `orm:"size(100)"`
    NumBeds       int
    NumBedR       int
    NumBaths      int
    Bedroom       int
    Guests        int
    PropertyType  string  `orm:"size(100)"`
    Amenities     string  `orm:"type(text)"`
}

type PropertyDetails struct {
    Id          int64          `orm:"auto;pk"`
    PropertyId  *RentalProperty `orm:"rel(fk)"`
    Description string         `orm:"type(text);null"`
    Images      string         `orm:"type(text);null"` // JSON array of image URLs
}

func init() {
    orm.RegisterModel(new(Location), new(Hotel), new(RentalProperty), new(PropertyDetails))
}