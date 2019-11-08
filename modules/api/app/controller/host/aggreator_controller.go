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

package host

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
	log "github.com/sirupsen/logrus"
)

func GetAggregatorListOfGrp(c *gin.Context) {
	var (
		limit int
		page  int
		err   error
	)
	pageTmp := c.DefaultQuery("page", "")
	limitTmp := c.DefaultQuery("limit", "")
	page, limit, err = h.PageParser(pageTmp, limitTmp)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	grpIDtmp := c.Params.ByName("host_group")
	if grpIDtmp == "" {
		h.JSONR(c, badstatus, "grp id is missing")
		return
	}
	grpID, err := strconv.Atoi(grpIDtmp)
	if err != nil {
		log.Debugf("grpIDtmp: %v", grpIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	aggregators := []f.Cluster{}
	var dt *gorm.DB
	if limit != -1 && page != -1 {
		dt = db.Falcon.Raw("SELECT * from cluster WHERE grp_id = ? limit ?,?", grpID, page, limit).Scan(&aggregators)
	} else {
		dt = db.Falcon.Where("grp_id = ?", grpID).Find(&aggregators)
	}
	if dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	hostgroupName := ""
	if len(aggregators) != 0 {
		hostgroupName, err = aggregators[0].HostGroupName()
		if err != nil {
			h.JSONR(c, badstatus, err)
			return
		}
	}

	h.JSONR(c, map[string]interface{}{
		"hostgroup":   hostgroupName,
		"aggregators": aggregators,
	})
	return
}

func GetAggregator(c *gin.Context) {
	aggIDtmp := c.Params.ByName("id")
	if aggIDtmp == "" {
		h.JSONR(c, badstatus, "agg id is missing")
		return
	}
	aggID, err := strconv.Atoi(aggIDtmp)
	if err != nil {
		log.Debugf("aggIDtmp: %v", aggIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	aggregator := f.Cluster{}
	if dt := db.Falcon.Where("id = ?", aggID).Find(&aggregator); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, aggregator)
	return
}

type APICreateAggregatorInput struct {
	GrpId       int64  `json:"hostgroup_id" binding:"required"`
	Numerator   string `json:"numerator" binding:"required"`
	Denominator string `json:"denominator" binding:"required"`
	Endpoint    string `json:"endpoint" binding:"required"`
	Metric      string `json:"metric" binding:"required"`
	Tags        string `json:"tags" binding:"exists"`
	Step        int    `json:"step" binding:"required"`
	// DsType      string `json:"ds_type" binding:"exists"`
}

func CreateAggregator(c *gin.Context) {
	var inputs APICreateAggregatorInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, fmt.Sprintf("binding error: %v", err))
		return
	}
	user, _ := h.GetUser(c)
	if !user.IsAdmin() {
		hostgroup := f.HostGroup{ID: inputs.GrpId}
		if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
			h.JSONR(c, expecstatus, fmt.Sprintf("find hostgroup error: %v", dt.Error.Error()))
			return
		}
		if hostgroup.CreateUser != user.Name {
			h.JSONR(c, badstatus, "You don't have permission!")
			return
		}
	}
	agg := f.Cluster{
		GrpId:       inputs.GrpId,
		Numerator:   inputs.Numerator,
		Denominator: inputs.Denominator,
		Endpoint:    inputs.Endpoint,
		Metric:      inputs.Metric,
		Tags:        inputs.Tags,
		DsType:      "GAUGE",
		Step:        inputs.Step,
		Creator:     user.Name}
	if dt := db.Falcon.Create(&agg); dt.Error != nil {
		h.JSONR(c, expecstatus, fmt.Sprintf("create aggregator got error: %v", dt.Error.Error()))
		return
	}
	h.JSONR(c, agg)
	return
}

type APIUpdateAggregatorInput struct {
	ID          int64  `json:"id" binding:"required"`
	Numerator   string `json:"numerator" binding:"required"`
	Denominator string `json:"denominator" binding:"required"`
	Endpoint    string `json:"endpoint" binding:"required"`
	Metric      string `json:"metric" binding:"required"`
	Tags        string `json:"tags" binding:"exists"`
	Step        int    `json:"step" binding:"required"`
	// DsType      string `json:"ds_type" binding:"exists"`
}

func UpdateAggregator(c *gin.Context) {
	var inputs APIUpdateAggregatorInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	aggregator := f.Cluster{ID: inputs.ID}
	if dt := db.Falcon.Find(&aggregator); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	user, _ := h.GetUser(c)
	if !user.IsAdmin() {
		hostgroup := f.HostGroup{ID: aggregator.GrpId}
		if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
			h.JSONR(c, expecstatus, fmt.Sprintf("find hostgroup got error: %v", dt.Error.Error()))
			return
		}
		//only admin & aggregator creator can update it
		if hostgroup.CreateUser != user.Name && aggregator.Creator != user.Name {
			h.JSONR(c, badstatus, "You don't have permission!")
			return
		}
	}
	uaggregator := map[string]interface{}{
		"Numerator":   inputs.Numerator,
		"Denominator": inputs.Denominator,
		"Endpoint":    inputs.Endpoint,
		"Metric":      inputs.Metric,
		"Tags":        inputs.Tags,
		"Step":        inputs.Step}
	if dt := db.Falcon.Model(&aggregator).Where("id = ?", aggregator.ID).Update(uaggregator).Find(&aggregator); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, aggregator)
	return
}

func DeleteAggregator(c *gin.Context) {
	aggIDtmp := c.Params.ByName("id")
	if aggIDtmp == "" {
		h.JSONR(c, badstatus, "agg id is missing")
		return
	}
	aggID, err := strconv.Atoi(aggIDtmp)
	if err != nil {
		log.Debugf("aggIDtmp: %v", aggIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	aggregator := f.Cluster{}
	if dt := db.Falcon.Where("id = ?", aggID).Find(&aggregator); dt.Error != nil {
		h.JSONR(c, expecstatus, fmt.Sprintf("find aggregator got error: %v", dt.Error.Error()))
		return
	}
	user, _ := h.GetUser(c)
	if !user.IsAdmin() {
		hostgroup := f.HostGroup{}
		if dt := db.Falcon.Where("id = ?", aggregator.GrpId).Find(&hostgroup); dt.Error != nil {
			h.JSONR(c, expecstatus, fmt.Sprintf("find hostgroup got error: %v", dt.Error.Error()))
			return
		}
		if hostgroup.CreateUser != user.Name {
			h.JSONR(c, badstatus, "You don't have permission!")
			return
		}
	}

	if dt := db.Falcon.Table("cluster").Where("id = ?", aggID).Delete(&aggregator); dt.Error != nil {
		h.JSONR(c, expecstatus, fmt.Sprintf("delete aggregator got error: %v", dt.Error))
		return
	}
	h.JSONR(c, fmt.Sprintf("aggregator:%v has been deleted", aggID))
	return
}
