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

package template

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest

func Routes(r *gin.Engine) {
	db = config.Con()
	tmpr := r.Group("/api/v1/template")
	tmpr.Use(utils.AuthSessionMidd)
	tmpr.GET("", GetTemplates)
	tmpr.POST("", CreateTemplate)
	tmpr.GET("/:tpl_id", GetATemplate)
	tmpr.PUT("", UpdateTemplate)
	tmpr.GET("/:tpl_id/hostgroup", GetATemplateHostgroup)
	tmpr.DELETE("/:tpl_id", DeleteTemplate)
	tmpr.POST("/action", CreateActionToTmplate)
	tmpr.PUT("/action", UpdateActionToTmplate)

	actr := r.Group("/api/v1/action")
	actr.GET("/:act_id", GetActionByID)

	//simple list for ajax use
	tmpr2 := r.Group("/api/v1/template_simple")
	tmpr.Use(utils.AuthSessionMidd)
	tmpr2.GET("", GetTemplatesSimple)
}
