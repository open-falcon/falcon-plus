package model

import (
	"fmt"
)

type Template struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parentId"`
	ActionId int    `json:"actionId"`
	Creator  string `json:"creator"`
}

func (this *Template) String() string {
	return fmt.Sprintf(
		"<Id:%d, Name:%s, ParentId:%d, ActionId:%d, Creator:%s>",
		this.Id,
		this.Name,
		this.ParentId,
		this.ActionId,
		this.Creator,
	)
}
