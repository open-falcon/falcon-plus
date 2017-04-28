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

// Operations about Relations
type RelController struct {
	BaseController
}

// @Title Get vue tag tree
// @Description get tags for vue tree
// @Param	id		body 	int64	true	"tag id"
// @Success 200 {object} []models.TreeNode All nodes under the current node
// @Failure 403 string error
// @router /treeNode [get]
func (c *RelController) GetTreeNodes() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	tag_id, _ := c.GetInt64("id", 0)

	nodes, err := op.GetTreeNodes(tag_id)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, nodes)
	}
}

// @Title Get vue tag tree
// @Description get tags for vue tree
// @Success 200 {object} []models.TreeNode all nodes of the tree(read)
// @Failure 403 string error
// @router /tree [get]
func (c *RelController) GetTree() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	tree, err := op.GetTree()
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, []models.TreeNode{*tree})
	}
}

// @Title Get vue tree node(operate)
// @Description get tags for vue tree
// @Success 200 {object} []int64 all ids of the node that can be operated
// @Failure 403 string error
// @router /tree/opnode [get]
func (c *RelController) GetOpNode() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	rel, _ := op.GetTreeOpNode()
	c.SendMsg(200, rel)
}

// @Title GetTagNodes
// @Description get tags for ztree
// @Param	id		body 	int64	true	"tag id"
// @Param	name		body 	string	true	"tag name"
// @Param	lv		body 	int	true	"tag level"
// @Param	otherParam	body 	string	true	"zTreeAsyncTest"
// @Success 200 {object} []models.TagNode All nodes under the current node for ztree
// @Failure 403 string error
// @router /zTreeNodes [post]
func (c *RelController) GetzTreeNodes() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	tag_id, _ := c.GetInt64("id", 0)
	//name := c.GetString("name")
	//lv, _ := c.GetInt("lv")

	nodes, err := op.GetTagTags(tag_id)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.Data["json"] = nodes
		c.ServeJSON()
	}
}

