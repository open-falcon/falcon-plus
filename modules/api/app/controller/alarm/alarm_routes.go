package alarm

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	authapi := r.Group("/api/v1/alarm")
	// authapi.Use(utils.AuthSessionMidd)
	authapi.POST("/eventcases", AlarmLists)
	authapi.GET("/eventcases", AlarmLists)
	authapi.POST("/events", EventsGet)
	authapi.GET("/events", EventsGet)
}
