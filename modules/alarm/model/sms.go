package model

import (
	"fmt"
)

type Sms struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}

func (this *Sms) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Content:%s>",
		this.Tos,
		this.Content,
	)
}
