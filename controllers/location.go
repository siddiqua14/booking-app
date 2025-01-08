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
    fmt.Printf("API Response: %+v\n", apiResponse)

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
                id, err := o.Insert(location)
                if err != nil {
                    fmt.Printf("Error inserting location: %v\n", err)
                } else {
                    fmt.Printf("Inserted location with ID: %d\n", id)
                }
            }
        }
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