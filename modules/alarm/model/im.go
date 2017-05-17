package model

import (
	"fmt"
)

type IM struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}

func (this *IM) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Content:%s>",
		this.Tos,
		this.Content,
	)
}
