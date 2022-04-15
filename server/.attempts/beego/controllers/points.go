package controllers

import (
	"drift-server/models"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/server/web"
)

// Operations for inflection points
type InflectionPointController struct {
	web.Controller
}

// @Title Create Inflection Point
// @Description create a new inflection point
// @Param	body		body 	models.InflectionPoint	true		"The inflection point content"
// @Success 200 {string} models.InflectionPoint.Id
// @Failure 403 body is empty
// @router / [post]
func (controller *InflectionPointController) Post() {

	// Create a new inflection point and unmarshall body into it
	logs.Info(fmt.Sprintf("%s", controller.Ctx.Input.RequestBody))

	// The same as an inflection point, but w/ json and only core fields
	var point models.InflectionPointRequest
	err := json.Unmarshal(controller.Ctx.Input.RequestBody, &point)
	if err != nil {
		logs.Warn(err)
	}

	// Get or Create the point
	inflectionPoint := models.GetOrCreateInflectionPoint(point)
	controller.Data["json"] = map[string]int{"Id": inflectionPoint.Id}
	controller.ServeJSON()
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [get]
//func (o *InflectionPointController) Get() {
//	objectId := o.Ctx.Input.Param(":objectId")
//	if objectId != "" {
//		ob, err := models.GetOne(objectId)
//		if err != nil {
//			o.Data["json"] = err.Error()
//		} else {
//			o.Data["json"] = ob
//		}
//	}
//	o.ServeJSON()
//}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
//func (o *InflectionPointController) GetAll() {
//	obs := models.GetAll()
//	o.Data["json"] = obs
//	o.ServeJSON()
//}

// @Title Update
// @Description update the object
// @Param	objectId		path 	string	true		"The objectid you want to update"
// @Param	body		body 	models.Object	true		"The body"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [put]
//func (o *InflectionPointController) Put() {
//	objectId := o.Ctx.Input.Param(":objectId")
//	var ob models.Object
//	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
//
//	err := models.Update(objectId, ob.Score)
//	if err != nil {
//		o.Data["json"] = err.Error()
//	} else {
//		o.Data["json"] = "update success!"
//	}
//	o.ServeJSON()
//}

// @Title Delete
// @Description delete the object
// @Param	objectId		path 	string	true		"The objectId you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 objectId is empty
// @router /:objectId [delete]
//func (o *InflectionPointController) Delete() {
//	objectId := o.Ctx.Input.Param(":objectId")
//	models.Delete(objectId)
//	o.Data["json"] = "delete success!"
//	o.ServeJSON()
//}
