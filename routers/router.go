package routers

import (
    "github.com/beego/beego/v2/server/web"
    "booking-app/controllers"
)

func init() {
    web.Router("/fetch_locations", &controllers.LocationController{}, "get:FetchAndStoreLocations")
	//web.Router("/locations", &controllers.LocationController{}, "get:GetLocations")
}