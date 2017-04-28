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
	"strings"

	"github.com/astaxie/beego"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
)

// Operations about Triggers
type TriggerController struct {
	BaseController
}

// @Title CreateTrigger
// @Description create triggers
// @Param	body	body 	models.Trigger	true	"body for trigger content"
// @Success 200 {object} models.Id Id
// @Failure 403 string error
// @router / [post]
func (c *TriggerController) CreateTrigger() {
	var trigger models.Trigger
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &trigger)
	trigger.Id = 0

	id, err := op.AddTrigger(&trigger)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetTriggersCnt
// @Description get Triggers number
// @Param   query     query   string  false    "trigger name"
// @Success 200 {object} models.Total trigger total number
// @Failure 403 string error
// @router /cnt [get]
func (c *TriggerController) GetTriggersCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetTriggersCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTriggers
// @Description get all Triggers
// @Param   query     query   string  false    "trigger name"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} []models.Trigger triggers info
// @Failure 403 string error
// @router /search [get]
func (c *TriggerController) GetTriggers() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	triggers, err := op.GetTriggers(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, triggers)
	}
}

// @Title Get
// @Description get trigger by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Trigger
// @Failure 403 string error
// @router /:id [get]
func (c *TriggerController) GetTrigger() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		trigger, err := op.GetTrigger(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, trigger)
		}
	}
}

// @Title UpdateTrigger
// @Description update the trigger
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Trigger	true		"body for trigger content"
// @Success 200 {object} models.Trigger trigger info
// @Failure 403 string error
// @router /:id [put]
func (c *TriggerController) UpdateTrigger() {
	var trigger models.Trigger

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &trigger)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if u, err := op.UpdateTrigger(id, &trigger); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteTrigger
// @Description delete the trigger
// @Param	id		path 	string	true		"The id you want to delete"
// @Success {code:200, data:"delete success!"} delete success!
// @Failure {code:403, msg:string}
// @router /:id [delete]
func (c *TriggerController) DeleteTrigger() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteTrigger(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
