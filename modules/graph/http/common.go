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
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
	"github.com/toolkits/file"

	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func configCommonRoutes() {
	// compatible with anteye
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	router.GET("/api/v2/health", func(c *gin.Context) {
		JSONR(c, 200, gin.H{"msg": "ok"})
	})

	router.GET("/api/v2/version", func(c *gin.Context) {
		JSONR(c, 200, gin.H{"value": g.VERSION})
	})

	router.GET("/api/v2/workdir", func(c *gin.Context) {
		JSONR(c, 200, gin.H{"value": file.SelfDir()})
	})

	router.GET("/api/v2/config", func(c *gin.Context) {
		JSONR(c, 200, gin.H{"value": g.Config()})
	})

	router.POST("/api/v2/config/reload", func(c *gin.Context) {
		g.ParseConfig(g.ConfigFile)
		JSONR(c, 200, gin.H{"msg": "ok"})
	})

	router.GET("/api/v2/stats/graph-queue-size", func(c *gin.Context) {
		rt := make(map[string]int)
		for i := 0; i < store.GraphItems.Size; i++ {
			keys := store.GraphItems.KeysByIndex(i)
			oneHourAgo := time.Now().Unix() - 3600

			count := 0
			for _, ckey := range keys {
				item := store.GraphItems.First(ckey)
				if item == nil {
					continue
				}

				if item.Timestamp > oneHourAgo {
					count++
				}
			}
			i_s := strconv.Itoa(i)
			rt[i_s] = count
		}
		JSONR(c, 200, rt)
	})
}
