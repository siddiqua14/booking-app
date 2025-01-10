package controllers

import (
    "strings"
    "github.com/beego/beego/v2/server/web"
    "github.com/beego/beego/v2/client/orm"
    "booking-app/models"
)

type PropertyController struct {
    web.Controller
}

// ListProperties handles the GET request for property listings
func (c *PropertyController) ListProperties() {
    location := c.GetString("location", "New York")

    // Initialize ORM
    o := orm.NewOrm()
    qs := o.QueryTable("rental_property")

    // Apply location filter
    qs = qs.Filter("Location__icontains", location)

    // Prepare result slice
    var properties []models.RentalProperty
    _, err := qs.All(&properties)

    if err != nil {
        c.Data["json"] = map[string]string{
            "error": "Failed to retrieve properties",
            "details": err.Error(),
        }
        c.ServeJSON()
        return
    }

    // Prepare enhanced property response
    var enhancedProperties []map[string]interface{}
    for _, prop := range properties {
        var details models.PropertyDetails
        err := o.QueryTable("property_details").Filter("HotelID", prop.HotelID).One(&details)

        // Handle the error for property details query
        if err != nil {
            // If there's an error, use default/empty values for details
            details = models.PropertyDetails{
                ImageUrl1:  "", // default empty image URL
                ImageUrl2:  "",
                ImageUrl3:  "",
                ImageUrl4:  "",
                ImageUrl5:  "",
                Description: "No description available",
            }
        }

        enhancedProp := map[string]interface{}{
            "ID":           prop.ID,
            "HotelID":      prop.HotelID,
            "HotelName":    prop.HotelName,
            "Location":     prop.Location,
            "Price":        prop.Price,
            "PropertyType": prop.PropertyType,
            "Guests":       prop.Guests,
            "Rating":       prop.Rating,
            "ReviewCount":  prop.ReviewCount,
            "NumBeds":      prop.NumBeds,
            "NumBedR":      prop.NumBedR,
            "NumBaths":     prop.NumBaths,
            "Bedroom":      prop.Bedroom,
            "Amenities":    strings.Split(prop.Amenities, ","),
            "HeroImage":    details.ImageUrl1,
        }
        enhancedProperties = append(enhancedProperties, enhancedProp)
    }

    // Pass data to the template
    c.Data["Properties"] = enhancedProperties
    c.TplName = "property_list.tpl"
    c.Render()
}

// GetPropertyDetails retrieves detailed information about a specific property
func (c *PropertyController) GetPropertyDetails() {
    propertyID := c.GetString("id")
    if propertyID == "" {
        c.Ctx.Output.Status = 400
        c.Data["json"] = map[string]string{"error": "Property ID is required"}
        c.ServeJSON()
        return
    }

    o := orm.NewOrm()
    
    var property models.RentalProperty
    err := o.QueryTable("rental_property").Filter("HotelID", propertyID).One(&property)
    if err != nil {
        c.Ctx.Output.Status = 500
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
        return
    }

    var details models.PropertyDetails
    err = o.QueryTable("property_details").Filter("HotelID", propertyID).One(&details)
    if err != nil {
        c.Ctx.Output.Status = 500
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
        return
    }

    enhancedProperty := map[string]interface{}{
        "HotelID":      property.HotelID,
        "HotelName":    property.HotelName,
        "Location":     property.Location,
        "Price":        property.Price,
        "PropertyType": property.PropertyType,
        "Guests":       property.Guests,
        "Rating":       property.Rating,
        "ReviewCount":  property.ReviewCount,
        "NumBeds":      property.NumBeds,
        "NumBedR":      property.NumBedR,
        "NumBaths":     property.NumBaths,
        "Bedroom":      property.Bedroom,
        "Amenities":    strings.Split(property.Amenities, ","),
        "HeroImage":    details.ImageUrl1,
        "Images": []string{
            details.ImageUrl1,
            details.ImageUrl2,
            details.ImageUrl3,
            details.ImageUrl4,
            details.ImageUrl5,
        },
        "Description": details.Description,
    }

    // Pass data to the template
    c.Data["Property"] = enhancedProperty
    c.TplName = "property_details.tpl"
    c.Render()
}