package model

import (
	"fmt"
)

type Mail struct {
	Tos     string `json:"tos"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func (this *Mail) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Subject:%s, Content:%s>",
		this.Tos,
		this.Subject,
		this.Content,
	)
}
