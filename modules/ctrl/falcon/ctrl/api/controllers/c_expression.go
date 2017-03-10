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

// Operations about Expressions
type ExpressionController struct {
	BaseController
}

// @Title CreateExpression
// @Description create expressions
// @Param	body	body 	models.ExpressionAction	true	"body for expression content"
// @Success 200 {object} models.Id models.Expression.Id
// @Failure 403 string error
// @router / [post]
func (c *ExpressionController) CreateExpression() {
	var ea models.ExpressionAction
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	json.Unmarshal(c.Ctx.Input.RequestBody, &ea)
	beego.Debug("pause", ea.Expression.Pause)

	id, err := op.AddAction(&ea.Action)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	ea.Expression.ActionId = id
	ea.Expression.CreateUserId = op.User.Id
	id, err = op.AddExpression(&ea.Expression)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title GetExpressionsCnt
// @Description get Expressions number
// @Param   query	query   string	false    "expression name"
// @Param   mine	query   bool	false    "only show mine expressions, default true"
// @Success 200 {object} models.Total expression number
// @Failure 403 string error
// @router /cnt [get]
func (c *ExpressionController) GetExpressionsCnt() {
	var user_id int64
	query := strings.TrimSpace(c.GetString("query"))
	mine, _ := c.GetBool("mine", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if mine {
		user_id = op.User.Id
	}
	cnt, err := op.GetExpressionsCnt(query, user_id)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetExpressions
// @Description get all Expressions
// @Param   query	query   string  false    "expression name"
// @Param   mine	query   bool	false    "only show mine expressions, default true"
// @Param   per		query   int     false    "per page number"
// @Param   offset	query   int     false    "offset  number"
// @Success 200 {object} []models.ExpressionUi expressionuis
// @Failure 403 string error
// @router /search [get]
func (c *ExpressionController) GetExpressions() {
	var user_id int64
	query := strings.TrimSpace(c.GetString("query"))
	mine, _ := c.GetBool("mine", true)
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	if mine {
		user_id = op.User.Id
	}
	ret, err := op.GetExpressions(query, user_id, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title Get
// @Description get expression by id
// @Param	id	path 	int	true	"The key for staticblock"
// @Success 200 {object} models.ExpressionAction expression and action info
// @Failure 403 string error
// @router /:id [get]
func (c *ExpressionController) GetExpressionAction() {

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if ret, err := op.GetExpressionAction(id); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title UpdateExpressionAction
// @Description update the expression
// @Param	id	path 	string	true	"The id you want to update"
// @Param	body	body 	models.Expression	true	"body for expression content"
// @Success 200 {object} models.Expression expression
// @Failure 403 string error
// @router /:id [put]
func (c *ExpressionController) UpdateExpressionAction() {
	var ea models.ExpressionAction

	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	json.Unmarshal(c.Ctx.Input.RequestBody, &ea)

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if u, err := op.UpdateExpressionAction(id, &ea); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, u)
	}
}

// @Title UpdateExpression
// @Description update the expression
// @Param	id	query 	string	true	"The id you want to update"
// @Param	pause	query 	int	true	"1: pause, 0: not pause"
// @Success 200 null success
// @Failure 403 string error
// @router /pause [put]
func (c *ExpressionController) PauseExpression() {
	var pause int

	id, err := c.GetInt64("id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}
	pause, err = c.GetInt("pause")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if _, err := op.PauseExpression(id, pause); err != nil {
		c.SendMsg(400, err.Error())
	} else {
		c.SendMsg(200, nil)
	}
}

// @Title DeleteExpression
// @Description delete the expression
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 string error
// @router /:id [delete]
func (c *ExpressionController) DeleteExpression() {
	id, err := c.GetInt64(":id")
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	err = op.DeleteExpression(id)
	if err != nil {
		c.SendMsg(403, err.Error())
		return
	}

	beego.Debug("delete success!")

	c.SendMsg(200, "delete success!")
}
