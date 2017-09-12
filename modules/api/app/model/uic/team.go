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

package uic

import (
	"errors"
	"fmt"

	"github.com/open-falcon/falcon-plus/modules/api/config"
)

type Team struct {
	ID      int64  `json:"id,"`
	Name    string `json:"name"`
	Resume  string `json:"resume"`
	Creator int64  `json:"creator"`
}

func (this Team) TableName() string {
	return "team"
}

func (this Team) Members() (users []User, err error) {
	db := config.Con()
	var tmapping []RelTeamUser
	if dt := db.Uic.Where("tid = ?", this.ID).Find(&tmapping); dt.Error != nil {
		err = dt.Error
		return
	}
	users = []User{}
	var uids []int64
	for _, t := range tmapping {
		uids = append(uids, t.Uid)
	}
	//no user bind to team
	if len(uids) == 0 {
		return
	}
	uidstr, err := arrIntToString(uids)
	if err != nil {
		return
	}

	if dt := db.Uic.Select("name, id, cnname").Where(fmt.Sprintf("id in (%s)", uidstr)).Find(&users); dt.Error != nil {
		err = dt.Error
		return
	}
	return
}

func (this Team) GetCreatorName() (userName string, err error) {
	userName = "unknown"
	db := config.Con()
	user := User{ID: this.Creator}
	if dt := db.Uic.Find(&user); dt.Error != nil {
		err = dt.Error
	} else {
		userName = user.Name
	}
	return
}

func arrIntToString(arr []int64) (result string, err error) {
	result = ""
	for indx, a := range arr {
		if indx == 0 {
			result = fmt.Sprintf("%v", a)
		} else {
			result = fmt.Sprintf("%v,%v", result, a)
		}
	}
	if result == "" {
		err = errors.New(fmt.Sprintf("array is empty, err: %v", arr))
	}
	return
}
