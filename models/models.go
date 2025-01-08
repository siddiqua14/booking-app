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
    Id        int64  `orm:"auto;pk"`
    HotelId   *Hotel `orm:"rel(fk)"`
    Name      string `orm:"size(255);null"`
    Type      string `orm:"size(100);null"`
    Bedrooms  int    
    Bathrooms int    
    Amenities string `orm:"type(text);null"`
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