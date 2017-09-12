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

package http

import (
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/graph/proc"
)

func configProcRoutes() {
	// counter
	router.GET("/counter/all", func(c *gin.Context) {
		JSONR(c, 200, proc.GetAll())
	})

	// compatible with falcon task monitor
	router.GET("/statistics/all", func(c *gin.Context) {
		ret := make(map[string]interface{})
		ret["msg"] = "success"
		ret["data"] = proc.GetAll()
		JSONR(c, 200, ret)
	})

}
