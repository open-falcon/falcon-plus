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
	"math"

	MUtils "github.com/open-falcon/falcon-plus/common/utils"
)

// DsType 即RRD中的Datasource的类型：GAUGE|COUNTER|DERIVE
type GraphItem struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	Tags      map[string]string `json:"tags"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	DsType    string            `json:"dstype"`
	Step      int               `json:"step"`
	Heartbeat int               `json:"heartbeat"`
	Min       string            `json:"min"`
	Max       string            `json:"max"`
}

func (this *GraphItem) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Tags:%v, Value:%v, TS:%d %v DsType:%s, Step:%d, Heartbeat:%d, Min:%s, Max:%s>",
		this.Endpoint,
		this.Metric,
		this.Tags,
		this.Value,
		this.Timestamp,
		MUtils.UnixTsFormat(this.Timestamp),
		this.DsType,
		this.Step,
		this.Heartbeat,
		this.Min,
		this.Max,
	)
}

func (this *GraphItem) PrimaryKey() string {
	return MUtils.PK(this.Endpoint, this.Metric, this.Tags)
}

func (t *GraphItem) Checksum() string {
	return MUtils.Checksum(t.Endpoint, t.Metric, t.Tags)
}

func (this *GraphItem) UUID() string {
	return MUtils.UUID(this.Endpoint, this.Metric, this.Tags, this.DsType, this.Step)
}

type GraphDeleteParam struct {
	Endpoint string `json:"endpoint"`
	Metric   string `json:"metric"`
	Step     int    `json:"step"`
	DsType   string `json:"dstype"`
	Tags     string `json:"tags"`
}

type GraphDeleteResp struct {
}

// ConsolFun 是RRD中的概念，比如：MIN|MAX|AVERAGE
type GraphQueryParam struct {
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	ConsolFun string `json:"consolFuc"`
	Endpoint  string `json:"endpoint"`
	Counter   string `json:"counter"`
	Step      int    `json:"step"`
}

type GraphQueryResponse struct {
	Endpoint string     `json:"endpoint"`
	Counter  string     `json:"counter"`
	DsType   string     `json:"dstype"`
	Step     int        `json:"step"`
	Values   []*RRDData `json:"Values"` //大写为了兼容已经再用这个api的用户
}

// 页面上已经可以看到DsType和Step了，直接带进查询条件，Graph更易处理
type GraphAccurateQueryParam struct {
	Checksum  string `json:"checksum"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	ConsolFun string `json:"consolFuc"`
	DsType    string `json:"dsType"`
	Step      int    `json:"step"`
}

type GraphAccurateQueryResponse struct {
	Values []*RRDData `json:"Values"`
}

type JsonFloat float64

func (v JsonFloat) MarshalJSON() ([]byte, error) {
	f := float64(v)
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return []byte("null"), nil
	} else {
		return []byte(fmt.Sprintf("%f", f)), nil
	}
}

type RRDData struct {
	Timestamp int64     `json:"timestamp"`
	Value     JsonFloat `json:"value"`
}

func NewRRDData(ts int64, val float64) *RRDData {
	return &RRDData{Timestamp: ts, Value: JsonFloat(val)}
}

func (this *RRDData) String() string {
	return fmt.Sprintf(
		"<RRDData:Value:%v TS:%d %v>",
		this.Value,
		this.Timestamp,
		MUtils.UnixTsFormat(this.Timestamp),
	)
}

type GraphInfoParam struct {
	Endpoint string `json:"endpoint"`
	Counter  string `json:"counter"`
}

type GraphInfoResp struct {
	ConsolFun string `json:"consolFun"`
	Step      int    `json:"step"`
	Filename  string `json:"filename"`
}

type GraphFullyInfo struct {
	Endpoint  string `json:"endpoint"`
	Counter   string `json:"counter"`
	ConsolFun string `json:"consolFun"`
	Step      int    `json:"step"`
	Filename  string `json:"filename"`
	Addr      string `json:"addr"`
}

type GraphLastParam struct {
	Endpoint string `json:"endpoint"`
	Counter  string `json:"counter"`
}

type GraphLastResp struct {
	Endpoint string   `json:"endpoint"`
	Counter  string   `json:"counter"`
	Value    *RRDData `json:"value"`
}
