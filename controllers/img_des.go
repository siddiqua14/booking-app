package controllers

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "log"
    "strings"

    "github.com/beego/beego/v2/client/orm"
    "github.com/beego/beego/v2/server/web"
    "booking-app/models"
)

type LocationController struct {
    web.Controller
}

func (c *LocationController) FetchHotelImagesAndDescription() {
    o := orm.NewOrm()
    var rentalProperties []models.RentalProperty

    // Fetch all rental properties from the database
    _, err := o.QueryTable(new(models.RentalProperty)).All(&rentalProperties)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Error fetching rental properties: %v", err)}
        c.ServeJSON()
        return
    }

    var hotelDetailsList []map[string]interface{}

    for _, rentalProperty := range rentalProperties {
        hotelID := rentalProperty.HotelID
        url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/web/stays/details?id=%s&checkIn=2025-01-09&checkOut=2025-01-15", hotelID)
        log.Printf("Making request to URL: %s", url)

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            log.Printf("Error creating request for hotel ID %s: %v", hotelID, err)
            continue
        }

        req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
        req.Header.Add("x-rapidapi-key", "3c935dc998msh43f2b397ec5205dp174193jsnb246f518916d")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            log.Printf("Error making request for hotel ID %s: %v", hotelID, err)
            continue
        }
        defer resp.Body.Close()

        log.Printf("API Response Status for hotel ID %s: %s", hotelID, resp.Status)

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Printf("Error reading response body for hotel ID %s: %v", hotelID, err)
            continue
        }

        log.Printf("Raw API Response for hotel ID %s: %s", hotelID, string(body))

        var apiResponse map[string]interface{}
        if err := json.Unmarshal(body, &apiResponse); err != nil {
            log.Printf("Error unmarshalling response for hotel ID %s: %v", hotelID, err)
            continue
        }

        data, ok := apiResponse["data"].(map[string]interface{})
        if !ok {
            log.Printf("Invalid data format in response for hotel ID %s", hotelID)
            continue
        }

        propertyImg := ""
        description := ""

        // Parse hotelPhotos to get all images of thumb_url
        if hotelPhotos, ok := data["hotelPhotos"].([]interface{}); ok {
            var images []string
            for _, photo := range hotelPhotos {
                if photoMap, ok := photo.(map[string]interface{}); ok {
                    if thumbURL, ok := photoMap["thumb_url"].(string); ok {
                        images = append(images, thumbURL)
                    }
                }
            }
            propertyImg = strings.Join(images, ", ")
        }

        // Parse hotelTranslation to get the description
        if hotelTranslation, ok := data["hotelTranslation"].([]interface{}); ok {
            if len(hotelTranslation) > 0 {
                if translationMap, ok := hotelTranslation[0].(map[string]interface{}); ok {
                    if desc, ok := translationMap["description"].(string); ok {
                        description = desc
                    }
                }
            }
        }

        // Update the existing rental property with additional details
        rentalProperty.PropertyImg = propertyImg
        rentalProperty.Description = description

        if _, err := o.Update(&rentalProperty); err != nil {
            log.Printf("Error updating rental property for hotel ID %s: %v", hotelID, err)
            continue
        }

        log.Printf("Rental property updated successfully for hotel ID %s", hotelID)

        // Collect details for JSON output
        hotelDetails := map[string]interface{}{
            "IDHotel":       rentalProperty.IDHotel,
            "HotelID":       rentalProperty.HotelID,
            "HotelName":     rentalProperty.HotelName,
            "DestID":        rentalProperty.DestID,
            "Location":      rentalProperty.Location,
            "Rating":        rentalProperty.Rating,
            "ReviewCount":   rentalProperty.ReviewCount,
            "Price":         rentalProperty.Price,
            "NumBeds":       rentalProperty.NumBeds,
            "NumBedR":       rentalProperty.NumBedR,
            "NumBaths":      rentalProperty.NumBaths,
            "Bedroom":       rentalProperty.Bedroom,
            "Guests":        rentalProperty.Guests,
            "PropertyType":  rentalProperty.PropertyType,
            "Amenities":     rentalProperty.Amenities,
            "PropertyImg":   rentalProperty.PropertyImg,
            "Description":   rentalProperty.Description,
        }
        hotelDetailsList = append(hotelDetailsList, hotelDetails)
    }

    c.Data["json"] = map[string]interface{}{
        "message": "Rental properties updated successfully",
        "details": hotelDetailsList,
    }
    c.ServeJSON()
}