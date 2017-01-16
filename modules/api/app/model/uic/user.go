package uic

import (
	con "github.com/masato25/owl_backend/config"
	"github.com/spf13/viper"
)

type User struct {
	ID     int64  `json:"id" `
	Name   string `json:"name"`
	Cnname string `json:"cnname"`
	Passwd string `json:"-"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	IM     string `json:"im" gorm:"column:im"`
	QQ     string `json:"qq" gorm:"column:qq"`
	Role   int    `json:"role"`
}

func skipAccessControll() bool {
	return !viper.GetBool("access_control")
}

func (this User) IsAdmin() bool {
	if skipAccessControll() {
		return true
	}
	if this.Role == 2 || this.Role == 1 {
		return true
	}
	return false
}

func (this User) IsSuperAdmin() bool {
	if skipAccessControll() {
		return true
	}
	if this.Role == 2 {
		return true
	}
	return false
}

func (this User) FindUser() (user User, err error) {
	db := con.Con()
	user = this
	dt := db.Uic.Find(&user)
	if dt.Error != nil {
		err = dt.Error
		return
	}
	return
}

type Session struct {
	ID      int64
	Uid     int64
	Sig     string
	Expired int
}

func (this Session) TableName() string {
	return "session"
}

func (this User) TableName() string {
	return "user"
}
