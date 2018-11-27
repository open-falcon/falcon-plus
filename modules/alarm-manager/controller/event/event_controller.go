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

package event

import (
	"net/http"

	"github.com/gin-gonic/gin"
	coommonModel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/controller"
	mevent "github.com/open-falcon/falcon-plus/modules/alarm-manager/model/event"
)

func RecvAlarmEvent(c *gin.Context) {
	var eve *coommonModel.Event
	if err := c.BindJSON(&eve); err != nil {
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if err := mevent.Store.InsertAlarmEvent(eve); err != nil {
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, controller.Resp{
		Code:    http.StatusOK,
		Message: "alarm event recv success",
		Data:    "",
	})
	return
}

// 单独获取event信息接口
func GetEvents(c *gin.Context) {
	var inputs mevent.EventApiInputs
	if err := c.Bind(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if err := checkEventInputs(inputs); err != nil {
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	// Default time interval: 1h, limit 100
	inputs = TimeQueryLimitFilters(inputs)
	event, err := mevent.Store.GetEventsInfo(inputs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "",
		Data:    event,
	})
	return
}

func GetEventsFaults(c *gin.Context) {
	var inputs mevent.EventApiInputs
	if err := c.Bind(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if err := checkEventInputs(inputs); err != nil {
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	// Default time interval: 1h, limit 100
	inputs = TimeQueryLimitFilters(inputs)
	eventfault, err := mevent.Store.GetEventsFaultsInfo(inputs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "",
		Data:    eventfault,
	})
	return
}

// 如果需要返回包含故障的数量，参数传入have_fault=true
func GetEventCount(c *gin.Context) {
	var inputs mevent.EventApiInputs
	if err := c.Bind(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if err := checkEventInputs(inputs); err != nil {
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	// Default time interval: 1h
	inputs = TimeParamFilters(inputs)
	count, err := mevent.Store.GetEventCount(inputs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "",
		Data:    count,
	})
	return
}
