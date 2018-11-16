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

package alarm_manager

import (
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

var server string

func Routes(r *gin.Engine) {
	server = config.GetAmAddr()
	amapi := r.Group("/api/v1")
	amapi.Use(utils.AuthSessionMidd)
	amapi.POST("/event", Forwarder)
	amapi.POST("/event/fault", Forwarder)
	amapi.POST("/event/count", Forwarder)

	amapi.POST("/fault", Forwarder)
	amapi.GET("/fault", Forwarder)
	amapi.GET("/fault/:id", Forwarder)

	amapi.GET("/fault/:id/timeline", Forwarder)

	amapi.GET("/fault/:id/event", Forwarder)
	amapi.PUT("/fault/:id/event", Forwarder)
	amapi.DELETE("/fault/:id/event", Forwarder)

	amapi.GET("/fault/:id/tag", Forwarder)
	amapi.PUT("/fault/:id/tag", Forwarder)
	amapi.DELETE("/fault/:id/tag", Forwarder)

	amapi.GET("/fault/:id/comment", Forwarder)
	amapi.PUT("/fault/:id/comment", Forwarder)
	amapi.DELETE("/fault/:id/comment", Forwarder)

	amapi.PUT("/fault/:id/owner", Forwarder)
	amapi.PUT("/fault/:id/state", Forwarder)
	amapi.PUT("/fault/:id/follower", Forwarder)
	amapi.PUT("/fault/:id/property", Forwarder)
}
