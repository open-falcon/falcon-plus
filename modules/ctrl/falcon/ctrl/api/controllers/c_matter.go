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

	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
	"time"
)

// Operations about Matters
type MatterController struct {
	BaseController
}

// @Title GetMatters
// @Description get Matters
// @Param   status    query   int     true    "matter status"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} []models.Matters matters
// @Failure 403 string error
// @router /search [get]
func (c *MatterController) GetMatters() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	//status := alarmModels.STATUS_PENDING
	status, err := c.GetInt("status")
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	matters, err := op.QueryMatters(status, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, matters)
	}
}

// @Title GetMattersCnt
// @Description get Matters number
// @Param   status     query   int  false    "matter status"
// @Success 200 {object} int matter total number
// @Failure 403 string error
// @router /cnt [get]
func (c *MatterController) GetMattersCnt() {
	status, _ := c.GetInt("status")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetMatterCnt(status)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, cnt)
	}
}

// @Title update Matter
// @Description update matter
// @Param	id	path 	int          true	"The id you want to update"
// @Param	body	body 	models.Matter       true	"body for matter content"
// @Success 200 {object} models.Matter matter  info
// @Failure 403 string error
// @router /:id [put]
func (c *MatterController) UpdateMatter() {
	var matter models.Matter
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &matter)
	fmt.Println(matter)

	if err := op.UpdateMatter(id, matter); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, "")
	}
}

// @Title GetEvents
// @Description get Events
// @Param   matter    query   int     true    "matter id"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} []models.Event matters
// @Failure 403 string error
// @router /event/search [get]
func (c *MatterController) GetEvents() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	//status := alarmModels.STATUS_PENDING
	matter, err := c.GetInt64("matter")
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	matters := op.QueryEventsByMatter(matter, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, matters)
	}
}

// @Title GetEventCnt
// @Description get Event number
// @Param   matter   query   int  false    "matter id"
// @Success 200 {object} int event total number
// @Failure 403 string error
// @router /event/cnt [get]
func (c *MatterController) GetEventCnt() {
	matter, _ := c.GetInt64("matter")
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.QueryEventsCntByMatter(matter)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, cnt)
	}
}

// @Title Claim matter
// @Description claim matter
// @Param   body    body    models.Claim    true    "body for team content"
// @Success 200 {object} ok
// @Failure 403 string error
// @router /claim [post]
func (c *MatterController) CreateClaim() {
	var claim models.Claim
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	json.Unmarshal(c.Ctx.Input.RequestBody, &claim)
	claim.Timestamp = time.Now().Unix()
	claim.User = op.User.Name
	err := op.AddClaim(claim)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, "ok")
	}
}
