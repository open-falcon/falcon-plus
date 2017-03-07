package g

import (
	"fmt"
	"sync"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/utils"
)

type EventDto struct {
	Id       string `json:"id"`
	Endpoint string `json:"endpoint"`
	Metric   string `json:"metric"`
	Counter  string `json:"counter"`

	Func       string `json:"func"`
	LeftValue  string `json:"leftValue"`
	Operator   string `json:"operator"`
	RightValue string `json:"rightValue"`

	Note string `json:"note"`

	MaxStep     int `json:"maxStep"`
	CurrentStep int `json:"currentStep"`
	Priority    int `json:"priority"`

	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`

	ExpressionId int `json:"expressionId"`
	StrategyId   int `json:"strategyId"`
	TemplateId   int `json:"templateId"`

	Link string `json:"link"`
}

type SafeEvents struct {
	sync.RWMutex
	M map[string]*EventDto
}

type OrderedEvents []*EventDto

func (this OrderedEvents) Len() int {
	return len(this)
}
func (this OrderedEvents) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this OrderedEvents) Less(i, j int) bool {
	return this[i].Timestamp < this[j].Timestamp
}

var Events = &SafeEvents{M: make(map[string]*EventDto)}

func (this *SafeEvents) Delete(id string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, id)
}

func (this *SafeEvents) Len() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.M)
}

func (this *SafeEvents) Clone() map[string]*EventDto {
	m := make(map[string]*EventDto)
	this.RLock()
	defer this.RUnlock()
	for key, val := range this.M {
		m[key] = val
	}
	return m
}

func (this *SafeEvents) Put(event *model.Event) {
	if event.Status == "OK" {
		this.Delete(event.Id)
		return
	}

	dto := &EventDto{}
	dto.Id = event.Id
	dto.Endpoint = event.Endpoint
	dto.Metric = event.Metric()
	dto.Counter = event.Counter()
	dto.Func = event.Func()
	dto.LeftValue = utils.ReadableFloat(event.LeftValue)
	dto.Operator = event.Operator()
	dto.RightValue = utils.ReadableFloat(event.RightValue())
	dto.Note = event.Note()

	dto.MaxStep = event.MaxStep()
	dto.CurrentStep = event.CurrentStep
	dto.Priority = event.Priority()

	dto.Status = event.Status
	dto.Timestamp = event.EventTime

	dto.ExpressionId = event.ExpressionId()
	dto.StrategyId = event.StrategyId()
	dto.TemplateId = event.TplId()

	dto.Link = Link(event)

	this.Lock()
	defer this.Unlock()
	this.M[dto.Id] = dto
}

func Link(event *model.Event) string {
	tplId := event.TplId()
	if tplId != 0 {
		return fmt.Sprintf("%s/template/view/%d", Config().Api.Portal, tplId)
	}

	eid := event.ExpressionId()
	if eid != 0 {
		return fmt.Sprintf("%s/expression/view/%d", Config().Api.Portal, eid)
	}

	return ""
}
