package routers

import (
    "github.com/beego/beego/v2/server/web"
    "booking-app/controllers"
)

func init() {
    web.Router("/fetch_locations", &controllers.LocationController{}, "get:FetchAndStoreLocations")
    web.Router("/fetch_stays_data", &controllers.LocationController{}, "get:FetchFilteredStaysData")
    web.Router("/fetch-hotel-details", &controllers.LocationController{}, "get:FetchHotelDetails")
    // in your router.go or wherever you define your routes
    web.Router("/fetch-hotel-images-and-description", &controllers.LocationController{}, "get:FetchHotelImagesAndDescription")
}