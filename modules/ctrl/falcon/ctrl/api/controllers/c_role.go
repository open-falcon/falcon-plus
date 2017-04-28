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

// Operations about Roles
type RoleController struct {
	BaseController
}

// @Title CreateRole
// @Description create roles
// @Param	body	body 	models.Role	true	"body for role content"
// @Success 200 {object} models.Id Id
// @Failure 403 string error
// @router / [post]
func (c *RoleController) CreateRole() {
	var role models.Role
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &role)
	role.Id = 0

	if id, err := op.AddRole(&role); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetRolesCnt
// @Description get Roles number
// @Param   query     query   string  false    "role name"
// @Success 200 {object} models.Total role number
// @Failure 403 string error
// @router /cnt [get]
func (c *RoleController) GetRolesCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetRolesCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetRoles
// @Description get all Roles
// @Param   query     query   string  false    "role name"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} []models.Role roles info
// @Failure 403 string error
// @router /search [get]
func (c *RoleController) GetRoles() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	roles, err := op.GetRoles(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, roles)
	}
}

// @Title Get
// @Description get role by id
// @Param	id	path 	int	true	"The key for staticblock"
// @Success 200 {object} models.Role role info
// @Failure 403 string error
// @router /:id [get]
func (c *RoleController) GetRole() {
	id, err := c.GetInt64(":id")

	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		role, err := op.GetRole(id)
		if err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, role)
		}
	}
}

// @Title UpdateRole
// @Description update the role
// @Param	id	path 	string	true	"The id you want to update"
// @Param	body	body 	models.Role	true	"body for role content"
// @Success 200 {object} models.Role role info
// @Failure 403 string error
// @router /:id [put]
func (c *RoleController) UpdateRole() {
	var role models.Role

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &role)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if o, err := op.UpdateRole(id, &role); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, o)
	}
}

// @Title DeleteRole
// @Description delete the role
// @Param	id	path 	string	true	"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 string error
// @router /:id [delete]
func (c *RoleController) DeleteRole() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteRole(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
