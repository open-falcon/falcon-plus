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

// Operations about Strategys
type StrategyController struct {
	BaseController
}

// @Title CreateStrategy
// @Description create strategys
// @Param	body	body 	models.Strategy	true	"body for strategy content"
// @Success 200 {object} models.Id Id
// @Failure 403 string error
// @router / [post]
func (c *StrategyController) CreateStrategy() {
	var strategy models.Strategy
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	json.Unmarshal(c.Ctx.Input.RequestBody, &strategy)
	strategy.Id = 0

	if id, err := op.AddStrategy(&strategy); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetStrategysCnt
// @Description get Strategys number
// @Param   tid     query   int     false       "template id"
// @Param   query   query   string  false       "strategy name"
// @Success 200 {object} models.Total strategy total number
// @Failure 403 string error
// @router /cnt [get]
func (c *StrategyController) GetStrategysCnt() {
	tid, _ := c.GetInt64("tid", 0)
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetStrategysCnt(tid, query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetStrategys
// @Description get all Strategys
// @Param   tid		query   int     false       "template id"
// @Param   query	query   string  false       "strategy name"
// @Param   per		query   int     false       "per page number"
// @Param   offset	query   int     false       "offset  number"
// @Success 200 {object} []models.Strategy strategys info
// @Failure 403 string error
// @router /search [get]
func (c *StrategyController) GetStrategys() {
	tid, _ := c.GetInt64("tid", 0)
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	strategys, err := op.GetStrategys(tid, query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, strategys)
	}
}

// @Title Get
// @Description get strategy by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Strategy strategy info
// @Failure 403 string error
// @router /:id [get]
func (c *StrategyController) GetStrategy() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		strategy, err := op.GetStrategy(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, strategy)
		}
	}
}

// @Title UpdateStrategy
// @Description update the strategy
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Strategy	true		"body for strategy content"
// @Success 200 {object} models.Strategy strategy info
// @Failure 403 string error
// @router /:id [put]
func (c *StrategyController) UpdateStrategy() {
	var strategy models.Strategy

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &strategy)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if u, err := op.UpdateStrategy(id, &strategy); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteStrategy
// @Description delete the strategy
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 string error
// @router /:id [delete]
func (c *StrategyController) DeleteStrategy() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteStrategy(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
