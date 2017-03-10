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
package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl"
	_ "github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models/session"
)

const (
	DB_PREFIX   = ""
	PAGE_PER    = 10
	SYS_R_TOKEN = "falcon_read"
	SYS_O_TOKEN = "falcon_operate"
	SYS_A_TOKEN = "falcon_admin"
)

const (
	SYS_F_R_TOKEN = 1 << iota
	SYS_F_O_TOKEN
	SYS_F_A_TOKEN
)

const (
	_ = iota
	SYS_IDX_R_TOKEN
	SYS_IDX_O_TOKEN
	SYS_IDX_A_TOKEN
)

var (
	dbTables = []string{
		"action",
		"expression",
		"host",
		"kv",
		"log",
		"role",
		"session",
		"strategy",
		"tag",
		"tag_host",
		"tag_rel",
		"tag_tpl",
		"team",
		"team_user",
		"template",
		"token",
		"tpl_rel",
		"trigger",
		"user",
	}
)

// ctl meta name
const (
	CTL_M_HOST = iota
	CTL_M_ROLE
	CTL_M_SYSTEM
	CTL_M_TAG
	CTL_M_USER
	CTL_M_TOKEN
	CTL_M_TPL
	CTL_M_RULE
	CTL_M_TEMPLATE
	CTL_M_TRIGGER
	CTL_M_EXPRESSION
	CTL_M_TEAM
	CTL_M_TAG_HOST
	CTL_M_TAG_TPL
	CTL_M_SIZE
)

// ctl method name
const (
	CTL_A_ADD = iota
	CTL_A_DEL
	CTL_A_SET
	CTL_A_GET
	CTL_A_SIZE
)

type Ids struct {
	Ids []int64 `json:"ids"`
}

type Id struct {
	Id int64 `json:"id"`
}

type Total struct {
	Total int64 `json:"total"`
}

var (
	moduleCache  [CTL_M_SIZE]cache
	sysTagSchema *TagSchema

	moduleName = [CTL_M_SIZE]string{
		"host", "role", "system", "tag", "user", "token",
		"template", "rule", "trigger", "expression", "team",
	}

	actionName = [CTL_A_SIZE]string{
		"add", "del", "set", "get",
	}
)

func initModels(conf *falcon.ConfCtrl) (err error) {
	if err = initConfig(conf); err != nil {
		panic(err)
	}
	if err = initAuth(conf); err != nil {
		panic(err)
	}
	if err = initCache(conf); err != nil {
		panic(err)
	}
	if err = initMetric(conf); err != nil {
		panic(err)
	}
	return nil
}

func initMetric(c *falcon.ConfCtrl) error {
	for _, m := range c.Metrics {
		metrics = append(metrics, &Metric{Name: m})
	}
	return nil
}

func initAuth(c *falcon.ConfCtrl) error {
	Auths = make(map[string]AuthInterface)
	for _, name := range strings.Split(c.Ctrl.Str(falcon.C_AUTH_MODULE), ",") {
		if auth, ok := allAuths[name]; ok {
			if auth.Init(c) == nil {
				Auths[name] = auth
			}
		}
	}
	return nil
}

// called by (p *Ctrl) Init()
// already load file config and def config
// will load db config
func initConfig(conf *falcon.ConfCtrl) error {

	beego.Debug(fmt.Sprintf("%s Init()", conf.Name))

	conf.Agent.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["agent"])
	conf.Loadbalance.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["loadbalance"])
	conf.Backend.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["backend"])
	conf.Ctrl.Set(falcon.APP_CONF_DEFAULT, falcon.ConfDefault["ctrl"])

	c := &conf.Ctrl
	dsn := c.Str(falcon.C_DSN)
	dbMaxConn, _ := c.Int(falcon.C_DB_MAX_CONN)
	dbMaxIdle, _ := c.Int(falcon.C_DB_MAX_IDLE)

	// config
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "falconSessionId"
	beego.BConfig.WebConfig.Session.SessionProvider = "mysql"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = dsn
	beego.BConfig.WebConfig.Session.SessionDisableHTTPOnly = false
	beego.BConfig.WebConfig.StaticDir["/"] = "static"
	beego.BConfig.WebConfig.StaticDir["/static"] = "static/static"

	// connect db, can not register db twice  :(
	orm.RegisterDataBase("default", "mysql", dsn, dbMaxIdle, dbMaxConn)

	// get config from db
	o := orm.NewOrm()
	if c, err := GetDbConfig(o, "ctrl"); err == nil {
		conf.Ctrl.Set(falcon.APP_CONF_DB, c)
	}

	// config -> beego config
	if addr := strings.Split(c.Str(falcon.C_HTTP_ADDR), ":"); len(addr) == 2 {
		beego.BConfig.Listen.HTTPAddr = addr[0]
		beego.BConfig.Listen.HTTPPort, _ = strconv.Atoi(addr[1])
	} else if len(addr) == 1 {
		beego.BConfig.Listen.HTTPPort, _ = strconv.Atoi(addr[0])
	}
	beego.BConfig.AppName = conf.Name
	beego.BConfig.RunMode = c.Str(falcon.C_RUN_MODE)
	beego.BConfig.WebConfig.EnableDocs, _ = c.Bool(falcon.C_ENABLE_DOCS)
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime, _ = c.Int64(falcon.C_SEESION_GC_MAX_LIFETIME)
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime, _ = c.Int(falcon.C_SESSION_COOKIE_LIFETIME)

	if beego.BConfig.RunMode == "dev" {
		beego.Debug("orm debug on")
		orm.Debug = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/doc"] = "swagger"
	}

	// tag
	var err error
	sysTagSchema, err = NewTagSchema(c.Str(falcon.C_TAG_SCHEMA))

	return err
}

func initCache(c *falcon.ConfCtrl) error {
	for _, module := range strings.Split(
		c.Ctrl.Str(falcon.C_CACHE_MODULE), ",") {
		for k, v := range moduleName {
			if v == module {
				moduleCache[k] = cache{
					enable: true,
					data:   make(map[int64]interface{}),
				}
				break
			}
		}
	}
	return nil
}
func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModelWithPrefix("",
		new(User), new(Host), new(Tag),
		new(Role), new(Token), new(Log),
		new(Tag_rel), new(Tpl_rel), new(Team),
		new(Template), new(Expression), new(Action),
		new(Strategy))

	ctrl.RegisterPrestart(initModels)
	ctrl.RegisterReload(initModels)
}
