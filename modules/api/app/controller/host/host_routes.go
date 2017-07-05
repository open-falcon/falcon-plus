package host

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	hostr := r.Group("/api/v1")
	hostr.Use(utils.AuthSessionMidd)
	//hostgroup
	hostr.GET("/hostgroup", GetHostGroups)
	hostr.POST("/hostgroup", CrateHostGroup)
	hostr.POST("/hostgroup/host", BindHostToHostGroup)
	hostr.PUT("/hostgroup/host", UnBindAHostToHostGroup)
	hostr.GET("/hostgroup/:host_group", GetHostGroup)
	hostr.DELETE("/hostgroup/:host_group", DeleteHostGroup)

	//plugins
	hostr.GET("/hostgroup/:host_group/plugins", GetPluginOfGrp)
	hostr.POST("/plugin", CreatePlugin)
	hostr.DELETE("/plugin/:id", DeletePlugin)

	//aggreator
	hostr.GET("/hostgroup/:host_group/aggregators", GetAggregatorListOfGrp)
	hostr.GET("/aggregator/:id", GetAggregator)
	hostr.POST("/aggregator", CreateAggregator)
	hostr.PUT("/aggregator", UpdateAggregator)
	hostr.DELETE("/aggregator/:id", DeleteAggregator)

	//template
	hostr.POST("/hostgroup/template", BindTemplateToGroup)
	hostr.PUT("/hostgroup/template", UnBindTemplateToGroup)
	hostr.GET("/hostgroup/:host_group/template", GetTemplateOfHostGroup)

	//host
	hostr.GET("/host/:host_id/template", GetTplsRelatedHost)
	hostr.GET("/host/:host_id/hostgroup", GetGrpsRelatedHost)

	//maintain
	hostr.POST("/host/maintain", SetMaintain)
	hostr.DELETE("/host/maintain", UnsetMaintain)
}
