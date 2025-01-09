package controllers

import (
    "github.com/beego/beego/v2/server/web"
    "github.com/beego/beego/v2/client/orm"
    "booking-app/models"
    
)

type PropertyController struct {
    web.Controller
}

// ListProperties handles the GET request for property listings
func (c *PropertyController) ListProperties() {
    // Get query parameters
    location := c.GetString("location")
    minPrice, _ := c.GetFloat("min_price")
    maxPrice, _ := c.GetFloat("max_price")
    guests, _ := c.GetInt("guests")
    propertyType := c.GetString("property_type")

    // Initialize ORM
    o := orm.NewOrm()
    qs := o.QueryTable("rental_property")

    // Apply filters
    if location != "" {
        qs = qs.Filter("Location__icontains", location)
    }

    if minPrice > 0 {
        qs = qs.Filter("Price__gte", minPrice)
    }

    if maxPrice > 0 {
        qs = qs.Filter("Price__lte", maxPrice)
    }

    if guests > 0 {
        qs = qs.Filter("Guests__gte", guests)
    }

    if propertyType != "" {
        qs = qs.Filter("PropertyType", propertyType)
    }

    // Prepare result slice
    var properties []models.RentalProperty
    _, err := qs.All(&properties)

    if err != nil {
        c.Data["json"] = map[string]string{
            "error": "Failed to retrieve properties",
            "details": err.Error(),
        }
        c.Ctx.Output.SetStatus(500)
        c.ServeJSON()
        return
    }

    // If no properties found
    if len(properties) == 0 {
        c.Data["json"] = []map[string]interface{}{}
        c.ServeJSON()
        return
    }

    // Prepare enhanced property response
    var enhancedProperties []map[string]interface{}
    for _, prop := range properties {
        enhancedProp := map[string]interface{}{
            "HotelID":       prop.ID,
            "HotelName":     prop.Name,
            "Location":      prop.Location,
            "Price":         prop.Price,
            "PropertyType":  prop.PropertyType,
            "Guests":        prop.Guests,
            "Rating":        prop.Rating,
            "ReviewCount":   prop.ReviewCount,
            "Amenities":     prop.Amenities,
            "Description":   prop.Description,
        }
        enhancedProperties = append(enhancedProperties, enhancedProp)
    }

    c.Data["json"] = enhancedProperties
    c.ServeJSON()
}

// GetPropertyDetails retrieves detailed information about a specific property
func (c *PropertyController) GetPropertyDetails() {
    hotelID, _ := c.GetInt("id")

    o := orm.NewOrm()
    property := &models.RentalProperty{ID: hotelID}

    err := o.Read(property)
    if err == orm.ErrNoRows {
        c.Data["json"] = map[string]string{
            "error": "Property not found",
        }
        c.Ctx.Output.SetStatus(404)
    } else if err != nil {
        c.Data["json"] = map[string]string{
            "error": "Failed to retrieve property details",
            "details": err.Error(),
        }
        c.Ctx.Output.SetStatus(500)
    } else {
        // Prepare detailed property response
        propertyDetails := map[string]interface{}{
            "HotelID":       property.ID,
            "HotelName":     property.Name,
            "Location":      property.Location,
            "Price":         property.Price,
            "PropertyType":  property.PropertyType,
            "Guests":        property.Guests,
            "Rating":        property.Rating,
            "ReviewCount":   property.ReviewCount,
            "Amenities":     property.Amenities,
            "Description":   property.Description,
            // Add more details as needed
        }
        c.Data["json"] = propertyDetails
    }

    c.ServeJSON()
}