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
	"time"
)

// +----------+---------------------+------+-----+-------------------+-----------------------------+
// | Field    | Type                | Null | Key | Default           | Extra                       |
// +----------+---------------------+------+-----+-------------------+-----------------------------+
// | id       | bigint(20) unsigned | NO   | PRI | NULL              | auto_increment              |
// | name     | varchar(255)        | NO   | UNI |                   |                             |
// | obj      | varchar(10240)      | NO   |     |                   |                             |
// | obj_type | varchar(255)        | NO   |     |                   |                             |
// | metric   | varchar(128)        | NO   |     |                   |                             |
// | tags     | varchar(1024)       | NO   |     |                   |                             |
// | dstype   | varchar(32)         | NO   |     | GAUGE             |                             |
// | step     | int(11) unsigned    | NO   |     | 60                |                             |
// | mock     | double              | NO   |     | 0                 |                             |
// | creator  | varchar(64)         | NO   |     |                   |                             |
// | t_create | datetime            | NO   |     | NULL              |                             |
// | t_modify | timestamp           | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
// +----------+---------------------+------+-----+-------------------+-----------------------------+

//no_data
type Mockcfg struct {
	ID   int64  `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
	Obj  string `json:"obj" gorm:"column:obj"`
	//group, host, other
	ObjType string    `json:"obj_type" gorm:"column:obj_type"`
	Metric  string    `json:"metric" gorm:"column:metric"`
	Tags    string    `json:"tags" gorm:"column:tags"`
	DsType  string    `json:"dstype" gorm:"column:dstype"`
	Step    int       `json:"step" gorm:"column:step"`
	Mock    float64   `json:"mock" gorm:"column:mock"`
	Creator string    `json:"creator" gorm:"column:creator"`
	TCreate time.Time `json:"-"`
}

func (this Mockcfg) TableName() string {
	return "mockcfg"
}
