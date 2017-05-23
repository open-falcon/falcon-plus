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

package strategy

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"io/ioutil"

	"github.com/gin-gonic/gin"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
	"github.com/spf13/viper"
)

func GetStrategys(c *gin.Context) {
	var strategys []f.Strategy
	tidtmp := c.DefaultQuery("tid", "")
	if tidtmp == "" {
		h.JSONR(c, badstatus, "tid is missing")
		return
	}
	tid, err := strconv.Atoi(tidtmp)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	dt := db.Falcon.Where("tpl_id = ?", tid).Find(&strategys)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, strategys)
	return
}

type APICreateStrategyInput struct {
	Metric     string `json:"metric" binding:"required"`
	Tags       string `json:"tags"`
	MaxStep    int    `json:"max_step" binding:"required"`
	Priority   int    `json:"priority" binding:"exists"`
	Func       string `json:"func" binding:"required"`
	Op         string `json:"op" binding:"required"`
	RightValue string `json:"right_value" binding:"required"`
	Note       string `json:"note"`
	RunBegin   string `json:"run_begin"`
	RunEnd     string `json:"run_end"`
	TplId      int64  `json:"tpl_id" binding:"required"`
}

func (this APICreateStrategyInput) CheckFormat() (err error) {
	validOp := regexp.MustCompile(`^(>|=|<|!)(=)?$`)
	validRightValue := regexp.MustCompile(`^\d+$`)
	validTime := regexp.MustCompile(`^\d{2}:\d{2}$`)
	switch {
	case !validOp.MatchString(this.Op):
		err = errors.New("op's formating is not vaild")
	case !validRightValue.MatchString(this.RightValue):
		err = errors.New("right_value's formating is not vaild")
	case !validTime.MatchString(this.RunBegin) && this.RunBegin != "":
		err = errors.New("run_begin's formating is not vaild, please refer ex. 00:00")
	case !validTime.MatchString(this.RunEnd) && this.RunEnd != "":
		err = errors.New("run_end's formating is not vaild, please refer ex. 24:00")
	}
	return
}

func CreateStrategy(c *gin.Context) {
	var inputs APICreateStrategyInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.CheckFormat(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	strategy := f.Strategy{
		Metric:     inputs.Metric,
		Tags:       inputs.Tags,
		MaxStep:    inputs.MaxStep,
		Priority:   inputs.Priority,
		Func:       inputs.Func,
		Op:         inputs.Op,
		RightValue: inputs.RightValue,
		Note:       inputs.Note,
		RunBegin:   inputs.RunBegin,
		RunEnd:     inputs.RunEnd,
		TplId:      inputs.TplId,
	}
	dt := db.Falcon.Save(&strategy)
	if dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, "stragtegy created")
	return
}

func GetStrategy(c *gin.Context) {
	sidtmp := c.Params.ByName("sid")
	if sidtmp == "" {
		h.JSONR(c, badstatus, "sid is missing")
		return
	}
	sid, err := strconv.Atoi(sidtmp)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	strategy := f.Strategy{ID: int64(sid)}
	if dt := db.Falcon.Find(&strategy); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, strategy)
	return
}

type APIUpdateStrategyInput struct {
	ID         int64  `json:"id" binding:"required"`
	Metric     string `json:"metric" binding:"required"`
	Tags       string `json:"tags"`
	MaxStep    int    `json:"max_step" binding:"required"`
	Priority   int    `json:"priority" binding:"exists"`
	Func       string `json:"func" binding:"required"`
	Op         string `json:"op" binding:"required"`
	RightValue string `json:"right_value" binding:"required"`
	Note       string `json:"note"`
	RunBegin   string `json:"run_begin"`
	RunEnd     string `json:"run_end"`
}

func (this APIUpdateStrategyInput) CheckFormat() (err error) {
	validOp := regexp.MustCompile(`^(>|=|<|!)(=)?$`)
	validRightValue := regexp.MustCompile(`^\d+$`)
	validTime := regexp.MustCompile(`^\d{2}:\d{2}$`)
	switch {
	case !validOp.MatchString(this.Op):
		err = errors.New("op's formating is not vaild")
	case !validRightValue.MatchString(this.RightValue):
		err = errors.New("right_value's formating is not vaild")
	case !validTime.MatchString(this.RunBegin) && this.RunBegin != "":
		err = errors.New("run_begin's formating is not vaild, please refer ex. 00:00")
	case !validTime.MatchString(this.RunEnd) && this.RunEnd != "":
		err = errors.New("run_end's formating is not vaild, please refer ex. 24:00")
	}
	return
}

func UpdateStrategy(c *gin.Context) {
	var inputs APIUpdateStrategyInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.CheckFormat(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	strategy := f.Strategy{
		ID: inputs.ID,
	}
	if dt := db.Falcon.Find(&strategy); dt.Error != nil {
		h.JSONR(c, expecstatus, fmt.Sprintf("find strategy got error:%v", dt.Error))
		return
	}
	ustrategy := map[string]interface{}{
		"Metric":     inputs.Metric,
		"Tags":       inputs.Tags,
		"MaxStep":    inputs.MaxStep,
		"Priority":   inputs.Priority,
		"Func":       inputs.Func,
		"Op":         inputs.Op,
		"RightValue": inputs.RightValue,
		"Note":       inputs.Note,
		"RunBegin":   inputs.RunBegin,
		"RunEnd":     inputs.RunEnd}
	if dt := db.Falcon.Model(&strategy).Where("id = ?", strategy.ID).Update(ustrategy); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("stragtegy:%d has been updated", strategy.ID))
	return
}

func DeleteStrategy(c *gin.Context) {
	sidtmp := c.Params.ByName("sid")
	if sidtmp == "" {
		h.JSONR(c, badstatus, "sid is missing")
		return
	}
	sid, err := strconv.Atoi(sidtmp)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	strategy := f.Strategy{ID: int64(sid)}
	if dt := db.Falcon.Delete(&strategy); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("strategy:%d has been deleted", sid))
	return
}

func MetricQuery(c *gin.Context) {
	filePath := viper.GetString("metric_list_file")
	if filePath == "" {
		filePath = "./data/metric"
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	metrics := strings.Split(string(data), "\n")
	h.JSONR(c, metrics)
	return
}
