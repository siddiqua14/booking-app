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
    "github.com/beego/beego/v2/server/web/context" 
    "booking-app/models"
)

type LocationController struct {
    web.Controller
    apiBaseUrl string
    apiKey     string
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


// Init function to initialize apiBaseUrl and apiKey
func (c *LocationController) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
    c.Controller.Init(ctx, controllerName, actionName, app)

    apiBaseUrl, err := web.AppConfig.String("API_BASE_URL")
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Error retrieving API Base URL: %v", err)}
        c.ServeJSON()
        return
    }

    apiKey, err := web.AppConfig.String("API_KEY")
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Error retrieving API Key: %v", err)}
        c.ServeJSON()
        return
    }

    c.apiBaseUrl = apiBaseUrl
    c.apiKey = apiKey
}

// FetchAndStoreLocations retrieves locations and stores them
func (c *LocationController) FetchAndStoreLocations() {
    query := "New%20York"
    url := fmt.Sprintf("%s/web/stays/auto-complete?query=%s", c.apiBaseUrl, query)

    // Create HTTP client
    client := &http.Client{}

    // Create request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Error creating request: %v", err)}
        c.ServeJSON()
        return
    }

    // Add headers with dynamic API key
    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", c.apiKey)

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

    // Return the filtered data in the required format
    c.Data["json"] = map[string]interface{}{
        "success": true,
        "query":   query,
        "count":   len(filteredLocations),
        "data": map[string]interface{}{
            "data": filteredLocations,
        },
    }
    c.ServeJSON()
}

func (c *LocationController) FetchFilteredStaysData() {
    // Fetch all locations from the database
    o := orm.NewOrm()
    var locations []models.Location
    _, err := o.QueryTable("location").All(&locations)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": "Failed to retrieve locations"}
        c.ServeJSON()
        return
    }

    var allHotels []HotelDetails
    for _, loc := range locations {
        // Fetch data for each location
        hotels, err := getHotelData(loc.DestId, loc.DestType, c.apiBaseUrl, c.apiKey)
        if err != nil {
            log.Printf("Error fetching data for destID: %s, destType: %s - %v", loc.DestId, loc.DestType, err)
            continue
        }

        allHotels = append(allHotels, hotels...)
    }

    if len(allHotels) == 0 {
        log.Print("No hotels found for any location")
        c.Data["json"] = map[string]interface{}{
            "message": "No stays found for any location.",
        }
        c.ServeJSON()
        return
    }

    c.Data["json"] = map[string]interface{}{
        "message": "Stays data fetched successfully",
        "results": allHotels,
    }
    c.ServeJSON()
}

func getHotelData(destID, destType, apiBaseUrl, apiKey string) ([]HotelDetails, error) {
    if destID == "" || destType == "" {
        return nil, fmt.Errorf("destID and destType cannot be empty")
    }

    url := fmt.Sprintf("%s/web/stays/search?destId=%s&destType=%s&checkIn=2025-01-12&checkOut=2025-01-31", apiBaseUrl, destID, destType)
    log.Printf("Making request to URL: %s", url)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", apiKey)

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

    log.Printf("Raw API Response: %s", string(body))

    var apiResponse APIResponse
    if err := json.Unmarshal(body, &apiResponse); err != nil {
        var rawResponse map[string]interface{}
        if jsonErr := json.Unmarshal(body, &rawResponse); jsonErr == nil {
            log.Printf("API Response Structure: %+v", rawResponse)
        }
        return nil, fmt.Errorf("error unmarshaling response: %v", err)
    }

    if !apiResponse.Status {
        return nil, fmt.Errorf("API returned error: %s", apiResponse.Message)
    }

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
            NumBedR:     result.MatchingUnitConfigurations.CommonConfiguration.NbAllBedR,
            NumBaths:    result.MatchingUnitConfigurations.CommonConfiguration.NbBathrooms,
        }

        hotels = append(hotels, hotel)

        existingRentalProperty := models.RentalProperty{IDHotel: hotel.IDHotel}
        err := o.Read(&existingRentalProperty, "IDHotel")
        if err == nil {
            log.Printf("Hotel with IDHotel %d already exists, skipping insertion.", hotel.IDHotel)
            continue
        }

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

    // Fetch top 25 records from the rental_property table
    var rentalProperties []models.RentalProperty
    _, err := o.QueryTable("rental_property").Limit(25).All(&rentalProperties)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": "Failed to retrieve rental properties"}
        c.ServeJSON()
        return
    }

    var hotelDetailsList []map[string]interface{}

    for _, rentalProperty := range rentalProperties {
        hotelID := rentalProperty.IDHotel
        url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/detail?hotelId=%d&checkinDate=2025-01-10&checkoutDate=2025-01-15&units=metric", hotelID)
        log.Printf("Making request to URL: %s", url)

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            log.Printf("Error creating request for hotel ID %d: %v", hotelID, err)
            continue
        }

        req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
        req.Header.Add("x-rapidapi-key", c.apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            log.Printf("Error making request for hotel ID %d: %v", hotelID, err)
            continue
        }

        // Ensure the response body is closed after reading
        body, err := io.ReadAll(resp.Body)
        resp.Body.Close()  // Close response body immediately after reading
        if err != nil {
            log.Printf("Error reading response body for hotel ID %d: %v", hotelID, err)
            continue
        }

        log.Printf("API Response Status for hotel ID %d: %s", hotelID, resp.Status)
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