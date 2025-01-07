package controllers

import (
	"booking-app/models"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

type LocationController struct {
	web.Controller
}

func (c *LocationController) AddLocation() {
	locationName := c.GetString("location")
	if locationName == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "Location name is required"
		c.ServeJSON()
		return
	}

	// Fetch data from Booking API
	url := "https://booking-com18.p.rapidapi.com/stays/auto-complete?query=" + locationName
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "52d384abecmshb0390e7c79d8689p1a8cd1jsn9c2ed180601b")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = "Error fetching data from Booking API"
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var hotels []map[string]interface{}
	json.Unmarshal(body, &hotels)

	// Store data in the database
	hotelsJSON, _ := json.Marshal(hotels)
	err = models.AddLocation(locationName, string(hotelsJSON))
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = "Error saving data to database"
		c.ServeJSON()
		return
	}

	c.Data["json"] = "Location added successfully"
	c.ServeJSON()
}
