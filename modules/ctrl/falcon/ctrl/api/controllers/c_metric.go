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
package controllers

import (
	"strings"

	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
)

// Operations about Metrics
type MetricController struct {
	BaseController
}

// @Title GetMetricsCnt
// @Description get Metrics number
// @Param   query     query   string  false       "metric name"
// @Success 200 {object} models.Total  Metric total number
// @Failure 403 string error
// @router /cnt [get]
func (c *MetricController) GetMetricsCnt() {
	query := strings.TrimSpace(c.GetString("query"))
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	cnt, err := op.GetMetricsCnt(query)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, totalObj(cnt))
	}
}

// @Title GetMetrics
// @Description get all Metrics
// @Param   query     query   string  false       "metric name"
// @Param   per       query   int     false       "per page number"
// @Param   offset    query   int     false       "offset  number"
// @Success 200 {object} []models.Metric metrics info
// @Failure 403 string error
// @router /search [get]
func (c *MetricController) GetMetrics() {
	query := strings.TrimSpace(c.GetString("query"))
	per, _ := c.GetInt("per", models.PAGE_PER)
	offset, _ := c.GetInt("offset", 0)
	op, _ := c.Ctx.Input.GetData("op").(*models.Operator)

	metrics, err := op.GetMetrics(query, per, offset)
	if err != nil {
		c.SendMsg(403, err.Error())
	} else {
		c.SendMsg(200, metrics)
	}
}
