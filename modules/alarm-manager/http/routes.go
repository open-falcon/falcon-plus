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
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/controller/event"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/controller/fault"
)

func ConfigApiRoutes(r *gin.RouterGroup) {
	r.POST("/event", event.GetEvents)
	r.POST("/event/fault", event.GetEventsFaults)
	r.POST("/event/count", event.GetEventCount)

	r.POST("/fault", fault.Create)
	r.GET("/fault", fault.List)
	r.GET("/fault/:id", fault.Get)

	r.GET("/fault/:id/timeline", fault.GetTimeLine)

	r.GET("/fault/:id/event", fault.GetEvent)
	r.PUT("/fault/:id/event", fault.AddEvent)
	r.DELETE("/fault/:id/event", fault.DeleteEvent)

	r.GET("/fault/:id/tag", fault.GetTag)
	r.PUT("/fault/:id/tag", fault.AddTag)
	r.DELETE("/fault/:id/tag", fault.DeleteTag)

	r.GET("/fault/:id/comment", fault.GetComment)
	r.PUT("/fault/:id/comment", fault.AddComment)
	r.DELETE("/fault/:id/comment", fault.DeleteComment)

	r.PUT("/fault/:id/owner", fault.UpdateOwner)
	r.PUT("/fault/:id/state", fault.UpdateState)
	r.PUT("/fault/:id/follower", fault.UpdateFollower)
	r.PUT("/fault/:id/property", fault.UpdateBasic)
}
