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

package mockcfg

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
)

func GetNoDataList(c *gin.Context) {
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
	var dt *gorm.DB
	mockcfgs := []f.Mockcfg{}
	if limit != -1 && page != -1 {
		dt = db.Falcon.Raw(fmt.Sprintf("SELECT * from mockcfg limit %d,%d", page, limit)).Scan(&mockcfgs)
	} else {
		dt = db.Falcon.Find(&mockcfgs)
	}
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, mockcfgs)
	return
}

func GetNoData(c *gin.Context) {
	nidtmp := c.Params.ByName("nid")
	if nidtmp == "" {
		h.JSONR(c, badstatus, "nid is missing")
		return
	}
	nid, err := strconv.Atoi(nidtmp)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	mockcfg := f.Mockcfg{ID: int64(nid)}
	if dt := db.Falcon.Find(&mockcfg); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, mockcfg)
	return
}

func CreateNoData(c *gin.Context) {
	var inputs APICreateNoDataInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.CheckFormat(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, _ := h.GetUser(c)
	mockcfg := f.Mockcfg{
		Name:    inputs.Name,
		Obj:     inputs.Obj,
		ObjType: inputs.ObjType,
		Metric:  inputs.Metric,
		Tags:    inputs.Tags,
		DsType:  inputs.DsType,
		Step:    inputs.Step,
		Mock:    inputs.Mock,
		Creator: user.Name,
	}
	if dt := db.Falcon.Save(&mockcfg); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, mockcfg)
	return
}

func UpdateNoData(c *gin.Context) {
	var inputs APIUpdateNoDataInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.CheckFormat(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	mockcfg := &f.Mockcfg{ID: inputs.ID}
	umockcfg := map[string]interface{}{
		"Obj":     inputs.Obj,
		"ObjType": inputs.ObjType,
		"Metric":  inputs.Metric,
		"Tags":    inputs.Tags,
		"DsType":  inputs.DsType,
		"Step":    inputs.Step,
		"Mock":    inputs.Mock,
	}
	if dt := db.Falcon.Model(&mockcfg).Where("id = ?", inputs.ID).Update(umockcfg).Find(&mockcfg); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, mockcfg)
	return
}

func DeleteNoData(c *gin.Context) {
	nidtmp := c.Params.ByName("nid")
	if nidtmp == "" {
		h.JSONR(c, badstatus, "nid is missing")
		return
	}
	nid, err := strconv.Atoi(nidtmp)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	mockcfg := f.Mockcfg{ID: int64(nid)}
	if dt := db.Falcon.Delete(&mockcfg); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("mockcfg:%d is deleted", nid))
	return
}
