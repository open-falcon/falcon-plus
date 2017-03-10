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
	"github.com/astaxie/beego"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
)

// Operations about Auth
type AuthController struct {
	BaseController
}

// @Title get support auth modules
// @Description get support auth modules
// @Success 200 {object} []string  modules list
// @Failure 405 string error
// @router /modules [get]
func (c *AuthController) Modules() {
	m := []string{}
	for k, _ := range models.Auths {
		m = append(m, k)
	}

	c.SendMsg(200, m)
}

// @Title OAuth Login
// @Description auth login
// @Param	module	path	string	true	"the module you want to use(github/google)"
// @Success 302 redirect
// @Failure 405 string error
// @router /login/:module [get]
func (c *AuthController) Authorize() {
	module := c.GetString(":module")

	auth, ok := models.Auths[module]
	if !ok {
		c.SendMsg(405, models.ErrNoModule.Error())
		return
	}

	URL := auth.AuthorizeUrl(c.Ctx)
	if URL == "" {
		c.SendMsg(405, nil)
		return
	}

	c.Ctx.Redirect(302, URL)
}

// @Title OAuth module callback handle
// @Description Auth module callback handle
// @Param	module	path	string	true	"the module you want to use"
// @Success 302 redirect to RedirectUrl(default "/")
// @Failure 406 not acceptable
// @router /callback/:module [get]
func (c *AuthController) Callback() {
	auth, ok := models.Auths[c.GetString(":module")]
	if !ok {
		c.SendMsg(406, models.ErrNoModule.Error())
		return
	}
	cb := c.GetString("cb")

	uuid, err := auth.CallBack(c.Ctx)
	if err != nil {
		c.SendMsg(406, err.Error())
		return
	}

	if _, err = c.Access(uuid); err != nil {
		c.SendMsg(406, err.Error())
		return
	}

	c.Ctx.Redirect(302, "/#"+cb)
}

// @Title AuthLogin
// @Description auth login, such as ldap auth
// @Success 200 {object} models.OperatorInfo operator info, reload user's tokens
// @Failure 406 not acceptable
// @router /info [get]
func (c *AuthController) Info() {
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)
	if op.User != nil {
		op.Token = op.UserTokens()
		c.SetSession("token", op.Token)
	}
	c.SendMsg(200, op.Info())
}

// @Title AuthLogin
// @Description auth login, such as ldap auth
// @Param	username	query	string	false	"username for login"
// @Param	password	query	string	false	"passworld for login"
// @Param	method		query	string	false	"login method"
// @Success 200 {object} models.OperatorInfo  operator info
// @Failure 406 not acceptable
// @router /login [post]
func (c *AuthController) PostLogin() {
	var (
		op           *models.Operator
		err          error
		ok           bool
		uuid, method string
		auth         models.AuthInterface
	)
	op, ok = c.Ctx.Input.GetData("op").(*models.Operator)
	if ok && op.User != nil {
		goto out
	}

	if method = c.GetString("method"); method == "" {
		err = models.ErrParam
		goto out_err
	}

	if auth, ok = models.Auths[method]; !ok {
		err = models.ErrNoExits
		goto out_err
	}

	if ok, uuid, err = auth.Verify(c); !ok {
		err = models.ErrLogin
		goto out_err
	}

	op, _ = c.Access(uuid)
out:
	c.SendMsg(200, op.Info())
	return

out_err:
	c.SendMsg(406, err.Error())
}

// @Title Auth Logout
// @Description user logout, reset cookie
// @Success 200 {string} logout success!
// @Failure 405 Method Not Allowed
// @router /logout [get]
func (c *AuthController) Logout() {
	if uid := c.GetSession("uid"); uid != nil {
		c.DelSession("uid")
		c.SendMsg(200, "logout success!")
	} else {
		c.SendMsg(405, models.ErrNoLogged.Error())
	}
}

func (c *AuthController) Access(uuid string) (op *models.Operator, err error) {
	op, _ = c.Ctx.Input.GetData("op").(*models.Operator)
	op.User, err = op.GetUserByUuid(uuid)
	if err != nil {
		// uuid no exist, create
		sys, _ := models.GetUser(1, op.O)
		sysOp := &models.Operator{
			User:  sys,
			O:     op.O,
			Token: models.SYS_F_A_TOKEN | models.SYS_F_O_TOKEN,
		}
		op.User, err = sysOp.AddUser(&models.User{Uuid: uuid, Name: uuid})
		if err != nil {
			beego.Debug("add user failed ", err.Error())
			return
		}
	}
	op.Token = op.UserTokens()

	beego.Debug("get login user ", op.User)
	c.SetSession("uid", op.User.Id)
	c.SetSession("token", op.Token)
	return
}
