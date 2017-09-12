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
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

type RelTeamUser struct {
	ID  int64
	Tid int64
	Uid int64
}

func (this RelTeamUser) TableName() string {
	return "rel_team_user"
}

func (this RelTeamUser) Me() {
	db := config.Con()
	db.Uic.Where("id = 1")
}
