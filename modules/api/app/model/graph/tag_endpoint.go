package graph

import "time"

type TagEndpoint struct {
	ID         uint `gorm:"primary_key"`
	Tag        string
	EndpointID int
	Ts         int
	TCreate    time.Time
	TModify    time.Time
}

func (TagEndpoint) TableName() string {
	return "tag_endpoint"
}
