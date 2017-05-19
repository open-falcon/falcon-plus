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

package model

import (
	"fmt"
)

type Template struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parentId"`
	ActionId int    `json:"actionId"`
	Creator  string `json:"creator"`
}

func (this *Template) String() string {
	return fmt.Sprintf(
		"<Id:%d, Name:%s, ParentId:%d, ActionId:%d, Creator:%s>",
		this.Id,
		this.Name,
		this.ParentId,
		this.ActionId,
		this.Creator,
	)
}
