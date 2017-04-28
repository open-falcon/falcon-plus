// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
)

// Operations about porfile/config/info
type AdminController struct {
	BaseController
}

// @Title Get config
// @Description get module config
// @Param	module	path	string	true	"module name"
// @Success 200 {object} map[string]string online host list
// @Failure 403 string error
// @router /online/:module [get]
func (c *AdminController) GetOnline() {
	var err error

	module := c.GetString(":module")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	conf, err := op.OnlineGet(module)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, conf)
	}
}

// @Title Get config
// @Description get module config
// @Param	module	path	string	true	"module name"
// @Success 200 {object} [3]map[string]string {defualt{}, conf{}, configfile{}}
// @Failure 403 string error
// @router /config/:module [get]
func (c *AdminController) GetConfig() {
	var err error

	module := c.GetString(":module")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	conf, err := op.ConfigGet(module)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, conf)
	}
}

// @Title update config
// @Description get tag role user
// @Param	module	path	string	true	"module"
// @Param	body	body	map[string]string	true	""
// @Success 200 {string} success
// @Failure 403 string error
// @router /config/:module [put]
func (c *AdminController) UpdateConfig() {
	var conf map[string]string

	module := c.GetString(":module")

	beego.Debug(string(c.Ctx.Input.RequestBody))
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &conf)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if err := op.ConfigSet(module, conf); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, "success")
	}
}

// @Title Get config
// @Description get tag role user
// @Param	action	path	string	true	"action"
// @Success 200 {string} result
// @Failure 403 string error
// @router /debug/:action [get]
func (c *AdminController) GetDebugAction() {
	var err error
	var obj interface{}
	action := c.GetString(":action")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	switch action {
	case "populate":
		obj, err = op.Populate()
	case "reset_db":
		obj, err = op.ResetDb()
	default:
		err = fmt.Errorf("%s %s", models.ErrUnsupported.Error(), action)
	}

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}

// @Title Get expansion status
// @Description graph expansion config
// @Param	module	path	string	true	"module name"
// @Success 200 {object} models.ExpansionStatus expansion status
// @Failure 403 string error
// @router /expansion/:module [get]
func (c *AdminController) GetExpansion() {
	var err error

	module := c.GetString(":module")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	obj, err := op.ExpansionGet(module)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, obj)
	}
}

// @Title set expansion (add endpoint)
// @Description expansion set
// @Param	body	body 	models.ExpansionStatus	true		"body for endpoints content"
// @Success 200 string success
// @Failure 403 string error
// @router /expansion/:module [put]
func (c *AdminController) SetExpansion() {
	var err error
	var status models.ExpansionStatus
	var cur_status *models.ExpansionStatus

	module := c.GetString(":module")
	json.Unmarshal(c.Ctx.Input.RequestBody, &status)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	cur_status, err = op.ExpansionGet(module)
	if err != nil {
		goto out
	}

	if !cur_status.Migrating && status.Migrating {
		err = op.ExpansionBegin(module, status.NewEndpoint)
		goto out
	}

	if cur_status.Migrating && !status.Migrating {
		err = op.ExpansionFinish(module)
		goto out
	}

	err = models.ErrParam
out:
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, "success")
	}
}
