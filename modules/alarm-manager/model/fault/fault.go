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
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/config"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/model/event"
)

// TODO:add context to error to make it specific

const (
	ActionFollow   = "FOLLOW"   // follow the fault
	ActionUnfollow = "UNFOLLOW" // unfollow the fault
)

const (
	StateUnprocessing = "UNPROCESSING" // state of fault
	StateProcessing   = "PROCESSING"   // state of fault
	StateClosed       = "CLOSED"       // state of fault
	StateDiscarded    = "DISCARDED"    // state of fault
)

// Store contains db connection info of alarm-manager. See Init.
var Store *FaultStore

func Init() {
	Store = NewFaultStore(config.Con().AM)
}

type FaultStore struct {
	AMDB   *gorm.DB
	Locker *sync.RWMutex
}

func NewFaultStore(db *gorm.DB) *FaultStore {
	return &FaultStore{AMDB: db, Locker: &sync.RWMutex{}}
}

func (fs *FaultStore) Create(createInfo CreateInfo) (FaultInfo, error) {
	tx := fs.AMDB.Begin()
	if tx.Error != nil {
		return FaultInfo{}, tx.Error
	}

	id, err := createFault(tx, &Fault{
		Title:   createInfo.Title,
		Note:    createInfo.Note,
		Creator: createInfo.Creator,
		Owner:   createInfo.Owner,
		State:   StateProcessing,
	})
	if err != nil {
		tx.Rollback()
		return FaultInfo{}, err
	}

	if err := addFollower(tx, id, createInfo.Creator); err != nil {
		tx.Rollback()
		return FaultInfo{}, err
	}

	if err := addTag(tx, id, createInfo.Tags); err != nil {
		tx.Rollback()
		return FaultInfo{}, err
	}

	if err := stateChangeLog(tx, &StateChangeLog{
		ReferTable: Fault{}.TableName(),
		ReferId:    id,
		ReferField: "owner",
		From:       "",
		To:         createInfo.Owner,
	}); err != nil {
		tx.Rollback()
		return FaultInfo{}, err
	}

	if err := stateChangeLog(tx, &StateChangeLog{
		ReferTable: Fault{}.TableName(),
		ReferId:    id,
		ReferField: "state",
		From:       StateUnprocessing,
		To:         StateProcessing,
	}); err != nil {
		tx.Rollback()
		return FaultInfo{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return FaultInfo{}, err
	}

	faultInfo, err := fs.Get(id)
	if err != nil {
		return FaultInfo{}, err
	}

	return faultInfo, nil
}

func createFault(tx *gorm.DB, f *Fault) (uint, error) {
	err := tx.Create(f).Error
	return f.ID, err
}

func addFollower(tx *gorm.DB, faultId uint, follower string) error {
	return tx.Create(&FaultFollower{FaultId: faultId, Follower: follower}).Error
}

func addTag(tx *gorm.DB, faultid uint, tags []string) error {
	for _, v := range tags {
		err := tx.Create(&FaultTag{FaultId: faultid, Tag: v}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func stateChangeLog(tx *gorm.DB, s *StateChangeLog) error {
	return tx.Create(s).Error
}

func (fs *FaultStore) AddEvent(faultid uint, eventids []uint) (FaultInfo, error) {
	tx := fs.AMDB.Begin()
	if tx.Error != nil {
		return FaultInfo{}, tx.Error
	}

	for _, v := range eventids {
		err := tx.Create(&FaultEvent{FaultId: faultid, EventId: v}).Error
		if err != nil {
			tx.Rollback()
			return FaultInfo{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return FaultInfo{}, err
	}

	faultInfo, err := fs.Get(faultid)
	if err != nil {
		return FaultInfo{}, err
	}

	return faultInfo, nil
}

func (fs *FaultStore) DeleteEvent(faultid uint, eventids []uint) (
	FaultInfo, error) {
	tx := fs.AMDB.Begin()
	if tx.Error != nil {
		return FaultInfo{}, tx.Error
	}
	for _, v := range eventids {
		// After removing some record by soft delete,
		// adding same record will fail if  the field exists unique key.
		// Thus,use hard delete rather than soft delete.
		err := tx.Unscoped().Where("fault_id = ? and event_id = ?",
			faultid, v).Delete(&FaultEvent{}).Error
		if err != nil {
			tx.Rollback()
			return FaultInfo{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return FaultInfo{}, err
	}

	faultInfo, err := fs.Get(faultid)
	if err != nil {
		return FaultInfo{}, err
	}

	return faultInfo, nil
}

func (fs *FaultStore) GetEvent(faultid uint) ([]event.EventInfo, error) {
	var faultEvents []FaultEvent
	if err := fs.AMDB.Order("event_id desc").Where("fault_id = ?",
		faultid).Find(&faultEvents).Error; err != nil {
		return []event.EventInfo{}, err
	}

	eventInfos := make([]event.EventInfo, 0)
	for _, v := range faultEvents {
		eventInfo, err := event.Store.GetEventByID(v.EventId)
		if err != nil {
			return []event.EventInfo{}, err
		}
		eventInfos = append(eventInfos, eventInfo)
	}

	return eventInfos, nil
}

func (fs *FaultStore) AddComment(faultid uint, creator, comment string) (
	FaultInfo, error) {
	faultComment := FaultComment{
		FaultId: faultid,
		Creator: creator,
		Comment: comment,
	}
	if err := fs.AMDB.Create(&faultComment).Error; err != nil {
		return FaultInfo{}, err
	}

	faultInfo, err := fs.Get(faultid)
	if err != nil {
		return FaultInfo{}, err
	}

	return faultInfo, nil
}

func (fs *FaultStore) DeleteComment(faultid uint, creator, comment string) (FaultInfo, error) {
	err := fs.AMDB.Unscoped().Where("fault_id = ? and creator = ? and comment = ?",
		faultid, creator, comment).Delete(&FaultComment{}).Error
	if err != nil {
		return FaultInfo{}, err
	}

	faultInfo, err := fs.Get(faultid)
	if err != nil {
		return FaultInfo{}, err
	}

	return faultInfo, nil
}

func (fs *FaultStore) GetComment(faultid uint) ([]CommentInfo, error) {
	var faultComment []FaultComment
	if err := fs.AMDB.Order("created_at desc").Where("fault_id = ?",
		faultid).Find(&faultComment).Error; err != nil {
		return []CommentInfo{}, err
	}

	commentInfos := make([]CommentInfo, 0, len(faultComment))
	for _, v := range faultComment {
		commentInfos = append(commentInfos, CommentInfo{
			CreatedAt: v.CreatedAt,
			Creator:   v.Creator,
			Comment:   v.Comment,
		})
	}
	return commentInfos, nil
}

func (fs *FaultStore) UpdateOwner(faultid uint, newOwner string) (
	FaultInfo, error) {

	// Lock to protect data from inconsistent.
	fs.Locker.Lock()
	defer fs.Locker.Unlock()

	fault, err := fs.GetFaultBasic(faultid)
	if err != nil {
		return FaultInfo{}, err
	}
	if is := ownerIsValid(fault.Owner, newOwner); !is {
		return FaultInfo{}, fmt.Errorf("new owner is the same as old owner")
	}

	tx := fs.AMDB.Begin()
	if tx.Error != nil {
		return FaultInfo{}, tx.Error
	}
	changeInfo := &StateChangeLog{
		ReferTable: Fault{}.TableName(),
		ReferId:    faultid,
		ReferField: "owner",
		From:       fault.Owner,
		To:         newOwner,
	}
	if err := update(tx, changeInfo); err != nil {
		return FaultInfo{}, err
	}

	faultInfo, err := fs.Get(faultid)
	if err != nil {
		return FaultInfo{}, err
	}

	return faultInfo, nil
}

func ownerIsValid(from string, to string) bool {
	return from != to
}

func update(tx *gorm.DB, changeInfo *StateChangeLog) error {
	if err := tx.Model(&Fault{}).Where("id = ?",
		changeInfo.ReferId).Update(changeInfo.ReferField,
		changeInfo.To).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(changeInfo).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (fs *FaultStore) UpdateState(faultid uint, newState string) (FaultInfo, error) {
	// Lock to protect data from inconsistent.
	fs.Locker.Lock()
	defer fs.Locker.Unlock()

	fault, err := fs.GetFaultBasic(faultid)
	if err != nil {
		return FaultInfo{}, err
	}
	if is := stateIsValid(fault.State, newState); !is {
		return FaultInfo{}, fmt.Errorf("new state is invalid")
	}

	tx := fs.AMDB.Begin()
	if tx.Error != nil {
		return FaultInfo{}, tx.Error
	}
	changeInfo := &StateChangeLog{
		ReferTable: Fault{}.TableName(),
		ReferId:    faultid,
		ReferField: "state",
		From:       fault.State,
		To:         newState,
	}
	if err := update(tx, changeInfo); err != nil {
		return FaultInfo{}, err
	}

	faultInfo, err := fs.Get(faultid)
	if err != nil {
		return FaultInfo{}, err
	}

	return faultInfo, nil
}

func stateIsValid(from string, to string) bool {
	var is bool
	switch from {
	case StateProcessing:
		is = (to == StateClosed) || (to == StateDiscarded)
	case StateClosed:
		is = (to == StateProcessing) || (to == StateDiscarded)
	case StateDiscarded:
		is = to == StateProcessing
	default:
		is = false
	}
	return is
}

func (fs *FaultStore) AddTag(faultid uint, tags []string) (FaultInfo, error) {
	tx := fs.AMDB.Begin()
	if tx.Error != nil {
		return FaultInfo{}, tx.Error
	}

	for _, v := range tags {
		err := tx.Create(&FaultTag{FaultId: faultid, Tag: v}).Error
		if err != nil {
			tx.Rollback()
			return FaultInfo{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return FaultInfo{}, err
	}

	faultInfo, err := fs.Get(faultid)
	if err != nil {
		return FaultInfo{}, err
	}

	return faultInfo, nil
}

func (fs *FaultStore) DeleteTag(faultid uint, tags []string) (
	FaultInfo, error) {

	tx := fs.AMDB.Begin()
	if tx.Error != nil {
		return FaultInfo{}, tx.Error
	}

	for _, v := range tags {
		// After removing some record by soft delete,
		// adding same record will fail if  the field exists unique key.
		// Thus,use hard delete rather than soft delete.
		err := tx.Unscoped().Where("fault_id = ? and tag = ?",
			faultid, v).Delete(&FaultTag{}).Error
		if err != nil {
			tx.Rollback()
			return FaultInfo{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return FaultInfo{}, err
	}

	faultInfo, err := fs.Get(faultid)
	if err != nil {
		return FaultInfo{}, err
	}

	return faultInfo, nil
}

func (fs *FaultStore) UpdateFollower(faultid uint, follower []string,
	action string) (FaultInfo, error) {

	if action != ActionFollow && action != ActionUnfollow {
		return FaultInfo{}, fmt.Errorf("action is invalid")
	}

	tx := fs.AMDB.Begin()
	if tx.Error != nil {
		return FaultInfo{}, tx.Error
	}

	if action == ActionFollow {
		for _, v := range follower {
			err := tx.Create(&FaultFollower{FaultId: faultid, Follower: v}).Error
			if err != nil {
				tx.Rollback()
				return FaultInfo{}, err
			}
		}
	}

	if action == ActionUnfollow {
		for _, v := range follower {
			// After removing some record by soft delete,
			// adding same record will fail if  the field exists unique key.
			// Thus,use hard delete rather than soft delete.
			if err := tx.Unscoped().Where("fault_id = ? and follower = ?",
				faultid, v).Delete(&FaultFollower{}).Error; err != nil {
				tx.Rollback()
				return FaultInfo{}, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return FaultInfo{}, err
	}

	faultInfo, err := fs.Get(faultid)
	if err != nil {
		return FaultInfo{}, err
	}

	return faultInfo, nil
}

func (fs *FaultStore) UpdateBasic(faultid uint, title, note string) (
	FaultInfo, error) {

	if err := fs.AMDB.Model(&Fault{}).Where("id = ?",
		faultid).Updates(Fault{Title: title, Note: note}).Error; err != nil {
		return FaultInfo{}, err
	}

	faultInfo, err := fs.Get(faultid)
	if err != nil {
		return FaultInfo{}, err
	}

	return faultInfo, nil
}

func (fs *FaultStore) List(filter Filter) ([]FaultInfo, uint, error) {
	countDB := fs.AMDB
	faultDB := fs.AMDB.Order("created_at desc").Limit(filter.Limit).Offset(filter.Offset)
	if filter.Follower != "" {
		countDB = countDB.Joins("JOIN fault_follower ON fault_follower.fault_id = fault.id")
		faultDB = faultDB.Joins("JOIN fault_follower ON fault_follower.fault_id = fault.id")
	}
	if filter.Tag != "" {
		countDB = countDB.Joins("JOIN fault_tag ON fault_tag.fault_id = fault.id")
		faultDB = faultDB.Joins("JOIN fault_tag ON fault_tag.fault_id = fault.id")
	}
	sql := SQLBuilder(filter)

	var count uint
	err := countDB.Where(sql).Find(&([]Fault{})).Count(&count).Error
	if err != nil {
		return []FaultInfo{}, 0, err
	}

	faults := make([]Fault, 0)
	err = faultDB.Where(sql).Find(&faults).Error
	if err != nil {
		return []FaultInfo{}, 0, err
	}

	faultinfos := make([]FaultInfo, 0, len(faults))
	for _, v := range faults {
		faultinfo, err := fs.Get(v.ID)
		if err != nil {
			return []FaultInfo{}, 0, err
		}
		faultinfos = append(faultinfos, faultinfo)
	}

	return faultinfos, count, nil
}

func (fs *FaultStore) Get(id uint) (FaultInfo, error) {
	fault, err := fs.GetFaultBasic(id)
	if err != nil {
		return FaultInfo{}, err
	}

	commentInfos, err := fs.GetComment(id)
	if err != nil {
		return FaultInfo{}, err
	}

	followers, err := fs.GetFaultFollower(id)
	if err != nil {
		return FaultInfo{}, err
	}

	tags, err := fs.GetTag(id)
	if err != nil {
		return FaultInfo{}, err
	}

	eventInfos, err := fs.GetEvent(id)
	if err != nil {
		return FaultInfo{}, err
	}

	return FaultInfo{
		Id:        fault.ID,
		CreatedAt: fault.CreatedAt,
		Title:     fault.Title,
		Note:      fault.Note,
		Creator:   fault.Creator,
		Owner:     fault.Owner,
		State:     fault.State,
		Tags:      tags,
		Events:    eventInfos,
		Followers: followers,
		Comments:  commentInfos,
	}, nil
}

func (fs *FaultStore) GetFaultBasic(faultid uint) (Fault, error) {
	var fault Fault
	err := fs.AMDB.Where("id = ?", faultid).First(&fault).Error
	if err != nil {
		return Fault{}, err
	}

	return fault, nil
}

func (fs *FaultStore) GetFaultFollower(faultid uint) ([]string, error) {
	var faultFollowers []FaultFollower
	err := fs.AMDB.Where("fault_id = ?", faultid).Find(&faultFollowers).Error
	if err != nil {
		return []string{}, err
	}

	followers := make([]string, 0, len(faultFollowers))
	for _, v := range faultFollowers {
		followers = append(followers, v.Follower)
	}

	return followers, nil
}

func (fs *FaultStore) GetTag(faultid uint) ([]string, error) {
	var faultTags []FaultTag
	err := fs.AMDB.Where("fault_id = ?", faultid).Find(&faultTags).Error
	if err != nil {
		return []string{}, err
	}

	tags := make([]string, 0, len(faultTags))
	for _, v := range faultTags {
		tags = append(tags, v.Tag)
	}

	return tags, nil
}

// TODO: enrich the timeline of fault.
func (fs *FaultStore) GetTimeLine(faultid uint) (TimeLine, error) {
	eventInfos, err := fs.GetEvent(faultid)
	if err != nil {
		return TimeLine{}, err
	}

	firstTime, lastTime := EventTimeInfo(eventInfos)

	createdAt, closedAt, err := fs.FaultTimeInfo(faultid)
	if err != nil {
		return TimeLine{}, err
	}

	return TimeLine{
		FirstEventTime: firstTime,
		LastEventTime:  lastTime,
		FaultCreatedAt: createdAt,
		FaultClosedAt:  closedAt,
	}, nil
}

func EventTimeInfo(eventInfos []event.EventInfo) (firstTime, lastTime string) {
	if len(eventInfos) == 0 {
		return firstTime, lastTime
	}

	createdTime := make([]string, 0, len(eventInfos))
	for _, v := range eventInfos {
		createdTime = append(createdTime, v.EventTs)
	}

	sort.Strings(createdTime)

	return createdTime[0], createdTime[len(createdTime)-1]
}

func (fs *FaultStore) FaultTimeInfo(faultid uint) (
	createdAt, closedAt time.Time, err error) {

	var changeLog []StateChangeLog
	if err := fs.AMDB.Where("refer_id = ? and refer_field = ?",
		faultid, "state").Find(&changeLog).Error; err != nil {
		return createdAt, closedAt, err
	}

	if len(changeLog) == 0 {
		return createdAt, closedAt, fmt.Errorf("fault does not exist")
	}

	first := changeLog[0]
	if first.To == StateProcessing {
		createdAt = first.CreatedAt
	}
	last := changeLog[len(changeLog)-1]
	if last.To == StateClosed {
		closedAt = last.CreatedAt
	}

	return createdAt, closedAt, nil
}

func SQLBuilder(filter Filter) string {
	result := make([]string, 0)

	if filter.Start != 0 {
		start := time.Unix(int64(filter.Start), 0)
		sql := fmt.Sprintf(" %s.created_at >= '%s' ", Fault{}.TableName(), start)
		result = append(result, sql)
	}

	if filter.End != 0 {
		end := time.Unix(int64(filter.End), 0)
		sql := fmt.Sprintf(" %s.created_at <= '%s' ", Fault{}.TableName(), end)
		result = append(result, sql)
	}

	if filter.Creator != "" {
		sql := fmt.Sprintf(" %s.creator = '%s' ",
			Fault{}.TableName(), filter.Creator)
		result = append(result, sql)
	}

	if filter.Owner != "" {
		sql := fmt.Sprintf(" %s.owner = '%s' ",
			Fault{}.TableName(), filter.Owner)
		result = append(result, sql)
	}

	if filter.State != "" {
		sql := fmt.Sprintf(" %s.state = '%s' ",
			Fault{}.TableName(), filter.State)
		result = append(result, sql)
	}

	if filter.Title != "" {
		sql := fmt.Sprintf(" %s.title like '%%%s%%' ",
			Fault{}.TableName(), filter.Title)
		result = append(result, sql)
	}

	if filter.Follower != "" {
		sql := fmt.Sprintf(" %s.follower = '%s' ",
			FaultFollower{}.TableName(), filter.Follower)
		result = append(result, sql)
	}

	if filter.Tag != "" {
		sql := fmt.Sprintf(" %s.tag = '%s' ",
			FaultTag{}.TableName(), filter.Tag)
		result = append(result, sql)
	}

	return strings.Join(result, "and")
}
