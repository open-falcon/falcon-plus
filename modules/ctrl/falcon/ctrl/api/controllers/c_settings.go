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

	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
)

// Operations about porfile/config/info
type SetController struct {
	BaseController
}

// @Title Get config
// @Description get modules config
// @Param	module	path	string	true	"module name"
// @Success 200 {object} [3]map[string]string {defualt{}, conf{}, configfile{}}
// @Failure 403 string error
// @router /config/:module [get]
func (c *SetController) GetConfig() {
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

// @Title Get
// @Description get profile
// @Success 200 {object} models.User user info
// @Failure 403 string error
// @router /profile [get]
func (c *SetController) GetUser() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if user, err := op.GetUser(op.User.Id); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, user)
	}
}

// @Title Profile update
// @Description update profile
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User user info
// @Failure 403 string error
// @router /profile [put]
func (c *SetController) UpdateUser() {
	var user models.User

	json.Unmarshal(c.Ctx.Input.RequestBody, &user)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if u, err := op.UpdateUser(op.User.Id, &user); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title GetLogsCnt
// @Description get logs number
// @Param   begin  query   string  false       "end time(YYYY-MM-DD HH:mm:ss)"
// @Param   end    query   string  false       "begin time(YYYY-MM-DD HH:mm:ss)"
// @Success 200 {object} models.Total total number
// @Failure 403 string error
// @router /log/cnt [get]
func (c *SetController) GetLogsCnt() {
	begin := c.GetString("begin")
	end := c.GetString("end")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetLogsCnt(begin, end)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetLogs
// @Description get all Logs
// @Param   begin      query  string  false  "end time(YYYY-MM-DD HH:mm:ss)"
// @Param   end        query  string  false  "begin time(YYYY-MM-DD HH:mm:ss)"
// @Param   per       query   int     false  "per page number"
// @Param   offset    query   int     false  "offset  number"
// @Success 200 {object} []models.Log logs info
// @Failure 403 string error
// @router /log/search [get]
func (c *SetController) GetLogs() {
	begin := c.GetString("begin")
	end := c.GetString("end")
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	metrics, err := op.GetLogs(begin, end, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, metrics)
	}
}
