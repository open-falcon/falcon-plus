package graph

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/open-falcon/modules/api/app/utils"
	"github.com/open-falcon/open-falcon/modules/api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	authapi := r.Group("/api/v1")
	authapi.Use(utils.AuthSessionMidd)
	authapi.GET("/graph/endpoint", EndpointRegexpQuery)
	authapi.GET("/graph/endpoint_counter", EndpointCounterRegexpQuery)
	authapi.POST("/graph/history", QueryGraphDrawData)
}
