package controllers

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/beego/beego/v2/client/orm"
    "github.com/beego/beego/v2/server/web"
    "booking-app/models"
)

type LocationController struct {
    web.Controller
}

type AutoCompleteResponse struct {
    Results []struct {
        DestId   string `json:"dest_id"`
        DestType string `json:"dest_type"`
        Value    string `json:"value"`
    } `json:"results"`
}

func (c *LocationController) FetchAndStoreLocations() {
    // Create HTTP client
    client := &http.Client{}
    
    // Create request
    url := "https://booking-com18.p.rapidapi.com/web/stays/auto-complete?query=New%20York"
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": err.Error()}
        c.ServeJSON()
        return
    }

    // Add headers
    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", "52d384abecmshb0390e7c79d8689p1a8cd1jsn9c2ed180601b")

    // Make request
    resp, err := client.Do(req)
    if err != nil {
        c.Data["json"] = map[string]interface{}{"error": err.Error()}
        c.ServeJSON()
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("API request failed with status code: %d", resp.StatusCode)}
        c.ServeJSON()
        return
    }

    var result map[string]interface{}
    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(&result); err != nil {
        c.Data["json"] = map[string]interface{}{"error": err.Error()}
        c.ServeJSON()
        return
    }

    // Log API response
    fmt.Printf("API Response: %+v\n", result)

    // Store relevant data in the database
    if results, ok := result["results"].([]interface{}); ok {
        o := orm.NewOrm()
        for _, loc := range results {
            if locMap, ok := loc.(map[string]interface{}); ok {
                location := &models.Location{
                    DestId:   locMap["dest_id"].(string),
                    DestType: locMap["dest_type"].(string),
                    Value:    locMap["label"].(string),
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

    // Return the full API response
    c.Data["json"] = map[string]interface{}{
        "message": "Locations fetched and stored successfully",
        "data":    result,
    }
    c.ServeJSON()
}
