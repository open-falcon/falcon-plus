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

package expression

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
)

func GetExpressionList(c *gin.Context) {
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
	expressions := []f.Expression{}
	if limit != -1 && page != -1 {
		dt = db.Falcon.Raw(fmt.Sprintf("SELECT * from expression limit %d,%d", page, limit)).Scan(&expressions)
	} else {
		dt = db.Falcon.Find(&expressions)
	}
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, expressions)
	return
}

func GetExpression(c *gin.Context) {
	eidtmp := c.Params.ByName("eid")
	if eidtmp == "" {
		h.JSONR(c, badstatus, "eid is missing")
		return
	}
	eid, err := strconv.Atoi(eidtmp)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	expression := f.Expression{ID: int64(eid)}
	if dt := db.Falcon.Find(&expression); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	action := f.Action{ID: expression.ActionId}
	if dt := db.Falcon.Find(&action); dt.Error != nil {
		h.JSONR(c, badstatus, fmt.Sprintf("find action got error: %v", dt.Error.Error()))
		return
	}
	h.JSONR(c, map[string]interface{}{
		"expression": expression,
		"action":     action,
	})
	return
}

type APICreateExrpessionInput struct {
	Expression string    `json:"expression" binding:"required"`
	Func       string    `json:"func" binding:"required"`
	Op         string    `json:"op" binding:"required"`
	RightValue string    `json:"right_value" binding:"required"`
	MaxStep    int       `json:"max_step" binding:"required"`
	Priority   int       `json:"priority" binding:"required"`
	Note       string    `json:"note" binding:"exists"`
	Pause      int       `json:"pause" binding:"exists"`
	Action     ActionTmp `json:"action" binding:"required"`
	// ActionId   string `json:"action_id" binding:"exists"`
}

type ActionTmp struct {
	UIC                []string `json:"uic" binding:"required"`
	URL                string   `json:"url" binding:"exists"`
	Callback           int      `json:"callback" binding:"exists"`
	BeforeCallbackSMS  int      `json:"before_callback_sms" binding:"exists"`
	AfterCallbackSMS   int      `json:"after_callback_sms" binding:"exists"`
	BeforeCallbackMail int      `json:"before_callback_mail" binding:"exists"`
	AfterCallbackMail  int      `json:"after_callback_mail" binding:"exists"`
}

func (this APICreateExrpessionInput) CheckFormat() (err error) {
	validOp := regexp.MustCompile(`^(>|=|<|!)(=)?$`)
	validRightValue := regexp.MustCompile(`^\-?\d+(\.\d+)?$`)
	switch {
	case !validOp.MatchString(this.Op):
		err = errors.New("op's formating is not vaild")
	case !validRightValue.MatchString(this.RightValue):
		err = errors.New("right_value's formating is not vaild")
	}
	return
}

func CreateExrpession(c *gin.Context) {
	var inputs APICreateExrpessionInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.CheckFormat(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, _ := h.GetUser(c)
	tx := db.Falcon.Begin()
	action := f.Action{
		UIC:                strings.Join(inputs.Action.UIC, ","),
		URL:                inputs.Action.URL,
		Callback:           inputs.Action.Callback,
		BeforeCallbackSMS:  inputs.Action.BeforeCallbackSMS,
		BeforeCallbackMail: inputs.Action.BeforeCallbackMail,
		AfterCallbackSMS:   inputs.Action.AfterCallbackSMS,
		AfterCallbackMail:  inputs.Action.AfterCallbackMail,
	}
	if dt := tx.Save(&action); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		tx.Rollback()
		return
	}
	expression := f.Expression{
		Expression: inputs.Expression,
		Func:       inputs.Func,
		Op:         inputs.Op,
		RightValue: inputs.RightValue,
		MaxStep:    inputs.MaxStep,
		Priority:   inputs.Priority,
		Note:       inputs.Note,
		Pause:      inputs.Pause,
		CreateUser: user.Name,
		ActionId:   action.ID,
	}
	dt := tx.Save(&expression)
	if dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, "expression created")
	return
}

type APIUpdateExrpessionInput struct {
	ID         int64      `json:"id"  binding:"required"`
	Expression string     `json:"expression" binding:"required"`
	Func       string     `json:"func" binding:"required"`
	Op         string     `json:"op" binding:"required"`
	RightValue string     `json:"right_value" binding:"required"`
	MaxStep    int        `json:"max_step" binding:"required"`
	Priority   int        `json:"priority" binding:"required"`
	Note       string     `json:"note" binding:"exists"`
	Pause      int        `json:"pause" binding:"exists"`
	Action     ActionTmpU `json:"action" binding:"required"`
}

