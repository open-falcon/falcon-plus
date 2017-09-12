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

package graph

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
	authapi := r.Group("/api/v1")
	authapi.Use(utils.AuthSessionMidd)
	authapi.GET("/graph/endpointobj", EndpointObjGet)
	authapi.GET("/graph/endpoint", EndpointRegexpQuery)
	authapi.GET("/graph/endpoint_counter", EndpointCounterRegexpQuery)
	authapi.POST("/graph/history", QueryGraphDrawData)
	authapi.POST("/graph/lastpoint", QueryGraphLastPoint)
	authapi.DELETE("/graph/endpoint", DeleteGraphEndpoint)
	authapi.DELETE("/graph/counter", DeleteGraphCounter)

	grfanaapi := r.Group("/api")
	grfanaapi.GET("/v1/grafana", GrafanaMainQuery)
	grfanaapi.GET("/v1/grafana/metrics/find", GrafanaMainQuery)
	grfanaapi.POST("/v1/grafana/render", GrafanaRender)
	grfanaapi.GET("/v1/grafana/render", GrafanaRender)

}
