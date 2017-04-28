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

// Operations about Tags
type TagController struct {
	BaseController
}

// @Title CreateTag
// @Description create tags
// @Param	body	body 	models.Tag	true	"body for tag content"
// @Success 200 {object} models.Id Id
// @Failure 403 string error
// @router / [post]
func (c *TagController) CreateTag() {
	var tag models.Tag
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &tag)

	if id, err := op.AddTag(&tag); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetTagsCnt
// @Description get Tags number
// @Param   query     query   string  false       "tag name"
// @Success 200 {object} models.Total tag total number
// @Failure 403 string error
// @router /cnt [get]
func (c *TagController) GetTagsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetTagsCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTags
// @Description get all Tags
// @Param   query     query   string  false       "tag name"
// @Param   per       query   int     false       "per page number"
// @Param   offset    query   int     false       "offset  number"
// @Success 200 {object} []models.Tag  tags info
// @Failure 403 string error
// @router /search [get]
func (c *TagController) GetTags() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	tags, err := op.GetTags(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, tags)
	}
}

// @Title Get
// @Description get tag by id
// @Param	id		path 	int	true		"The key for staticblock"
// @Success 200 {object} models.Tag tag info
// @Failure 403 string error
// @router /:id [get]
func (c *TagController) GetTag() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		if tag, err := op.GetTag(id); err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, tag)
		}
	}
}

// @Title UpdateTag
// @Description update the tag
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Tag	true		"body for tag content"
// @Success 200 {object} models.Tag tag info
// @Failure 403 string error
// @router /:id [put]
func (c *TagController) UpdateTag() {
	var tag models.Tag
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &tag)

	if u, err := op.UpdateTag(id, &tag); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title DeleteTag
// @Description delete the tag
// @Param	id	path	string	true	"The id you want to delete"
// @Success 200 {string} "delete success!"
// @Failure 403 string error
// @router /:id [delete]
func (c *TagController) DeleteTag() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteTag(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
