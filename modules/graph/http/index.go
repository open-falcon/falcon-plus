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
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
)

func configIndexRoutes() {
	// 触发索引全量更新, 同步操作
	router.GET("/index/updateAll", func(c *gin.Context) {
		go index.UpdateIndexAllByDefaultStep()
		JSONR(c, 200, gin.H{"msg": "ok"})
	})

	// 获取索引全量更新的并行数
	router.GET("/index/updateAll/concurrent", func(c *gin.Context) {
		JSONR(c, 200, gin.H{"msg": "ok", "value": index.GetConcurrentOfUpdateIndexAll()})
	})

	type APIIndexItemInput struct {
		Endpoint string `json:"endpoint" form:"endpoint" binding:"required"`
		Metric   string `json:"metric" form:"metric" binding:"required"`
		Step     int    `json:"step" form:"step" binding:"required"`
		Dstype   string `json:"dstype" form:"dstype" binding:"required"`
		Tags     string `json:"tags" form:"tags"`
	}

	// 更新一条索引数据,用于手动建立索引 endpoint metric step dstype tags
	router.POST("/api/v2/index", func(c *gin.Context) {
		inputs := []*APIIndexItemInput{}
		if err := c.Bind(&inputs); err != nil {
			c.AbortWithError(500, err)
			return
		}

		for _, in := range inputs {
			err, tags := cutils.SplitTagsString(in.Tags)
			if err != nil {
				log.Error("split tags:", in.Tags, "error:", err)
				continue
			}

			err = index.UpdateIndexOne(in.Endpoint, in.Metric, tags, in.Dstype, in.Step)
			if err != nil {
				log.Error("build index fail, item:", in, "error:", err)
			} else {
				log.Debug("build index manually", in)
			}
		}

		JSONR(c, 200, gin.H{"msg": "ok"})
	})
}
