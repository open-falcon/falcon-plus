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
package models

import "fmt"

const (
	STATUS_CREATING = iota
	STATUS_PENDING
	STATUS_PROCESSING
	STATUS_SOLVED
	STATUS_UPGRADED
	STATUS_IGNORED
)

type Matter struct {
	Id        int64
	Status    int
	Tag       string
	Starttime int64
	Endtime   int64
}

type Events struct {
	Eventid    string
	Endpoint   string
	Hostname   string
	Pdl        string
	Metric     string
	Pushedtags string
}

type Claim struct {
	Matter    int64
	User      string
	Timestamp int64
	Commit    string
}
func (op *Operator) QueryMatters(status int, per int, offset int) ([]Matter, error) {
	var matters []Matter
	_, err := op.O.Raw("SELECT a.matter AS id , b.tag, b.uic, b.status, b.starttime , b.endtime FROM alarm_event.user_matter a ,alarm_event.matter b WHERE a.user=? AND b.status=? AND b.id=a.matter order by b.id desc limit ? offset  ?", op.User.Name, status, per, offset).QueryRows(&matters)
	return matters, err
}

func (op *Operator) QueryEventsByMatter(matterID int64, per, offset int) []Events {
	var events []Events
	_, err := op.O.Raw("SELECT distinct a.event AS eventid, b.endpoint, b.hostname AS hostname, b.pdl, b.metric,b.pushed_tags AS pushedtags FROM alarm_event.matter_event a, alarm_event.event b WHERE a.matter=? AND a.event =b.event_id order by a.id desc limit ? offset ?", matterID, per, offset).QueryRows(&events)
	if err != nil {
		fmt.Println(err)
	}
	encountered := map[Events]bool{}
	result := []Events{}
	for v := range events {
		if encountered[events[v]] == true {
		} else {
			encountered[events[v]] = true
			result = append(result, events[v])
		}
	}
	return result
}

func (op *Operator) QueryEventsCntByMatter(matterID int64) (int, error) {
	var events []Events
	_, err := op.O.Raw("SELECT distinct a.event AS eventid FROM alarm_event.matter_event a, alarm_event.event b WHERE a.matter=? AND a.event =b.event_id order by a.id desc", matterID).QueryRows(&events)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	encountered := map[Events]bool{}
	result := []Events{}
	for v := range events {
		if encountered[events[v]] == true {
		} else {
			encountered[events[v]] = true
			result = append(result, events[v])
		}
	}
	cnt := len(result)
	return cnt, err
}

func (op *Operator) UpdateMatter(id int64, _m Matter) error {
	if _, err := op.GetMatter(id); err != nil {
		return err
	}
	_, err := op.O.Raw("UPDATE alarm_event.matter SET status = ? WHERE id =?", _m.Status, id).Exec()
	return err

}
func (op *Operator) GetMatter(id int64) (Matter, error) {
	var matter Matter
	err := op.O.Raw("SELECT id , status, tag, uic , starttime , endtime FROM alarm_event.matter WHERE id=? ", id).QueryRow(&matter)
	return matter, err
}

func (op *Operator) GetMatterCnt(status int) (int64, error) {
	var matterCnt int64
	err := op.O.Raw("SELECT count(*) AS cnt FROM alarm_event.user_matter a ,alarm_event.matter b WHERE a.user=? AND b.status=? AND b.id=a.matter order by b.id desc", op.User.Name, status).QueryRow(&matterCnt)
	return matterCnt, err
}

func (op *Operator) AddClaim(claim Claim) error {
	op.O.Begin()
	_, err := op.O.Raw("UPDATE alarm_event.matter SET status = ? WHERE id =?", STATUS_PROCESSING, claim.Matter).Exec()
	_, err = op.O.Raw("INSERT INTO alarm_event.matter_claim (`matter`,`user`,`comment`,`timestamp`) VALUES (?,?,?,?)", claim.Matter, claim.User, claim.Commit, claim.Timestamp).Exec()
	if err != nil {
		err = op.O.Rollback()
	} else {
		err = op.O.Commit()
	}
	return err
}
