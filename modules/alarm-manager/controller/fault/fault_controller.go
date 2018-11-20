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

package fault

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/controller"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/model/fault"
)

// Create reads params from request and creates fault.
// If successful, fault created will be returned.
func Create(c *gin.Context) {
	var createInfo fault.CreateInfo
	if err := c.BindJSON(&createInfo); err != nil {
		log.Errorf("read params fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	user, err := GetUser(c)
	if err != nil {
		log.Errorf("get user info from session fails: %v", err)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	createInfo.Creator = user

	fault, err := fault.Store.Create(createInfo)
	if err != nil {
		log.Errorf("create fault fails: %v", err)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "fault create succeed",
		Data:    fault,
	})
	return
}

type WebSession struct {
	Name string
	Sig  string
}

func GetUser(c *gin.Context) (string, error) {
	apiToken := c.Request.Header.Get("Apitoken")

	var websession WebSession
	if err := json.Unmarshal([]byte(apiToken), &websession); err != nil {
		return "", err
	}

	return websession.Name, nil
}

// Get gets fault by fault id.
// If successful, fault info in detail will be returned.
func Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	faultinfo, err := fault.Store.Get(uint(id))
	if err != nil {
		log.Errorf("get fault by id fails: %v,faultid: %v", err, id)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Get fault successfully",
		Data:    faultinfo,
	})
	return
}

// GetEvent gets event in fault by fault id.
// If successful, event will be returned.
func GetEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		log.Errorf("convert event id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	eventinfo, err := fault.Store.GetEvent(uint(id))
	if err != nil {
		log.Errorf("get event fails: %v,faultid: %v", err, id)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Get event of fault successfully",
		Data:    eventinfo,
	})
	return
}

