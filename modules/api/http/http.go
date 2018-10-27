package http

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	yaag_gin "github.com/masato25/yaag/gin"
	"github.com/masato25/yaag/yaag"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/alarm"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/dashboard_graph"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/dashboard_screen"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/expression"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/graph"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/host"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/mockcfg"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/strategy"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/template"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/uic"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/spf13/viper"
)

var routes *gin.Engine

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func configRoutes() {
	routes.Use(utils.CORS())
	routes.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, I'm Falcon+ (｡A｡)")
	})
	graph.Routes(routes)
	uic.Routes(routes)
	template.Routes(routes)
	strategy.Routes(routes)
	host.Routes(routes)
	expression.Routes(routes)
	mockcfg.Routes(routes)
	dashboard_graph.Routes(routes)
	dashboard_screen.Routes(routes)
	alarm.Routes(routes)
	configCommonRoutes()
}

func Start(vip *viper.Viper) {
	go startHttpServer(vip)
}

func startHttpServer(vip *viper.Viper) {
	if vip.GetString("log_level") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	routes = gin.Default()
	if vip.GetBool("gen_doc") {
		yaag.Init(&yaag.Config{
			On:       true,
			DocTitle: "Gin",
			DocPath:  vip.GetString("gen_doc_path"),
			BaseUrls: map[string]string{"Production": "/api/v1", "Staging": "/api/v1"},
		})
		routes.Use(yaag_gin.Document())
	}
	//start gin server
	addr := vip.GetString("web_port")
	log.Debugf("will start with port:%v", addr)

	configRoutes()

	go routes.Run(addr)
}
