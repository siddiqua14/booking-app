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

type FilteredLocation struct {
    DestId   string `json:"dest_id"`
    DestType string `json:"dest_type"`
    Value    string `json:"value"`
}
type HotelDetails struct {
    IDHotel     int     `json:"id_hotel"`
    HotelID     string  `json:"hotel_id"`
    HotelName   string  `json:"hotel_name"`
    DestID      string  `json:"dest_id"`
    Location    string  `json:"location"`
    Rating      float64 `json:"rating"`
    ReviewCount int     `json:"review_count"`
    Price       string  `json:"price"`
    NumBeds     int     `json:"num_beds"`
    NumBedR     int     `json:"num_bedrooms"`
    NumBaths    int     `json:"num_bathrooms"`
}
type APIResponse struct {
    Data struct {
        Results []struct {
            ID               string `json:"id"`
            BasicPropertyData struct {
                ID      int `json:"id"`
                Reviews struct {
                    TotalScore   float64 `json:"totalScore"`
                    ReviewsCount int     `json:"reviewsCount"`
                } `json:"reviews"`
                
            } `json:"basicPropertyData"`
            Location struct {
                DisplayLocation string `json:"displayLocation"`
            } `json:"location"`
            DisplayName struct {
                Text string `json:"text"`
            } `json:"displayName"`
            MatchingUnitConfigurations struct {
                CommonConfiguration struct {
                    NbAllBeds   int `json:"nbAllBeds"`
                    NbAllBedR   int `json:"nbBedrooms"`
                    NbBathrooms int `json:"nbBathrooms"`
                } `json:"commonConfiguration"`
            } `json:"matchingUnitConfigurations"`
            PriceDisplayInfoIrene struct {
                DisplayPrice struct {
                    AmountPerStay struct {
                        Amount string `json:"amount"`
                    } `json:"amountPerStay"`
                } `json:"displayPrice"`
            } `json:"priceDisplayInfoIrene"`
        } `json:"results"`
    } `json:"data"`
    Meta struct {
        CurrentPage   int `json:"currentPage"`
        Limit        int `json:"limit"`
        TotalRecords int `json:"totalRecords"`
        TotalPage    int `json:"totalPage"`
    } `json:"meta"`
    Status  bool   `json:"status"`
    Message string `json:"message"`
}

// In-memory storage for dest_id and dest_type
var storedDestId string
var storedDestType string
var storedHotelIDs []int


func (c *LocationController) FetchAndStoreLocations() {
    // Create HTTP client
    client := &http.Client{}

    // Create request
    url := "https://booking-com18.p.rapidapi.com/web/stays/auto-complete?query=New%20York"
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Error creating request: %v", err)}
        c.ServeJSON()
        return
    }

    // Add headers
    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", "52d384abecmshb0390e7c79d8689p1a8cd1jsn9c2ed180601b")

    // Make request
    resp, err := client.Do(req)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Error making request: %v", err)}
        c.ServeJSON()
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(resp.Body)
        bodyString := string(bodyBytes)
        c.Data["json"] = map[string]interface{}{
            "error": fmt.Sprintf("API request failed with status code: %d", resp.StatusCode),
            "response": bodyString,
        }
        c.ServeJSON()
        return
    }

    var apiResponse map[string]interface{}
    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(&apiResponse); err != nil {
        c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Error decoding response: %v", err)}
        c.ServeJSON()
        return
    }

    // Log API response
    fmt.Printf("Auto-complete API Response: %+v\n", apiResponse)

    // Extract relevant data from the API response
    var filteredLocations []FilteredLocation
    if data, ok := apiResponse["data"].([]interface{}); ok {
        o := orm.NewOrm()
        for _, item := range data {
            if itemMap, ok := item.(map[string]interface{}); ok {
                filteredLocation := FilteredLocation{
                    DestId:   itemMap["dest_id"].(string),
                    DestType: itemMap["dest_type"].(string),
                    Value:    itemMap["label"].(string),
                }
                filteredLocations = append(filteredLocations, filteredLocation)
                // Store in the database
                location := &models.Location{
                    DestId:   itemMap["dest_id"].(string),
                    DestType: itemMap["dest_type"].(string),
                    Value:    itemMap["label"].(string),
                }
                 // Check if the dest_id already exists in the database before inserting
                existingLocation := models.Location{DestId: location.DestId}
                err := o.Read(&existingLocation, "DestId")
                if err != nil {
                     // Only insert if not already present
                    id, err := o.Insert(location)
                    if err != nil {
                        fmt.Printf("Error inserting location: %v\n", err)
                    } else {
                        fmt.Printf("Inserted location with ID: %d\n", id)
                    }
                } else {
                    fmt.Println("Location already exists, skipping insert.")
                }
            }
        }
    }
    // Store the first location's dest_id and dest_type for fetching hotel data later
	if len(filteredLocations) > 0 {
		storedDestId = filteredLocations[0].DestId
		storedDestType = filteredLocations[0].DestType
	}
    // Return the filtered data in the required format
    c.Data["json"] = map[string]interface{}{
        "success": true,
        "query": "New York",
        "count": len(filteredLocations),
        "data": map[string]interface{}{
            "data": filteredLocations,
        },
    }
    c.ServeJSON()
}

