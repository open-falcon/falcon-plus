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
