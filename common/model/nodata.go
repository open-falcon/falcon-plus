package model

import (
	"fmt"

	"github.com/open-falcon/falcon-plus/common/utils"
	ttime "github.com/toolkits/time"
)

type NodataItem struct {
	Counter string `json:"counter"`
	Ts      int64  `json:"ts"`
	FStatus string `json:"fstatus"`
	FTs     int64  `json:"fts"`
}

func (this *NodataItem) String() string {
	return fmt.Sprintf("{NodataItem counter:%s ts:%s fecthStatus:%s fetchTs:%s}",
		this.Counter, ttime.FormatTs(this.Ts), this.FStatus, ttime.FormatTs(this.FTs))
}

type NodataConfig struct {
	Id       int               `json:"id"`
	Name     string            `json:"name"`
	ObjType  string            `json:"objType"`
	Endpoint string            `json:"endpoint"`
	Metric   string            `json:"metric"`
	Tags     map[string]string `json:"tags"`
	Type     string            `json:"type"`
	Step     int64             `json:"step"`
	Mock     float64           `json:"mock"`
}

func NewNodataConfig(id int, name string, objType string, endpoint string, metric string, tags map[string]string, dstype string, step int64, mock float64) *NodataConfig {
	return &NodataConfig{id, name, objType, endpoint, metric, tags, dstype, step, mock}
}

func (this *NodataConfig) String() string {
	return fmt.Sprintf("{NodataConfig id:%d, name:%s, objType:%s, endpoint:%s, metric:%s, tags:%s, type:%s, step:%d, mock:%f}",
		this.Id, this.Name, this.ObjType, this.Endpoint, this.Metric, utils.SortedTags(this.Tags), this.Type, this.Step, this.Mock)
}
