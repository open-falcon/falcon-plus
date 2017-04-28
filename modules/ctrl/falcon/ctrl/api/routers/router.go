// @APIVersion 1.0.0
// @Title falcon ctrl API
// @Description Open-Falcon 是小米运维部开源的一款互联网企业级监控系统解决方案.
// @Contact yubo@xiaomi.com
// @TermsOfServiceUrl http://open-falcon.org/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html

package routers

import (
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/controllers"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
)

const (
	ACL = true
)

func init() {
	beego.InsertFilter("/v1.0/*", beego.BeforeRouter, profileFilter)
	beego.InsertFilter("/v1.0/*", beego.BeforeRouter, accessFilter)

	if ACL {
		beego.InsertFilter("/v1.0/*", beego.BeforeRouter, accessFilter)
	}
	ns := beego.NewNamespace("/v1.0",
		beego.NSNamespace("/auth", beego.NSInclude(&controllers.AuthController{})),
		beego.NSNamespace("/host", beego.NSInclude(&controllers.HostController{})),
		beego.NSNamespace("/role", beego.NSInclude(&controllers.RoleController{})),
		beego.NSNamespace("/tag", beego.NSInclude(&controllers.TagController{})),
		beego.NSNamespace("/user", beego.NSInclude(&controllers.UserController{})),
		beego.NSNamespace("/token", beego.NSInclude(&controllers.TokenController{})),
		beego.NSNamespace("/rel", beego.NSInclude(&controllers.RelController{})),
		beego.NSNamespace("/team", beego.NSInclude(&controllers.TeamController{})),
		beego.NSNamespace("/template", beego.NSInclude(&controllers.TemplateController{})),
		beego.NSNamespace("/expression", beego.NSInclude(&controllers.ExpressionController{})),
		beego.NSNamespace("/strategy", beego.NSInclude(&controllers.StrategyController{})),
		beego.NSNamespace("/settings", beego.NSInclude(&controllers.SetController{})),
		beego.NSNamespace("/metric", beego.NSInclude(&controllers.MetricController{})),
		beego.NSNamespace("/admin", beego.NSInclude(&controllers.AdminController{})),
		beego.NSNamespace("/matter", beego.NSInclude(&controllers.MatterController{})),
	)
	beego.AddNamespace(ns)
}

func accessFilter(ctx *context.Context) {
	if strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/auth") {
		return
	}

	op, ok := ctx.Input.GetData("op").(*models.Operator)
	if !ok || op.User == nil {
		http.Error(ctx.ResponseWriter, "Unauthorized", 401)
		return
	}

	if strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/settings") {
		return
	}

	if strings.HasPrefix(ctx.Request.RequestURI, "/v1.0/admin") {
		if !op.IsAdmin() {
			http.Error(ctx.ResponseWriter, "permission denied, admin only", 403)
		}
		return
	}

	switch ctx.Request.Method {
	case "GET":
		if !op.IsReader() {
			http.Error(ctx.ResponseWriter, "permission denied, read", 403)
		}
	case "POST", "PUT", "DELETE":
		if !op.IsOperator() {
			beego.Debug("TOKEN ", op.Token)
			http.Error(ctx.ResponseWriter, "permission denied, operate", 403)
		}
	default:
		http.Error(ctx.ResponseWriter, "Method Not Allowed", 405)
	}
}

func profileFilter(ctx *context.Context) {
	op := &models.Operator{O: orm.NewOrm()}
	ctx.Input.SetData("op", op)
	if id, ok := ctx.Input.Session("uid").(int64); ok {
		u, err := models.GetUser(id, op.O)
		if err != nil {
			beego.Debug("login, but can not found user")
			return
		}

		op.User = u
		op.Token, _ = ctx.Input.Session("token").(int)
	} else {
		beego.Debug("not login 2")
	}
}
