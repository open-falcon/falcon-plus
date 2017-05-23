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

// +---------+------------------+------+-----+---------+-------+
// | Field   | Type             | Null | Key | Default | Extra |
// +---------+------------------+------+-----+---------+-------+
// | grp_id  | int(10) unsigned | NO   | PRI | NULL    |       |
// | host_id | int(11)          | NO   | PRI | NULL    |       |
// +---------+------------------+------+-----+---------+-------+

type GrpHost struct {
	GrpID  int64 `json:"grp_id" gorm:"column:grp_id"`
	HostID int64 `json:"host_id" gorm:"column:host_id"`
}

func (this GrpHost) TableName() string {
	return "grp_host"
}
