package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/expression"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/graph"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/host"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/mockcfg"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/strategy"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/template"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller/uic"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
)

func StartGin(port string, r *gin.Engine) {
	r.Use(utils.CORS())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, I'm Falcon+ (｡A｡)")
	})
	graph.Routes(r)
	uic.Routes(r)
	template.Routes(r)
	strategy.Routes(r)
	host.Routes(r)
	expression.Routes(r)
	mockcfg.Routes(r)
	r.Run(port)
}
