package graph

import "time"

type Endpoint struct {
	ID               uint              `gorm:"primary_key"`
	Endpoint         string            `json:"endpoint"`
	Ts               int               `json:"-"`
	TCreate          time.Time         `json:"-"`
	TModify          time.Time         `json:"-"`
	EndpointCounters []EndpointCounter `gorm:"ForeignKey:EndpointIDE"`
}

func (Endpoint) TableName() string {
	return "endpoint"
}
