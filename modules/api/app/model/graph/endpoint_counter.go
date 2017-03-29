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
