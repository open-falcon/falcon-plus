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
package falcon

import (
	"fmt"
	"math"
	"strings"
)

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// code == 0 => success
// code == 1 => bad request
type RpcResp struct {
	Code int `json:"code"`
}

func (p *RpcResp) String() string {
	return fmt.Sprintf("<Code: %d>", p.Code)
}

type Null struct {
}

/* agent/lb */
type LbResp struct {
	Message string
	Total   int
	Invalid int
}

func (p *LbResp) String() string {
	return fmt.Sprintf("Total:%v Invalid:%v Message:%s>",
		p.Total, p.Invalid, p.Message)
}

type MetaData struct {
	Host  string  `json:"host"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Ts    int64   `json:"ts"`
	Step  int64   `json:"step"`
	Type  string  `json:"type"`
	Tags  string  `json:"tags"`
}

func (t *MetaData) String() string {
	return fmt.Sprintf("MetaData host:%s metric:%s Timestamp:%d Step:%d Value:%f type:%s Tags:%v",
		t.Host, t.Name, t.Ts, t.Step, t.Value, t.Type, t.Tags)
}

func (p *MetaData) Id() string {
	return fmt.Sprintf("%s/%s/%s/%s/%d", p.Host, p.Name, p.Tags, p.Type, p.Step)
}

func (p *MetaData) Rrd() (*RrdItem, error) {
	e := &RrdItem{}

	e.Host = p.Host
	e.Name = p.Name
	e.Tags = p.Tags
	e.TimeStemp = p.Ts
	e.Value = p.Value
	e.Step = int(p.Step)
	if e.Step < MIN_STEP {
		e.Step = MIN_STEP
	}
	e.Heartbeat = e.Step * 2

	if p.Type == GAUGE {
		e.Type = p.Type
		e.Min = "U"
		e.Max = "U"
	} else if p.Type == COUNTER {
		e.Type = DERIVE
		e.Min = "0"
		e.Max = "U"
	} else if p.Type == DERIVE {
		e.Type = DERIVE
		e.Min = "0"
		e.Max = "U"
	} else {
		return e, fmt.Errorf("not_supported_counter_type")
	}

	//move to backend
	//e.TimeStemp = e.TimeStemp - e.TimeStemp%int64(e.Step)

	return e, nil
}

func (p *MetaData) Tsdb() *TsdbItem {
	t := TsdbItem{Tags: make(map[string]string)}

	if p.Tags != "" {
		tags := strings.Split(p.Tags, ",")
		for _, tag := range tags {
			kv := strings.SplitN(tag, "=", 2)
			if len(kv) == 2 {
				t.Tags[kv[0]] = kv[1]
			}
		}
	}
	t.Tags["host"] = p.Host
	t.Metric = p.Name
	t.Timestamp = p.Ts
	t.Value = p.Value
	return &t
}

type TsdbItem struct {
	Metric    string            `json:"metric"`
	Tags      map[string]string `json:"tags"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
}

func (p *TsdbItem) String() string {
	return fmt.Sprintf("Metric:%s, Tags:%v, Value:%v, TS:%d",
		p.Metric, p.Tags, p.Value, p.Timestamp)
}

func (p *TsdbItem) TsdbString() (s string) {
	s = fmt.Sprintf("put %s %d %.3f ", p.Metric, p.Timestamp, p.Value)

	for k, v := range p.Tags {
		key := strings.ToLower(strings.Replace(k, " ", "_", -1))
		value := strings.Replace(v, " ", "_", -1)
		s += key + "=" + value + " "
	}

	return s
}

/* lb/storage */
// Type: GAUGE|COUNTER|DERIVE
type RrdItem struct {
	Host      string  `json:"host"`
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	TimeStemp int64   `json:"ts"`
	Step      int     `json:"step"`
	Type      string  `json:"type"`
	Tags      string  `json:"tags"`
	Heartbeat int     `json:"hb"`
	Min       string  `json:"min"`
	Max       string  `json:"max"`
}

func (p *RrdItem) String() string {
	return fmt.Sprintf("Host:%s, Key:%s, Tags:%v, Value:%v, "+
		"TS:%d %v Type:%s, Step:%d, Heartbeat:%d, Min:%s, Max:%s",
		p.Host,
		p.Name,
		p.Tags,
		p.Value,
		p.TimeStemp,
		FmtTs(p.TimeStemp),
		p.Type,
		p.Step,
		p.Heartbeat,
		p.Min,
		p.Max,
	)
}

func (p *RrdItem) Csum() string {
	return Md5sum(p.Id())
}

func (p *RrdItem) Id() string {
	return fmt.Sprintf("%s/%s/%s/%s/%d", p.Host, p.Name, p.Tags, p.Type, p.Step)
}

// ConsolFun 是RRD中的概念，比如：MIN|MAX|AVERAGE
type RrdQuery struct {
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	Host      string `json:"host"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Step      int    `json:"step"`
	ConsolFun string `json:"consolFuc"`
}

func (p *RrdQuery) Csum() string {
	return Md5sum(p.Id())
}

func (p *RrdQuery) Id() string {
	return fmt.Sprintf("%s/%s/%s/%d", p.Host, p.Name, p.Type, p.Step)
}

type RrdResp struct {
	Host string     `json:"host"`
	Name string     `json:"name"`
	Type string     `json:"type"`
	Step int        `json:"step"`
	Vs   []*RRDData `json:"Vs"` //大写为了兼容已经再用这个api的用户
}

func (p *RrdResp) Csum() string {
	return Md5sum(p.Id())
}

func (p *RrdResp) Id() string {
	return fmt.Sprintf("%s/%s/%s/%d", p.Host, p.Name, p.Type, p.Step)
}

type RrdQueryCsum struct {
	Csum      string `json:"csum"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	ConsolFun string `json:"consolFuc"`
}

type RrdRespCsum struct {
	Values []*RRDData `json:"values"`
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
	Ts int64     `json:"ts"`
	V  JsonFloat `json:"v"`
}

func (p *RRDData) String() string {
	return fmt.Sprintf(
		"RRDData:Value:%v TS:%d %v",
		p.V,
		p.Ts,
		FmtTs(p.Ts),
	)
}

type File struct {
	Filename string
	Data     []byte
}