// @Title GetTagHostCnt
// @Description get Tag-Host number
// @Param	query	query   string  false	"host name"
// @Param	tag_id	query   int	true	"tag id"
// @Param	deep	query   bool	false	"search sub tag"
// @Success 200 {object} models.Total total number
// @Failure 403 string error
// @router /tag/host/cnt [get]
func (c *RelController) GetTagHostCnt() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	n, err := op.GetTagHostCnt(tag_id, query, deep)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetHost
// @Description get all Host
// @Param	tag_id		query	int	true	"tag id"
// @Param	query	query	string	false	"host name"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	per		query	int	false	"per page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.RelTagHost tag host info
// @Failure 403 string error
// @router /tag/host/search [get]
func (c *RelController) GetTagHost() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	ret, err := op.GetTagHost(tag_id, query, deep, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag host relation
// @Description create tag/host relation
// @Param	body	body	models.RelTagHost	true	""
// @Success 200 {object} models.Id Id
// @Failure 403 string error
// @router /tag/host [post]
func (c *RelController) CreateTagHost() {
	var rel models.RelTagHost

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if id, err := op.CreateTagHost(rel); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title create tag host relation
// @Description create tag/hosts relation
// @Param	body	body	models.RelTagHosts	true	""
// @Success 200 {object} models.Id Id
// @Failure 403 string error
// @router /tag/hosts [post]
func (c *RelController) CreateTagHosts() {
	var rel models.RelTagHosts

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if id, err := op.CreateTagHosts(rel); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title delete tag host relation
// @Description delete tag/host relation
// @Param	body		body 	models.Id	true	"relation id"
// @Success 200 {object} models.Total affected number
// @Failure 403 string error
// @router /tag/host [delete]
func (c *RelController) DelTagHost() {
	var rel models.Id
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := op.DeleteTagHost(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag host relation
// @Description delete tag/hosts relation
// @Param	body	body 	models.ids	true	"unbind multiple tag-host relation"
// @Success 200 {object} models.Total affected number
// @Failure 403 string error
// @router /tag/hosts [delete]
func (c *RelController) DelTagHosts() {
	var rel models.Ids
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := op.DeleteTagHosts(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTagTemplateCnt
// @Description get Tag-Template number
// @Param	query	query   string  false	"template name"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	mine	query   bool	false	"search mine template"
// @Param	tag_id	query   int	true	"tag id"
// @Success 200 {object} models.Total total number
// @Failure 403 string error
// @router /tag/template/cnt [get]
func (c *RelController) GetTagTplCnt() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	mine, _ := c.GetBool("mine", true)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	n, err := op.GetTagTplCnt(tag_id, query, deep, mine)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTemplate
// @Description get all Template
// @Param	tag_id	query	int	true	"tag id"
// @Param	query	query	string	false	"template name"
// @Param	deep	query   bool	false	"search sub tag"
// @Param	mine	query   bool	false	"search mine template"
// @Param	per	query	int	false	"per page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.RelTagTplUi templates info
// @Failure 403 string error
// @router /tag/template/search [get]
func (c *RelController) GetTagTpl() {
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	deep, _ := c.GetBool("deep", true)
	mine, _ := c.GetBool("mine", true)
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	ret, err := op.GetTagTpl(tag_id, query, deep, mine, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag template relation
// @Description create tag/template relation
// @Param	body	body	models.RelTagTpl	true	""
// @Success 200 {object} models.Id Id
// @Failure 403 string error
// @router /tag/template [post]
func (c *RelController) CreateTagTpl() {
	var rel models.RelTagTpl

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if id, err := op.CreateTagTpl(rel); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title create tag template relation
// @Description create tag/templates relation
// @Param	body	body	models.RelTagTpls	true	""
// @Success 200 {object} models.Id Id
// @Failure 403 string error
// @router /tag/templates [post]
func (c *RelController) CreateTagTpls() {
	var rel models.RelTagTpls

	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	if id, err := op.CreateTagTpls(rel); err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(id))
	}
}

// @Title delete tag template relation
// @Description delete tag/template relation
// @Param	body		body 	models.RelTagTpl	true	""
// @Success 200 {object} models.Total affected number
// @Failure 403 string error
// @router /tag/template [delete]
func (c *RelController) DelTagTpl() {
	var rel models.RelTagTpl
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := op.DeleteTagTpl(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag template relation
// @Description delete tag/templates relation
// @Param	body	body 	models.RelTagTpls	true	""
// @Success 200 {object} models.Total affected number
// @Failure 403 string error
// @router /tag/templates [delete]
func (c *RelController) DelTagTpls() {
	var rel models.RelTagTpls
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := op.DeleteTagTpls(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTagRoleUserCnt
// @Description get tag role user number
// @Param	query	query   string  false	"user name"
// @Param	tag_id	query   int	true	"tag id"
// @Success 200 {object} models.Total user total number
// @Failure 403 string error
// @router /tag/role/user/cnt [get]
func (c *RelController) GetTagRoleUserCnt() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	n, err := op.GetTagRoleUserCnt(global, tag_id, query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title GetTagRoleUser
// @Description get tag role user
// @Param	tag_id	query	int	true	"tag id"
// @Param	query	query	string	false	"user name"
// @Param	per	query	int	false	"per page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.Host hosts info
// @Failure 403 string error
// @router /tag/role/user/search [get]
func (c *RelController) GetTagRoleUser() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	ret, err := op.GetTagRoleUser(global, tag_id, query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag role users relation
// @Description create tag/role/users relation
// @Param	body	body 	models.RelTagRoleUser	true	""
// @Success 200 {object} models.Id affected number
// @Failure 403 string error
// @router /tag/role/user [post]
func (c *RelController) CreateTagRoleUser() {
	var rel models.RelTagRoleUser
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := op.CreateTagRoleUser(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag role user relation
// @Description delete tag/role/user relation
// @Param	body		body 	models.RelTagHost	true	""
// @Success 200 {object} models.Id affected id
// @Failure 403 string error
// @router /tag/role/user [delete]
func (c *RelController) DelTagRoleUser() {
	var rel models.RelTagRoleUser
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := op.DeleteTagRoleUser(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(n))
	}
}

// @Title GetTagRoleTokenCnt
// @Description get tag role token number
// @Param	query	query   string  false	"token name"
// @Param	tag_id	query   int	true	"tag id"
// @Success 200 {object} models.Id affected id
// @Failure 403 string error
// @router /tag/role/token/cnt [get]
func (c *RelController) GetTagRoleTokenCnt() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	n, err := op.GetTagRoleTokenCnt(global, tag_id, query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(n))
	}
}

// @Title GetTagRoleToken
// @Description get tag role token
// @Param	tag_id	query	int	true	"tag id"
// @Param	query	query	string	false	"token name"
// @Param	per	query	int	false	"per page number"
// @Param	offset	query	int	false	"offset  number"
// @Success 200 {object} []models.Host hosts info
// @Failure 403 string error
// @router /tag/role/token/search [get]
func (c *RelController) GetTagRoleToken() {
	global, _ := c.GetBool("global", false)
	tag_id, _ := c.GetInt64("tag_id", 0)
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	ret, err := op.GetTagRoleToken(global, tag_id, query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, ret)
	}
}

// @Title create tag role tokens relation
// @Description create tag/role/tokens relation
// @Param	body	body 	models.RelTagRoleToken	true	""
// @Success 200 {object} models.Total affected number
// @Failure 403 string error
// @router /tag/role/token [post]
func (c *RelController) CreateTagRoleToken() {
	var rel models.RelTagRoleToken
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := op.CreateTagRoleToken(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(n))
	}
}

// @Title delete tag role token relation
// @Description delete tag/role/token relation
// @Param	body		body 	models.RelTagHost	true	""
// @Success 200 {object} models.Id affected id
// @Failure 403 string error
// @router /tag/role/token [delete]
func (c *RelController) DelTagRoleToken() {
	var rel models.RelTagRoleToken
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	json.Unmarshal(c.Ctx.Input.RequestBody, &rel)

	n, err := op.DeleteTagRoleToken(rel)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, idObj(n))
	}
}
