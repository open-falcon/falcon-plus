package model

import (
	"fmt"
)

type Chat struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}

func (this *Chat) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Content:%s>",
		this.Tos,
		this.Content,
	)
}
