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

package falcon_portal

import (
	"github.com/open-falcon/falcon-plus/modules/api/app/model/uic"
	con "github.com/open-falcon/falcon-plus/modules/api/config"
	log "github.com/sirupsen/logrus"
)

type Template struct {
	ID         int64  `json:"id" gorm:"column:id"`
	Name       string `json:"tpl_name" gorm:"column:tpl_name"`
	ParentID   int64  `json:"parent_id" orm:"column:parent_id"`
	ActionID   int64  `json:"action_id" orm:"column:action_id"`
	CreateUser string `json:"create_user" orm:"column:create_user"`
}

func (this Template) TableName() string {
	return "tpl"
}

func (this Template) FindUserName() (name string, err error) {
	var user uic.User
	user.Name = this.CreateUser
	db := con.Con()
	dt := db.Uic.Find(&user)
	if dt.Error != nil {
		err = dt.Error
		return
	}
	name = user.Name
	return
}

func (this Template) FindParentName() (name string, err error) {
	var ptpl Template
	if this.ParentID == 0 {
		return
	}
	ptpl.ID = this.ParentID
	db := con.Con()
	dt := db.Falcon.Find(&ptpl)
	if dt.Error != nil {
		log.Debugf("tpl_id: %v find parent: %v with error: %s", this.ID, ptpl.ID, dt.Error.Error())
		return
	}
	name = ptpl.Name
	return
}