// AddEvent adds event into fault by event id and fault id.
// If successful, the newest fault info in detail will be returned.
func AddEvent(c *gin.Context) {
	id, eventids := Params(c, "id", "eventids", ",")

	if is := IsEmpty(id); is {
		msg := "id is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}
	faultid, err := strconv.Atoi(id)
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if is := IsEmpty(eventids); is {
		msg := "eventids is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}
	eids := make([]uint, 0, len(eventids))
	for _, v := range eventids {
		eid, err := strconv.Atoi(v)
		if err != nil {
			log.Errorf("convert event id fails: %v", err)
			c.JSON(http.StatusBadRequest, controller.Resp{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}
		eids = append(eids, uint(eid))
	}

	faultinfo, err := fault.Store.AddEvent(uint(faultid), eids)
	if err != nil {
		log.Errorf("add event fails: %v,faultid: %v", err, faultid)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Add event into fault successfully",
		Data:    faultinfo,
	})
	return
}

// DeleteEvent deletes event from fault by event id and fault id.
// If successful, the newest fault info in detail will be returned.
func DeleteEvent(c *gin.Context) {
	id, eventids := Params(c, "id", "eventids", ",")

	if is := IsEmpty(id); is {
		msg := "id is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}
	faultid, err := strconv.Atoi(id)
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if is := IsEmpty(eventids); is {
		msg := "eventids is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}
	eids := make([]uint, 0, len(eventids))
	for _, v := range eventids {
		eid, err := strconv.Atoi(v)
		if err != nil {
			log.Errorf("convert event id fails: %v", err)
			c.JSON(http.StatusBadRequest, controller.Resp{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}
		eids = append(eids, uint(eid))
	}

	faultinfo, err := fault.Store.DeleteEvent(uint(faultid), eids)
	if err != nil {
		log.Errorf("delete event fails: %v,faultid: %v", err, faultid)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Delete event from fault successfully",
		Data:    faultinfo,
	})
	return
}

// GetComment gets comment of fault by fault id.
// If successful, comment info will be returned.
func GetComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	commentInfo, err := fault.Store.GetComment(uint(id))
	if err != nil {
		log.Errorf("get comment fails: %v,faultid: %v", err, id)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Get comment of fault successfully",
		Data:    commentInfo,
	})
	return
}

// AddComment adds one comment into fault by fault id and comment.
// If successful, the newest fault info in detail will be returned.
func AddComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	comment := c.Query("comment")
	if comment == "" {
		msg := "comment is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}

	user, err := GetUser(c)
	if err != nil {
		log.Errorf("get user info from session fails: %v", err)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	faultInfo, err := fault.Store.AddComment(uint(id), user, comment)
	if err != nil {
		log.Errorf("add comment fails: %v,faultid: %v", err, id)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Add comment successfully",
		Data:    faultInfo,
	})
	return
}

// DeleteComment deletes comment from fault by fault id and comment and creator.
// If successful, the newest fault info in detail will be returned.
func DeleteComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	comment := c.Query("comment")
	if comment == "" {
		msg := "comment is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}

	creator := c.Query("creator")
	if creator == "" {
		msg := "creator is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}

	faultInfo, err := fault.Store.DeleteComment(uint(id), creator, comment)
	if err != nil {
		log.Errorf("delete comment fails: %v,faultid: %v", err, id)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Delete comment successfully",
		Data:    faultInfo,
	})
	return
}

// GetTag gets tag info of fault by fault id.
// If successful, tags of fault will be returned.
func GetTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	taginfo, err := fault.Store.GetTag(uint(id))
	if err != nil {
		log.Errorf("get tag fails: %v,faultid: %v", err, id)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Get tag of fault successfully",
		Data:    taginfo,
	})
	return
}

// AddTag adds tag into fault by fault id and tags.
// If successful, the newest fault info in detail will be returned.
func AddTag(c *gin.Context) {
	id, tag := Params(c, "id", "tags", ",")

	if is := IsEmpty(id); is {
		msg := "id is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}
	faultid, err := strconv.Atoi(id)
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	tags := make([]string, 0)
	for _, v := range tag {
		if v != "" {
			tags = append(tags, v)
		}
	}
	if is := IsEmpty(tags); is {
		msg := "tags is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}

	faultinfo, err := fault.Store.AddTag(uint(faultid), tags)
	if err != nil {
		log.Errorf("add tags fails: %v,faultid: %v", err, faultid)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Add tag into fault successfully",
		Data:    faultinfo,
	})
	return
}

// DeleteTag deletes tags from fault by fault id and tags.
// If successful, the newest fault in detail will be returned.
func DeleteTag(c *gin.Context) {
	id, tag := Params(c, "id", "tags", ",")

	if is := IsEmpty(id); is {
		msg := "id is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}
	faultid, err := strconv.Atoi(id)
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	tags := make([]string, 0)
	for _, v := range tag {
		if v != "" {
			tags = append(tags, v)
		}
	}
	if is := IsEmpty(tags); is {
		msg := "tags is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}

	faultinfo, err := fault.Store.DeleteTag(uint(faultid), tags)
	if err != nil {
		log.Errorf("delete tag fails: %v,faultid: %v", err, faultid)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Delete tag from fault successfully",
		Data:    faultinfo,
	})
	return
}

// Params gets param in url path and query string.
func Params(c *gin.Context, pathkey, querykey, separator string) (string, []string) {
	pathValue := c.Params.ByName(pathkey)

	queryValue := c.Query(querykey)
	if queryValue == "" {
		return pathValue, []string{}
	}
	queryValues := strings.Split(queryValue, separator)
	return pathValue, queryValues
}

// UpdateOwner updates owner of fault by fault id.
// If successful, the newest fault in detail will be returned.
func UpdateOwner(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	owner := c.Query("owner")
	if owner == "" {
		msg := "owner is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}

	fault, err := fault.Store.UpdateOwner(uint(id), owner)
	if err != nil {
		log.Errorf("update owner fails: %v,faultid: %v", err, id)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Update owner successfully",
		Data:    fault,
	})
	return
}

// UpdateState updates state of fault by fault id.
// If successful, the newest fault in detail will be returned.
func UpdateState(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	state := c.Query("state")
	if state == "" {
		msg := "state is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}

	fault, err := fault.Store.UpdateState(uint(id), state)
	if err != nil {
		log.Errorf("update state fails: %v,faultid: %v", err, id)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Update state successfully",
		Data:    fault,
	})
	return
}

// UpdateFollwer updates follower of fault by fault id.
// If successful, the newest fault in detail will be returned.
func UpdateFollower(c *gin.Context) {
	action := c.Query("action")
	if action == "" {
		msg := "action is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}

	id, follower := Params(c, "id", "follower", ",")
	if is := IsEmpty(id); is {
		msg := "id is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}
	faultid, err := strconv.Atoi(id)
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	followers := make([]string, 0)
	for _, v := range follower {
		if v != "" {
			followers = append(followers, v)
		}
	}
	if is := IsEmpty(followers); is {
		msg := "follower is empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}

	fault, err := fault.Store.UpdateFollower(uint(faultid), followers, action)
	if err != nil {
		log.Errorf("update follower fails: %v,faultid: %v", err, faultid)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Update follower successfully",
		Data:    fault,
	})
	return
}

// List gets fault by filter in request.
// If successful, matched fault and count will be returned.
// Count is the number of fault which meets the filter rather than
// the number of fault in response body.
func List(c *gin.Context) {
	filter, err := FilterInfo(c)
	if err != nil {
		log.Errorf("filter is invalid: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	faults, count, err := fault.Store.List(filter)
	if err != nil {
		log.Errorf("list fault fails: %v", err)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "List fault successfully",
		Data: struct {
			Faults []fault.FaultInfo
			Count  uint
		}{
			faults,
			count,
		},
	})
	return
}

// GetTimeLine gets timeline of fault by fault id.
// If successful, timeline of fault will be returned.
func GetTimeLine(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	timeLine, err := fault.Store.GetTimeLine(uint(id))
	if err != nil {
		log.Errorf("get fault timeline fails: %v,faultid: %v", err, id)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Get fault timeline successfully",
		Data:    timeLine,
	})
	return
}

// UpdateBasic updates title and note of fault by fault id.
// If successful, the newest fault in detail is returned.
func UpdateBasic(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		log.Errorf("convert fault id fails: %v", err)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	title := c.Query("title")
	note := c.Query("note")
	if title == "" && note == "" {
		msg := "title and note can not be all empty"
		log.Errorf(msg)
		c.JSON(http.StatusBadRequest, controller.Resp{
			Code:    http.StatusBadRequest,
			Message: msg,
		})
		return
	}

	faultinfo, err := fault.Store.UpdateBasic(uint(id), title, note)
	if err != nil {
		log.Errorf("update fault basic fails: %v,faultid: %v", err, id)
		c.JSON(http.StatusInternalServerError, controller.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, controller.Resp{
		Code:    http.StatusOK,
		Message: "Update fault basic successfully",
		Data:    faultinfo,
	})
	return
}

func FilterInfo(c *gin.Context) (fault.Filter, error) {
	defaultEnd := time.Now().Unix()
	defaultStart := defaultEnd - 3600*24

	startParam := c.DefaultQuery("start", strconv.FormatInt(defaultStart, 10))
	start, err := strconv.Atoi(startParam)
	if err != nil {
		return fault.Filter{}, err
	}

	endParam := c.DefaultQuery("end", strconv.FormatInt(defaultEnd, 10))
	end, err := strconv.Atoi(endParam)
	if err != nil {
		return fault.Filter{}, err
	}

	creator := c.Query("creator")
	owner := c.Query("owner")
	state := c.Query("state")
	title := c.Query("title")
	follower := c.Query("follower")
	tag := c.Query("tag")

	defaultLimit := "10"
	limitParam := c.DefaultQuery("limit", defaultLimit)
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return fault.Filter{}, err
	}
	if limit <= 0 || limit >= 50 {
		limit = 10
	}

	defaultOffset := "0"
	offsetParam := c.DefaultQuery("offset", defaultOffset)
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		return fault.Filter{}, err
	}

	return fault.Filter{
		Start:    uint(start),
		End:      uint(end),
		Creator:  creator,
		Owner:    owner,
		State:    state,
		Title:    title,
		Follower: follower,
		Tag:      tag,
		Limit:    uint(limit),
		Offset:   uint(offset),
	}, nil
}