func (c *LocationController) FetchFilteredStaysData() {
    // Log the stored values
    log.Printf("Fetching data with destID: %s, destType: %s", storedDestId, storedDestType)
    
    hotels, err := getHotelData(storedDestId, storedDestType)
    if err != nil {
        log.Printf("Error in getHotelData: %v", err)
        c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Error fetching stays data: %v", err)}
        c.ServeJSON()
        return
    }

    if len(hotels) == 0 {
        log.Print("No hotels found in the response")
        c.Data["json"] = map[string]interface{}{
            "message": "No stays found for the given location.",
        }
        c.ServeJSON()
        return
    }
    // Store all hotel IDs for fetching hotel details later
    for _, hotel := range hotels {
        storedHotelIDs = append(storedHotelIDs, hotel.IDHotel)
    }

    c.Data["json"] = map[string]interface{}{
        "message": "Stays data fetched successfully",
        "results": hotels,
    }
    c.ServeJSON()
}

func getHotelData(destID, destType string) ([]HotelDetails, error) {
    if destID == "" || destType == "" {
        return nil, fmt.Errorf("destID and destType cannot be empty")
    }

    url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/web/stays/search?destId=%s&destType=%s&checkIn=2025-01-12&checkOut=2025-01-31", destID, destType)
    log.Printf("Making request to URL: %s", url)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", "52d384abecmshb0390e7c79d8689p1a8cd1jsn9c2ed180601b")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    log.Printf("API Response Status: %s", resp.Status)

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %v", err)
    }

    // Log the raw response for debugging
    log.Printf("Raw API Response: %s", string(body))

    var apiResponse APIResponse
    if err := json.Unmarshal(body, &apiResponse); err != nil {
        // Try to unmarshal into a map to see the actual structure
        var rawResponse map[string]interface{}
        if jsonErr := json.Unmarshal(body, &rawResponse); jsonErr == nil {
            log.Printf("API Response Structure: %+v", rawResponse)
        }
        return nil, fmt.Errorf("error unmarshaling response: %v", err)
    }

    // Check if the response was successful
    if !apiResponse.Status {
        return nil, fmt.Errorf("API returned error: %s", apiResponse.Message)
    }

    // Check if we have any results in the data
    if len(apiResponse.Data.Results) == 0 {
        log.Printf("API returned success but no results. Total records: %d", apiResponse.Meta.TotalRecords)
        return nil, nil
    }

    var hotels []HotelDetails
    o := orm.NewOrm()

    for i, result := range apiResponse.Data.Results {
        log.Printf("Processing result %d: %+v", i, result)
        
        hotel := HotelDetails{
            IDHotel:     result.BasicPropertyData.ID,
            HotelID:     result.ID,
            HotelName:   result.DisplayName.Text,
            DestID:      destID,
            Location:    result.Location.DisplayLocation,
            Rating:      result.BasicPropertyData.Reviews.TotalScore,
            ReviewCount: result.BasicPropertyData.Reviews.ReviewsCount,
            Price:       result.PriceDisplayInfoIrene.DisplayPrice.AmountPerStay.Amount,
            NumBeds:     result.MatchingUnitConfigurations.CommonConfiguration.NbAllBeds,
            NumBedR:     result.MatchingUnitConfigurations.CommonConfiguration.NbAllBeds, // Assuming NbAllBedR is a typo
            NumBaths:    result.MatchingUnitConfigurations.CommonConfiguration.NbBathrooms,
        }

        // Append to hotels slice
        hotels = append(hotels, hotel)

        // Check if the hotel already exists in the database
        existingRentalProperty := models.RentalProperty{IDHotel: hotel.IDHotel}
        err := o.Read(&existingRentalProperty, "IDHotel")
        if err == nil {
            // Hotel already exists, skip insertion
            log.Printf("Hotel with IDHotel %d already exists, skipping insertion.", hotel.IDHotel)
            continue
        }

        // Map HotelDetails to RentalProperty
        rentalProperty := models.RentalProperty{
            IDHotel:     hotel.IDHotel,
            HotelID:     hotel.HotelID,
            HotelName:   hotel.HotelName,
            DestID:      hotel.DestID,
            Location:    hotel.Location,
            Rating:      hotel.Rating,
            ReviewCount: hotel.ReviewCount,
            Price:       hotel.Price,
            NumBeds:     hotel.NumBeds,
            NumBedR:     hotel.NumBedR,
            NumBaths:    hotel.NumBaths,
        }

        // Insert rentalProperty into the database
        _, err = o.Insert(&rentalProperty)
        if err != nil {
            log.Printf("Error inserting rental property: %v", err)
        } else {
            log.Printf("Inserted rental property with ID: %d", rentalProperty.ID)
        }
    }

    log.Printf("Processed %d rental properties successfully", len(hotels))
    return hotels, nil
}

