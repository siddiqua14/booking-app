package routers

import (
	"booking-app/controllers"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/location", &controllers.LocationController{}, "post:AddLocation")
}
