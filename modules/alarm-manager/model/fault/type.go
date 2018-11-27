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
	"time"

	"github.com/jinzhu/gorm"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/model/event"
)

type FaultInfo struct {
	Id        uint
	CreatedAt time.Time
	Title     string
	Note      string
	Creator   string
	Owner     string
	State     string
	Tags      []string
	Events    []event.EventInfo
	Followers []string
	Comments  []CommentInfo
}

type CommentInfo struct {
	CreatedAt time.Time
	Creator   string
	Comment   string
}

type CreateInfo struct {
	Title   string   `json:"title" form:"title" binding:"required"`
	Note    string   `json:"note" form:"note"`
	Creator string   `json:"creator" form:"creator"`
	Owner   string   `json:"owner" form:"owner" binding:"required"`
	Tags    []string `json:"tags" form:"tags"`
}

type TimeLine struct {
	FirstEventTime string
	LastEventTime  string
	FaultCreatedAt time.Time
	FaultClosedAt  time.Time
}

type Filter struct {
	Start    uint
	End      uint
	Creator  string
	Owner    string
	State    string
	Title    string
	Follower string
	Tag      string
	Limit    uint
	Offset   uint
}

type Fault struct {
	gorm.Model
	Title   string
	Note    string
	Creator string
	Owner   string
	State   string
}

func (Fault) TableName() string {
	return "fault"
}

type FaultComment struct {
	gorm.Model
	FaultId uint
	Creator string
	Comment string `gorm:"size:999"`
}

func (FaultComment) TableName() string {
	return "fault_comment"
}

type FaultEvent struct {
	gorm.Model
	FaultId uint `gorm:"unique_index:idx_faultid_eventid"`
	EventId uint `gorm:"unique_index:idx_faultid_eventid"`
}

func (FaultEvent) TableName() string {
	return "fault_event"
}

type FaultFollower struct {
	gorm.Model
	FaultId  uint   `gorm:"unique_index:idx_faultid_follower"`
	Follower string `gorm:"unique_index:idx_faultid_follower"`
}

func (FaultFollower) TableName() string {
	return "fault_follower"
}

type FaultTag struct {
	gorm.Model
	FaultId uint   `gorm:"unique_index:idx_faultid_tag"`
	Tag     string `gorm:"unique_index:idx_faultid_tag"`
}

func (FaultTag) TableName() string {
	return "fault_tag"
}

type StateChangeLog struct {
	gorm.Model
	ReferTable string
	ReferId    uint
	ReferField string
	From       string
	To         string
}

func (StateChangeLog) TableName() string {
	return "state_change_log"
}
