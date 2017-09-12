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
	con "github.com/open-falcon/falcon-plus/modules/api/config"
	"github.com/spf13/viper"
)

type User struct {
	ID     int64  `json:"id" `
	Name   string `json:"name"`
	Cnname string `json:"cnname"`
	Passwd string `json:"-"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	IM     string `json:"im" gorm:"column:im"`
	QQ     string `json:"qq" gorm:"column:qq"`
	Role   int    `json:"role"`
}

func skipAccessControll() bool {
	return !viper.GetBool("access_control")
}

func (this User) IsAdmin() bool {
	if skipAccessControll() {
		return true
	}
	if this.Role == 2 || this.Role == 1 {
		return true
	}
	return false
}

func (this User) IsSuperAdmin() bool {
	if skipAccessControll() {
		return true
	}
	if this.Role == 2 {
		return true
	}
	return false
}

func (this User) FindUser() (user User, err error) {
	db := con.Con()
	user = this
	dt := db.Uic.Find(&user)
	if dt.Error != nil {
		err = dt.Error
		return
	}
	return
}

type Session struct {
	ID      int64
	Uid     int64
	Sig     string
	Expired int
}

func (this Session) TableName() string {
	return "session"
}

func (this User) TableName() string {
	return "user"
}
