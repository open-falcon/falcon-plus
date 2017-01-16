package falcon_portal

import (
	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/open-falcon/modules/api/app/model/uic"
	con "github.com/open-falcon/open-falcon/modules/api/config"
)

type Template struct {
	ID         int64  `json:"id" gorm:"column:id"`
	Name       string `json:"tpl_name" gorm:"column:tpl_name"`
	ParentID   int64  `json:"parent_id" orm:"column:parent_id"`
	ActionID   int64  `json:"action_id" orm:"column:action_id"`
	CreateUser string `json:"create_user" orm:"column:create_user"`
}

func (this Template) TableName() string {
	return "tpl"
}

func (this Template) FindUserName() (name string, err error) {
	var user uic.User
	user.Name = this.CreateUser
	db := con.Con()
	dt := db.Uic.Find(&user)
	if dt.Error != nil {
		err = dt.Error
		return
	}
	name = user.Name
	return
}

func (this Template) FindParentName() (name string, err error) {
	var ptpl Template
	if this.ParentID == 0 {
		return
	}
	ptpl.ID = this.ParentID
	db := con.Con()
	dt := db.Falcon.Find(&ptpl)
	if dt.Error != nil {
		log.Debugf("tpl_id: %v find parent: %v with error: %s", this.ID, ptpl.ID, dt.Error.Error())
		return
	}
	name = ptpl.Name
	return
}
