package routers

import (
    "github.com/beego/beego/v2/server/web"
    "booking-app/controllers"
)

func init() {
    web.Router("/", &controllers.MainController{})

    web.Router("/fetch_locations", &controllers.LocationController{}, "get:FetchAndStoreLocations")
    web.Router("/fetch_stays_data", &controllers.LocationController{}, "get:FetchFilteredStaysData")
    web.Router("/fetch-hotel-details", &controllers.LocationController{}, "get:FetchHotelDetails")
    web.Router("/fetch-hotel-images-and-description", &controllers.ImageDescriptionController{}, "get:FetchHotelImagesAndDescriptions")
    // Page route
    // Property listing endpoint
    web.Router("/v1/property/list", &controllers.PropertyController{}, "get:ListProperties")
    
    // API route for AJAX calls
    web.Router("/v1/property/details", &controllers.PropertyController{}, "get:GetPropertyDetails")
}