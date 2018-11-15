// Copyright 2018 Xiaomi, Inc.
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

package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/controller/event"
)

func Start(logLevel string, webPort string) {
	if logLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	configRoutes(router)
	if err := router.Run(webPort); err != nil {
		log.Fatal("am start failed, ", err)
	}
}

func configRoutes(router *gin.Engine) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})

	router.POST("/v1/recv", event.RecvAlarmEvent)

	api := router.Group("/api/v1")
	ConfigApiRoutes(api)
}
