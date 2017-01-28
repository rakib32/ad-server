package routers

import (
	"ad-server/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/v1/ads", &controllers.AdController{}, "get:ServeAds")
	beego.Router("/api/v1/ad/impression/:adid/:deliveryid", &controllers.AdController{}, "post:SaveImpression")
	beego.Router("/api/v1/ad/click/:adid/:deliveryid", &controllers.AdController{}, "post:SaveClick")
}
