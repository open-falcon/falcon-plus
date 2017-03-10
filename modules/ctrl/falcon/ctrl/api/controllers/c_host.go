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

	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
)

// Operations about Hosts
type HostController struct {
	BaseController
}

// @Title CreateHost
// @Description create hosts
// @Param	body	body 	models.Host	true	"body for host content"
// @Success 200 {object} models.Id id
// @Failure 403 string error
// @router / [post]
func (c *HostController) CreateHost() {
	var host models.Host

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &host)
	host.Id = 0

	if id, err := op.AddHost(&host); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetHostsCnt
// @Description get Hosts number
// @Param   query     query   string  false       "host name"
// @Success 200 {object} models.Total  host total number
// @Failure 403 string error
// @router /cnt [get]
func (c *HostController) GetHostsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetHostsCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetHosts
// @Description get all Hosts
// @Param   query     query   string  false       "host name"
// @Param   per       query   int     false       "per page number"
// @Param   offset    query   int     false       "offset  number"
// @Success 200 {object} []models.Host hosts info
// @Failure 403 string error
// @router /search [get]
func (c *HostController) GetHosts() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	hosts, err := op.GetHosts(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, hosts)
	}
}

// @Title GetHost
// @Description get host by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Host host info
// @Failure 403 string error
// @router /:id [get]
func (c *HostController) GetHost() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		host, err := op.GetHost(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, host)
		}
	}
}

// @Title UpdateHost
// @Description update the host
// @Param	id		path 	string	true	"The id you want to update"
// @Param	body		body 	models.Host	true	"body for host content"
// @Success 200 {object} models.Host host info
// @Failure 403 string error
// @router /:id [put]
func (c *HostController) UpdateHost() {
	var host models.Host
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &host)

	if u, err := op.UpdateHost(id, &host); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteHost
// @Description delete the host
// @Param	id	path	string	true	"The id you want to delete"
// @Success 200 {string} "delete success!"
// @Failure 403 string error
// @router /:id [delete]
func (c *HostController) DeleteHost() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteHost(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	c.SendMsg(200, "delete success!")
}
