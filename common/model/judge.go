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

	"github.com/open-falcon/falcon-plus/common/utils"
)

type JudgeItem struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	Value     float64           `json:"value"`
	ValueRaw  string            `json:"valueRaw"`
	Timestamp int64             `json:"timestamp"`
	JudgeType string            `json:"judgeType"`
	Tags      map[string]string `json:"tags"`
}

func (this *JudgeItem) String() string {
	return fmt.Sprintf("<Endpoint:%s, Metric:%s, Value:%.2f ValueRaw:%s, Timestamp:%d, JudgeType:%s Tags:%v>",
		this.Endpoint,
		this.Metric,
		this.Value,
		this.ValueRaw,
		this.Timestamp,
		this.JudgeType,
		this.Tags)
}

func (this *JudgeItem) PrimaryKey() string {
	return utils.Md5(utils.PK(this.Endpoint, this.Metric, this.Tags))
}

type HistoryData struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
	ValueRaw  string  `json:"valueRaw"`
}
