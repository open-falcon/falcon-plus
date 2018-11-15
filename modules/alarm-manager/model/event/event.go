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
	"fmt"
	"strconv"
	"strings"
	"time"
	//log "github.com/Sirupsen/logrus"
)

const (
	eventsJoin        = "INNER JOIN events as b on a.event_case_id = b.`event_case_id`"
	eventReceiverJoin = "INNER JOIN event_receiver as c on b.`id` = c.`event_id` "
	faultJoin         = "INNER JOIN fault_event as d on b.`id` = d.`event_id`"
)

func GetSelectConFilter(in EventApiInputs) string {
	var sql string
	if in.Status != "" && in.NowStatus {
		//此处a.`status`为当前event状态,即event_case表中实时更新的event_case状态
		sql = "b.`id`, a.`event_case_id`, a.`endpoint`, a.`counter`, a.`func`, a.`cond`, a.`note`, " +
			"a.`max_step`, a.`current_step`, a.`priority`, a.`status`, a.`event_ts`, " +
			"a.`tpl_creator`, a.`expression_id`, a.`strategy_id`, a.`template_id`"
	}

	sql = "b.`id`, a.`event_case_id`, a.`endpoint`, a.`counter`, a.`func`, b.`cond`, a.`note`, " +
		"a.`max_step`, b.`current_step`, a.`priority`, b.`status`, b.`event_ts`, " +
		"a.`tpl_creator`, a.`expression_id`, a.`strategy_id`, a.`template_id`"
	return sql
}

func GetQueryConFilter(in EventApiInputs) string {
	cond := []string{}
	if in.UserName != "" {
		cond = append(cond, fmt.Sprintf("c.`user` = '%s'", in.UserName))
	}
	if in.Uic != nil && len(in.Uic) != 0 {
		for i, u := range in.Uic {
			in.Uic[i] = fmt.Sprintf("'%s'", u)
		}
		cond = append(cond, fmt.Sprintf("c.`uic` in (%s)", strings.Join(in.Uic, ", ")))
	}

	priority := []string{}
	if in.Priority != nil && len(in.Priority) != 0 {
		for _, po := range in.Priority {
			priority = append(priority, fmt.Sprintf("'%s'", strconv.Itoa(po)))
		}
		cond = append(cond, fmt.Sprintf("a.`priority` in (%s)", strings.Join(priority, ", ")))
	}

	if in.StartTime != 0 {
		cond = append(cond, fmt.Sprintf("c.`event_ts` >= '%s'", time.Unix(in.StartTime, 0).Format(timeLayout)))
	}
	if in.EndTime != 0 {
		cond = append(cond, fmt.Sprintf("c.`event_ts` <= '%s'", time.Unix(in.EndTime, 0).Format(timeLayout)))
	}
	if in.Status != "" {
		if in.NowStatus {
			//此处status为当前event状态,即event_case表中实时更新的event_case状态
			cond = append(cond, fmt.Sprintf("a.`status` = '%s'", in.Status))
		} else {
			//此处status为历史event状态,即events中历史状态
			cond = append(cond, fmt.Sprintf("b.`status` = '%s'", in.Status))
		}
	}

	if in.Endpoint != "" {
		cond = append(cond, fmt.Sprintf("a.`endpoint` regexp '%s'", in.Endpoint))
	}
	if in.Counter != "" {
		cond = append(cond, fmt.Sprintf("a.`counter` regexp '%s'", in.Counter))
	}

	wherefilter := strings.Join(cond, " AND ")
	return wherefilter
}

func (store *EventStore) GetEventsInfo(in EventApiInputs) ([]EventInfo, error) {
	var err error
	selectsql := GetSelectConFilter(in)
	wheresql := GetQueryConFilter(in)

	event := []EventInfo{}
	if in.HaveFault {
		err = store.AMDB.Table("event_cases as a").Select(selectsql).Joins(eventsJoin).
			Joins(eventReceiverJoin).Joins(faultJoin).Where(wheresql).Group("b.`id`").Order("c.`event_ts` desc").
			Limit(in.Limit).Offset(in.Offset).Find(&event).Error
	} else {
		err = store.AMDB.Table("event_cases as a").Select(selectsql).Joins(eventsJoin).
			Joins(eventReceiverJoin).Where(wheresql).Group("b.`id`").Order("c.`event_ts` desc").
			Limit(in.Limit).Offset(in.Offset).Find(&event).Error
	}
	if err != nil {
		return []EventInfo{}, err
	}
	return event, nil
}

func (store *EventStore) GetFaultsInfo(id uint) ([]FaultInfos, error) {
	fault := []FaultInfos{}
	selectsql := "a.`title`, b.`fault_id`, a.`updated_at`"
	faultEventJoin := "INNER JOIN `fault_event` as b on a.`id` = b.`fault_id`"

	err := store.AMDB.Table("fault as a").Select(selectsql).Joins(faultEventJoin).
		Where("b.`event_id` = ? and b.`deleted_at` IS NULL", id).Order("a.`updated_at` desc").Find(&fault).Error

	if err != nil {
		return []FaultInfos{}, err
	}
	return fault, nil
}

func (store *EventStore) GetEventsFaultsInfo(in EventApiInputs) ([]EventFaultInfos, error) {
	event, err := store.GetEventsInfo(in)
	if err != nil {
		return []EventFaultInfos{}, err
	}
	eventfault := []EventFaultInfos{}
	for _, e := range event {
		fault, err := store.GetFaultsInfo(e.ID)
		if err != nil {
			return []EventFaultInfos{}, err
		}
		eventfault = append(eventfault, EventFaultInfos{Event: e, Faults: fault})
	}
	return eventfault, nil
}

func (store *EventStore) GetEventCount(in EventApiInputs) (EventCount, error) {
	var err error
	wheresql := GetQueryConFilter(in)

	event := []EventInfo{}
	if in.HaveFault {
		err = store.AMDB.Table("event_cases as a").Select("b.`id`").Joins(eventsJoin).Joins(eventReceiverJoin).
			Joins(faultJoin).Where(wheresql).Group("b.`id`").Find(&event).Error
	} else {
		err = store.AMDB.Table("event_cases as a").Select("b.`id`").Joins(eventsJoin).Joins(eventReceiverJoin).
			Where(wheresql).Group("b.`id`").Find(&event).Error
	}

	if err != nil {
		return EventCount{Count: 0}, err
	}
	return EventCount{Count: len(event)}, nil
}

func (store *EventStore) GetEventByID(id uint) (EventInfo, error) {
	selectsql := "b.`id`, a.`event_case_id`, a.`endpoint`, a.`counter`, " +
		"a.`func`, b.`cond`, a.`note`, a.`max_step`, b.`current_step`, a.`priority`, b.`status`, " +
		"b.`event_ts`, a.`tpl_creator`, a.`expression_id`, a.`strategy_id`, a.`template_id`"

	event := EventInfo{}
	err := store.AMDB.Table("event_cases as a").Select(selectsql).Joins(eventsJoin).
		Where("b.`id` = ?", id).Find(&event).Error

	if err != nil {
		return EventInfo{}, err
	}
	return event, nil
}
