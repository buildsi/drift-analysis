package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	web.GlobalControllerRouter["drift-server/controllers:InflectionPointController"] = append(web.GlobalControllerRouter["drift-server/controllers:InflectionPointController"],
		web.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	//    beego.GlobalControllerRouter["drift-server/controllers:ObjectController"] = append(beego.GlobalControllerRouter["drift-server/controllers:ObjectController"],
	//        beego.ControllerComments{
	//            Method: "GetAll",
	//            Router: "/",
	//            AllowHTTPMethods: []string{"get"},
	//            MethodParams: param.Make(),
	//            Filters: nil,
	//            Params: nil})

	//    beego.GlobalControllerRouter["drift-server/controllers:ObjectController"] = append(beego.GlobalControllerRouter["drift-server/controllers:ObjectController"],
	//        beego.ControllerComments{
	//            Method: "Get",
	//            Router: "/:objectId",
	//            AllowHTTPMethods: []string{"get"},
	//            MethodParams: param.Make(),
	//            Filters: nil,
	//            Params: nil})

	//    beego.GlobalControllerRouter["drift-server/controllers:ObjectController"] = append(beego.GlobalControllerRouter["drift-server/controllers:ObjectController"],
	//        beego.ControllerComments{
	//            Method: "Put",
	//            Router: "/:objectId",
	//            AllowHTTPMethods: []string{"put"},
	//            MethodParams: param.Make(),
	//            Filters: nil,
	//            Params: nil})
	//
	//    beego.GlobalControllerRouter["drift-server/controllers:ObjectController"] = append(beego.GlobalControllerRouter["drift-server/controllers:ObjectController"],
	//        beego.ControllerComments{
	//            Method: "Delete",
	//            Router: "/:objectId",
	//            AllowHTTPMethods: []string{"delete"},
	//            MethodParams: param.Make(),
	//            Filters: nil,
	//            Params: nil})
}
