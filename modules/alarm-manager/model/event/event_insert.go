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
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"

	coommonModel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/api"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/config"
)

const timeLayout = "2006-01-02 15:04:05"

var Store *EventStore

func Init() {
	Store = NewEventStore(config.Con().AM)
}

type EventStore struct {
	AMDB   *gorm.DB
	Locker *sync.RWMutex
}

func NewEventStore(db *gorm.DB) *EventStore {
	return &EventStore{AMDB: db, Locker: &sync.RWMutex{}}
}

func (store *EventStore) insertEventReceiver(eventid uint, event *coommonModel.Event) error {
	action := api.GetAction(event)
	if action == nil {
		return fmt.Errorf("action is nil")
	}

	teams := strings.Split(action.UIC, ",")
	if len(teams) == 0 {
		return fmt.Errorf("team is empty")
	}
	for _, team := range teams {
		if team == "" {
			continue
		}
		users := api.UsersOf(team)
		for _, user := range users {
			eventReceiver := &EventReceiver{
				EventID:        eventid,
				ActionConfigID: int64(action.ID),
				Uic:            team,
				User:           user.Name,
				EventTs:        time.Unix(event.EventTime, 0).Format(timeLayout),
			}

			if err := store.AMDB.Create(&eventReceiver).Error; err != nil {
				return fmt.Errorf("create event_receiver err:%v", err)
			}
		}
	}
	return nil
}

func (store *EventStore) insertEvent(eve *coommonModel.Event) (uint, error) {
	event := &Events{
		EventCaseID: eve.Id,
		CurrentStep: eve.CurrentStep,
		Cond:        fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
		Status:      eve.Status,
		EventTs:     time.Unix(eve.EventTime, 0).Format(timeLayout),
	}

	if err := store.AMDB.Create(&event).Error; err != nil {
		return 0, err
	}

	return event.ID, nil
}

func (store *EventStore) updateEventCase(eve *coommonModel.Event, event []EventCases) error {
	var err error
	//TODO: 使用gorm UPDATE, 不使用EXEC
	sqltemplete := `UPDATE event_cases SET
				update_at = ?,
				max_step = ?,
				current_step = ?,
				note = ?,
				cond = ?,
				status = ?,
				func = ?,
				priority = ?,
				tpl_creator = ?,
				action_config_id = ?,
				expression_id = ?,
				strategy_id = ?,
				template_id = ?`
	// reopen case
	if event[0].ProcessStatus == "resolved" || event[0].ProcessStatus == "ignored" {
		sqltemplete = fmt.Sprintf("%v ,process_status = '%s', process_note = %d", sqltemplete, "unresolved", 0)
	}

	tpl_creator := ""
	if eve.Tpl() != nil {
		tpl_creator = eve.Tpl().Creator
	}

	if eve.CurrentStep == 1 {
		// update start time of cases
		sqltemplete = fmt.Sprintf("%v , event_ts = ? WHERE event_case_id = ?", sqltemplete)
		err = store.AMDB.Exec(
			sqltemplete,
			time.Unix(eve.EventTime, 0).Format(timeLayout),
			eve.MaxStep(),
			eve.CurrentStep,
			eve.Note(),
			fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
			eve.Status,
			eve.Func(),
			eve.Priority(),
			tpl_creator,
			api.GetAction(eve).ID,
			eve.ExpressionId(),
			eve.StrategyId(),
			eve.TplId(),
			time.Unix(eve.EventTime, 0).Format(timeLayout),
			eve.Id,
		).Error
	} else {
		sqltemplete = fmt.Sprintf("%v WHERE event_case_id = ?", sqltemplete)
		err = store.AMDB.Exec(
			sqltemplete,
			time.Unix(eve.EventTime, 0).Format(timeLayout),
			eve.MaxStep(),
			eve.CurrentStep,
			eve.Note(),
			fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
			eve.Status,
			eve.Func(),
			eve.Priority(),
			tpl_creator,
			api.GetAction(eve).ID,
			eve.ExpressionId(),
			eve.StrategyId(),
			eve.TplId(),
			eve.Id,
		).Error
	}
	return err
}

func counterGenerate(metric string, tags string) (counter string) {
	counter = metric
	if tags != "" {
		counter = fmt.Sprintf("%s/%s", metric, tags)
	}
	return
}

func (store *EventStore) insertEventCase(eve *coommonModel.Event) error {
	tpl_creator := ""
	if eve.Tpl() != nil {
		tpl_creator = eve.Tpl().Creator
	}

	eventcase := &EventCases{
		EventCaseID:    eve.Id,
		Endpoint:       eve.Endpoint,
		Counter:        counterGenerate(eve.Metric(), utils.SortedTags(eve.PushedTags)),
		Func:           eve.Func(),
		Cond:           fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
		Note:           eve.Note(),
		MaxStep:        eve.MaxStep(),
		CurrentStep:    eve.CurrentStep,
		Priority:       eve.Priority(),
		Status:         eve.Status,
		EventTs:        time.Unix(eve.EventTime, 0).Format(timeLayout),
		UpdateAt:       time.Unix(eve.EventTime, 0).Format(timeLayout),
		TplCreator:     tpl_creator,
		ActionConfigID: api.GetAction(eve).ID,
		ExpressionID:   eve.ExpressionId(),
		StrategyID:     eve.StrategyId(),
		TemplateID:     eve.TplId(),
	}

	if err := store.AMDB.Create(&eventcase).Error; err != nil {
		return err
	}
	return nil
}

func (store *EventStore) isExistOfEventCases(id string) []EventCases {
	var event []EventCases
	store.AMDB.Where("event_case_id = ?", id).Find(&event)
	return event
}

func (store *EventStore) InsertAlarmEvent(eve *coommonModel.Event) error {
	event := store.isExistOfEventCases(eve.Id)
	if len(event) == 0 {
		// if not exist, insert eventcase
		if err := store.insertEventCase(eve); err != nil {
			log.Errorf("insert eventcases fail: %v", err)
			return err
		}

	} else {
		// if exist, update eventcase
		if err := store.updateEventCase(eve, event); err != nil {
			log.Errorf("update eventcases fail: %v", err)
			return err
		}
	}

	// insert event, store event history
	eventid, err := store.insertEvent(eve)
	if err != nil {
		log.Errorf("insert event fail: %v", err)
		return err
	}

	// insert event receiver
	if err := store.insertEventReceiver(eventid, eve); err != nil {
		log.Errorf("insert event receiver fail: %v", err)
		return err
	}
	return nil
}
