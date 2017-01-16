package uic

import (
	"errors"
	"fmt"

	"github.com/masato25/owl_backend/config"
)

type Team struct {
	ID      int64  `json:"id,"`
	Name    string `json:"name"`
	Resume  string `json:"resume"`
	Creator int64  `json:"creator"`
}

func (this Team) TableName() string {
	return "team"
}

func (this Team) Members() (users []User, err error) {
	db := config.Con()
	var tmapping []RelTeamUser
	if dt := db.Uic.Where("tid = ?", this.ID).Find(&tmapping); dt.Error != nil {
		err = dt.Error
		return
	}
	users = []User{}
	var uids []int64
	for _, t := range tmapping {
		uids = append(uids, t.Uid)
	}
	//no user bind to team
	if len(uids) == 0 {
		return
	}
	uidstr, err := arrIntToString(uids)
	if err != nil {
		return
	}

	if dt := db.Uic.Select("name, id, cnname").Where(fmt.Sprintf("id in (%s)", uidstr)).Find(&users); dt.Error != nil {
		err = dt.Error
		return
	}
	return
}

func (this Team) GetCreatorName() (userName string, err error) {
	userName = "unknown"
	db := config.Con()
	user := User{ID: this.Creator}
	if dt := db.Uic.Find(&user); dt.Error != nil {
		err = dt.Error
	} else {
		userName = user.Name
	}
	return
}

func arrIntToString(arr []int64) (result string, err error) {
	result = ""
	for indx, a := range arr {
		if indx == 0 {
			result = fmt.Sprintf("%v", a)
		} else {
			result = fmt.Sprintf("%v,%v", result, a)
		}
	}
	if result == "" {
		err = errors.New(fmt.Sprintf("array is empty, err: %v", arr))
	}
	return
}
