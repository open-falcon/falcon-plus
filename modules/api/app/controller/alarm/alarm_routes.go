package alarm

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
	alarmapi := r.Group("/api/v1/alarm")
	alarmapi.Use(utils.AuthSessionMidd)
	alarmapi.POST("/eventcases", AlarmLists)
	alarmapi.GET("/eventcases", AlarmLists)
	alarmapi.POST("/events", EventsGet)
	alarmapi.GET("/events", EventsGet)
	alarmapi.POST("/event_note", AddNotesToAlarm)
	alarmapi.GET("/event_note", GetNotesOfAlarm)
}
