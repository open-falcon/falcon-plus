package model

import (
	"fmt"

	"github.com/open-falcon/falcon-plus/common/utils"
)

type EJudgeItem struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	JudgeType string            `json:"judgeType"`
	Values    map[string]string `json:"values"`
	Filters   map[string]string `json:"filters"`
	Timestamp int64             `json:"timestamp"`
}

func (this *EJudgeItem) String() string {
	return fmt.Sprintf("<Endpoint:%s, Metric:%s, Timestamp:%d, JudgeType:%s Values:%v Filters:%v>",
		this.Endpoint,
		this.Metric,
		this.Timestamp,
		this.JudgeType,
		this.Values,
		this.Filters)
}

func (this *EJudgeItem) PrimaryKey() string {
	return utils.Md5(this.Metric)
}

type EHistoryData struct {
	Timestamp int64              `json:"timestamp"`
	Values    map[string]float64 `json:"values"`
}
