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

// Operations about Teams
type TeamController struct {
	BaseController
}

// @Title CreateTeam
// @Description create teams
// @Param	body	body 	models.Team	true	"body for team content"
// @Success 200 {object} models.Id Id
// @Failure 403 string error
// @router / [post]
func (c *TeamController) CreateTeam() {
	var team models.Team
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	json.Unmarshal(c.Ctx.Input.RequestBody, &team)
	team.Creator = op.User.Id
	id, err := op.AddTeam(&team)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetTeamsCnt
// @Description get Teams number
// @Param   query     query   string  false    "team name"
// @Success 200 {object} models.Total team total number
// @Failure 403 string error
// @router /cnt [get]
func (c *TeamController) GetTeamsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	own, _ := c.GetBool("own", false)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetTeamsCnt(query, own)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetTeams
// @Description get all Teams
// @Param   query     query   string  false    "team name"
// @Param   own       query   bool    false    "check if the creator is yourself"
// @Param   per       query   int     false    "per page number"
// @Param   offset    query   int     false    "offset  number"
// @Success 200 {object} []models.Team teams info
// @Failure 403 string error
// @router /search [get]
func (c *TeamController) GetTeams() {
	query := strings.TrimSpace(c.GetString("query"))
	own, _ := c.GetBool("own", false)
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	teams, err := op.GetTeams(query, own, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, teams)
	}
}

// @Title Get team by id
// @Description get team by id
// @Param	id	path 	int	true	"team id"
// @Success 200 {object} models.Team team info
// @Failure 403 string error
// @router /:id [get]
func (c *TeamController) GetTeam() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		if t, err := op.GetTeam(id); err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, t)
		}
	}
}

// @Title Get team op.ber
// @Description get team by id
// @Param	id	path 	int	true	"team id"
// @Success 200 {object} models.TeamMembers user info
// @Failure 403 string error
// @router /:id/member [get]
func (c *TeamController) GetMember() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
		if obj, err := op.GetMember(id); err != nil {
			c.SendMsg(403, err.Error())
		} else {
			c.SendMsg(200, obj)
		}
	}
}

// @Title UpdateTeam
// @Description update the team
// @Param	id	path 	string		true	"The id you want to update"
// @Param	body	body 	models.Team	true	"body for team content"
// @Success 200 {object} models.Team team info
// @Failure 403 string error
// @router /:id [put]
func (c *TeamController) UpdateTeam() {
	var team models.Team

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &team)

	if t, err := op.UpdateTeam(id, &team); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, t)
	}
}

// @Title update Team op.bers
// @Description create teams
// @Param	body	body 	models.TeamMemberIds	true	"body for team content"
// @Success 200 {object} models.TeamMemberIds member info
// @Failure 403 string error
// @router /:id/member [put]
func (c *TeamController) UpdateMember() {
	var member models.TeamMemberIds
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &member)

	if m, err := op.UpdateMember(id, &member); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, m)
	}
}

// @Title DeleteTeam
// @Description delete the team
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} "delete success!"
// @Failure 403 string error
// @router /:id [delete]
func (c *TeamController) DeleteTeam() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteTeam(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	c.SendMsg(200, "delete success!")
}
