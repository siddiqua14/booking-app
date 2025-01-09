package controllers

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"

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
	HotelID    string  `json:"hotel_id"`
	HotelName  string  `json:"hotel_name"`
	DestID     string  `json:"dest_id"`
    Location   string  `json:"location"`
	Rating     float64 `json:"rating"`
	ReviewCount int    `json:"review_count"`
}
// In-memory storage for dest_id and dest_type
var storedDestId string
var storedDestType string

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

// FetchFilteredStaysData fetches hotel data based on dest_id and dest_type
func (c *LocationController) FetchFilteredStaysData() {
	hotels, err := getHotelData(storedDestId, storedDestType)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Error fetching stays data: %v", err)}
		c.ServeJSON()
		return
	}

	if len(hotels) == 0 {
		c.Data["json"] = map[string]interface{}{
			"message": "No stays found for the given location.",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"message": "Stays data fetched successfully",
		"results": hotels,
	}
	c.ServeJSON()
}

// getHotelData fetches hotel data from the Booking.com API
func getHotelData(destID, destType string) ([]HotelDetails, error) {
    url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/web/stays/search?destId=%s&destType=%s&checkIn=2025-01-12&checkOut=2025-01-31", destID, destType)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", "52d384abecmshb0390e7c79d8689p1a8cd1jsn9c2ed180601b")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var apiResponse map[string]interface{}
    if err := json.Unmarshal(body, &apiResponse); err != nil {
        return nil, err
    }

    // Log the full API response
    fmt.Printf("Hotel Search API Response: %+v\n", apiResponse)

    // Extract the relevant data from the API response
    var hotels []HotelDetails
    if data, ok := apiResponse["data"].(map[string]interface{}); ok {
        results, ok := data["results"].([]interface{})
        if !ok {
            return nil, fmt.Errorf("unexpected type for results field, got: %T", data["results"])
        }

        // Loop through each hotel and extract relevant data
        for _, item := range results {
            if itemMap, ok := item.(map[string]interface{}); ok {

                // Extract hotel ID
                hotelID, ok := itemMap["id"]
                if !ok {
                    continue
                }

                // Extract hotel name from displayName
                displayName, ok := itemMap["displayName"].(map[string]interface{})
                if !ok {
                    continue
                }
                hotelName, ok := displayName["text"]
                if !ok {
                    continue
                }
                
                // Extract location from basicPropertyData -> location -> displayLocation
                // Extract location
                basicPropertyData, ok := itemMap["basicPropertyData"].(map[string]interface{})
                if !ok {
                    continue
                }
                locationData, ok := basicPropertyData["location"].(map[string]interface{})
                if !ok {
                    continue
                }
                location, ok := locationData["displayLocation"].(string)
                if !ok {
                    location = "Unknown"
                }

                // Extract rating
                reviews, ok := basicPropertyData["reviews"].(map[string]interface{})
                var rating float64
                if ok {
                    rating, _ = reviews["totalScore"].(float64)
                }

                // Extract review count
                var reviewCount int
                if ok {
                    if rc, ok := reviews["reviewCount"].(float64); ok {
                        reviewCount = int(rc)
                    }
                }
                // Convert hotelID and hotelName to string
                hotelIDStr := fmt.Sprintf("%v", hotelID)
                hotelNameStr := fmt.Sprintf("%v", hotelName)

                // Append the hotel details to the list
                hotel := HotelDetails{
                    HotelID:     hotelIDStr,
                    HotelName:   hotelNameStr,
                    DestID:      destID,
                    Location:    location,
                    Rating:      rating,
                    ReviewCount: reviewCount,
                }
                hotels = append(hotels, hotel)
            }
        }
    } else {
        return nil, fmt.Errorf("unexpected type for data field, got: %T", apiResponse["data"])
    }

    return hotels, nil
}
