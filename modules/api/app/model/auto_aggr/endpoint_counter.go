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

package graph

import "time"

type EndpointCounter struct {
	ID         uint `gorm:"primary_key"`
	EndpointID int
	Counter    string
	Step       int
	Type       string
	Ts         int
	TCreate    time.Time
	TModify    time.Time
}

func (EndpointCounter) TableName() string {
	return "endpoint_counter"
}
