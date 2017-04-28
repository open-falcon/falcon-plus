package controllers

import (
	"github.com/astaxie/beego"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
)

type Search struct {
	Name        string
	Placeholder string
}

type BaseController struct {
	beego.Controller
}

func init() {
	// The hookfuncs will run in beego.Run()
	// beego.AddAPPStartHook(start)
}

func (c *BaseController) SendMsg(code int, msg interface{}) {
	c.Ctx.ResponseWriter.WriteHeader(code)
	c.Data["json"] = msg
	c.ServeJSON()
}

func totalObj(n int64) models.Total {
	return models.Total{Total: n}
}

func idObj(n int64) models.Id {
	return models.Id{Id: n}
}
