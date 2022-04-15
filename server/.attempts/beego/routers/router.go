package routers

import (
	"drift-server/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	ns := web.NewNamespace("/v1",
		web.NSNamespace("/point",
			web.NSInclude(
				&controllers.InflectionPointController{},
			),
		),
		//		web.NSNamespace("/user",
		//			web.NSInclude(
		//				&controllers.UserController{},
		//			),
		//		),
	)
	web.AddNamespace(ns)
}
