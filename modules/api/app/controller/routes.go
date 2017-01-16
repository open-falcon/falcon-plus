package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masato25/owl_backend/app/controller/expression"
	"github.com/masato25/owl_backend/app/controller/graph"
	"github.com/masato25/owl_backend/app/controller/host"
	"github.com/masato25/owl_backend/app/controller/mockcfg"
	"github.com/masato25/owl_backend/app/controller/strategy"
	"github.com/masato25/owl_backend/app/controller/template"
	"github.com/masato25/owl_backend/app/controller/uic"
	"github.com/masato25/owl_backend/app/utils"
)

func StartGin(port string, r *gin.Engine) {
	r.Use(utils.CORS())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, I'm OWL (｡A｡)")
	})
	graph.Routes(r)
	uic.Routes(r)
	template.Routes(r)
	strategy.Routes(r)
	host.Routes(r)
	expression.Routes(r)
	mockcfg.Routes(r)
	r.Run()
}