type ActionTmpU struct {
	UIC                []string `json:"uic" binding:"required"`
	URL                string   `json:"url" binding:"exists"`
	Callback           int      `json:"callback" binding:"exists"`
	BeforeCallbackSMS  int      `json:"before_callback_sms" binding:"exists"`
	AfterCallbackSMS   int      `json:"after_callback_sms" binding:"exists"`
	BeforeCallbackMail int      `json:"before_callback_mail" binding:"exists"`
	AfterCallbackMail  int      `json:"after_callback_mail" binding:"exists"`
}

func (this APIUpdateExrpessionInput) CheckFormat() (err error) {
	validOp := regexp.MustCompile(`^(>|=|<|!)(=)?$`)
	validRightValue := regexp.MustCompile(`^\d+$`)
	switch {
	case !validOp.MatchString(this.Op):
		err = errors.New("op's formating is not vaild")
	case !validRightValue.MatchString(this.RightValue):
		err = errors.New("right_value's formating is not vaild")
	}
	return
}

func UpdateExrpession(c *gin.Context) {
	var inputs APIUpdateExrpessionInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if err := inputs.CheckFormat(); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	tx := db.Falcon.Begin()
	user, _ := h.GetUser(c)
	expression := f.Expression{ID: inputs.ID}
	if dt := tx.Find(&expression); dt.Error != nil {
		h.JSONR(c, expecstatus, fmt.Sprintf(
			"find expression got error:%v", dt.Error.Error()))
		tx.Rollback()
		return
	}
	if !user.IsAdmin() {
		if expression.CreateUser != user.Name {
			h.JSONR(c, badstatus, "You don't have permission!")
			tx.Rollback()
			return
		}
	}
	uexpression := map[string]interface{}{
		"ID":         expression.ID,
		"Expression": inputs.Expression,
		"Func":       inputs.Func,
		"Op":         inputs.Op,
		"RightValue": inputs.RightValue,
		"MaxStep":    inputs.MaxStep,
		"Priority":   inputs.Priority,
		"Note":       inputs.Note,
		"Pause":      inputs.Pause,
	}
	dt := tx.Model(&expression).Where("id = ?", expression.ID).Update(uexpression).Find(&expression)
	if dt.Error != nil {
		h.JSONR(c, expecstatus, fmt.Sprintf(
			"update expression got error: %v", dt.Error))
		tx.Rollback()
		return
	}
	actionTmp := f.Action{ID: expression.ActionId}
	uaction := map[string]interface{}{
		"ID":                 actionTmp.ID,
		"UIC":                strings.Join(inputs.Action.UIC, ","),
		"URL":                inputs.Action.URL,
		"Callback":           inputs.Action.Callback,
		"BeforeCallbackSMS":  inputs.Action.BeforeCallbackSMS,
		"BeforeCallbackMail": inputs.Action.BeforeCallbackMail,
		"AfterCallbackSMS":   inputs.Action.AfterCallbackSMS,
		"AfterCallbackMail":  inputs.Action.AfterCallbackMail,
	}
	if dt = tx.Find(&actionTmp, expression.ActionId); dt.Error != nil {
		h.JSONR(c, expecstatus, fmt.Sprintf(
			"find action got error: %v", dt.Error))
		tx.Rollback()
		return
	}
	dt = tx.Model(&actionTmp).Where("id = ?", actionTmp.ID).Update(uaction)
	if dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("expression:%v has been updated", inputs.ID))
	return
}

func DeleteExpression(c *gin.Context) {
	eidtmp := c.Params.ByName("eid")
	if eidtmp == "" {
		h.JSONR(c, badstatus, "eid is missing")
		return
	}
	eid, err := strconv.Atoi(eidtmp)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	tx := db.Falcon.Begin()
	user, _ := h.GetUser(c)
	expression := f.Expression{ID: int64(eid)}
	if !user.IsAdmin() {
		tx.Find(&expression)
		if expression.CreateUser != user.Name {
			h.JSONR(c, badstatus, "You don't have permission!")
			tx.Rollback()
			return
		}
	}
	dt := tx.Table("action").Where("id = ?", expression.ActionId).Delete(&f.Action{})
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		tx.Rollback()
		return
	}
	if dt := tx.Delete(&expression); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("expression:%d has been deleted", eid))
	return
}