func (c *LocationController) FetchHotelDetails() {
    o := orm.NewOrm()

    var hotelDetailsList []map[string]interface{}

    for _, hotelID := range storedHotelIDs {
        url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/detail?hotelId=%d&checkinDate=2025-01-09&checkoutDate=2025-01-15&units=metric", hotelID)
        log.Printf("Making request to URL: %s", url)

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            log.Printf("Error creating request for hotel ID %d: %v", hotelID, err)
            continue
        }

        req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
        req.Header.Add("x-rapidapi-key", "52d384abecmshb0390e7c79d8689p1a8cd1jsn9c2ed180601b")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            log.Printf("Error making request for hotel ID %d: %v", hotelID, err)
            continue
        }
        defer resp.Body.Close()

        log.Printf("API Response Status for hotel ID %d: %s", hotelID, resp.Status)

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Printf("Error reading response body for hotel ID %d: %v", hotelID, err)
            continue
        }

        log.Printf("Raw API Response for hotel ID %d: %s", hotelID, string(body))

        var apiResponse map[string]interface{}
        if err := json.Unmarshal(body, &apiResponse); err != nil {
            log.Printf("Error unmarshalling response for hotel ID %d: %v", hotelID, err)
            continue
        }

        data, ok := apiResponse["data"].(map[string]interface{})
        if !ok {
            log.Printf("Invalid data format in response for hotel ID %d", hotelID)
            continue
        }

        availableRooms, ok := data["available_rooms"].(float64)
        if !ok {
            log.Printf("Invalid available_rooms format in response for hotel ID %d", hotelID)
            continue
        }

        roomRecommendations, ok := data["room_recommendation"].([]interface{})
        if !ok {
            log.Printf("Invalid room_recommendation format in response for hotel ID %d", hotelID)
            continue
        }

        var totalAdults, totalChildren float64
        for _, recommendation := range roomRecommendations {
            if recommendationMap, ok := recommendation.(map[string]interface{}); ok {
                if adults, ok := recommendationMap["adults"].(float64); ok {
                    totalAdults += adults
                }
                if children, ok := recommendationMap["children"].(float64); ok {
                    totalChildren += children
                }
            }
        }

        accommodationTypeName, ok := data["accommodation_type_name"].(string)
        if !ok {
            log.Printf("Invalid accommodation_type_name format in response for hotel ID %d", hotelID)
            continue
        }

        facilitiesBlock, ok := data["facilities_block"].(map[string]interface{})
        if !ok {
            log.Printf("Invalid facilities_block format in response for hotel ID %d", hotelID)
            continue
        }

        facilities, ok := facilitiesBlock["facilities"].([]interface{})
        if !ok {
            log.Printf("Invalid facilities format in response for hotel ID %d", hotelID)
            continue
        }

        var amenities []string
        for i, facility := range facilities {
            if i >= 3 {
                break
            }
            facilityMap, ok := facility.(map[string]interface{})
            if !ok {
                log.Printf("Invalid facility format in response for hotel ID %d", hotelID)
                continue
            }
            name, ok := facilityMap["name"].(string)
            if !ok {
                log.Printf("Invalid facility name format in response for hotel ID %d", hotelID)
                continue
            }
            amenities = append(amenities, name)
        }

        // Update the existing rental property with additional details
        rentalProperty := models.RentalProperty{IDHotel: hotelID}
        if err := o.Read(&rentalProperty, "IDHotel"); err != nil {
            log.Printf("Rental property not found for hotel ID %d", hotelID)
            continue
        }

        rentalProperty.Bedroom = int(availableRooms)
        rentalProperty.Guests = int(totalAdults + totalChildren)
        rentalProperty.PropertyType = accommodationTypeName
        rentalProperty.Amenities = strings.Join(amenities, ", ")

        if _, err := o.Update(&rentalProperty); err != nil {
            log.Printf("Error updating rental property for hotel ID %d: %v", hotelID, err)
            continue
        }

        log.Printf("Rental property updated successfully for hotel ID %d", hotelID)

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
        }
        hotelDetailsList = append(hotelDetailsList, hotelDetails)
    }

    // Send the hotel details as JSON response
    c.Data["json"] = map[string]interface{}{
        "message": "Rental properties updated successfully",
        "details": hotelDetailsList,
    }
    c.ServeJSON()
}